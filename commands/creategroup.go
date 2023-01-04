package commands

import (
	"fmt"
	"strings"

	pb "github.com/varunjain0606/grpc-simple-chat/protos"
	"github.com/varunjain0606/grpc-simple-chat/types"
)

type CreateGroup struct {}

func (c *CreateGroup) Execute(client *types.Client, args []string) {
	if !Islogged {
		fmt.Println("\nMust be logged in to create a group")
        return
	}
	if len(args) != 1 {
		fmt.Println("\nRequired only group name")
		return
	}
	user := &pb.User{
		Id:   ID,
		Name: UserName,
	}
	s := strings.TrimSpace(args[0])
	if len(s) == 0 {
		fmt.Println("\nInvalid group name")
		return
	}

	_, err := client.CreateGroup(client.Context, &pb.Group{User: user, Group: args[0]})
	if err != nil {
		fmt.Printf("\nunable to create group: %v\n", err)
	}
	fmt.Printf("\nCreated group: %v\n", args[0])

	err = connect(client, Group, user, args[0], "")
	if err!= nil {
        fmt.Println(err)
	}
}
