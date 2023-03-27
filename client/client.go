package client

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"raspifan/config"
	"strconv"
)

var client = mqtt.NewClient(mqtt.NewClientOptions().AddBroker(config.Config.Broker.URL))

func Connect() error {
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func Disconnect() {
	client.Disconnect(1000)
}

func SubscribeToTemperatureTopic() (chan float64, error) {
	temperatures := make(chan float64)
	token := client.Subscribe(config.Config.Broker.Topic, 2, func(_ mqtt.Client, message mqtt.Message) {
		value, err := strconv.Atoi(string(message.Payload()))
		if err != nil {
			return
		}

		temperatures <- float64(value) / 1000
	})

	if err := token.Error(); token.Wait() && err != nil {
		return nil, err
	}

	return temperatures, nil
}
