package main

import (
	"fmt"

	"log"
	"sync"

	"github.com/varunjain0606/grpc-simple-chat/commands"
	pb "github.com/varunjain0606/grpc-simple-chat/protos"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var client pb.BroadcastClient
var wait *sync.WaitGroup

const (
	One2One = "one2one"
	Group = "group"
)

func init() {
	wait = &sync.WaitGroup{}
}

func connect(comm string, user *pb.User, grp, friend string) error {
	var streamerror error
	stream, err := client.CreateStream(context.Background(), &pb.Connect{
		User:   user,
		Active: true,
		Group:  grp,
		Recipient: friend,
		Type: comm,
	})
	if err != nil {
		return fmt.Errorf("connection failed: %v", err)
	}
	if comm == Group {
		fmt.Printf("Created stream successfully. Now listening on group %s\n", grp)
	}
	if comm == One2One {
		fmt.Printf("Created stream successfully with friend %s\n", friend)
	}
	wait.Add(1)
	go func(str pb.Broadcast_CreateStreamClient) {
		defer wait.Done()
		for {
			msg, err := str.Recv()
			if err != nil {
				streamerror = fmt.Errorf("Error reading message: %v", err)
				break
			}
			if msg.Group == grp {
				fmt.Printf("%v : %s\n", msg.User.Name, msg.Content)
			}
			if comm == One2One && msg.To == user.Name && msg.To == friend {
				fmt.Printf("%v : %s\n", msg.User.Name, msg.Content)
			}
		}
	}(stream)
	return streamerror
}

func login(ctx context.Context, name string) (string, error) {
	res, err := client.Login(ctx, &pb.LoginRequest{
		Name:     name,
	})
	if err != nil {
		return "", err
	}
	return res.Id, nil
}

func logout(ctx context.Context, id string) error {
	_, err := client.Logout(ctx, &pb.LogoutRequest{Id: id})
	if s, ok := status.FromError(err); ok && s.Code() == codes.Unavailable {
		log.Println("unable to logout (connection already closed)")
		return nil
	}
	return err
}

func joinGroup(ctx context.Context, user *pb.User, group string) error {
	return connect(Group, user, group, "")
}

func createGroup(ctx context.Context, user *pb.User, group string) error {
	_, err := client.CreateGroup(ctx, &pb.Group{User: user, Group: group})
	return err
}

func leaveGroup(ctx context.Context, user *pb.User, group string) error {
	_, err := client.LeaveGroup(ctx, &pb.Group{User: user, Group: group})
	return err
}

func main() {
   commands.HandleInterativeCommands()
}

// func main() {
// 	done := make(chan int)
// 	name := flag.String("n", "default", "The name of the user")
// 	group := flag.String("g", "default", "Group to connect")
// 	friend := flag.String("f", "default", "Friend to connect")
// 	comm := flag.String("c", "group", "Type of communication")
// 	flag.Parse()

// 	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
// 	if err != nil {
// 		log.Fatalf("Couldnt connect to service: %v", err)
// 	}

// 	client = pb.NewBroadcastClient(conn)

// 	ctx := context.Background()

// 	id, err := login(ctx, *name)
// 	if err != nil {
// 		log.Fatalf("Login failed: %v", err)
// 	}

// 	log.Println("User logged in")

// 	user := &pb.User{
// 		Id:   id,
// 		Name: *name,
// 	}

// 	if *group == "" || *group == "default" {
// 		fmt.Println("Invalid group")
// 	} else {
// 		err = createGroup(ctx, user, *group)
// 		if err != nil {
// 			log.Fatalf("Create group failed. Maybe already exists: %v", err)
// 		} else {
// 			fmt.Printf("Created group: %s\n", *group)
// 		}
// 	}

// 	if *comm == Group {
// 		err = joinGroup(ctx, user, *group)
// 		if err!= nil {
// 			log.Fatalf("Join group failed: %v", err)
// 		}
// 	} else {
// 		err = connect(One2One, user, "", *friend)
// 		if err != nil {
// 			log.Fatalf("Couldnt connect to friend: %v", err)
// 		}
// 		fmt.Printf("Connected to friend: %s\n", *friend)
// 	}
// 	wait.Add(1)
// 	go func() {
// 		defer wait.Done()
// 		scanner := bufio.NewScanner(os.Stdin)
// 		for scanner.Scan() {
// 			msg := &pb.Message{
// 				User:        user,
// 				Content:   scanner.Text(),
// 				Group: *group,
// 				Timestamp: time.Now().String(),
// 				To: *friend,
// 			}
// 			_, err := client.BroadcastMessage(context.Background(), &pb.Type{
// 				Type: Group,
// 				Message: msg,
// 			})
// 			if err != nil {
// 				fmt.Printf("Error Sending Message: %v", err)
// 				break
// 			}
// 		}
// 	}()

// 	go func() {
// 		wait.Wait()
// 		close(done)
// 	}()

// 	<-done
// }
