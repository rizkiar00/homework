package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/rizkiar00/homework/pkg/config"
	"github.com/rizkiar00/homework/pkg/constant"
	httpresponse "github.com/rizkiar00/homework/pkg/response"
)

func Timeout(cfg config.HTTPConfig) echoMiddleware.ContextTimeoutConfig {
	return echoMiddleware.ContextTimeoutConfig{
		Timeout: cfg.Timeout(),
		ErrorHandler: func(err error, ctx echo.Context) error {
			if err != nil && !errors.Is(err, context.DeadlineExceeded) {
				return err
			}

			return httpresponse.Error(ctx, http.StatusServiceUnavailable, constant.CodeServiceUnavailable, constant.MessageServiceUnavailable)
		},
	}
}
