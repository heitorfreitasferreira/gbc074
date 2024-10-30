package server

import (
	"context"
	"encoding/json"
	"log"

	"library-manager/cad-server/internal/database"
	"library-manager/cad-server/internal/utils"
	api_cad "library-manager/shared/api/cad"
)

var qos byte = 2

type Server struct {
	api_cad.UnimplementedPortalCadastroServer

	userRepo database.UserRepo
	bookRepo database.BookRepo
}

func NewServer(userRepo database.UserRepo, bookRepo database.BookRepo) *Server {
	return &Server{
		userRepo: userRepo,
		bookRepo: bookRepo,
	}
}

func (s *Server) NovoUsuario(ctx context.Context, usuario *api_cad.Usuario) (*api_cad.Status, error) {
	cpf := utils.CPF(usuario.Cpf)
	if !cpf.Validate() {
		log.Printf("CPF inválido")
		return &api_cad.Status{Status: 1, Msg: "CPF inválido"}, nil
	}

	jsonData, err := json.Marshal(usuario)
	if err != nil {
		log.Printf("Erro ao converter dados do usuário para JSON: %v", err)
		return &api_cad.Status{Status: 1, Msg: "Erro ao converter dados para JSON"}, nil
	}
	// TODO: Salvar usuário no banco de dados
	_ = jsonData
	return &api_cad.Status{Status: 0}, nil
}

func (s *Server) ObtemUsuario(ctx context.Context, request *api_cad.Identificador) (*api_cad.Usuario, error) {
	user, err := s.userRepo.ObtemUsuario(utils.CPF(request.Id))

	if err != nil {
		return nil, err
	}

	return &api_cad.Usuario{
		Cpf:  request.Id,
		Nome: user.Nome,
	}, nil
}

func (s *Server) ObtemTodosUsuarios(request *api_cad.Vazia, stream api_cad.PortalCadastro_ObtemTodosUsuariosServer) error {
	users, err := s.userRepo.ObtemTodosUsuarios()
	if err != nil {
		return err
	}
	for _, user := range users {
		protoUser := database.UserToProto(user)
		err := stream.Send(&protoUser)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) EditaUsuario(ctx context.Context, usuario *api_cad.Usuario) (*api_cad.Status, error) {
	log.Default().Printf("Editando usuário %s", usuario.Cpf)
	user := database.User{
		Cpf:  utils.CPF(usuario.Cpf),
		Nome: usuario.Nome,
	}
	if !user.Cpf.Validate() {
		log.Printf("CPF inválido")
		return &api_cad.Status{Status: 1, Msg: "CPF inválido"}, nil
	}

	status, err := s.userRepo.EditaUsuario(user)
	if err != nil {
		return &api_cad.Status{Status: 1, Msg: "Erro ao atualizar usuário"}, err
	}

	// Publicar mensagem de atualização no tópico MQTT
	jsonData, err := json.Marshal(usuario)
	if err != nil {
		log.Printf("Erro ao converter dados do usuário para JSON: %v", err)
		return &api_cad.Status{Status: 1, Msg: "Erro ao converter dados para JSON"}, nil
	}

	// TODO: Editar usuário no banco de dados
	_ = jsonData

	return &api_cad.Status{Status: status.Status}, nil
}

func (s *Server) RemoveUsuario(ctx context.Context, request *api_cad.Identificador) (*api_cad.Status, error) {
	status, err := s.userRepo.RemoveUsuario(utils.CPF(request.Id))
	if err != nil {
		return &api_cad.Status{
			Status: status.Status,
			Msg:    err.Error(),
		}, err
	}

	// Publicar mensagem de deleção no tópico MQTT
	jsonData, err := json.Marshal(request)
	if err != nil {
		log.Printf("Erro ao converter dados do usuário para JSON: %v", err)
		return &api_cad.Status{Status: 1, Msg: "Erro ao converter dados para JSON"}, nil
	}
	// TODO: Remover usuário do banco de dados
	_ = jsonData
	return &api_cad.Status{Status: status.Status, Msg: status.Msg}, nil
}

func (s *Server) NovoLivro(ctx context.Context, livro *api_cad.Livro) (*api_cad.Status, error) {
	isbn := utils.ISBN(livro.Isbn)
	if !isbn.Validate() {
		return &api_cad.Status{Status: 1, Msg: "ISBN inválido"}, nil
	}

	jsonData, err := json.Marshal(livro)
	if err != nil {
		log.Printf("Erro ao converter dados do livro para JSON: %v", err)
		return &api_cad.Status{Status: 1, Msg: "Erro ao converter dados para JSON"}, nil
	}
	//TODO: Salvar livro no banco de dados
	_ = jsonData
	return &api_cad.Status{Status: 0}, nil
}

func (s *Server) ObtemLivro(ctx context.Context, request *api_cad.Identificador) (*api_cad.Livro, error) {
	book, err := s.bookRepo.ObtemLivro(utils.ISBN(request.Id))

	if err != nil {
		return nil, err
	}

	return &api_cad.Livro{
		Isbn:   string(book.Isbn),
		Titulo: book.Titulo,
		Autor:  book.Autor,
		Total:  book.Total,
	}, nil
}

func (s *Server) ObtemTodosLivros(request *api_cad.Vazia, stream api_cad.PortalCadastro_ObtemTodosLivrosServer) error {
	books, err := s.bookRepo.ObtemTodosLivros()
	if err != nil {
		return err
	}
	for _, book := range books {
		protoBook := database.BookToProto(book)
		err := stream.Send(&protoBook)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) EditaLivro(ctx context.Context, livro *api_cad.Livro) (*api_cad.Status, error) {
	book := database.Book{
		Isbn:   utils.ISBN(livro.Isbn),
		Titulo: livro.Titulo,
		Autor:  livro.Autor,
		Total:  livro.Total,
	}
	// Validar ISBN aqui, se necessário
	// if !book.Isbn.Validate() { ... }

	status, err := s.bookRepo.EditaLivro(book)
	if err != nil {
		return &api_cad.Status{Status: 1, Msg: "Erro ao atualizar livro"}, err
	}

	// Publicar mensagem de atualização no tópico MQTT
	jsonData, err := json.Marshal(livro)
	if err != nil {
		log.Printf("Erro ao converter dados do livro para JSON: %v", err)
		return &api_cad.Status{Status: 1, Msg: "Erro ao converter dados para JSON"}, nil
	}
	// TODO: Editar livro no banco de dados
	_ = jsonData
	return &api_cad.Status{Status: status.Status}, nil
}

func (s *Server) RemoveLivro(ctx context.Context, request *api_cad.Identificador) (*api_cad.Status, error) {
	status, err := s.bookRepo.RemoveLivro(utils.ISBN(request.Id))
	if err != nil {
		return &api_cad.Status{
			Status: status.Status,
			Msg:    err.Error(),
		}, err
	}

	// Publicar mensagem de deleção no tópico MQTT
	jsonData, err := json.Marshal(request)
	if err != nil {
		log.Printf("Erro ao converter dados do livro para JSON: %v", err)
		return &api_cad.Status{Status: 1, Msg: "Erro ao converter dados para JSON"}, nil
	}

	// TODO: Remover livro do banco de dados
	_ = jsonData
	return &api_cad.Status{Status: status.Status, Msg: status.Msg}, nil
}
