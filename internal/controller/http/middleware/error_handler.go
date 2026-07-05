package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rizkiar00/homework/pkg/constant"
	httpresponse "github.com/rizkiar00/homework/pkg/response"
	"github.com/sirupsen/logrus"
)

func HTTPErrorHandler(logger *logrus.Logger) echo.HTTPErrorHandler {
	return func(err error, ctx echo.Context) {
		if ctx.Response().Committed {
			return
		}

		status := http.StatusInternalServerError
		code := constant.CodeInternalServer
		message := constant.MessageInternalServer

		if echoErr, ok := err.(*echo.HTTPError); ok {
			status = echoErr.Code
			message = http.StatusText(echoErr.Code)
			if text, ok := echoErr.Message.(string); ok && text != "" {
				message = text
			}
		}

		switch status {
		case http.StatusBadRequest:
			code = constant.CodeBadRequest
		case http.StatusUnauthorized:
			code = constant.CodeUnauthorized
		case http.StatusNotFound:
			code = constant.CodeNotFound
		case http.StatusConflict:
			code = constant.CodeConflict
		case http.StatusServiceUnavailable:
			code = constant.CodeServiceUnavailable
		}

		logger.WithFields(logrus.Fields{
			"request_id": GetRequestID(ctx),
			"method":     ctx.Request().Method,
			"path":       ctx.Request().URL.Path,
			"status":     status,
			"error":      err.Error(),
		}).Error("http error")

		if writeErr := httpresponse.Error(ctx, status, code, message); writeErr != nil {
			logger.WithError(writeErr).Error("failed to write error response")
		}
	}
}
