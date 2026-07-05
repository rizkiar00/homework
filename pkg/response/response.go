package response

import "github.com/labstack/echo/v4"

type Body struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

func JSON(ctx echo.Context, status int, code string, message string, data interface{}) error {
	return ctx.JSON(status, Body{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func JSONWithMeta(ctx echo.Context, status int, code string, message string, data interface{}, meta interface{}) error {
	return ctx.JSON(status, Body{
		Code:    code,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func Error(ctx echo.Context, status int, code string, message string) error {
	return ctx.JSON(status, Body{
		Code:    code,
		Message: message,
	})
}
