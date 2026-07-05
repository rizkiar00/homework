package resource

import (
	"os"
	"path/filepath"

	"github.com/rizkiar00/homework/pkg/config"
	"github.com/sirupsen/logrus"
)

func NewLogger(cfg config.Config) *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	if cfg.AppConfig.Env == "local" {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	if cfg.Log.FilePath != "" {
		if err := os.MkdirAll(filepath.Dir(cfg.Log.FilePath), 0755); err == nil {
			if file, err := os.OpenFile(cfg.Log.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); err == nil {
				logger.SetOutput(file)
			}
		}
	}

	return logger
}
