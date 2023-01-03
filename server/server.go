package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	pb "github.com/varunjain0606/grpc-simple-chat/protos"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	glog "google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

var grpcLog glog.LoggerV2

var channels map[string][]string
var users map[string]string

func init() {
	grpcLog = glog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout)
}

type GroupConnection struct {
	stream pb.Broadcast_CreateStreamServer
	id     string
	name string
	active bool
	group string
	error  chan error
}

type One2OneConnection struct {
	stream pb.Broadcast_CreateStreamServer
	id     string
	name string
	active bool
	user string
	error  chan error
}

type Server struct {
	Group map[string][]*GroupConnection
	One2One map[string][]*One2OneConnection
	pb.UnimplementedBroadcastServer
}

func (s *Server) Login(_ context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	timestamp := time.Now()
	id := sha256.Sum256([]byte(timestamp.String() + req.Name))
	token := hex.EncodeToString(id[:])
	//users[token] = req.Name
	s.setName(token, req.Name)
	return &pb.LoginResponse{Id: token}, nil
}

func (s *Server) Logout(_ context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	ok := s.delName(req.Id)
	if !ok {
		return nil, status.Error(codes.NotFound, "id not found")
	}
	return new(pb.LogoutResponse), nil
}

func (s *Server) setName(id string, name string) {
	s.One2One[id] = []*One2OneConnection{}
	users[id] = name
}

func (s *Server) delName(id string) (ok bool) {
	delete(users, id)
	delete(s.One2One, id)
	return true
}

func (s *Server) ListChannels(ctx context.Context, q *pb.ItemQuery) (*pb.ItemResponse, error) {
	//fmt.Println("Groups:")
	var res pb.ItemResponse
	if q.Type == "groups" {
		for k := range channels {
			res.Items = append(res.Items, &pb.Item{
				Name: k,
			})
		}
	} else if q.Type == "users" {
		for k := range users {
			res.Items = append(res.Items, &pb.Item{
				Name: users[k],
			})
		}
	} else {
		return nil, errors.New("invalid type")
	}
	return &res, nil
}

func remove(s []string, i int) []string {
    s[i] = s[len(s)-1]
    return s[:len(s)-1]
}

func (s *Server) CreateGroup(ctx context.Context, pgrp *pb.Group) (*pb.Close, error) {
	_, ok := channels[pgrp.Group]
	// If the key exists
	if ok {
		fmt.Println("Group exists. You may connect to it")
	} else {
		fmt.Printf("Creating new group %s\n", pgrp.Group)
		channels[pgrp.Group] = []string{pgrp.User.Id}
		s.Group[pgrp.Group] = []*GroupConnection{}
	}
	return &pb.Close{}, nil
}

func (s *Server) LeaveGroup(ctx context.Context, pgrp *pb.Group) (*pb.Close, error) {
	val, ok := channels[pgrp.Group]
    // If the key exists
    if ok {
        fmt.Println("Group exists. Checking if user exists in group")
		for k , v := range val {
			if v == pgrp.User.Id {
				val = remove(val, k)
			}
		}
		channels[pgrp.Group] = val
    }
	// to do. remove also the group connection
	return &pb.Close{}, nil
}

