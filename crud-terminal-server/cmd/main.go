package main

import (
	"flag"
	"fmt"
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
		panic(fmt.Errorf("erro ao conectar ao broker MQTT: %v", token.Error()))
	}
	defer mqttClient.Disconnect(250)

	list, err := net.Listen("tcp", fmt.Sprintf("%s:%s", *host, *port))
	if err != nil {
		panic(fmt.Errorf("erro ao abrir a porta %s: %v", *port, err))
	}
	s := grpc.NewServer()

	br_ufu_facom_gbc074_projeto_cadastro.RegisterPortalCadastroServer(s, server.NewServer(database.ConcreteUserRepo, mqttClient))

	go func() {
		<-ch
		fmt.Println("Gracefully stopping the server...")
		s.GracefulStop()
	}()

	fmt.Printf("Server listening at %v\n", list.Addr())
	if err := s.Serve(list); err != nil {
		panic(err)
	}
}
