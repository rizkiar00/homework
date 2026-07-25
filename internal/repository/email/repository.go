package email

import "context"

type Repository interface {
	SendVerificationCode(ctx context.Context, to string, name string, code string) error
	SendPasswordResetCode(ctx context.Context, to string, name string, code string) error
}
