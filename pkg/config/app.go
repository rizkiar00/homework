package config

type AppConfig struct {
	Name                   string `env:"APP_NAME" envDefault:"homework-api"`
	Env                    string `env:"APP_ENV" envDefault:"local"`
	Host                   string `env:"APP_HOST" envDefault:"0.0.0.0"`
	Port                   string `env:"APP_PORT" envDefault:"8080"`
	ShutdownTimeoutSeconds int    `env:"APP_SHUTDOWN_TIMEOUT_SECONDS" envDefault:"10"`
}

func (c AppConfig) Address() string {
	return c.Host + ":" + c.Port
}
