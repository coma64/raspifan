package main

import (
	"fmt"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/rpi"
	"raspifan/client"
	"raspifan/config"
)

var (
	pinAlwaysOn = rpi.P1_33
	pinPwm      = rpi.P1_32
)

func turnFanOn() {
	if err := pinPwm.PWM(gpio.DutyMax/100*gpio.Duty(config.Config.FanSpeed), physic.KiloHertz); err != nil {
		panic(err)
	}
}

func turnFanOff() {
	if err := pinPwm.PWM(gpio.Duty(0), physic.KiloHertz); err != nil {
		panic(err)
	}
}

func main() {
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

	isFanOn := false
	average := <-temperatures
	for temperature := range temperatures {
		average -= average / 10
		average += temperature / 10

		fmt.Printf("Received temperature: %v C; Current average: %v C\n", temperature, average)

		if !isFanOn && average >= 25 {
			isFanOn = true
			turnFanOn()
		} else if isFanOn && average < 25 {
			isFanOn = false
			turnFanOff()
		}
	}
}
