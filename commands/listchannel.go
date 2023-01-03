package commands

import (
	"fmt"

	pb "github.com/varunjain0606/grpc-simple-chat/protos"
	"github.com/varunjain0606/grpc-simple-chat/types"
)

type ListChannel struct {}

func (c *ListChannel) Execute(client *types.Client, args []string) {
	val, err := client.ListChannels(client.Context, &pb.ItemQuery{
		Type: "groups",
	})
	if err != nil {
		fmt.Printf("Could not list channels: %v", err)
	}
	fmt.Println("Groups:")
	for _, r := range val.Items {
		fmt.Println("\t" + r.Name)
	}
	val, err = client.ListChannels(client.Context, &pb.ItemQuery{
		Type: "users",
	})
	if err != nil {
		fmt.Printf("Could not list channels: %v", err)
	}
	fmt.Println("Users:")
	for _, r := range val.Items {
		fmt.Println("\t" + r.Name)
	}
}
