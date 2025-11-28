package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/codepnw/simple-bank/pkg/utils/validate"
	"github.com/joho/godotenv"
)

type EnvConfig struct {
	DB DBConfig `envPrefix:"DB_"`
}

type DBConfig struct {
	User     string `env:"USER" validate:"required"`
	Password string `env:"PASSWORD"`
	Host     string `env:"HOST" envDefault:"localhost"`
	Port     int    `env:"PORT" envDefault:"5432"`
	DBName   string `env:"NAME" validate:"required"`
	SSLMode  string `env:"SSL_MODE" envDefault:"disable"`
}

func LoadEnv(path string) (*EnvConfig, error) {
	godotenv.Load(path)

	cfg := new(EnvConfig)
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parse env failed: %w", err)
	}

	if err := validate.Struct(cfg); err != nil {
		return nil, fmt.Errorf("validate env failed: %w", err)
	}
	return cfg, nil
}
