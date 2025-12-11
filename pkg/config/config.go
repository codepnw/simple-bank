package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/codepnw/simple-bank/pkg/utils/helper"
	"github.com/joho/godotenv"
)

type EnvConfig struct {
	Server ServerConfig `envPrefix:"SERVER_"`
	DB     DBConfig     `envPrefix:"DB_"`
	JWT    JWTConfig    `envPrefix:"JWT_"`
	Paseto PasetoConfig `envPrefix:"PASETO_"`
}

type ServerConfig struct {
	HTTPAddr   string `env:"HTTP_ADDRESS" envDefault:"localhost:8080"`
	HTTPPrefix string `env:"HTTP_PREFIX" envDefault:"/api/v1"`
	GrpcAddr   string `env:"GRPC_ADDRESS" envDefault:"localhost:9090"`
}

type DBConfig struct {
	User     string `env:"USER" validate:"required"`
	Password string `env:"PASSWORD"`
	Host     string `env:"HOST" envDefault:"localhost"`
	Port     int    `env:"PORT" envDefault:"5432"`
	DBName   string `env:"NAME" validate:"required"`
	SSLMode  string `env:"SSL_MODE" envDefault:"disable"`
}

type JWTConfig struct {
	SecretKey  string `env:"SECRET_KEY" validate:"required"`
	RefreshKey string `env:"REFRESH_KEY" validate:"required"`
}

type PasetoConfig struct {
	SymmetricKey string `env:"SYMMETRIC_KEY" validate:"required,len=32"`
}

func LoadEnv(path string) (*EnvConfig, error) {
	godotenv.Load(path)

	cfg := new(EnvConfig)
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parse env failed: %w", err)
	}

	if err := helper.Validate(cfg); err != nil {
		return nil, fmt.Errorf("validate env failed: %w", err)
	}
	return cfg, nil
}
