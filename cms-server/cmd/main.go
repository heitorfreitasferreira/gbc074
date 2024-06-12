package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	br_ufu_facom_gbc074_projeto_biblioteca "github.com/rpc-mqtt-library-manager/cms-server/api"
	"google.golang.org/grpc"
)

type Server struct {
	br_ufu_facom_gbc074_projeto_biblioteca.UnimplementedPortalBibliotecaServer
}

func (s *Server) RealizaEmprestimo(stream br_ufu_facom_gbc074_projeto_biblioteca.PortalBiblioteca_RealizaEmprestimoServer) error {
	return nil
}

func (s *Server) RealizaDevolucao(stream br_ufu_facom_gbc074_projeto_biblioteca.PortalBiblioteca_RealizaDevolucaoServer) error {
	return nil
}

func (s *Server) BloqueiaUsuarios(ctx context.Context, req *br_ufu_facom_gbc074_projeto_biblioteca.Vazia) (*br_ufu_facom_gbc074_projeto_biblioteca.Status, error) {
	return nil, nil
}

func (s *Server) LiberaUsuarios(ctx context.Context, req *br_ufu_facom_gbc074_projeto_biblioteca.Vazia) (*br_ufu_facom_gbc074_projeto_biblioteca.Status, error) {
	return nil, nil
}

func (s *Server) ListaUsuariosBloqueados(req *br_ufu_facom_gbc074_projeto_biblioteca.Vazia, stream br_ufu_facom_gbc074_projeto_biblioteca.PortalBiblioteca_ListaUsuariosBloqueadosServer) error {
	return nil
}

func (s *Server) ListaLivrosEmprestados(req *br_ufu_facom_gbc074_projeto_biblioteca.Vazia, stream br_ufu_facom_gbc074_projeto_biblioteca.PortalBiblioteca_ListaLivrosEmprestadosServer) error {
	return nil
}

func (s *Server) ListaLivrosEmFalta(req *br_ufu_facom_gbc074_projeto_biblioteca.Vazia, stream br_ufu_facom_gbc074_projeto_biblioteca.PortalBiblioteca_ListaLivrosEmFaltaServer) error {
	return nil
}

func (s *Server) PesquisaLivro(req *br_ufu_facom_gbc074_projeto_biblioteca.Criterio, stream br_ufu_facom_gbc074_projeto_biblioteca.PortalBiblioteca_PesquisaLivroServer) error {
	return nil
}

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

	// Create gRPC server instance
	server := grpc.NewServer()

	// Register server handlers
	br_ufu_facom_gbc074_projeto_biblioteca.RegisterPortalBibliotecaServer(server, &Server{})

	log.Printf("Server listening at %v\n", list.Addr())
	// Make server run to be used by gRPC
	err = server.Serve(list)
	if err != nil {
		log.Fatal(err)
	}

}
