package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
)

type (
	Config struct {
		RabbitConfig    RabbitConfig    `envPrefix:"RMQ_"`
		LoggerConfig    LoggerConfig    `envPrefix:"LOG_"`
		PostreSqlConfig PostreSqlConfig `envPrefix:"PG_"`
	}

	RabbitConfig struct {
		Host              string `env:"HOST" envDefault:"localhost"`
		Port              int    `env:"PORT"  envDefault:"5672"`
		Username          string `env:"USERNAME" envDefault:"guest"`
		Password          string `env:"PASSWORD" envDefault:"guest"`
		VirtualHost       string `env:"VIRTUAL_HOST" envDefault:"/"`
		ReconnectAttempts uint   `env:"RECONNECT_ATTEMPTS" envDefault:"5"`
	}

	LoggerConfig struct {
		Level string `env:"LEVEL" envDefault:"INFO"`
	}

	PostreSqlConfig struct {
		Host     string `env:"HOST" envDefault:"localhost"`
		Port     int    `env:"PORT"  envDefault:"5432"`
		Username string `env:"USERNAME" envDefault:"postgres"`
		Password string `env:"PASSWORD" envDefault:"password"`
		DbName   string `env:"DBNAME" envDefault:"BalanceService"`
	}
)

func GetConfig() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, errors.New("Can't parse config")
	}
	return &cfg, nil
}
