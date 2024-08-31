package handlers

import (
	"encoding/json"
	"library-manager/bib-server/internal/database"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func logTopicMessage(topic UserBookTopic, msg []byte) {
	log.Printf("Mensagem recebida no tópico %s: %s\n", topic, msg)
}

func BookLoanHandler(client mqtt.Client, msg mqtt.Message) {
	payload := msg.Payload()
	logTopicMessage(BookLoanTopic, payload)
	var userBook database.UserBook
	err := json.Unmarshal(payload, &userBook)
	if err != nil {
		log.Printf("Erro ao converter dados do payload: %v", err)
		return
	}

	book, err := database.ConcreteBookRepo.GetById(userBook.BookISBN)
	if err != nil {
		log.Printf("Livro inexistente: %v", err)
		return
	}
	if book.Total == 0 || book.Remaining < 1 {
		log.Printf("Livro sem exemplares disponíveis")
		return
	}

	user, err := database.ConcreteUserRepo.GetById(userBook.UserCPF)
	if err != nil {
		log.Printf("Usuário inexistente: %v", err)
		return
	}

	if user.Blocked {
		log.Printf("Usuário bloqueado")
		return
	}

	database.ConcreteUserBookRepo.Create(userBook)
}

func BookReturnHandler(client mqtt.Client, msg mqtt.Message) {
	payload := msg.Payload()
	logTopicMessage(BookReturnTopic, payload)

	var userBook database.UserBook
	json.Unmarshal(payload, &userBook)
	err := json.Unmarshal(payload, &userBook)
	if err != nil {
		log.Printf("Erro ao converter dados do payload: %v", err)
		return
	}

	book, err := database.ConcreteBookRepo.GetById(userBook.BookISBN)
	if err != nil {
		log.Printf("Erro ao obter livro: %v", err)
		return
	}

	err = database.ConcreteBookRepo.Update(database.Book{
		ISBN:      book.ISBN,
		Title:     book.Title,
		Author:    book.Author,
		Total:     book.Total,
		Remaining: book.Remaining + 1,
	})
	if err != nil {
		log.Printf("Erro ao atualizar livro: %v", err)
		return
	}

	err = database.ConcreteUserBookRepo.RemoveUserLoan(userBook)
	if err != nil {
		log.Printf("Erro ao remover empréstimo: %v", err)
	}

}

func BookListBorrowedHandler(client mqtt.Client, msg mqtt.Message) {
	logTopicMessage(BookListBorrowedTopic, msg.Payload())

}

func BookListMissingHandler(client mqtt.Client, msg mqtt.Message) {
	logTopicMessage(BookListMissingTopic, msg.Payload())
}

func BookSearchHandler(client mqtt.Client, msg mqtt.Message) {
	logTopicMessage(BookSearchTopic, msg.Payload())
}

func UserBlockHandler(client mqtt.Client, msg mqtt.Message) {
	logTopicMessage(UserBlockTopic, msg.Payload())

	currTime := time.Now().UnixMilli()
	userBookList := database.ConcreteUserBookRepo.GetAll()
	for _, user := range userBookList {
		userLoans := database.ConcreteUserBookRepo.GetUserLoans(user.UserCPF)
		for _, loan := range userLoans {
			// Bloqueira usuário com empréstimo maior que 10 segundos
			if currTime-loan.Timestamp > 10*1000 {
				user, err := database.ConcreteUserRepo.GetById(user.UserCPF)
				if err != nil {
					log.Printf("Erro ao obter usuário: %v", err)
					return
				}

				database.ConcreteUserRepo.Update(database.User{
					CPF:     user.CPF,
					Name:    user.Name,
					Blocked: true,
				})
				break
			}
		}
	}
}

func UserFreeHandler(client mqtt.Client, msg mqtt.Message) {
	logTopicMessage(UserFreeTopic, msg.Payload())

	currTime := time.Now().UnixMilli()
	userList := database.ConcreteUserRepo.GetAll()
	for _, user := range userList {
		if user.Blocked {
			userLoans := database.ConcreteUserBookRepo.GetUserLoans(user.CPF)
			shouldUnblock := true
			for _, loan := range userLoans {
				if currTime-loan.Timestamp > 10*1000 {
					shouldUnblock = false
					break
				}
			}

			if shouldUnblock {
				database.ConcreteUserRepo.Update(database.User{
					CPF:     user.CPF,
					Name:    user.Name,
					Blocked: false,
				})
			}
		}
	}
}

func UserListBlockedHandler(client mqtt.Client, msg mqtt.Message) {
	logTopicMessage(UserListBlockedTopic, msg.Payload())
}
