package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/varunjain0606/grpc-simple-chat/types"
)

type Command interface {
	Execute(grpcClient *types.Client, args []string)
}

const (
	One2One = "one2one"
	Group = "group"
)

var commands = map[string]Command{
	"help": &PrintCommands{},
	"login": &Login{},
	"logout": &Logout{},
	"create_group": &CreateGroup{},
	"join_group": &JoinGroup{},
	"leave_group": &LeaveGroup{},
	"send_message": &Connect{},
	"list_channels": &ListChannel{},
	"stream": &ReceiveMessage{},
	"exit": &Logout{},
}

func GetCommandByName(name string) Command{
	if command := commands[name]; command == nil {
		fmt.Println("Write a message")
		fmt.Println("Press exit to come to the main menu")
		return commands["stream"]
	} else {
		 return command
	}
}

func HandleInterativeCommands(){
	for {
		handleCommand("help")
		if Islogged {
			fmt.Printf("<%s>$ ", UserName)
		} else {
			fmt.Print("$ ")
		}
		reader := bufio.NewReader(os.Stdin)
		cmdString, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		cmdString = strings.TrimSuffix(cmdString, "\n")
		handleCommand(cmdString)
		if cmdString == "exit" {
			os.Exit(1)
		}
	}
}


func handleCommand( command string  ){
	//println(command)
	client := types.GetGrpcClient()
	commandSplit := strings.Split(command, " ")
	cmd := commandSplit[0]
	GetCommandByName(cmd).Execute(client , commandSplit[1:])
}