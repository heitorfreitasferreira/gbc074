package server

import (
	"context"
	"log"

	"github.com/rpc-mqtt-library-manager/cms-server/api"
	"github.com/rpc-mqtt-library-manager/cms-server/internal/database"
)

type Server struct {
	br_ufu_facom_gbc074_projeto_biblioteca.UnimplementedPortalBibliotecaServer

	userRepo database.UserRepo
	bookRepo database.BookRepo
}

func NewServer(userRepo database.UserRepo, bookRepo database.BookRepo) *Server {
	return &Server{
		userRepo: userRepo,
		bookRepo: bookRepo,
	}
}

func (s *Server) RealizaEmprestimo(stream br_ufu_facom_gbc074_projeto_biblioteca.PortalBiblioteca_RealizaEmprestimoServer) error {
	log.Println("cms-server/realiza-emprestimo")
	return nil
}

func (s *Server) RealizaDevolucao(stream br_ufu_facom_gbc074_projeto_biblioteca.PortalBiblioteca_RealizaDevolucaoServer) error {
	log.Println("cms-server/realiza-devolucao")
	return nil
}

func (s *Server) BloqueiaUsuarios(ctx context.Context, req *br_ufu_facom_gbc074_projeto_biblioteca.Vazia) (*br_ufu_facom_gbc074_projeto_biblioteca.Status, error) {
	log.Println("cms-server/bloqueia-usuarios")
	return nil, nil
}

func (s *Server) LiberaUsuarios(ctx context.Context, req *br_ufu_facom_gbc074_projeto_biblioteca.Vazia) (*br_ufu_facom_gbc074_projeto_biblioteca.Status, error) {
	log.Println("cms-server/libera-usuarios")
	return nil, nil
}

func (s *Server) ListaUsuariosBloqueados(req *br_ufu_facom_gbc074_projeto_biblioteca.Vazia, stream br_ufu_facom_gbc074_projeto_biblioteca.PortalBiblioteca_ListaUsuariosBloqueadosServer) error {
	log.Println("cms-server/lista-usuarios-bloqueados")
	return nil
}

func (s *Server) ListaLivrosEmprestados(req *br_ufu_facom_gbc074_projeto_biblioteca.Vazia, stream br_ufu_facom_gbc074_projeto_biblioteca.PortalBiblioteca_ListaLivrosEmprestadosServer) error {
	log.Println("cms-server/lista-livors-emprestados")
	return nil
}

func (s *Server) ListaLivrosEmFalta(req *br_ufu_facom_gbc074_projeto_biblioteca.Vazia, stream br_ufu_facom_gbc074_projeto_biblioteca.PortalBiblioteca_ListaLivrosEmFaltaServer) error {
	log.Println("cms-server/lista-livros-em-falta")
	return nil
}

func (s *Server) PesquisaLivro(req *br_ufu_facom_gbc074_projeto_biblioteca.Criterio, stream br_ufu_facom_gbc074_projeto_biblioteca.PortalBiblioteca_PesquisaLivroServer) error {
	log.Println("cms-server/pesquisa-livro:", req.Criterio)
	return nil
}
