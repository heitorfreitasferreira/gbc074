package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/rpc-mqtt-library-manager/cms-server/api"
	"github.com/rpc-mqtt-library-manager/cms-server/internal/database"
	"github.com/rpc-mqtt-library-manager/cms-server/internal/server"
	"google.golang.org/grpc"
)

func main() {
	log.Printf("Hello, World! CMS Server")

	// Get CLI arguments
	port := flag.String("port", "50051", "Port to listen on")
	host := flag.String("host", "127.0.0.1", "Host to listen on")
	flag.Parse()

	// Create TCP port listener
	list, err := net.Listen("tcp", fmt.Sprintf("%s:%s", *host, *port))
	if err != nil {
		log.Fatalf("error opening port %s: %v", *port, err)
	}

	// Create gRPC grpcServer instance
	grpcServer := grpc.NewServer()

	// Register server handlers
	br_ufu_facom_gbc074_projeto_biblioteca.RegisterPortalBibliotecaServer(
		grpcServer,
		server.NewServer(
			database.ConcreteUserRepo,
			database.ConcreteBookRepo,
		),
	)

	log.Printf("Server listening at %v\n", list.Addr())
	// Make server run to be used by gRPC
	err = grpcServer.Serve(list)
	if err != nil {
		log.Fatal(err)
	}

}
