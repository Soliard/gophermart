package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerHost      string `env:"ADDRESS"`
	LogLevel        string `env:"LOG_LEVEL"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	TokenSecret     string `env:"TOKEN_SECRET"`
	TokenExpMinutes int    `env:"TOKEN_EXP"`
}

func New() (*Config, error) {
	config := &Config{}

	flag.StringVar(&config.ServerHost, "a", "localhost:8080", "server addres")
	flag.StringVar(&config.LogLevel, "l", "debug", "log level")
	//postgres://postgres:postgres@localhost:5432/gotplmetrics?sslmode=disable
	flag.StringVar(&config.DatabaseDSN, "d", "", "database connection string")
	flag.StringVar(&config.TokenSecret, "s", "gigasecret", "key will be used for jwt")
	flag.IntVar(&config.TokenExpMinutes, "e", 1, "time in minutes to token expiring")
	flag.Parse()

	err := env.Parse(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
