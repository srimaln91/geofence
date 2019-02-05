package mqtt

import (
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/xid"
)

type MQTTClient struct {
	Id     string
	Client mqtt.Client
}

func New(brokerAddr string) MQTTClient {

	guid := xid.New().String()
	opts := createClientOptions(guid, brokerAddr)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}

	return MQTTClient{
		Id:     guid,
		Client: client,
	}

}

func createClientOptions(clientId string, uri string) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(uri)
	// opts.SetUsername(uri.User.Username())
	// password, _ := uri.User.Password()
	// opts.SetPassword(password)
	opts.SetClientID(clientId)
	return opts
}

func (mc *MQTTClient) PublishMessage(topic string, message string) mqtt.Token {
	return mc.Client.Publish(topic, 0, false, message)
}

func (mc *MQTTClient) SubscribeMessage(topic string, msgChan chan mqtt.Message) {

	mc.Client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		// fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
		msgChan <- msg
	})

}

func (mc *MQTTClient) GetMessageChannel() chan mqtt.Message {
	return make(chan mqtt.Message)
}
