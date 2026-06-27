package config

import "fmt"

type DatabaseConfig struct {
	Name     string `env:"DB_NAME"`
	Driver   string `env:"DB_DRIVER" envDefault:"postgres"`
	Host     string `env:"DB_HOST"`
	Port     int    `env:"DB_PORT" envDefault:"5432"`
	Username string `env:"DB_USERNAME"`
	Password string `env:"DB_PASSWORD"`
	Schema   string `env:"DB_SCHEMA" envDefault:"public"`
}

func (c DatabaseConfig) IsConfigured() bool {
	return c.Host != "" && c.Name != "" && c.Username != ""
}

func (c DatabaseConfig) GenerateConnectionString() string {
	if !c.IsConfigured() {
		return ""
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?search_path=%s&sslmode=disable",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
		c.Schema,
	)
}
