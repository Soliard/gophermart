package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerHost      string `env:"RUN_ADDRESS"`
	LogLevel        string `env:"LOG_LEVEL"`
	DatabaseDSN     string `env:"DATABASE_URI"`
	TokenSecret     string `env:"TOKEN_SECRET"`
	TokenExpMinutes int    `env:"TOKEN_EXP"`
	AccrualAddress  string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func New() (*Config, error) {
	config := &Config{}

	flag.StringVar(&config.ServerHost, "a", "localhost:8080", "server addres")
	flag.StringVar(&config.LogLevel, "l", "debug", "log level")
	flag.StringVar(&config.DatabaseDSN, "d", "postgres://postgres:postgres@localhost:5432/gophermart?sslmode=disable", "database connection string")
	flag.StringVar(&config.TokenSecret, "s", "gigasecret", "key will be used for jwt")
	flag.IntVar(&config.TokenExpMinutes, "e", 10, "time in minutes to token expiring")
	flag.StringVar(&config.AccrualAddress, "r", "localhost:5050", "address accural system")
	flag.Parse()

	err := env.Parse(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
