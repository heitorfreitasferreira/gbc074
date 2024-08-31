package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"library-manager/bib-server/internal/api"
	"library-manager/bib-server/internal/database"
	"library-manager/bib-server/internal/queue"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Server struct {
	api_bib.UnimplementedPortalBibliotecaServer

	userRepo     database.UserRepo
	bookRepo     database.BookRepo
	userBookRepo database.UserBookRepo

	mqttClient mqtt.Client
}

func NewServer(userRepo database.UserRepo, bookRepo database.BookRepo, userBookRepo database.UserBookRepo, mqttClient mqtt.Client) *Server {
	return &Server{
		userRepo:     userRepo,
		bookRepo:     bookRepo,
		userBookRepo: userBookRepo,
		mqttClient:   mqttClient,
	}
}

var qos byte = 2

func (s *Server) publishMessage(topic queue.UserBookTopic, message []byte) error {
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

func (s *Server) publishWithEmptyMessage(topic queue.UserBookTopic) error {
	var errorMessage string

	if s.mqttClient.IsConnected() {
		token := s.mqttClient.Publish(string(topic), qos, false, nil)
		token.Wait()
		if token.Error() == nil {
			log.Printf("Mensagem vazia publicada no tópico %v", topic)
			return nil
		} else {
			errorMessage = fmt.Sprintf("Erro ao publicar mensagem vazia no tópico MQTT: %v", token.Error())
		}
	} else {
		errorMessage = "Cliente MQTT não está conectado"
	}

	return errors.New("Server.publishWithEmptyMessage: " + errorMessage)
}

func (s *Server) RealizaEmprestimo(stream api_bib.PortalBiblioteca_RealizaEmprestimoServer) error {
	data, err := stream.Recv()
	log.Printf("Recebendo dados %v", data)

	if err != nil {
		return stream.SendAndClose(
			&api_bib.Status{
				Status: 1,
				Msg:    "Erro ao receber dados",
			},
		)
	}

	if data.Usuario == nil || data.Livro == nil {
		return stream.SendAndClose(
			&api_bib.Status{
				Status: 1,
				Msg:    "Dados inválidos",
			},
		)
	}

	userBook := database.NewUserBook(data.Usuario.Id, data.Livro.Id)
	jsonData, err := json.Marshal(userBook)
	if err != nil {
		errMsg := fmt.Sprintf("Erro ao converter dados para JSON: %v", err)
		log.Println(errMsg)
		return stream.SendAndClose(
			&api_bib.Status{
				Status: 1,
				Msg:    errMsg,
			},
		)
	}

	err = s.publishMessage(queue.BookLoanTopic, jsonData)
	if err != nil {
		return stream.SendAndClose(
			&api_bib.Status{
				Status: 1,
				Msg:    err.Error(),
			},
		)
	}

	return stream.SendAndClose(
		&api_bib.Status{
			Status: 0,
			Msg:    "Solicitação de empréstimo realizada!",
		},
	)
}

func (s *Server) RealizaDevolucao(stream api_bib.PortalBiblioteca_RealizaDevolucaoServer) error {
	data, err := stream.Recv()
	log.Printf("Recebendo dados %v", data)

	if err != nil {
		return stream.SendAndClose(
			&api_bib.Status{
				Status: 1,
				Msg:    "Erro ao receber dados",
			},
		)
	}

	if data.Usuario == nil || data.Livro == nil {
		return stream.SendAndClose(
			&api_bib.Status{
				Status: 1,
				Msg:    "Dados inválidos",
			},
		)
	}

	userBook := database.NewUserBook(data.Usuario.Id, data.Livro.Id)
	jsonData, err := json.Marshal(userBook)
	if err != nil {
		errMsg := fmt.Sprintf("Erro ao converter dados para JSON: %v", err)
		log.Println(errMsg)
		return stream.SendAndClose(
			&api_bib.Status{
				Status: 1,
				Msg:    errMsg,
			},
		)
	}

	err = s.publishMessage(queue.BookReturnTopic, jsonData)
	if err != nil {
		return stream.SendAndClose(
			&api_bib.Status{
				Status: 1,
				Msg:    err.Error(),
			},
		)
	}

	return stream.SendAndClose(
		&api_bib.Status{
			Status: 0,
			Msg:    "Solicitação de empréstimo realizada!",
		},
	)

}

func (s *Server) BloqueiaUsuarios(ctx context.Context, req *api_bib.Vazia) (*api_bib.Status, error) {
	err := s.publishWithEmptyMessage(queue.UserBlockTopic)
	if err != nil {
		return &api_bib.Status{
			Status: 0,
			Msg:    err.Error(),
		}, err
	}

	return &api_bib.Status{
		Status: 0,
		Msg:    "Solicitação de bloqueio de usuários realizada!",
	}, nil
}

func (s *Server) LiberaUsuarios(ctx context.Context, req *api_bib.Vazia) (*api_bib.Status, error) {
	err := s.publishWithEmptyMessage(queue.UserFreeTopic)
	if err != nil {
		return &api_bib.Status{
			Status: 0,
			Msg:    err.Error(),
		}, err
	}

	return &api_bib.Status{
		Status: 0,
		Msg:    "Solicitação de liberação de usuários realizada!",
	}, nil

}

func (s *Server) ListaUsuariosBloqueados(req *api_bib.Vazia, stream api_bib.PortalBiblioteca_ListaUsuariosBloqueadosServer) error {
	// Get books if loaned more than 10 seconds ago
	return nil
}

func (s *Server) ListaLivrosEmprestados(req *api_bib.Vazia, stream api_bib.PortalBiblioteca_ListaLivrosEmprestadosServer) error {
	return nil
}

func (s *Server) ListaLivrosEmFalta(req *api_bib.Vazia, stream api_bib.PortalBiblioteca_ListaLivrosEmFaltaServer) error {
	return nil
}

func (s *Server) PesquisaLivro(req *api_bib.Criterio, stream api_bib.PortalBiblioteca_PesquisaLivroServer) error {
	return nil
}
