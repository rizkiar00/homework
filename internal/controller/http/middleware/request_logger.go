package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func RequestLogger(logger *logrus.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			start := time.Now()
			err := next(ctx)
			if err != nil {
				ctx.Error(err)
			}

			req := ctx.Request()
			res := ctx.Response()
			logger.WithFields(logrus.Fields{
				"request_id": GetRequestID(ctx),
				"method":     req.Method,
				"path":       req.URL.Path,
				"status":     res.Status,
				"latency_ms": time.Since(start).Milliseconds(),
				"remote_ip":  ctx.RealIP(),
				"user_agent": req.UserAgent(),
			}).Info("http request completed")

			return nil
		}
	}
}
