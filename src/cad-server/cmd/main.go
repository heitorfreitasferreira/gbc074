package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"library-manager/cad-server/internal/server"
	api_cad "library-manager/shared/api/cad"

	"google.golang.org/grpc"
)

func main() {
	port := flag.String("port", "50052", "Port to listen on")
	host := flag.String("host", "127.0.0.1", "Host to listen on")
	databaseAddr := flag.String("replica", "http://localhost:21000", "Address of the database server")

	flag.Parse()
	ch := make(chan os.Signal, 1)
	// Receive notifications on chanel if receive os.Interrupt or syscall.SIGTERM.
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	list, err := net.Listen("tcp", fmt.Sprintf("%s:%s", *host, *port))
	if err != nil {
		log.Fatalf("error opening port %s: %v", *port, err)
	}
	s := grpc.NewServer()

	api_cad.RegisterPortalCadastroServer(s, server.NewServer(*databaseAddr))

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
