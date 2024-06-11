package main

import (
	"flag"

	"github.com/rpc-mqtt-library-manager/crud-terminal-client/cmd/dfa.go"
	"github.com/rpc-mqtt-library-manager/crud-terminal-client/cmd/transition"
	"github.com/rpc-mqtt-library-manager/crud-terminal-client/cmd/view"
	"github.com/rpc-mqtt-library-manager/crud-terminal-client/utils"
)

func main() {
	sm := dfa.NewViewPicker()
	port := flag.String("port", "50051", "Port to listen on")

	flag.Parse()

	utils.Port = *port
	for {
		view.View(sm)
		transition.Transition(sm)
	}
}
