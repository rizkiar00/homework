package auth

import (
	"context"

	"github.com/rizkiar00/homework/internal/model"
	"github.com/rizkiar00/homework/pkg/token"
)

type Usecase interface {
	Register(ctx context.Context, request model.RegisterRequest) (model.AuthUserResponse, error)
	Login(ctx context.Context, request model.LoginRequest) (model.LoginResponse, error)
	VerifyEmail(ctx context.Context, request model.VerifyEmailRequest) error
	ResendVerification(ctx context.Context, request model.ResendVerificationRequest) error
	ForgotPassword(ctx context.Context, request model.ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, request model.ResetPasswordRequest) error
	Logout(ctx context.Context, claims token.Claims) error
	Me(ctx context.Context, userID string) (model.AuthUserResponse, error)
}
