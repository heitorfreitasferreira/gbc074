package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	// Import para os repositórios
	"library-manager/bib-server/internal/server" // Import para o pacote do servidor
	api_bib "library-manager/shared/api/bib"

	"google.golang.org/grpc"
)

func main() {
	port := flag.String("port", "50051", "Port to listen on")
	host := flag.String("host", "127.0.0.1", "Host to listen on")
	// Endereços dos bancos de dados
	bookDatabaseAddr := flag.String("cluster1", "http://localhost:11000", "Address of the book database server")
	userDatabaseAddr := flag.String("cluster0", "http://localhost:13000", "Address of the user database server")

	flag.Parse()
	ch := make(chan os.Signal, 1)
	// Recebe notificações de interrupção e termina graciosamente
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	// Cria um listener TCP na porta especificada
	list, err := net.Listen("tcp", fmt.Sprintf("%s:%s", *host, *port))
	if err != nil {
		log.Fatalf("error opening port %s: %v", *port, err)
	}
	s := grpc.NewServer()

	api_bib.RegisterPortalVBibliotecaServer(s, server.NewServer(*userDatabaseAddr, *bookDatabaseAddr))

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