func (s *Server) CreateStream(pconn *pb.Connect, stream pb.Broadcast_CreateStreamServer) error {
	if pconn.Type == "one2one" {
		reverseConn := false
		// check if connection already exists. If not then add. From - to and to-from
		if len(s.One2One[pconn.User.Id]) == 0 {
			reverseConn = true
		}
		for k, c := range s.One2One {
			if k == pconn.User.Id {
				for _, v := range c {
					if v.name == pconn.User.Name && v.user == pconn.Recipient{
						log.Println("Connection already exists")
						return nil
					}
					if v.name == pconn.Recipient && v.user != pconn.User.Name {
						log.Println("Reverse connection doesn't exist. Will add")
						reverseConn = true
					}
				}
            }
        }
		conn := &One2OneConnection{
			stream: stream,
			id:     pconn.User.Id,
			name:   pconn.User.Name,
			active: true,
			error:  make(chan error),
			user: pconn.Recipient,
		}
		c := s.One2One[pconn.User.Id]
		c = append(c, conn)
		s.One2One[pconn.User.Id] = c
		users[pconn.User.Id] = pconn.User.Name

		log.Printf("Creating stream between %s and %s", pconn.User.Name, pconn.Recipient)
		if reverseConn {
			conn := &One2OneConnection{
				stream: stream,
				id:     pconn.User.Id,
				name:   pconn.Recipient,
				active: true,
				error:  make(chan error),
				user: pconn.User.Name,
			}
			log.Printf("Creating stream between %s and %s", pconn.Recipient, pconn.User.Name)
			c = s.One2One[pconn.User.Id]
			c = append(c, conn)
			s.One2One[pconn.User.Id] = c
		}
		return <-conn.error
	} else if pconn.Type == "group" {
		_, ok := s.Group[pconn.Group] 
		if !ok {
			log.Println("Group does not exist. Please create it")
			return errors.New("group does not exist")
		}
		conn := &GroupConnection{
			stream: stream,
			id:     pconn.User.Id,
			name:   pconn.User.Name,
			active: true,
			error:  make(chan error),
			group: pconn.Group,
		}
		c := s.Group[pconn.Group]
		c = append(c, conn)
		s.Group[pconn.Group] = c
		if pconn.Group != "" && pconn.Group != "default" {
			val, ok := channels[pconn.Group]
			// If the key exists
			if ok {
				fmt.Printf("Group exists. Connecting user :%s to %s\n",pconn.User.Name, pconn.Group)
				val = append(val, pconn.User.Id)
				channels[pconn.Group] = val
			} else {
				fmt.Printf("Creating new group %s", pconn.Group)
				channels[pconn.Group] = []string{pconn.User.Id}
			}
		}
		return <-conn.error
	}
	return nil
}

func (s *Server) BroadcastMessage(ctx context.Context, msg *pb.Type) (*pb.Close, error) {
	wait := sync.WaitGroup{}
	done := make(chan int)
	if msg.Type == "one2one" {
		for _, o := range s.One2One {
			for _, conn := range o {
				if conn.name != msg.Message.To && conn.user != msg.Message.User.Name {
					continue
				}
				wait.Add(1)
				go func(msg *pb.Type, conn *One2OneConnection) {
					defer wait.Done()
					if conn.active {
						err := conn.stream.Send(msg.Message)
						if err != nil {
							grpcLog.Errorf("Error with Stream: %v - Error: %v", conn.stream, err)
							conn.active = false
							conn.error <- err
						}
						grpcLog.Info("Sending message from:", conn.name, " to: ", conn.user)
					}
				}(msg, conn)
			}
		}
		go func() {
			wait.Wait()
			close(done)
		}()
	} else if msg.Type == "group" {
		for k, o := range s.Group {
			if k == msg.Message.Group {
				for _, conn := range o {
						wait.Add(1)
						go func(msg *pb.Type, conn *GroupConnection) {
							defer wait.Done()
							if conn.active {
								err := conn.stream.Send(msg.Message)
								if err != nil {
									grpcLog.Errorf("Error with Stream: %v - Error: %v", conn.stream, err)
									conn.active = false
									conn.error <- err
								}
								grpcLog.Info("Sending message from:", conn.name, " to: ", conn.group)
							}
					}(msg, conn)
				}
			}
		}
		go func() {
			wait.Wait()
			close(done)
		}()
	}

	<-done
	return &pb.Close{}, nil
}

func main() {
	groupConnections := make(map[string][]*GroupConnection)
	one2oneConnections :=  make(map[string][]*One2OneConnection)

	channels = make(map[string][]string)
	users = make(map[string]string)

	server := &Server{groupConnections, one2oneConnections, pb.UnimplementedBroadcastServer{}}

	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("error creating the server %v", err)
	}

	grpcLog.Info("Starting server at port :8080")

	pb.RegisterBroadcastServer(grpcServer, server)
	grpcServer.Serve(listener)
}
