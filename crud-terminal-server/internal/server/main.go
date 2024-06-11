package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	br_ufu_facom_gbc074_projeto_cadastro "github.com/rpc-mqtt-library-manager/crud-terminal-server/api"
	"github.com/rpc-mqtt-library-manager/crud-terminal-server/internal/database"
	"github.com/rpc-mqtt-library-manager/crud-terminal-server/internal/utils"
)

type Server struct {
	br_ufu_facom_gbc074_projeto_cadastro.UnimplementedPortalCadastroServer

	userRepo database.UserRepo

	mqttClient mqtt.Client
}

func NewServer() *Server {
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883")
	opts.SetClientID("cadastro_server")
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Erro ao conectar ao broker MQTT: %v", token.Error())
	}

	return &Server{
		mqttClient: client,
		userRepo:   database.NewInMemoryUserRepo(),
	}
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Conectado ao broker MQTT")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Conexão perdida: %v", err)
}

func (s *Server) NovoUsuario(ctx context.Context, usuario *br_ufu_facom_gbc074_projeto_cadastro.Usuario) (*br_ufu_facom_gbc074_projeto_cadastro.Status, error) {
	user := database.User{
		Cpf:  utils.CPF(usuario.Cpf),
		Nome: usuario.Nome,
	}
	if !user.Cpf.Validate() {
		return &br_ufu_facom_gbc074_projeto_cadastro.Status{Status: 1, Msg: "CPF inválido"}, nil
	}

	// Publicar mensagem no tópico MQTT
	jsonData, err := json.Marshal(user)
	if err != nil {
		log.Printf("Erro ao converter dados do usuário para JSON: %v", err)
		return &br_ufu_facom_gbc074_projeto_cadastro.Status{Status: 1, Msg: "Erro ao converter dados para JSON"}, nil
	}
	token := s.mqttClient.Publish("user/create", 2, false, jsonData)
	token.Wait()

	return &br_ufu_facom_gbc074_projeto_cadastro.Status{Status: 0}, nil
}

// ObterUsuario: obtém um usuário pelo CPF
func (s *Server) ObterUsuario(ctx context.Context, request *br_ufu_facom_gbc074_projeto_cadastro.Identificador) (*br_ufu_facom_gbc074_projeto_cadastro.Usuario, error) {
	user, err := s.userRepo.ObtemUsuario(utils.CPF(request.Id))

	if err != nil {
		return nil, err
	}

	return &br_ufu_facom_gbc074_projeto_cadastro.Usuario{
		Cpf:  request.Id,
		Nome: user.Nome,
	}, nil
}

// AtualizarUsuario: atualiza os dados de um usuário
func (s *Server) AtualizarUsuario(ctx context.Context, usuario *br_ufu_facom_gbc074_projeto_cadastro.Usuario) (*br_ufu_facom_gbc074_projeto_cadastro.Status, error) {
	user := database.User{
		Cpf:  utils.CPF(usuario.Cpf),
		Nome: usuario.Nome,
	}
	if !user.Cpf.Validate() {
		return &br_ufu_facom_gbc074_projeto_cadastro.Status{Status: 1, Msg: "CPF inválido"}, nil
	}

	status, err := s.userRepo.EditaUsuario(user)
	if err != nil {
		return &br_ufu_facom_gbc074_projeto_cadastro.Status{Status: 1, Msg: "Erro ao atualizar usuário"}, err
	}

	// Publicar mensagem de atualização no tópico MQTT
	jsonData, err := json.Marshal(user)
	if err != nil {
		log.Printf("Erro ao converter dados do usuário para JSON: %v", err)
		return &br_ufu_facom_gbc074_projeto_cadastro.Status{Status: 1, Msg: "Erro ao converter dados para JSON"}, nil
	}
	token := s.mqttClient.Publish("user/update", 2, false, jsonData)
	token.Wait()

	return &br_ufu_facom_gbc074_projeto_cadastro.Status{Status: status.Status}, nil
}

// DeletarUsuario: deleta um usuário pelo CPF
func (s *Server) DeletarUsuario(ctx context.Context, request *br_ufu_facom_gbc074_projeto_cadastro.Identificador) (*br_ufu_facom_gbc074_projeto_cadastro.Status, error) {
	status, err := s.userRepo.RemoveUsuario(utils.CPF(request.Id))
	if err != nil {
		return &br_ufu_facom_gbc074_projeto_cadastro.Status{
			Status: status.Status,
			Msg:    err.Error(),
		}, err
	}

	// Publicar mensagem de deleção no tópico MQTT
	jsonData, err := json.Marshal(request)
	if err != nil {
		log.Printf("Erro ao converter dados do usuário para JSON: %v", err)
		return &br_ufu_facom_gbc074_projeto_cadastro.Status{Status: 1, Msg: "Erro ao converter dados para JSON"}, nil
	}

	token := s.mqttClient.Publish("user/remove", 2, false, jsonData)
	token.Wait()

	return &br_ufu_facom_gbc074_projeto_cadastro.Status{Status: status.Status, Msg: status.Msg}, nil
}

// Função para escutar mensagens MQTT
func (s *Server) ListenMQTT() {
	s.mqttClient.Subscribe("user/create", 0, func(client mqtt.Client, msg mqtt.Message) {
		var usuario database.User
		if err := json.Unmarshal(msg.Payload(), &usuario); err != nil {
			log.Printf("Erro ao decodificar mensagem MQTT: %v", err)
			return
		}

		// Atualiza ou cria o usuário no repositório
		_, err := s.userRepo.EditaUsuario(usuario)
		if err != nil {
			log.Fatalf("Erro ao atualizar usuário via MQTT: %v", err)
		} else {
			log.Printf("Usuário atualizado via MQTT: %s", usuario.Cpf)
		}
	})

	s.mqttClient.Subscribe("user/remove", 0, func(client mqtt.Client, msg mqtt.Message) {
		var request br_ufu_facom_gbc074_projeto_cadastro.Identificador
		if err := json.Unmarshal(msg.Payload(), &request); err != nil {
			log.Printf("Erro ao decodificar mensagem MQTT: %v", err)
			return
		}

		cpf := utils.CPF(request.Id)
		_, err := s.userRepo.RemoveUsuario(cpf)
		if err != nil {
			log.Printf("Erro ao deletar usuário via MQTT: %v", err)
			return
		}
		log.Printf("Usuário deletado via MQTT: %s", request.Id)
	})
}
