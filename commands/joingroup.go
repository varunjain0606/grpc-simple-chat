package commands

import (
	"context"
	"fmt"
	"io"
	"strings"

	pb "github.com/varunjain0606/grpc-simple-chat/protos"
	"github.com/varunjain0606/grpc-simple-chat/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type JoinGroup struct {}

var GroupName string
var FriendName string

func connect(client *types.Client, comm string, user *pb.User, grp, friend string) error {
	// check to see if group exists or friend exists
	var streamerror error
	stream, err := client.CreateStream(context.Background(), &pb.Connect{
		User:   user,
		Active: true,
		Group:  grp,
		Recipient: friend,
		Type: comm,
	})
	if err != nil {
		fmt.Println(err)
	}
	GroupName = grp
	FriendName = friend
	if comm == Group {
		fmt.Printf("Created stream for group:%s\n", grp)
	}
	if comm == One2One {
		fmt.Printf("Created stream successfully with friend %s\n", friend)
	}
	client.Wg.Add(1)
	go func(str pb.Broadcast_CreateStreamClient) {
		defer client.Wg.Done()
		for {
			msg, err := str.Recv()
			if s, ok := status.FromError(err); ok && s.Code() == codes.Canceled {
				streamerror = fmt.Errorf("stream cancelled (shutting down)")
			} else if s, ok := status.FromError(err); ok && s.Code() == codes.NotFound {
				fmt.Println("group does not exist. Press enter to continue")
				streamerror = fmt.Errorf("group does not exist")
				break
			} else if s, ok := status.FromError(err); ok && s.Code() == codes.InvalidArgument {
				fmt.Println("User does not exist. Press enter to continue")
				streamerror = fmt.Errorf("user does not exist")
				break
			} else if err == io.EOF {
				streamerror = fmt.Errorf("stream closed by server")
				break
			} else if err != nil {
				streamerror = fmt.Errorf("error reading message: %v", err)
				break
			}
			if msg.Group == grp {
				fmt.Printf("%v : %s\n", msg.User.Name, msg.Content)
			}
			if comm == One2One && msg.To == user.Name && msg.User.Name != friend{
				fmt.Printf("%v : %s\n", msg.User.Name, msg.Content)
			}
		}
	}(stream)
	return streamerror
}

func (c *JoinGroup) Execute(client *types.Client, args []string) {
	if !Islogged {
		fmt.Println("\nMust be logged in to join a group")
        return
	}
	if len(args) !=  1{
		fmt.Println("\nRequired only group name")
		return
	}
	s := strings.TrimSpace(args[0])
	if len(s) == 0 {
		fmt.Println("\nInvalid group name")
		return
	}
	user := &pb.User{
		Id:   ID,
		Name: UserName,
	}
	err := connect(client, Group, user, args[0], "")
	if err!= nil {
        fmt.Println(err)
	}
}
