package user

import (
	"context"
	"time"

	"github.com/rizkiar00/homework/internal/entity"
)

type Repository interface {
	Create(ctx context.Context, data entity.User) (entity.User, error)
	FindByUsername(ctx context.Context, username string) (entity.User, error)
	FindByEmail(ctx context.Context, email string) (entity.User, error)
	FindByID(ctx context.Context, id string) (entity.User, error)
	CreatePendingRegistration(ctx context.Context, user entity.User, verification entity.UserEmailVerification) (entity.User, error)
	UpdatePendingRegistration(ctx context.Context, userID string, data entity.User, verification entity.UserEmailVerification) (entity.User, error)
	FindActiveVerificationByEmail(ctx context.Context, email string) (entity.UserEmailVerification, error)
	FindLatestVerificationByEmail(ctx context.Context, email string) (entity.UserEmailVerification, error)
	CountVerificationsByEmailSince(ctx context.Context, email string, since time.Time) (int64, error)
	IncrementVerificationAttempt(ctx context.Context, verificationID string) error
	VerifyEmail(ctx context.Context, userID string, verificationID string) error
	CreatePasswordReset(ctx context.Context, reset entity.UserPasswordReset) error
	FindActivePasswordResetByEmail(ctx context.Context, email string) (entity.UserPasswordReset, error)
	FindLatestPasswordResetByEmail(ctx context.Context, email string) (entity.UserPasswordReset, error)
	CountPasswordResetsByEmailSince(ctx context.Context, email string, since time.Time) (int64, error)
	IncrementPasswordResetAttempt(ctx context.Context, resetID string) error
	ResetPassword(ctx context.Context, userID string, resetID string, passwordHash string) error
}
