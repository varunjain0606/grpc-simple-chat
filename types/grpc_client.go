package types

import (
	"context"
	"log"
	"sync"

	pb "github.com/varunjain0606/grpc-simple-chat/protos"
	"google.golang.org/grpc"
)

var client *Client
var once sync.Once
var wg *sync.WaitGroup

type Client struct {
	pb.BroadcastClient
	context.Context
	Wg *sync.WaitGroup
}

func GetGrpcClient() *Client {
	wg = &sync.WaitGroup{}
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Couldnt connect to service: %v", err)
	}

	cl := pb.NewBroadcastClient(conn)
	once.Do(func() {
		client = &Client{cl, context.Background(), wg}
	})
	return client
}