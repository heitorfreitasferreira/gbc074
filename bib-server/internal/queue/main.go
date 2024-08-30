package queue

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"

	"library-manager/shared/database"
	"library-manager/shared/queue/handlers"
	"library-manager/shared/utils"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const qos byte = 2

var topicHandlers map[string]func(client mqtt.Client, msg mqtt.Message) = map[string]func(client mqtt.Client, msg mqtt.Message){
	string(handlers.BookLoanTopic): func(client mqtt.Client, msg mqtt.Message) {
		payload := msg.Payload()
		handlers.LogTopicMessage(handlers.BookLoanTopic, payload)
		var userBook database.UserBook
		err := json.Unmarshal(payload, &userBook)
		if err != nil {
			log.Printf("Erro ao converter dados do payload: %v", err)
			return
		}

		book, err := database.ConcreteBookRepo.ObtemLivro(utils.ISBN(userBook.BookISNB))
		if err != nil {
			log.Printf("Livro inexistente: %v", err)
			return
		}
		if book.Total == 0 {
			log.Printf("Livro sem exemplares disponíveis")
			return
		}

		_, err = database.ConcreteUserRepo.ObtemUsuario(utils.CPF(userBook.UserId))
		if err != nil {
			log.Printf("Usuário inexistente: %v", err)
			return
		}

		// TODO: Verificar se usuário está bloqueado
		database.ConcreteUserBookRepo.LoanBook(userBook)
	},
	string(handlers.BookReturnTopic): func(client mqtt.Client, msg mqtt.Message) {
		handlers.LogTopicMessage(handlers.BookReturnTopic, msg.Payload())
	},
	string(handlers.BookListBorrowedTopic): func(client mqtt.Client, msg mqtt.Message) {
		handlers.LogTopicMessage(handlers.BookListBorrowedTopic, msg.Payload())


	},
	string(handlers.BookListMissingTopic): func(client mqtt.Client, msg mqtt.Message) {
		handlers.LogTopicMessage(handlers.BookListMissingTopic, msg.Payload())
	},
	string(handlers.BookSearchTopic): func(client mqtt.Client, msg mqtt.Message) {
		handlers.LogTopicMessage(handlers.BookSearchTopic, msg.Payload())
	},
	string(handlers.UserBlockTopic): func(client mqtt.Client, msg mqtt.Message) {
		handlers.LogTopicMessage(handlers.UserBlockTopic, msg.Payload())
	},
	string(handlers.UserFreeTopic): func(client mqtt.Client, msg mqtt.Message) {
		handlers.LogTopicMessage(handlers.UserFreeTopic, msg.Payload())
	},
	string(handlers.UserListBlockedTopic): func(client mqtt.Client, msg mqtt.Message) {
		handlers.LogTopicMessage(handlers.UserListBlockedTopic, msg.Payload())
	},
}

var onConnect mqtt.OnConnectHandler = func(client mqtt.Client) {
	for topic, handler := range topicHandlers {
		if token := client.Subscribe(topic, qos, handler); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}
}

func GetMqttBroker(host string, port int, id int) mqtt.Client {
	// Use id (received as PID) to generate unique client identifier
	clientId := fmt.Sprintf("bib-server-%d", id)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ClientAuth:         tls.NoClientCert,
	}
	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:%d", host, port)).SetClientID(clientId).SetOnConnectHandler(onConnect).SetTLSConfig(tlsConfig)
	return mqtt.NewClient(opts)
}
