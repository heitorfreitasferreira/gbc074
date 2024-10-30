package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"library-manager/cad-server/internal/database"
	"library-manager/cad-server/internal/server"
	api_cad "library-manager/shared/api/cad"
	sharedDatabase "library-manager/shared/database"

	"google.golang.org/grpc"
)

func main() {
	port := flag.String("port", "50052", "Port to listen on")
	host := flag.String("host", "127.0.0.1", "Host to listen on")

	flag.Parse()
	ch := make(chan os.Signal, 1)
	// Receive notifications on chanel if receive os.Interrupt or syscall.SIGTERM.
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	list, err := net.Listen("tcp", fmt.Sprintf("%s:%s", *host, *port))
	if err != nil {
		log.Fatalf("error opening port %s: %v", *port, err)
	}
	s := grpc.NewServer()
	// TODO Resolver o banco de dados distribuido
	bookRepo, err := database.NewInMemoryBookRepo(sharedDatabase.Cluster1Replica0Path)
	if err != nil {
		log.Fatalf("error creating book repository: %v", err)
	}
	defer bookRepo.Close()

	userRepo, err := database.NewInMemoryUserRepo(sharedDatabase.Cluster0Replica0Path)
	if err != nil {
		log.Fatalf("error creating user repository: %v", err)
	}
	defer userRepo.Close()
	api_cad.RegisterPortalCadastroServer(s, server.NewServer(userRepo, bookRepo))

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
