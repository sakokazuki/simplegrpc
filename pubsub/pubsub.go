package pubsub

import (
	"github.com/sakokazuki/simplegrpc/event"

	"github.com/mattn/go-pubsub"
	"github.com/rs/zerolog/log"
)

type PubSuber interface {
	Publish(payload event.Payload)
	Subscribe(f func(payload event.Payload)) error
}

type PubSub struct {
	pubsub *pubsub.PubSub
}

func NewPubSub() PubSuber {
	return &PubSub{
		pubsub: pubsub.New(),
	}
}

func (d *PubSub) Publish(payload event.Payload) {
	log.Info().
		Str("data", payload.String()).
		Msg("publish payload")
	d.pubsub.Pub(payload)

}

func (d *PubSub) Subscribe(f func(payload event.Payload)) error {
	return d.pubsub.Sub(f)
}
