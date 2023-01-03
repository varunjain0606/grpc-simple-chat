package commands

import (
	"context"
	"fmt"
	"log"
	"strings"

	pb "github.com/varunjain0606/grpc-simple-chat/protos"
	"github.com/varunjain0606/grpc-simple-chat/types"
)

type Login struct {}

var ID string
var UserName string
var Islogged bool

func login(client *types.Client, ctx context.Context, name string) (string, error) {
	res, err := client.Login(ctx, &pb.LoginRequest{
		Name:     name,
	})
	if err != nil {
		return "", err
	}
	return res.Id, nil
}

func (c *Login) Execute(client *types.Client, args []string) {
	if len(args) < 1 || len(args) > 1{
		fmt.Println("\nRequired only username")
		return
	}
	s := strings.TrimSpace(args[0])
	if len(s) == 0 {
		fmt.Println("\nInvalid user name")
		return
	}
	var err error
	ID, err = login(client, client.Context, args[0])
	if err != nil {
		log.Fatalf("\nLogin failed: %v", err)
	}
	UserName = args[0]
	Islogged = true
	fmt.Println("\nUser logged in successfully")
}
