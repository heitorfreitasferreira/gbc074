package server

import (
	"context"
	"errors"
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rpc-mqtt-library-manager/cms-server/api"
	"github.com/rpc-mqtt-library-manager/cms-server/internal/database"
	"github.com/rpc-mqtt-library-manager/cms-server/internal/queue"
)

type Server struct {
	br_ufu_facom_gbc074_projeto_biblioteca.UnimplementedPortalBibliotecaServer

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

func (s *Server) publishMessage(topic queue.Topic, message string) error {
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

func (s *Server) RealizaEmprestimo(stream br_ufu_facom_gbc074_projeto_biblioteca.PortalBiblioteca_RealizaEmprestimoServer) error {
	s.publishMessage(queue.BookLoanTopic, "Emprestimo realizado")
	return nil
}

func (s *Server) RealizaDevolucao(stream br_ufu_facom_gbc074_projeto_biblioteca.PortalBiblioteca_RealizaDevolucaoServer) error {
	s.publishMessage(queue.BookReturnTopic, "Devolução realizada")
	return nil
}

func (s *Server) BloqueiaUsuarios(ctx context.Context, req *br_ufu_facom_gbc074_projeto_biblioteca.Vazia) (*br_ufu_facom_gbc074_projeto_biblioteca.Status, error) {
	s.publishMessage(queue.UserBlockTopic, "Usuário bloqueado")
	return nil, nil
}

func (s *Server) LiberaUsuarios(ctx context.Context, req *br_ufu_facom_gbc074_projeto_biblioteca.Vazia) (*br_ufu_facom_gbc074_projeto_biblioteca.Status, error) {
	s.publishMessage(queue.UserFreeTopic, "Usuário liberado")
	return nil, nil
}

func (s *Server) ListaUsuariosBloqueados(req *br_ufu_facom_gbc074_projeto_biblioteca.Vazia, stream br_ufu_facom_gbc074_projeto_biblioteca.PortalBiblioteca_ListaUsuariosBloqueadosServer) error {
	s.publishMessage(queue.UserListBlockedTopic, "Lista de usuários bloqueados")
	return nil
}

func (s *Server) ListaLivrosEmprestados(req *br_ufu_facom_gbc074_projeto_biblioteca.Vazia, stream br_ufu_facom_gbc074_projeto_biblioteca.PortalBiblioteca_ListaLivrosEmprestadosServer) error {
	s.publishMessage(queue.BookListBorrowedTopic, "Lista de livros emprestados")
	return nil
}

func (s *Server) ListaLivrosEmFalta(req *br_ufu_facom_gbc074_projeto_biblioteca.Vazia, stream br_ufu_facom_gbc074_projeto_biblioteca.PortalBiblioteca_ListaLivrosEmFaltaServer) error {
	s.publishMessage(queue.BookListMissingTopic, "Lista de livros em falta")
	return nil
}

func (s *Server) PesquisaLivro(req *br_ufu_facom_gbc074_projeto_biblioteca.Criterio, stream br_ufu_facom_gbc074_projeto_biblioteca.PortalBiblioteca_PesquisaLivroServer) error {
	s.publishMessage(queue.BookSearchTopic, "Livro encontrado")
	return nil
}
