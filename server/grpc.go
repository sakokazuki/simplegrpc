package server

import (
	"context"
	"encoding/json"
	"io"

	"github.com/sakokazuki/simplegrpc/config"
	"github.com/sakokazuki/simplegrpc/event"
	"github.com/sakokazuki/simplegrpc/manager"
	pb "github.com/sakokazuki/simplegrpc/protobuf"
	"github.com/sakokazuki/simplegrpc/pubsub"
	"google.golang.org/grpc/codes"

	"github.com/pkg/errors"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

var (
	jsonDBFile = "config/route_guid_db.json"
)

type GRPCServer struct {
	*grpc.Server
	config config.Config
}

type refreshEvents struct {
	client *manager.Client
	events []string
}

type StreamServer struct {
	clientManager *manager.ClientManager
	newClients    chan manager.Client
	removeClients chan manager.Client
	payloads      chan event.Payload
	refreshEvents chan refreshEvents
	pubsub        pubsub.PubSuber
}

// NewGRPCServer setup
func NewGRPCServer(opt Option) (*GRPCServer, error) {
	gs := &GRPCServer{
		config: opt.Config,
	}

	var opts []grpc.ServerOption

	gs.Server = grpc.NewServer(opts...)

	ss, err := NewStreamServer(opt)

	if err != nil {
		return nil, errors.Wrap(err, "failed to NewStereamServer")
	}
	pb.RegisterStreamServiceServer(gs.Server, ss)

	return gs, nil

}

func NewStreamServer(opt Option) (*StreamServer, error) {
	ss := &StreamServer{
		clientManager: manager.NewClientManager(),
		newClients:    make(chan manager.Client, 20),
		removeClients: make(chan manager.Client, 20),
		payloads:      make(chan event.Payload, 20),
		refreshEvents: make(chan refreshEvents, 20),
		pubsub:        opt.PubSuber,
	}

	if err := ss.pubsub.Subscribe(func(payload event.Payload) {
		ss.payloads <- payload
	}); err != nil {
		return nil, errors.Wrap(err, "failed to subscribe")
	}

	ss.Run()
	return ss, nil
}

func (ss *StreamServer) Run() {
	go func() {
		for {
			select {
			case client := <-ss.newClients:
				ss.clientManager.AddClient(client)
			case client := <-ss.removeClients:
				ss.clientManager.RemoveClient(client)
			case payload := <-ss.payloads:
				ss.clientManager.SendPayload(payload)
			case re := <-ss.refreshEvents:
				ss.clientManager.DeleteEvents(re.client)
				re.client.SetEvents(re.events)
				ss.clientManager.AddClient(*re.client)
			}
		}
	}()
}

func (ss *StreamServer) Events(es pb.StreamService_EventsServer) error {
	log.Info().Msg("Client Connection Start")
	client := manager.NewClient([]string{})
	ss.newClients <- client
	defer func() {
		log.Info().Msg("Client Connection End")
		ss.removeClients <- client
	}()

	//removeClientされたときにchan event.Payloadもcloseするのでそのときにgoroutine抜ける
	go func() {
		//clientのpayloadへの書き込みを待ち受けている。書き込みはmanagerのSendPayload()で行われる
		for pl := range client.ReceivePayload() {
			eventType := pb.EventType{Type: pl.Meta.Type}
			p := &pb.Payload{
				EventType: &eventType,
				Data:      string(pl.Data),
			}
			if err := es.Send(p); err != nil {
				log.Fatal().Err(err).Msg("failed to send message")
			} else {
			}
		}
	}()

	for {
		request, err := es.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			if grpc.Code(err) != codes.Canceled {
				return errors.Wrap(err, "Recv error")
			}
			return nil
		}

		if request.ForceClose {
			return nil
		}

		var events []string
		if request.Events == nil {
			events = make([]string, 0)
		} else {
			l := len(request.Events)
			events = make([]string, l)
			for i := 0; i < l; i++ {
				events[i] = request.Events[i].Type
			}
		}

		ss.refreshEvents <- refreshEvents{
			client: &client,
			events: events,
		}

	}
}

func (ss *StreamServer) Publish(ctx context.Context, j *pb.Json) (*pb.Success, error) {

	pubsub := ss.pubsub
	var pl event.Payload
	if err := json.Unmarshal([]byte(j.Data), &pl); err != nil {
		log.Error().Err(err).Msg("failed create payloads")
	}
	pubsub.Publish(pl)

	return &pb.Success{IsSuccess: true}, nil
}
