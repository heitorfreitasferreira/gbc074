package handlers

import (
	"encoding/json"
	"library-manager/shared/database"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func CreateBook(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Mensagem recebida no tópico book/create: %s\n", msg.Payload())
	var book database.Book
	err := json.Unmarshal(msg.Payload(), &book)
	if err != nil {
		log.Printf("Erro ao converter dados do livro para JSON: %v", err)
		return
	}
	database.ConcreteBookRepo.CreateBook(book)
}

func UpdateBook(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Mensagem recebida no tópico book/update: %s\n", msg.Payload())
	var book database.Book
	err := json.Unmarshal(msg.Payload(), &book)
	if err != nil {
		log.Printf("Erro ao converter dados do livro para JSON: %v", err)
		return
	}
	database.ConcreteBookRepo.EditaLivro(book)
}

func RemoveBook(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Mensagem recebida no tópico book/remove: %s\n", msg.Payload())
	var book database.Book
	err := json.Unmarshal(msg.Payload(), &book)
	if err != nil {
		log.Printf("Erro ao converter dados do livro para JSON: %v", err)
		return
	}
	database.ConcreteBookRepo.RemoveLivro(book.Isbn)
}
