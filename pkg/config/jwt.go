package config

import "time"

type JWTConfig struct {
	Secret           string `env:"JWT_SECRET" envDefault:"local-secret"`
	ExpiresInSeconds int    `env:"JWT_EXPIRES_IN_SECONDS" envDefault:"3600"`
}

func (c JWTConfig) ExpiresIn() time.Duration {
	return time.Duration(c.ExpiresInSeconds) * time.Second
}
