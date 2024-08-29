package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"library-manager/cad-server/api"
	"library-manager/cad-server/internal/database"
	"library-manager/cad-server/internal/queue"
	"library-manager/cad-server/internal/server"

	"google.golang.org/grpc"
)

func main() {
	port := flag.String("port", "50051", "Port to listen on")
	host := flag.String("host", "127.0.0.1", "Host to listen on")

	flag.Parse()
	ch := make(chan os.Signal, 1)
	// Receive notifications on chanel if receive os.Interrupt or syscall.SIGTERM.
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	mqttClient := queue.GetMqttBroker(*host, 1883, os.Getpid())
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("error connecting to MQTT broker: %v", token.Error())
	}

	// O defer garante que esse comando seja executado quando a função sair da pilha.
	defer mqttClient.Disconnect(250)

	list, err := net.Listen("tcp", fmt.Sprintf("%s:%s", *host, *port))
	if err != nil {
		log.Fatalf("error opening port %s: %v", *port, err)
	}
	s := grpc.NewServer()

	api.RegisterPortalCadastroServer(s, server.NewServer(database.ConcreteUserRepo, mqttClient, database.ConcreteBookRepo))

	go func() {
		<-ch
		log.Println("Gracefully stopping the server...")
		s.GracefulStop()
	}()

	log.Printf("Server listening at %v\n", list.Addr())
	if err := s.Serve(list); err != nil {
		log.Fatal(err)
	}
}
