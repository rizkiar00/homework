package model

import "errors"

var (
	ErrInvalidCredential     = errors.New("invalid username or password")
	ErrUsernameAlreadyExists = errors.New("username already exists")
)
