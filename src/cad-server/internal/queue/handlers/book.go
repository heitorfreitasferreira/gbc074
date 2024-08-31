package handlers

import (
	"encoding/json"
	"log"

	"library-manager/cad-server/internal/database"
	"library-manager/cad-server/internal/utils"
	api_cad "library-manager/shared/api/cad"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func CreateBook(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Mensagem recebida no tópico book/create: %s\n", msg.Payload())
	var protoBook api_cad.Livro
	err := json.Unmarshal(msg.Payload(), &protoBook)
	if err != nil {
		log.Printf("Erro ao converter dados do livro para JSON: %v", err)
		return
	}
	book := database.ProtoToBook(&protoBook)
	database.ConcreteBookRepo.CreateBook(book)
}

func UpdateBook(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Mensagem recebida no tópico book/update: %s\n", msg.Payload())
	var protoBook api_cad.Livro
	err := json.Unmarshal(msg.Payload(), &protoBook)
	if err != nil {
		log.Printf("Erro ao converter dados do livro para JSON: %v", err)
		return
	}

	book := database.ProtoToBook(&protoBook)
	database.ConcreteBookRepo.EditaLivro(book)
}

func RemoveBook(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Mensagem recebida no tópico book/remove: %s\n", msg.Payload())
	var protoId api_cad.Identificador
	err := json.Unmarshal(msg.Payload(), &protoId)
	if err != nil {
		log.Printf("Erro ao converter dados do livro para JSON: %v", err)
		return
	}
	database.ConcreteBookRepo.RemoveLivro(utils.ISBN(protoId.Id))
}
