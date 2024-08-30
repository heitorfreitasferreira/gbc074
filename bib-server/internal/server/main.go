package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	api "library-manager/shared/api/bib"
	"library-manager/shared/database"
	"library-manager/shared/queue"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Server struct {
	api.UnimplementedPortalBibliotecaServer

	userRepo database.UserRepo
	bookRepo database.BookRepo

	mqttClient mqtt.Client
}

func NewServer(userRepo database.UserRepo, bookRepo database.BookRepo, mqttClient mqtt.Client) *Server {
	return &Server{
		userRepo:   userRepo,
		bookRepo:   bookRepo,
		mqttClient: mqttClient,
	}
}

var qos byte = 2

func (s *Server) publishMessage(topic queue.Topic, message []byte) error {
	var errorMessage string

	if s.mqttClient.IsConnected() {
		token := s.mqttClient.Publish(string(topic), qos, false, message)
		token.Wait()
		if token.Error() == nil {
			log.Printf("Mensagem publicada no tópico %v", topic)
			return nil
		} else {
			errorMessage = fmt.Sprintf("Erro ao publicar mensagem no tópico MQTT: %v", token.Error())
		}
	} else {
		errorMessage = "Cliente MQTT não está conectado"
	}

	return errors.New("Server.publishMessage: " + errorMessage)
}

func (s *Server) RealizaEmprestimo(stream api.PortalBiblioteca_RealizaEmprestimoServer) error {
	data, err := stream.Recv()
	log.Printf("Recebendo dados %v", data)

	if err != nil {
		return stream.SendAndClose(
			&api.Status{
				Status: 1,
				Msg:    "Erro ao receber dados",
			},
		)
	}

	if data.Usuario == nil || data.Livro == nil {
		return stream.SendAndClose(
			&api.Status{
				Status: 1,
				Msg:    "Dados inválidos",
			},
		)
	}

	userBook := database.UserBook{
		UserId:   data.Usuario.Id,
		BookISNB: data.Livro.Id,
	}

	jsonData, err := json.Marshal(userBook)
	if err != nil {
		errMsg := fmt.Sprintf("Erro ao converter dados para JSON: %v", err)
		log.Println(errMsg)
		return stream.SendAndClose(
			&api.Status{
				Status: 1,
				Msg:    errMsg,
			},
		)
	}

	err = s.publishMessage(queue.BookLoanTopic, jsonData)
	if err != nil {
		return stream.SendAndClose(
			&api.Status{
				Status: 1,
				Msg:    err.Error(),
			},
		)
	}

	return stream.SendAndClose(
		&api.Status{
			Status: 0,
			Msg:    "Empréstimo realizado",
		},
	)
}

func (s *Server) RealizaDevolucao(stream api.PortalBiblioteca_RealizaDevolucaoServer) error {
	return nil
}

func (s *Server) BloqueiaUsuarios(ctx context.Context, req *api.Vazia) (*api.Status, error) {
	return nil, nil
}

func (s *Server) LiberaUsuarios(ctx context.Context, req *api.Vazia) (*api.Status, error) {
	return nil, nil
}

func (s *Server) ListaUsuariosBloqueados(req *api.Vazia, stream api.PortalBiblioteca_ListaUsuariosBloqueadosServer) error {
	return nil
}

func (s *Server) ListaLivrosEmprestados(req *api.Vazia, stream api.PortalBiblioteca_ListaLivrosEmprestadosServer) error {
	return nil
}

func (s *Server) ListaLivrosEmFalta(req *api.Vazia, stream api.PortalBiblioteca_ListaLivrosEmFaltaServer) error {
	return nil
}

func (s *Server) PesquisaLivro(req *api.Criterio, stream api.PortalBiblioteca_PesquisaLivroServer) error {
	return nil
}
