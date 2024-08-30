package server

import (
	"context"
	"encoding/json"
	"log"

	api "library-manager/shared/api/cad"
	"library-manager/shared/database"
	"library-manager/shared/utils"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var qos byte = 2

type Server struct {
	api.UnimplementedPortalCadastroServer

	userRepo database.UserRepo
	bookRepo database.BookRepo

	mqttClient mqtt.Client
}

func NewServer(userRepo database.UserRepo, mqttClient mqtt.Client, bookRepo database.BookRepo) *Server {
	return &Server{
		userRepo: userRepo,
		bookRepo: bookRepo,

		mqttClient: mqttClient,
	}
}

func (s *Server) NovoUsuario(ctx context.Context, usuario *api.Usuario) (*api.Status, error) {
	user := database.User{
		Cpf:  utils.CPF(usuario.Cpf),
		Nome: usuario.Nome,
	}
	if !user.Cpf.Validate() {
		log.Printf("CPF inválido")
		return &api.Status{Status: 1, Msg: "CPF inválido"}, nil
	}

	jsonData, err := json.Marshal(user)
	if err != nil {
		log.Printf("Erro ao converter dados do usuário para JSON: %v", err)
		return &api.Status{Status: 1, Msg: "Erro ao converter dados para JSON"}, nil
	}

	if s.mqttClient.IsConnected() {
		token := s.mqttClient.Publish("user/create", qos, false, jsonData)
		token.Wait()
		if token.Error() != nil {
			log.Printf("Erro ao publicar mensagem no tópico MQTT: %v", token.Error())
			return &api.Status{Status: 1, Msg: "Erro ao publicar mensagem no MQTT"}, nil
		} else {
			log.Println("Mensagem publicada no tópico user/create")
		}
	} else {
		log.Println("Cliente MQTT não está conectado")
		return &api.Status{Status: 1, Msg: "MQTT não conectado"}, nil
	}
	return &api.Status{Status: 0}, nil
}

func (s *Server) ObtemUsuario(ctx context.Context, request *api.Identificador) (*api.Usuario, error) {
	user, err := s.userRepo.ObtemUsuario(utils.CPF(request.Id))

	if err != nil {
		return nil, err
	}

	return &api.Usuario{
		Cpf:  request.Id,
		Nome: user.Nome,
	}, nil
}

func (s *Server) ObtemTodosUsuarios(request *api.Vazia, stream api.PortalCadastro_ObtemTodosUsuariosServer) error {
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

func (s *Server) EditaUsuario(ctx context.Context, usuario *api.Usuario) (*api.Status, error) {
	log.Default().Printf("Editando usuário %s", usuario.Cpf)
	user := database.User{
		Cpf:  utils.CPF(usuario.Cpf),
		Nome: usuario.Nome,
	}
	if !user.Cpf.Validate() {
		log.Printf("CPF inválido")
		return &api.Status{Status: 1, Msg: "CPF inválido"}, nil
	}

	status, err := s.userRepo.EditaUsuario(user)
	if err != nil {
		return &api.Status{Status: 1, Msg: "Erro ao atualizar usuário"}, err
	}

	// Publicar mensagem de atualização no tópico MQTT
	jsonData, err := json.Marshal(user)
	if err != nil {
		log.Printf("Erro ao converter dados do usuário para JSON: %v", err)
		return &api.Status{Status: 1, Msg: "Erro ao converter dados para JSON"}, nil
	}
	if s.mqttClient.IsConnected() {
		token := s.mqttClient.Publish("user/update", qos, false, jsonData)
		token.Wait()
		if token.Error() != nil {
			log.Printf("Erro ao publicar mensagem no tópico MQTT: %v", token.Error())
			return &api.Status{Status: 1, Msg: "Erro ao publicar mensagem no MQTT"}, nil
		}
	} else {
		log.Println("Cliente MQTT não está conectado")
		return &api.Status{Status: 1, Msg: "MQTT não conectado"}, nil
	}

	return &api.Status{Status: status.Status}, nil
}

// RemoveUsuario: deleta um usuário pelo CPF
func (s *Server) RemoveUsuario(ctx context.Context, request *api.Identificador) (*api.Status, error) {
	status, err := s.userRepo.RemoveUsuario(utils.CPF(request.Id))
	if err != nil {
		return &api.Status{
			Status: status.Status,
			Msg:    err.Error(),
		}, err
	}

	// Publicar mensagem de deleção no tópico MQTT
	jsonData, err := json.Marshal(request)
	if err != nil {
		log.Printf("Erro ao converter dados do usuário para JSON: %v", err)
		return &api.Status{Status: 1, Msg: "Erro ao converter dados para JSON"}, nil
	}

	if s.mqttClient.IsConnected() {
		token := s.mqttClient.Publish("user/remove", 2, false, jsonData)
		token.Wait()
		if token.Error() != nil {
			log.Printf("Erro ao publicar mensagem no tópico MQTT: %v", token.Error())
			return &api.Status{Status: 1, Msg: "Erro ao publicar mensagem no MQTT"}, nil
		}
	} else {
		log.Println("Cliente MQTT não está conectado")
		return &api.Status{Status: 1, Msg: "MQTT não conectado"}, nil
	}
	return &api.Status{Status: status.Status, Msg: status.Msg}, nil
}

func (s *Server) NovoLivro(ctx context.Context, livro *api.Livro) (*api.Status, error) {
	book := database.Book{
		Isbn:   utils.ISBN(livro.Isbn),
		Titulo: livro.Titulo,
		Autor:  livro.Autor,
	}
	if !book.Isbn.Validate() {
		return &api.Status{Status: 1, Msg: "ISBN inválido"}, nil
	}

	jsonData, err := json.Marshal(book)
	if err != nil {
		log.Printf("Erro ao converter dados do livro para JSON: %v", err)
		return &api.Status{Status: 1, Msg: "Erro ao converter dados para JSON"}, nil
	}

	if s.mqttClient.IsConnected() {
		token := s.mqttClient.Publish("book/create", qos, false, jsonData)
		token.Wait()
		if token.Error() != nil {
			log.Printf("Erro ao publicar mensagem no tópico MQTT: %v", token.Error())
			return &api.Status{Status: 1, Msg: "Erro ao publicar mensagem no MQTT"}, nil
		} else {
			log.Println("Mensagem publicada no tópico book/create")
		}
	} else {
		log.Println("Cliente MQTT não está conectado")
		return &api.Status{Status: 1, Msg: "MQTT não conectado"}, nil
	}
	return &api.Status{Status: 0}, nil
}

func (s *Server) ObtemLivro(ctx context.Context, request *api.Identificador) (*api.Livro, error) {
	book, err := s.bookRepo.ObtemLivro(utils.ISBN(request.Id))

	if err != nil {
		return nil, err
	}

	return &api.Livro{
		Isbn:   request.Id,
		Titulo: book.Titulo,
		Autor:  book.Autor,
	}, nil
}

func (s *Server) ObtemTodosLivros(request *api.Vazia, stream api.PortalCadastro_ObtemTodosLivrosServer) error {
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

func (s *Server) EditarLivro(ctx context.Context, livro *api.Livro) (*api.Status, error) {
	book := database.Book{
		Isbn:   utils.ISBN(livro.Isbn),
		Titulo: livro.Titulo,
		Autor:  livro.Autor,
	}
	// Validar ISBN aqui, se necessário
	// if !book.Isbn.Validate() { ... }

	status, err := s.bookRepo.EditaLivro(book)
	if err != nil {
		return &api.Status{Status: 1, Msg: "Erro ao atualizar livro"}, err
	}

	// Publicar mensagem de atualização no tópico MQTT
	jsonData, err := json.Marshal(book)
	if err != nil {
		log.Printf("Erro ao converter dados do livro para JSON: %v", err)
		return &api.Status{Status: 1, Msg: "Erro ao converter dados para JSON"}, nil
	}
	if s.mqttClient.IsConnected() {
		token := s.mqttClient.Publish("book/update", qos, false, jsonData)
		token.Wait()
		if token.Error() != nil {
			log.Printf("Erro ao publicar mensagem no tópico MQTT: %v", token.Error())
			return &api.Status{Status: 1, Msg: "Erro ao publicar mensagem no MQTT"}, nil
		}
	} else {
		log.Println("Cliente MQTT não está conectado")
		return &api.Status{Status: 1, Msg: "MQTT não conectado"}, nil
	}

	return &api.Status{Status: status.Status}, nil
}

func (s *Server) RemoveLivro(ctx context.Context, request *api.Identificador) (*api.Status, error) {
	status, err := s.bookRepo.RemoveLivro(utils.ISBN(request.Id))
	if err != nil {
		return &api.Status{
			Status: status.Status,
			Msg:    err.Error(),
		}, err
	}

	// Publicar mensagem de deleção no tópico MQTT
	jsonData, err := json.Marshal(request)
	if err != nil {
		log.Printf("Erro ao converter dados do livro para JSON: %v", err)
		return &api.Status{Status: 1, Msg: "Erro ao converter dados para JSON"}, nil
	}

	if s.mqttClient.IsConnected() {
		token := s.mqttClient.Publish("book/remove", 2, false, jsonData)
		token.Wait()
		if token.Error() != nil {
			log.Printf("Erro ao publicar mensagem no tópico MQTT: %v", token.Error())
			return &api.Status{Status: 1, Msg: "Erro ao publicar mensagem no MQTT"}, nil
		}
	} else {
		log.Println("Cliente MQTT não está conectado")
		return &api.Status{Status: 1, Msg: "MQTT não conectado"}, nil
	}
	return &api.Status{Status: status.Status, Msg: status.Msg}, nil
}
