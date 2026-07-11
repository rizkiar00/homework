package config

import (
	"fmt"
	"time"
)

type RedisConfig struct {
	Enabled        bool   `env:"REDIS_ENABLED" envDefault:"false"`
	Host           string `env:"REDIS_HOST" envDefault:"127.0.0.1"`
	Port           int    `env:"REDIS_PORT" envDefault:"6379"`
	Username       string `env:"REDIS_USERNAME"`
	Password       string `env:"REDIS_PASSWORD"`
	DB             int    `env:"REDIS_DB" envDefault:"0"`
	DialTimeoutSec int    `env:"REDIS_DIAL_TIMEOUT_SECONDS" envDefault:"5"`
}

func (c RedisConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c RedisConfig) DialTimeout() time.Duration {
	if c.DialTimeoutSec <= 0 {
		return 5 * time.Second
	}

	return time.Duration(c.DialTimeoutSec) * time.Second
}
