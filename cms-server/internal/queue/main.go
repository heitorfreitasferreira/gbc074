package queue

import (
	"crypto/tls"
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const qos byte = 2

var topicHandlers map[string]func(client mqtt.Client, msg mqtt.Message) = map[string]func(client mqtt.Client, msg mqtt.Message){
	"book/loan": func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("Mensagem recebida no tópico book/loan: %s\n", msg.Payload())
	},
	"book/return": func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("Mensagem recebida no tópico book/return: %s\n", msg.Payload())
	},
	"book/list-borrowed": func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("Mensagem recebida no tópico book/list-borrowed: %s\n", msg.Payload())
	},
	"book/list-missing": func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("Mensagem recebida no tópico book/missing: %s\n", msg.Payload())
	},
	"book/search": func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("Mensagem recebida no tópico book/search: %s\n", msg.Payload())
	},
	"user/block": func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("Mensagem recebida no tópico user/block: %s\n", msg.Payload())
	},
	"user/free": func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("Mensagem recebida no tópico user/free: %s\n", msg.Payload())
	},
	"user/list-blocked": func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("Mensagem recebida no tópico user/list-blocked: %s\n", msg.Payload())
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
