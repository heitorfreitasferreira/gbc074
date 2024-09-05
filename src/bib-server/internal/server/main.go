package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"library-manager/bib-server/internal/database"
	"library-manager/bib-server/internal/queue/handlers"
	"library-manager/bib-server/internal/server/utils"
	"library-manager/shared/api/bib"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Server struct {
	api_bib.UnimplementedPortalBibliotecaServer

	userRepo     database.UserRepo
	bookRepo     database.BookRepo
	userBookRepo database.UserBookRepo

	mqttClient mqtt.Client
}

func NewServer(userRepo database.UserRepo, bookRepo database.BookRepo, userBookRepo database.UserBookRepo, mqttClient mqtt.Client) api_bib.PortalBibliotecaServer {
	return &Server{
		userRepo:     userRepo,
		bookRepo:     bookRepo,
		userBookRepo: userBookRepo,
		mqttClient:   mqttClient,
	}
}

var qos byte = 2

func (s *Server) publishMessage(topic handlers.UserBookTopic, message []byte) error {
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

func (s *Server) publishWithEmptyMessage(topic handlers.UserBookTopic) error {
	var errorMessage string

	if s.mqttClient.IsConnected() {
		token := s.mqttClient.Publish(string(topic), qos, false, "")
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
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			break
		}

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

		err = s.publishMessage(handlers.BookLoanTopic, jsonData)
		if err != nil {
			return stream.SendAndClose(
				&api_bib.Status{
					Status: 1,
					Msg:    err.Error(),
				},
			)
		}
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

	userBook := database.ProtoToUserBook(data)
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

	err = s.publishMessage(handlers.BookReturnTopic, jsonData)
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
	err := s.publishWithEmptyMessage(handlers.UserBlockTopic)
	if err != nil {
		return &api_bib.Status{
			Status: 1,
			Msg:    err.Error(),
		}, err
	}

	return &api_bib.Status{
		Status: 0,
		Msg:    "Solicitação de bloqueio de usuários realizada!",
	}, nil
}

func (s *Server) LiberaUsuarios(ctx context.Context, req *api_bib.Vazia) (*api_bib.Status, error) {
	err := s.publishWithEmptyMessage(handlers.UserFreeTopic)
	if err != nil {
		return &api_bib.Status{
			Status: 1,
			Msg:    err.Error(),
		}, err
	}

	return &api_bib.Status{
		Status: 0,
		Msg:    "Solicitação de liberação de usuários realizada!",
	}, nil

}

func (s *Server) ListaUsuariosBloqueados(req *api_bib.Vazia, stream api_bib.PortalBiblioteca_ListaUsuariosBloqueadosServer) error {
	userList := s.userRepo.GetAll()

	currTime := time.Now().UnixMilli()
	for _, user := range userList {
		if user.Blocked {
			userLoans := s.userBookRepo.GetUserLoans(user.CPF)
			overdueLoanBooks := make([]*api_bib.Livro, 0)
			for _, loan := range userLoans {
				// Verifica se o empréstimo tem mais de 10 segundos
				if (currTime - loan.Timestamp) > 10*1000 {
					book, err := s.bookRepo.GetById(loan.BookISBN)
					if err != nil {
						stream.SendMsg(err)
						return err
					}

					protoBook := database.BookToProto(book)
					overdueLoanBooks = append(overdueLoanBooks, &protoBook)
				}
			}

			protoUser := database.UserToProto(user)

			blockedUser := &api_bib.UsuarioBloqueado{
				Usuario: &protoUser,
				Livros:  overdueLoanBooks,
			}

			stream.Send(blockedUser)
		}
	}

	return nil
}

func (s *Server) ListaLivrosEmprestados(req *api_bib.Vazia, stream api_bib.PortalBiblioteca_ListaLivrosEmprestadosServer) error {
	userList := s.userRepo.GetAll()

	for _, user := range userList {
		userLoans := s.userBookRepo.GetUserLoans(user.CPF)
		for _, loan := range userLoans {
			book, err := s.bookRepo.GetById(loan.BookISBN)
			if err != nil {
				stream.SendMsg(err)
				return err
			}
			protoBook := database.BookToProto(book)

			stream.Send(&protoBook)
		}

	}

	return nil
}

func (s *Server) ListaLivrosEmFalta(req *api_bib.Vazia, stream api_bib.PortalBiblioteca_ListaLivrosEmFaltaServer) error {
	bookList := s.bookRepo.GetAll()

	for _, book := range bookList {
		if book.Remaining == 0 {
			protoBook := database.BookToProto(book)
			stream.Send(&protoBook)
		}
	}

	return nil
}

func (s *Server) PesquisaLivro(req *api_bib.Criterio, stream api_bib.PortalBiblioteca_PesquisaLivroServer) error {
	allBooks := s.bookRepo.GetAll()

	filteredBooks, err := utils.FilterBooks(allBooks, req.Criterio)
	if err != nil {
		return err
	}

	for _, book := range filteredBooks {
		protoBook := database.BookToProto(book)
		stream.Send(&protoBook)
	}

	return nil
}
