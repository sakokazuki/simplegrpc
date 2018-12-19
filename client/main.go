package main

import (
	"context"
	"fmt"
	"log"
	pb "github.com/sakokazuki/simplegrpc/protobuf"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

var (
	serverAddr = "127.0.0.1:50151"
)

func eventType(et string) *pb.EventType {
	return &pb.EventType{Type: et}
}

func publish(client pb.StreamServiceClient, data string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	fmt.Printf("publish %s\n", data)
	_, err := client.Publish(ctx, &pb.Json{Data: data})
	if err != nil {
		log.Fatalf("fail to publish: %v", err)
	}
}

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewStreamServiceClient(conn)

	count := 0
	go func() {
		for {
			time.Sleep(time.Second * 5)
			publish(client, `{"meta":{"type":"unity:test"},"data":{"data":`+strconv.Itoa(count)+`}}`)
			count++
		}
	}()

	ctx := context.Background()
	req := pb.Request{
		Events: []*pb.EventType{
			eventType("program:1234:poll"),
			eventType("program:1234:views"),
		},
	}

	ss, err := client.Events(ctx)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := ss.Send(&req); err != nil {
			log.Fatal(err)
		}
	}()

	for {
		resp, err := ss.Recv()
		if err != nil {
			log.Println(err)
			continue
		}
		if resp == nil {
			log.Println("payload is nil")
			continue
		}
		fmt.Printf("Meta: %s\tData: %s\n", resp.EventType.Type, resp.Data)
	}

}
