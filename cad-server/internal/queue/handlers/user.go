package handlers

import (
	"encoding/json"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"library-manager/cad-server/internal/database"
)

func CreateUser(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Mensagem recebida no tópico user/create: %s\n", msg.Payload())
	var user database.User
	err := json.Unmarshal(msg.Payload(), &user)
	if err != nil {
		log.Printf("Erro ao converter dados do usuário para JSON: %v", err)
		return
	}
	database.ConcreteUserRepo.CreateUser(user)
}
func UpdateUser(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Mensagem recebida no tópico user/update: %s\n", msg.Payload())
	var user database.User
	err := json.Unmarshal(msg.Payload(), &user)
	if err != nil {
		log.Printf("Erro ao converter dados do usuário para JSON: %v", err)
		return
	}
	database.ConcreteUserRepo.EditaUsuario(user)
}
func RemoveUser(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Mensagem recebida no tópico user/remove: %s\n", msg.Payload())
	var user database.User
	err := json.Unmarshal(msg.Payload(), &user)
	if err != nil {
		log.Printf("Erro ao converter dados do usuário para JSON: %v", err)
		return
	}
	database.ConcreteUserRepo.RemoveUsuario(user.Cpf)
}
