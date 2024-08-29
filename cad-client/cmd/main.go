package main

import (
	"flag"
	"fmt"
	"library-manager/cad-client/api"
	"library-manager/cad-client/internal/handlers"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	port := flag.String("port", "50051", "Port to listen on")

	flag.Parse()

	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%s", *port), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Não foi possível conectar: %v", err)
		panic(err)
	}

	defer conn.Close()
	client := api.NewPortalCadastroClient(conn)
	for {
		handler := handlers.Choose()

		if err := handler(client); err != nil {
			log.Fatalf("Erro ao executar handler: %v", err)
		}
	}
}
