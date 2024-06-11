package server

import (
	"context"
	"encoding/json"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	br_ufu_facom_gbc074_projeto_cadastro "github.com/rpc-mqtt-library-manager/crud-terminal-server/api"
	"github.com/rpc-mqtt-library-manager/crud-terminal-server/internal/database"
	"github.com/rpc-mqtt-library-manager/crud-terminal-server/internal/utils"
)

var qos byte = 2

type Server struct {
	br_ufu_facom_gbc074_projeto_cadastro.UnimplementedPortalCadastroServer

	userRepo database.UserRepo

	mqttClient mqtt.Client
}

func NewServer(userRepo database.UserRepo, mqttClient mqtt.Client) *Server {
	return &Server{
		userRepo:   userRepo,
		mqttClient: mqttClient,
	}
}

func (s *Server) NovoUsuario(ctx context.Context, usuario *br_ufu_facom_gbc074_projeto_cadastro.Usuario) (*br_ufu_facom_gbc074_projeto_cadastro.Status, error) {
	user := database.User{
		Cpf:  utils.CPF(usuario.Cpf),
		Nome: usuario.Nome,
	}
	if !user.Cpf.Validate() {
		log.Printf("CPF inválido")
		return &br_ufu_facom_gbc074_projeto_cadastro.Status{Status: 1, Msg: "CPF inválido"}, nil
	}

	jsonData, err := json.Marshal(user)
	if err != nil {
		log.Printf("Erro ao converter dados do usuário para JSON: %v", err)
		return &br_ufu_facom_gbc074_projeto_cadastro.Status{Status: 1, Msg: "Erro ao converter dados para JSON"}, nil
	}

	if s.mqttClient.IsConnected() {
		token := s.mqttClient.Publish("user/create", qos, false, jsonData)
		token.Wait()
		if token.Error() != nil {
			log.Printf("Erro ao publicar mensagem no tópico MQTT: %v", token.Error())
			return &br_ufu_facom_gbc074_projeto_cadastro.Status{Status: 1, Msg: "Erro ao publicar mensagem no MQTT"}, nil
		} else {
			log.Println("Mensagem publicada no tópico user/create")
		}
	} else {
		log.Println("Cliente MQTT não está conectado")
		return &br_ufu_facom_gbc074_projeto_cadastro.Status{Status: 1, Msg: "MQTT não conectado"}, nil
	}
	return &br_ufu_facom_gbc074_projeto_cadastro.Status{Status: 0}, nil
}

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

func (s *Server) AtualizarUsuario(ctx context.Context, usuario *br_ufu_facom_gbc074_projeto_cadastro.Usuario) (*br_ufu_facom_gbc074_projeto_cadastro.Status, error) {
	user := database.User{
		Cpf:  utils.CPF(usuario.Cpf),
		Nome: usuario.Nome,
	}
	if !user.Cpf.Validate() {
		log.Printf("CPF inválido")
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
	if s.mqttClient.IsConnected() {
		token := s.mqttClient.Publish("user/update", qos, false, jsonData)
		token.Wait()
		if token.Error() != nil {
			log.Printf("Erro ao publicar mensagem no tópico MQTT: %v", token.Error())
			return &br_ufu_facom_gbc074_projeto_cadastro.Status{Status: 1, Msg: "Erro ao publicar mensagem no MQTT"}, nil
		}
	} else {
		log.Println("Cliente MQTT não está conectado")
		return &br_ufu_facom_gbc074_projeto_cadastro.Status{Status: 1, Msg: "MQTT não conectado"}, nil
	}

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

	if s.mqttClient.IsConnected() {
		token := s.mqttClient.Publish("user/remove", 2, false, jsonData)
		token.Wait()
		if token.Error() != nil {
			log.Printf("Erro ao publicar mensagem no tópico MQTT: %v", token.Error())
			return &br_ufu_facom_gbc074_projeto_cadastro.Status{Status: 1, Msg: "Erro ao publicar mensagem no MQTT"}, nil
		}
	} else {
		log.Println("Cliente MQTT não está conectado")
		return &br_ufu_facom_gbc074_projeto_cadastro.Status{Status: 1, Msg: "MQTT não conectado"}, nil
	}
	return &br_ufu_facom_gbc074_projeto_cadastro.Status{Status: status.Status, Msg: status.Msg}, nil
}
