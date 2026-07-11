package config

import "strings"

type AppConfig struct {
	Name                   string `env:"APP_NAME" envDefault:"homework-api"`
	Env                    string `env:"APP_ENV" envDefault:"local"`
	Host                   string `env:"APP_HOST" envDefault:"0.0.0.0"`
	Port                   string `env:"APP_PORT"`
	PlatformPort           string `env:"PORT"`
	ShutdownTimeoutSeconds int    `env:"APP_SHUTDOWN_TIMEOUT_SECONDS" envDefault:"10"`
}

func (c AppConfig) Address() string {
	return c.Host + ":" + c.ListenPort()
}

func (c AppConfig) ListenPort() string {
	if strings.TrimSpace(c.Port) != "" {
		return c.Port
	}
	if strings.TrimSpace(c.PlatformPort) != "" {
		return c.PlatformPort
	}

	return "8080"
}
