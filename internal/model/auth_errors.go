package model

import "errors"

import "github.com/rizkiar00/homework/pkg/constant"

var (
	ErrInvalidCredential     = errors.New(constant.MessageInvalidAuthCredential)
	ErrUsernameAlreadyExists = errors.New(constant.MessageUsernameAlreadyExists)
)
