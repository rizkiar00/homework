package middleware

import (
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/rizkiar00/homework/pkg/config"
)

func BodyLimit(cfg config.HTTPConfig) echoMiddleware.BodyLimitConfig {
	return echoMiddleware.BodyLimitConfig{
		Limit: cfg.BodyLimitValue(),
	}
}
