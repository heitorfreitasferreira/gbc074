package handlers

import (
	"encoding/json"
	"log"

	"library-manager/bib-server/internal/database"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func UserCreateHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Mensagem recebida no tópico user/create: %s\n", msg.Payload())
	var user database.User
	err := json.Unmarshal(msg.Payload(), &user)
	if err != nil {
		log.Printf("Erro ao converter dados do usuário para JSON: %v", err)
		return
	}
	database.ConcreteUserRepo.Create(user)
}
func UserUpdateHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Mensagem recebida no tópico user/update: %s\n", msg.Payload())
	var user database.User
	err := json.Unmarshal(msg.Payload(), &user)
	if err != nil {
		log.Printf("Erro ao converter dados do usuário para JSON: %v", err)
		return
	}
	database.ConcreteUserRepo.Update(database.User{
		CPF:     user.CPF,
		Name:    user.Name,
		Blocked: user.Blocked,
	})
}
func UserRemoveHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Mensagem recebida no tópico user/remove: %s\n", msg.Payload())
	var user database.User
	err := json.Unmarshal(msg.Payload(), &user)
	if err != nil {
		log.Printf("Erro ao converter dados do usuário para JSON: %v", err)
		return
	}
	database.ConcreteUserRepo.Remove(user.CPF)
}
