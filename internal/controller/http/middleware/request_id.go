package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rizkiar00/homework/pkg/constant"
)

func RequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			requestID := ctx.Request().Header.Get(constant.HeaderRequestID)
			if requestID == "" {
				requestID = uuid.NewString()
			}

			ctx.Set(constant.ContextRequestIDKey, requestID)
			ctx.Response().Header().Set(constant.HeaderRequestID, requestID)

			return next(ctx)
		}
	}
}

func GetRequestID(ctx echo.Context) string {
	requestID, ok := ctx.Get(constant.ContextRequestIDKey).(string)
	if !ok {
		return ""
	}

	return requestID
}
