package main

import (
	"flag"
	"fmt"
	"net"

	br_ufu_facom_gbc074_projeto_cadastro "github.com/rpc-mqtt-library-manager/crud-terminal-server/api"
	"github.com/rpc-mqtt-library-manager/crud-terminal-server/internal/server"
	"google.golang.org/grpc"
)

func main() {
	port := flag.String("port", "50051", "Port to listen on")

	flag.Parse()

	list, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()

	br_ufu_facom_gbc074_projeto_cadastro.RegisterPortalCadastroServer(s, &server.Server{})

	fmt.Printf("Server listening at %v\n", list.Addr())
	if err := s.Serve(list); err != nil {
		panic(err)
	}
}
