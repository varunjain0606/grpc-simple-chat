package commands

import (
	"fmt"

	"github.com/varunjain0606/grpc-simple-chat/types"
)

type PrintCommands struct {}

func (c *PrintCommands) Execute(client *types.Client, args []string) {
	fmt.Println("----------------------------------------------------------------")
	fmt.Println("\t\t\tAvailable commands")
	fmt.Println("----------------------------------------------------------------")
	for k := range commands {
        fmt.Println("\t\t\t" + k)
    }
}
