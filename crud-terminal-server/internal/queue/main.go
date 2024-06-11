package queue

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rpc-mqtt-library-manager/crud-terminal-server/internal/database"
)

const qos byte = 2

var topicHandlers map[string]func(client mqtt.Client, msg mqtt.Message) = map[string]func(client mqtt.Client, msg mqtt.Message){
	"user/create": func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("Mensagem recebida no tópico user/create: %s\n", msg.Payload())
		var user database.User
		err := json.Unmarshal(msg.Payload(), &user)
		if err != nil {
			log.Printf("Erro ao converter dados do usuário para JSON: %v", err)
			return
		}
		database.ConcreteUserRepo.CreateUser(user)
	},
	"user/update": func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("Mensagem recebida no tópico user/update: %s\n", msg.Payload())
		var user database.User
		err := json.Unmarshal(msg.Payload(), &user)
		if err != nil {
			log.Printf("Erro ao converter dados do usuário para JSON: %v", err)
			return
		}
		database.ConcreteUserRepo.EditaUsuario(user)
	},
	"user/remove": func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("Mensagem recebida no tópico user/remove: %s\n", msg.Payload())
		var user database.User
		err := json.Unmarshal(msg.Payload(), &user)
		if err != nil {
			log.Printf("Erro ao converter dados do usuário para JSON: %v", err)
			return
		}
		database.ConcreteUserRepo.RemoveUsuario(user.Cpf)
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
	clientId := fmt.Sprintf("cad-server-%d", id)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ClientAuth:         tls.NoClientCert,
	}
	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:%d", host, port)).SetClientID(clientId).SetOnConnectHandler(onConnect).SetTLSConfig(tlsConfig)
	return mqtt.NewClient(opts)
}
