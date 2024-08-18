package queue

import (
	"crypto/tls"
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const qos byte = 2

type Topic string

const (
	BookLoanTopic         Topic = "book/loan"
	BookReturnTopic       Topic = "book/return"
	BookListBorrowedTopic Topic = "book/list-borrowed"
	BookListMissingTopic  Topic = "book/list-missing"
	BookSearchTopic       Topic = "book/search"
	UserBlockTopic        Topic = "user/block"
	UserFreeTopic         Topic = "user/free"
	UserListBlockedTopic  Topic = "user/list-blocked"
)

func logTopicMessage(topic Topic, msg []byte) {
	log.Printf("Mensagem recebida no tópico %s: %s\n", topic, msg)
}

var topicHandlers map[string]func(client mqtt.Client, msg mqtt.Message) = map[string]func(client mqtt.Client, msg mqtt.Message){
	string(BookLoanTopic): func(client mqtt.Client, msg mqtt.Message) {
		logTopicMessage(BookLoanTopic, msg.Payload())
	},
	string(BookReturnTopic): func(client mqtt.Client, msg mqtt.Message) {
		logTopicMessage(BookReturnTopic, msg.Payload())
	},
	string(BookListBorrowedTopic): func(client mqtt.Client, msg mqtt.Message) {
		logTopicMessage(BookListBorrowedTopic, msg.Payload())
	},
	string(BookListMissingTopic): func(client mqtt.Client, msg mqtt.Message) {
		logTopicMessage(BookListMissingTopic, msg.Payload())
	},
	string(BookSearchTopic): func(client mqtt.Client, msg mqtt.Message) {
		logTopicMessage(BookSearchTopic, msg.Payload())
	},
	string(UserBlockTopic): func(client mqtt.Client, msg mqtt.Message) {
		logTopicMessage(UserBlockTopic, msg.Payload())
	},
	string(UserFreeTopic): func(client mqtt.Client, msg mqtt.Message) {
		logTopicMessage(UserFreeTopic, msg.Payload())
	},
	string(UserListBlockedTopic): func(client mqtt.Client, msg mqtt.Message) {
		logTopicMessage(UserListBlockedTopic, msg.Payload())
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
