package handlers

import (
	"encoding/json"
	"log"

	"library-manager/bib-server/internal/database"
	api_bib "library-manager/shared/api/bib"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func BookCreateHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Mensagem recebida no tópico book/create: %s\n", msg.Payload())
	var protoBook api_bib.Livro
	err := json.Unmarshal(msg.Payload(), &protoBook)
	if err != nil {
		log.Printf("Erro ao converter dados do livro para JSON: %v", err)
		return
	}
	database.ConcreteBookRepo.Create(database.Book{
		ISBN:      protoBook.Isbn,
		Title:     protoBook.Titulo,
		Author:    protoBook.Autor,
		Total:     protoBook.Total,
		Remaining: protoBook.Total,
	})
}

func BookUpdateHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Mensagem recebida no tópico book/update: %s\n", msg.Payload())
	var protoBook api_bib.Livro
	err := json.Unmarshal(msg.Payload(), &protoBook)
	if err != nil {
		log.Printf("Erro ao converter dados do livro para JSON: %v", err)
		return
	}

	// Obtem livros restantes pois mensagem pode vir do cad server e não conter o campo Remaining
	bookToUpdate, err := database.ConcreteBookRepo.GetById(protoBook.Isbn)
	if err != nil {
		log.Printf("Erro ao obter livro pelo ISBN: %v", err)
		return
	}

	database.ConcreteBookRepo.Update(database.Book{
		ISBN:      protoBook.Isbn,
		Author:    protoBook.Autor,
		Total:     protoBook.Total,
		Title:     protoBook.Titulo,
		Remaining: bookToUpdate.Remaining,
	})
}

func BookRemoveHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Mensagem recebida no tópico book/remove: %s\n", msg.Payload())
	var protoBook api_bib.Livro
	err := json.Unmarshal(msg.Payload(), &protoBook)
	if err != nil {
		log.Printf("Erro ao converter dados do livro para JSON: %v", err)
		return
	}

	database.ConcreteBookRepo.Remove(string(protoBook.Isbn))
}
