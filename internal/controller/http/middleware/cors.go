package middleware

import (
	"net/http"

	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/rizkiar00/homework/pkg/config"
	"github.com/rizkiar00/homework/pkg/constant"
)

func CORS(cfg config.HTTPConfig) echoMiddleware.CORSConfig {
	return echoMiddleware.CORSConfig{
		AllowOrigins: cfg.AllowedOrigins(),
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			"Authorization",
			"Content-Type",
			constant.HeaderRequestID,
		},
		ExposeHeaders: []string{
			constant.HeaderRequestID,
		},
	}
}
