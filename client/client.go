package client

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
	"raspifan/config"
	"strconv"
)

type Celsius int64

const OneDegree = Celsius(1_000)

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

func SubscribeToTemperatureTopic() (chan Celsius, error) {
	temperatures := make(chan Celsius)
	token := client.Subscribe(config.Config.Broker.Topic, 2, func(_ mqtt.Client, message mqtt.Message) {
		value, err := strconv.ParseInt(string(message.Payload()), 10, 0)
		if err != nil {
			return
		}

		temperatures <- Celsius(value)
	})

	if err := token.Error(); token.Wait() && err != nil {
		return nil, err
	}
	log.Info().
		Str("broker", config.Config.Broker.URL).
		Str("topic", config.Config.Broker.Topic).
		Msg("Subscribed to topic")

	return temperatures, nil
}
