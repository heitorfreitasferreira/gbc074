package queue

import (
	"crypto/tls"
	"fmt"
	"library-manager/bib-server/internal/queue/handlers"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const qos byte = 2

var topicHandlers map[string]func(client mqtt.Client, msg mqtt.Message) = map[string]func(client mqtt.Client, msg mqtt.Message){
	string(handlers.BookLoanTopic):         handlers.BookLoanHandler,
	string(handlers.BookReturnTopic):       handlers.BookReturnHandler,
	string(handlers.BookListBorrowedTopic): handlers.BookListBorrowedHandler,
	string(handlers.BookListMissingTopic):  handlers.BookListMissingHandler,
	string(handlers.BookSearchTopic):       handlers.BookSearchHandler,
	string(handlers.UserBlockTopic):        handlers.UserBlockHandler,
	string(handlers.UserFreeTopic):         handlers.UserFreeHandler,
	string(handlers.UserListBlockedTopic):  handlers.UserListBlockedHandler,
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
