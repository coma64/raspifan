package main

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/bcm283x"
	"raspifan/client"
)

var (
	pinAlwaysOn = bcm283x.GPIO2
	pinPwm      = bcm283x.GPIO12
)

func setFanSpeed(fanSpeed gpio.Duty) {
	if err := pinPwm.PWM(fanSpeed, physic.KiloHertz); err != nil {
		panic(err)
	}
}

func calcFanSpeed(temperature client.Celsius) int64 {
	speed := int64((temperature - 20*client.OneDegree) * 5)
	if speed > 100*1_000 {
		speed = 100 * 1_000
	} else if speed < 0 {
		speed = 0
	}

	return speed
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if state, err := host.Init(); err != nil {
		panic(err)
	} else {
		for _, driver := range state.Loaded {
			fmt.Printf("Loaded driver: %v\n", driver.String())
		}
	}

	if err := client.Connect(); err != nil {
		panic(err)
	}
	defer client.Disconnect()

	temperatures, err := client.SubscribeToTemperatureTopic()
	if err != nil {
		panic(err)
	}

	if err = pinAlwaysOn.Out(gpio.High); err != nil {
		panic(err)
	}

	for temperature := range temperatures {
		fanSpeed := calcFanSpeed(temperature)

		log.Debug().
			Int64("temperature", int64(temperature)).
			Int64("newFanSpeed", fanSpeed).
			Msg("Received temperature")

		setFanSpeed(gpio.DutyMax / 100_000 * gpio.Duty(fanSpeed))
	}
}
