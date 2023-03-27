package config

import "github.com/jinzhu/configor"

type config struct {
	Broker struct {
		URL   string
		Topic string
	}
	FanSpeed int `yaml:"fan-speed"`
}

var Config = &config{}

func init() {
	if err := configor.Load(Config, "config.yml"); err != nil {
		panic(err)
	}
}
