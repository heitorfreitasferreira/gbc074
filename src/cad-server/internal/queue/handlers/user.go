package handlers

import (
	"encoding/json"
	"log"

	"library-manager/cad-server/internal/database"
	"library-manager/cad-server/internal/utils"
	api_cad "library-manager/shared/api/cad"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func CreateUser(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Mensagem recebida no tópico user/create: %s\n", msg.Payload())
	var protoUser api_cad.Usuario
	err := json.Unmarshal(msg.Payload(), &protoUser)
	if err != nil {
		log.Printf("Erro ao converter dados do usuário para JSON: %v", err)
		return
	}

	user := database.ProtoToUser(&protoUser)
	database.ConcreteUserRepo.CreateUser(user)
}
func UpdateUser(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Mensagem recebida no tópico user/update: %s\n", msg.Payload())
	var protoUser api_cad.Usuario
	err := json.Unmarshal(msg.Payload(), &protoUser)
	if err != nil {
		log.Printf("Erro ao converter dados do usuário para JSON: %v", err)
		return
	}

	user := database.ProtoToUser(&protoUser)
	database.ConcreteUserRepo.EditaUsuario(user)
}
func RemoveUser(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Mensagem recebida no tópico user/remove: %s\n", msg.Payload())
	var protoId api_cad.Identificador
	err := json.Unmarshal(msg.Payload(), &protoId)
	if err != nil {
		log.Printf("Erro ao converter dados do usuário para JSON: %v", err)
		return
	}

	database.ConcreteUserRepo.RemoveUsuario(utils.CPF(protoId.Id))
}
