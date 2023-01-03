package commands

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	pb "github.com/varunjain0606/grpc-simple-chat/protos"
	"github.com/varunjain0606/grpc-simple-chat/types"
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
				fmt.Printf("Error Sending Message: %v", err)
				break
			}
		}
	}()

	if back {
		HandleInterativeCommands()
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
	fmt.Println("Press enter to connect and start messaging")
	if GroupName == "" || GroupName == "default" {
		send(client, One2One , user, "", FriendName)
	} else {
		send(client, Group , user, GroupName, "")
	}
}
