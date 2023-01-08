package commands

import (
	"bufio"
	"context"
	"log"
	"os"
	"time"

	pb "github.com/varunjain0606/grpc-simple-chat/protos"
	"github.com/varunjain0606/grpc-simple-chat/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ReceiveMessage struct {}

func send(client *types.Client, comm string, user *pb.User, grp, friend string) {
	done := make(chan int)
	back := false
	client.Wg.Add(1)
	go func() {
		defer client.Wg.Done()
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			if scanner.Text() == "exit" || scanner.Text() == "commands" {
				back = true
				close(done)
				break
			}
			msg := &pb.Message{
				User:        user,
				Content:   scanner.Text(),
				Group: grp,
				Timestamp: time.Now().String(),
				To: friend,
			}
			_, err := client.BroadcastMessage(context.Background(), &pb.Type{
				Type: comm,
				Message: msg,
			})
			if err != nil {
				if s, ok := status.FromError(err); ok && s.Code() == codes.Unavailable {
					log.Fatalf("Unable to communicate with server. Exiting")
				}
				log.Println(err)
				break
			}
		}
	}()

	if back {
		HandleCommands()
	}
	go func() {
		client.Wg.Wait()
		if !back {
			close(done)
		}
	}()

	<-done
}

func (c *ReceiveMessage) Execute(client *types.Client, args []string) {
	user := &pb.User{
		Id:   ID,
		Name: UserName,
	}
	if GroupName == "" || GroupName == "default" {
		send(client, One2One , user, "", FriendName)
	} else {
		send(client, Group , user, GroupName, "")
	}
}
