package main

import (
	"fmt"

	"github.com/rpc-mqtt-library-manager/crud-terminal-client/cmd/transition"
	"github.com/rpc-mqtt-library-manager/crud-terminal-client/cmd/types"
	"github.com/rpc-mqtt-library-manager/crud-terminal-client/cmd/view"
)

func main() {
	var state types.State
	for {
		view.View(state)
		state = transition.Transition(state)
		fmt.Printf("--------------------------------\n\n")
	}
}
