package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	AppConfig AppConfig
	Database  DatabaseConfig
	HTTP      HTTPConfig
	JWT       JWTConfig
	Log       LogConfig
}

func New() (Config, error) {
	_ = godotenv.Load()

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
