package config

import "github.com/jinzhu/configor"

type config struct {
	Broker struct {
		URL   string
		Topic string
	}
}

var Config = &config{}

func init() {
	if err := configor.Load(Config, "config.yml"); err != nil {
		panic(err)
	}
}
