package main

import (
	"flag"

	"library-manager/cad-client/cmd/dfa.go"
	"library-manager/cad-client/cmd/transition"
	"library-manager/cad-client/cmd/view"
	"library-manager/cad-client/utils"
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
