package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

type (
	Config struct {
		GRPC GRPC
		PG   PG
	}

	GRPC struct {
		Port        string `env:"GRPC_PORT" envDefault:"9090"`
		GatewayPort string `env:"GRPC_GATEWAY_PORT" envDefault:"8080"`
	}

	PG struct {
		Host     string `env:"POSTGRES_HOST" envDefault:"localhost"`
		Port     string `env:"POSTGRES_PORT" envDefault:"5432"`
		DB       string `env:"POSTGRES_DB" envDefault:"Library"`
		User     string `env:"POSTGRES_USER" envDefault:"user"`
		Password string `env:"POSTGRES_PASSWORD" envDefault:"00000000"`
		MaxConn  string `env:"POSTGRES_MAX_CONN" envDefault:"10"`
		URL      string
	}
)

func New() (*Config, error) {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse env params: %w", err)
	}

	cfg.PG.URL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.PG.User,
		cfg.PG.Password,
		cfg.PG.Host,
		cfg.PG.Port,
		cfg.PG.DB,
	)

	return cfg, nil
}
