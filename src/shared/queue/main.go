package queue

import (
	"crypto/tls"
	"fmt"

	"library-manager/shared/queue/handlers"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const qos byte = 2

var topicHandlers map[string]func(client mqtt.Client, msg mqtt.Message) = map[string]func(client mqtt.Client, msg mqtt.Message){
	"user/create": handlers.CreateUser,
	"user/update": handlers.UpdateUser,
	"user/remove": handlers.RemoveUser,
	"book/create": handlers.CreateBook,
	"book/update": handlers.UpdateBook,
	"book/remove": handlers.RemoveBook,
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
	clientId := fmt.Sprintf("cad-server-%d", id)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ClientAuth:         tls.NoClientCert,
	}
	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:%d", host, port)).SetClientID(clientId).SetOnConnectHandler(onConnect).SetTLSConfig(tlsConfig)
	return mqtt.NewClient(opts)
}
