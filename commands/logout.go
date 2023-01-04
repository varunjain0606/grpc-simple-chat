package commands

import (
	"context"
	"fmt"
	"log"

	pb "github.com/varunjain0606/grpc-simple-chat/protos"
	"github.com/varunjain0606/grpc-simple-chat/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Logout struct {}

func logout(client *types.Client, ctx context.Context) error {
	user := &pb.User{
		Id:   ID,
		Name: UserName,
	}
	_, err := client.LeaveGroup(client.Context, &pb.Group{User: user, Group: GroupName})
	if err != nil {
		fmt.Printf("\nunable to leave group. %s\n", err.Error())
	}
	_, err = client.Logout(ctx, &pb.LogoutRequest{Id: ID})
	if s, ok := status.FromError(err); ok && s.Code() == codes.Unavailable {
		log.Println("unable to logout (connection already closed)")
		return nil
	}
	Islogged = false
	UserName = ""
	return err
}

func (c *Logout) Execute(client *types.Client, args []string) {
	if !Islogged {
		fmt.Println("\nUser not logged in")
        return
	}
	err := logout(client, client.Context)
	if err != nil {
		log.Fatalf("Logout failed: %v", err)
	}
	fmt.Println("\nuser logged out successfully")
}