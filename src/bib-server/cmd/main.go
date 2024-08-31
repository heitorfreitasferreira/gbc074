package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"library-manager/bib-server/internal/database"
	"library-manager/bib-server/internal/queue"
	"library-manager/bib-server/internal/server"
	"library-manager/shared/api/bib"

	"google.golang.org/grpc"
)

func main() {
	log.Printf("Hello, World! CMS Server")

	// Get CLI arguments
	port := flag.String("port", "50051", "Port to listen on")
	host := flag.String("host", "127.0.0.1", "Host to listen on")
	flag.Parse()

	mqttClient := queue.GetMqttBroker(*host, 1883, os.Getpid())
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("error connecting to MQTT broker: %v", token.Error())
	} else {
		log.Println("Connected to MQTT broker")
	}

	defer mqttClient.Disconnect(250)

	// Create TCP port listener
	list, err := net.Listen("tcp", fmt.Sprintf("%s:%s", *host, *port))
	if err != nil {
		log.Fatalf("error opening port %s: %v", *port, err)
	}

	// Create gRPC grpcServer instance
	grpcServer := grpc.NewServer()

	// Register server handlers
	api_bib.RegisterPortalBibliotecaServer(
		grpcServer,
		server.NewServer(
			database.ConcreteUserRepo,
			database.ConcreteBookRepo,
			database.ConcreteUserBookRepo,
			mqttClient,
		),
	)

	// Create channel to capture SIGERM and gracefully stop grpc server.
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-channel
		log.Println("Gracefully stopping the server...")
		grpcServer.GracefulStop()
	}()

	log.Printf("Server listening at %v\n", list.Addr())
	// Make server run to be used by gRPC
	err = grpcServer.Serve(list)
	if err != nil {
		log.Fatal(err)
	}

}
