package resource

import (
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

	return logger
}
