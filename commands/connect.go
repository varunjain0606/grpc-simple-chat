package commands

import (
	"fmt"
	"strings"

	pb "github.com/varunjain0606/grpc-simple-chat/protos"
	"github.com/varunjain0606/grpc-simple-chat/types"
)

type Connect struct {}

func (c *Connect) Execute(client *types.Client, args []string) {
	if len(args) < 1 {
		fmt.Println("\nRequired friend name")
		return
	}
	s := strings.TrimSpace(args[0])
	if len(s) == 0 {
		fmt.Println("\nInvalid friend name")
		return
	}
	user := &pb.User{
		Id:   ID,
		Name: UserName,
	}
	FriendName = s
	connect(client, One2One, user, "", args[0])
}
