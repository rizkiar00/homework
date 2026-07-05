package response

import (
	"github.com/labstack/echo/v4"
	"github.com/rizkiar00/homework/pkg/customerror"
)

func CustomError(ctx echo.Context, err error) error {
	customErr := customerror.From(err)
	return Error(ctx, customErr.HTTPStatus, customErr.Code, customErr.Message)
}
