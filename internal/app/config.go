package app

import (
	"github.com/caarlos0/env/v7"
	"github.com/joho/godotenv"
)

type Config struct {
	Port        int    `env:"PORT"`
	PostgresDSN string `env:"POSTGRES_DSN"`
	APIURL      string `env:"API_URL"`
}

func NewConfig(envPath string) (c Config, err error) {
	err = godotenv.Load(envPath)
	if err != nil {
		return c, err
	}

	err = env.Parse(&c)
	if err != nil {
		return c, err
	}

	return c, nil
}
