package auth

import (
	"context"

	"github.com/rizkiar00/homework/internal/model"
)

type Usecase interface {
	Register(ctx context.Context, request model.RegisterRequest) (model.AuthUserResponse, error)
	Login(ctx context.Context, request model.LoginRequest) (model.LoginResponse, error)
	Me(ctx context.Context, userID string) (model.AuthUserResponse, error)
}
