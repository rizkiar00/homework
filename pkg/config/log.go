package config

type LogConfig struct {
	FilePath string `env:"LOG_FILE_PATH" envDefault:"logs/app.log"`
}
