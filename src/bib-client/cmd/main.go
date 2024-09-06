package main

import (
	"bufio"
	"flag"
	"fmt"
	"library-manager/bib-client/internal/handlers"
	"log"
	"os"

	api "library-manager/shared/api/bib"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	host := flag.String("host", "localhost", "Host to connect to")
	port := flag.String("port", "50052", "Port to listen on")

	flag.Parse()

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", *host, *port), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Não foi possível conectar: %v", err)
		panic(err)
	}

	defer conn.Close()
	client := api.NewPortalBibliotecaClient(conn)
	for {
		handler := handlers.Choose(*host, *port)

		if err := handler(client); err != nil {
			fmt.Printf("\nErro ao executar handler: %v", err)
			fmt.Println("Pressione ENTER...")

			bufio.NewReader(os.Stdin).ReadBytes('\n')
		}
	}
}
