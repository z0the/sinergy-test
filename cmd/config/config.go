package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string `env:"SERVER_PORT"`
	ServerHost string `env:"SERVER_HOST"`
	ClientPort string `env:"CLIENT_PORT"`
}

func (c *Config) init(envPath string) {
	isDev := flag.Bool("DEV", false, "")
	flag.Parse()
	if *isDev {
		err := godotenv.Load(envPath)
		if err != nil {
			panic(err)
		}
	}

	err := env.Parse(c)
	if err != nil {
		panic(err)
	}
}

func GetConfig(envPath string) Config {
	cfg := Config{}
	cfg.init(envPath)

	return cfg
}
