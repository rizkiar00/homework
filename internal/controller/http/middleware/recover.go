package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rizkiar00/homework/pkg/constant"
	httpresponse "github.com/rizkiar00/homework/pkg/response"
	"github.com/sirupsen/logrus"
)

func Recover(logger *logrus.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			defer func() {
				if recovered := recover(); recovered != nil {
					logger.WithFields(logrus.Fields{
						"request_id": GetRequestID(ctx),
						"method":     ctx.Request().Method,
						"path":       ctx.Request().URL.Path,
						"panic":      recovered,
					}).Error("panic recovered")

					if ctx.Response().Committed {
						err = nil
						return
					}

					err = httpresponse.Error(ctx, http.StatusInternalServerError, constant.CodeInternalServer, constant.MessageInternalServer)
				}
			}()

			return next(ctx)
		}
	}
}
