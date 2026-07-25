package customerror

import (
	"errors"
	"net/http"

	"github.com/rizkiar00/homework/pkg/constant"
	"gorm.io/gorm"
)

type Error struct {
	Code       string
	Message    string
	HTTPStatus int
	Err        error
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

func New(code string, message string, httpStatus int) *Error {
	return &Error{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

func Wrap(err error, code string, message string, httpStatus int) *Error {
	return &Error{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		Err:        err,
	}
}

func BadRequest(message string) *Error {
	return New(constant.CodeBadRequest, message, http.StatusBadRequest)
}

func Unauthorized(message string) *Error {
	return New(constant.CodeUnauthorized, message, http.StatusUnauthorized)
}

func NotFound(message string) *Error {
	return New(constant.CodeNotFound, message, http.StatusNotFound)
}

func Conflict(message string) *Error {
	return New(constant.CodeConflict, message, http.StatusConflict)
}

func TooManyRequests(message string) *Error {
	return New(constant.CodeTooManyRequests, message, http.StatusTooManyRequests)
}

func Internal(err error) *Error {
	return Wrap(err, constant.CodeInternalServer, constant.MessageInternalServer, http.StatusInternalServerError)
}

func From(err error) *Error {
	var customErr *Error
	if errors.As(err, &customErr) {
		return customErr
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return NotFound(constant.MessageDataNotFound)
	}

	return Internal(err)
}
