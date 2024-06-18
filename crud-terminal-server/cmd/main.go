package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	br_ufu_facom_gbc074_projeto_cadastro "github.com/rpc-mqtt-library-manager/crud-terminal-server/api"
	"github.com/rpc-mqtt-library-manager/crud-terminal-server/internal/database"
	"github.com/rpc-mqtt-library-manager/crud-terminal-server/internal/queue"
	"github.com/rpc-mqtt-library-manager/crud-terminal-server/internal/server"
	"google.golang.org/grpc"
)

func main() {
	port := flag.String("port", "50051", "Port to listen on")
	host := flag.String("host", "127.0.0.1", "Host to listen on")

	flag.Parse()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	mqttClient := queue.GetMqttBroker(*host, 1883, os.Getpid())
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("error connecting to MQTT broker: %v", token.Error())
	}
	defer mqttClient.Disconnect(250)

	list, err := net.Listen("tcp", fmt.Sprintf("%s:%s", *host, *port))
	if err != nil {
		log.Fatalf("error opening port %s: %v", *port, err)
	}
	s := grpc.NewServer()

	br_ufu_facom_gbc074_projeto_cadastro.RegisterPortalCadastroServer(s, server.NewServer(database.ConcreteUserRepo, mqttClient, database.ConcreteBookRepo))

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
