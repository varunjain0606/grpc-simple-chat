package commands

import (
	"fmt"
	"strings"

	pb "github.com/varunjain0606/grpc-simple-chat/protos"
	"github.com/varunjain0606/grpc-simple-chat/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LeaveGroup struct {}

func (c *LeaveGroup) Execute(client *types.Client, args []string) {
	if !Islogged {
		fmt.Println("\nMust be logged in to leave a group")
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

	_, err := client.LeaveGroup(client.Context, &pb.Group{User: user, Group: args[0]})
	if err != nil {
		if s, ok := status.FromError(err); ok && s.Code() == codes.NotFound {
			fmt.Printf("\nunable to leave group. %s\n", err.Error())
			return
		}
	}
	fmt.Printf("\nUser %s has left the group %s\n", UserName, args[0])
}
