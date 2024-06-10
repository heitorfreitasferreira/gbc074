package main

import (
	"flag"
	"fmt"

	"github.com/rpc-mqtt-library-manager/crud-terminal-client/cmd/transition"
	"github.com/rpc-mqtt-library-manager/crud-terminal-client/cmd/types"
	"github.com/rpc-mqtt-library-manager/crud-terminal-client/cmd/view"
	"github.com/rpc-mqtt-library-manager/crud-terminal-client/utils"
)

func main() {
	var state types.State
	port := flag.String("port", "50051", "Port to listen on")

	flag.Parse()

	utils.Port = *port
	for {
		view.View(state)
		state = transition.Transition(state)
		fmt.Printf("--------------------------------\n\n")
	}
}
