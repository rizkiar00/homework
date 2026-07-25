package impl

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rizkiar00/homework/internal/entity"
	"github.com/rizkiar00/homework/internal/model"
	userRepo "github.com/rizkiar00/homework/internal/repository/db/user"
	emailRepo "github.com/rizkiar00/homework/internal/repository/email"
	"github.com/rizkiar00/homework/pkg/config"
	"github.com/rizkiar00/homework/pkg/constant"
	"github.com/rizkiar00/homework/pkg/customerror"
	"github.com/rizkiar00/homework/pkg/token"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	verificationCodeTTL     = 15 * time.Minute
	passwordResetCodeTTL    = 15 * time.Minute
	maxVerificationAttempts = 5
	emailCodeLimitWindow    = 15 * time.Minute
	emailCodeCooldown       = 60 * time.Second
	maxEmailCodeRequests    = 3
)

type usecase struct {
	cfg          config.Config
	repo         userRepo.Repository
	emailRepo    emailRepo.Repository
	tokenService *token.Service
	blacklist    *token.Blacklist
}

func New(cfg config.Config, repo userRepo.Repository, emailRepo emailRepo.Repository, tokenService *token.Service, blacklist *token.Blacklist) *usecase {
	return &usecase{
		cfg:          cfg,
		repo:         repo,
		emailRepo:    emailRepo,
		tokenService: tokenService,
		blacklist:    blacklist,
	}
}

func (u *usecase) Register(ctx context.Context, request model.RegisterRequest) (model.AuthUserResponse, error) {
	request.FullName = strings.TrimSpace(request.FullName)
	request.Username = strings.TrimSpace(request.Username)
	request.Email = strings.ToLower(strings.TrimSpace(request.Email))
	if request.FullName == "" || request.Username == "" || request.Email == "" || len(request.Password) < 8 {
		return model.AuthUserResponse{}, customerror.BadRequest(constant.MessageInvalidRegisterCredential)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.AuthUserResponse{}, err
	}

	code, err := generateVerificationCode()
	if err != nil {
		return model.AuthUserResponse{}, err
	}

	now := time.Now()
	user := entity.User{
		UserID:        uuid.NewString(),
		FullName:      request.FullName,
		Username:      request.Username,
		Email:         request.Email,
		PasswordHash:  string(hash),
		Role:          constant.RoleUser,
		RoleID:        int64Pointer(2),
		IsActive:      false,
		EmailVerified: false,
		CreatedAt:     now,
	}
	verification := entity.UserEmailVerification{
		VerificationID: uuid.NewString(),
		UserID:         user.UserID,
		Email:          request.Email,
		CodeHash:       u.hashVerificationCode(request.Email, code),
		ExpiresAt:      now.Add(verificationCodeTTL),
		CreatedAt:      now,
	}

	existingByEmail, err := u.repo.FindByEmail(ctx, request.Email)
	if err == nil {
		if existingByEmail.EmailVerified {
			return model.AuthUserResponse{}, customerror.Conflict(constant.MessageEmailAlreadyExists)
		}

		existingByUsername, usernameErr := u.repo.FindByUsername(ctx, request.Username)
		if usernameErr == nil && existingByUsername.UserID != existingByEmail.UserID {
			return model.AuthUserResponse{}, customerror.Conflict(constant.MessageUsernameAlreadyExists)
		}
		if usernameErr != nil && !errors.Is(usernameErr, gorm.ErrRecordNotFound) {
			return model.AuthUserResponse{}, usernameErr
		}
		if err := u.enforceVerificationEmailLimit(ctx, request.Email); err != nil {
			return model.AuthUserResponse{}, err
		}

		user.UserID = existingByEmail.UserID
		verification.UserID = existingByEmail.UserID
		row, updateErr := u.repo.UpdatePendingRegistration(ctx, existingByEmail.UserID, user, verification)
		if updateErr != nil {
			return model.AuthUserResponse{}, updateErr
		}
		if sendErr := u.emailRepo.SendVerificationCode(ctx, request.Email, request.FullName, code); sendErr != nil {
			return model.AuthUserResponse{}, sendErr
		}

		return toAuthUserResponse(row), nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return model.AuthUserResponse{}, err
	}

	existingByUsername, err := u.repo.FindByUsername(ctx, request.Username)
	if err == nil {
		if existingByUsername.EmailVerified || existingByUsername.Email != request.Email {
			return model.AuthUserResponse{}, customerror.Conflict(constant.MessageUsernameAlreadyExists)
		}
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return model.AuthUserResponse{}, err
	}

	row, err := u.repo.CreatePendingRegistration(ctx, user, verification)
	if err != nil {
		if strings.Contains(err.Error(), "users_username") {
			return model.AuthUserResponse{}, customerror.Conflict(constant.MessageUsernameAlreadyExists)
		}
		if strings.Contains(err.Error(), "users_email") {
			return model.AuthUserResponse{}, customerror.Conflict(constant.MessageEmailAlreadyExists)
		}
		return model.AuthUserResponse{}, err
	}
	if err := u.emailRepo.SendVerificationCode(ctx, request.Email, request.FullName, code); err != nil {
		return model.AuthUserResponse{}, err
	}

	return toAuthUserResponse(row), nil
}

func int64Pointer(value int64) *int64 {
	return &value
}

func (u *usecase) VerifyEmail(ctx context.Context, request model.VerifyEmailRequest) error {
	request.Email = strings.ToLower(strings.TrimSpace(request.Email))
	request.Code = strings.TrimSpace(request.Code)
	if request.Email == "" || request.Code == "" {
		return customerror.BadRequest(constant.MessageInvalidVerificationCode)
	}

	verification, err := u.repo.FindActiveVerificationByEmail(ctx, request.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return customerror.BadRequest(constant.MessageInvalidVerificationCode)
		}
		return err
	}
	if verification.AttemptCount >= maxVerificationAttempts {
		return customerror.BadRequest(constant.MessageInvalidVerificationCode)
	}

	expected := u.hashVerificationCode(request.Email, request.Code)
	if !hmac.Equal([]byte(expected), []byte(verification.CodeHash)) {
		_ = u.repo.IncrementVerificationAttempt(ctx, verification.VerificationID)
		return customerror.BadRequest(constant.MessageInvalidVerificationCode)
	}

	return u.repo.VerifyEmail(ctx, verification.UserID, verification.VerificationID)
}

func (u *usecase) ResendVerification(ctx context.Context, request model.ResendVerificationRequest) error {
	request.Email = strings.ToLower(strings.TrimSpace(request.Email))
	if request.Email == "" {
		return customerror.BadRequest(constant.MessageInvalidVerificationCode)
	}

	row, err := u.repo.FindByEmail(ctx, request.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	if row.EmailVerified {
		return nil
	}
	if err := u.enforceVerificationEmailLimit(ctx, request.Email); err != nil {
		return err
	}

	code, err := generateVerificationCode()
	if err != nil {
		return err
	}
	now := time.Now()
	verification := entity.UserEmailVerification{
		VerificationID: uuid.NewString(),
		UserID:         row.UserID,
		Email:          row.Email,
		CodeHash:       u.hashVerificationCode(row.Email, code),
		ExpiresAt:      now.Add(verificationCodeTTL),
		CreatedAt:      now,
	}
	if _, err := u.repo.UpdatePendingRegistration(ctx, row.UserID, row, verification); err != nil {
		return err
	}

	return u.emailRepo.SendVerificationCode(ctx, row.Email, row.FullName, code)
}

func (u *usecase) ForgotPassword(ctx context.Context, request model.ForgotPasswordRequest) error {
	request.Email = strings.ToLower(strings.TrimSpace(request.Email))
	if request.Email == "" {
		return nil
	}

	row, err := u.repo.FindByEmail(ctx, request.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	if !row.IsActive || !row.EmailVerified {
		return nil
	}
	if err := u.enforcePasswordResetEmailLimit(ctx, request.Email); err != nil {
		return err
	}

	code, err := generateVerificationCode()
	if err != nil {
		return err
	}

	now := time.Now()
	reset := entity.UserPasswordReset{
		ResetID:   uuid.NewString(),
		UserID:    row.UserID,
		Email:     row.Email,
		CodeHash:  u.hashPasswordResetCode(row.Email, code),
		ExpiresAt: now.Add(passwordResetCodeTTL),
		CreatedAt: now,
	}
	if err := u.repo.CreatePasswordReset(ctx, reset); err != nil {
		return err
	}

	return u.emailRepo.SendPasswordResetCode(ctx, row.Email, row.FullName, code)
}

func (u *usecase) ResetPassword(ctx context.Context, request model.ResetPasswordRequest) error {
	request.Email = strings.ToLower(strings.TrimSpace(request.Email))
	request.Code = strings.TrimSpace(request.Code)
	if request.Email == "" || request.Code == "" || len(request.NewPassword) < 8 {
		return customerror.BadRequest(constant.MessageInvalidResetPassword)
	}

	reset, err := u.repo.FindActivePasswordResetByEmail(ctx, request.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return customerror.BadRequest(constant.MessageInvalidResetPasswordCode)
		}
		return err
	}
	if reset.AttemptCount >= maxVerificationAttempts {
		return customerror.BadRequest(constant.MessageInvalidResetPasswordCode)
	}

	expected := u.hashPasswordResetCode(request.Email, request.Code)
	if !hmac.Equal([]byte(expected), []byte(reset.CodeHash)) {
		_ = u.repo.IncrementPasswordResetAttempt(ctx, reset.ResetID)
		return customerror.BadRequest(constant.MessageInvalidResetPasswordCode)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return u.repo.ResetPassword(ctx, reset.UserID, reset.ResetID, string(hash))
}

func (u *usecase) Login(ctx context.Context, request model.LoginRequest) (model.LoginResponse, error) {
	row, err := u.repo.FindByUsername(ctx, strings.TrimSpace(request.Username))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.LoginResponse{}, customerror.Unauthorized(constant.MessageInvalidAuthCredential)
		}
		return model.LoginResponse{}, err
	}
	if !row.IsActive || !row.EmailVerified {
		return model.LoginResponse{}, customerror.Unauthorized(constant.MessageEmailNotVerified)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(row.PasswordHash), []byte(request.Password)); err != nil {
		return model.LoginResponse{}, customerror.Unauthorized(constant.MessageInvalidAuthCredential)
	}

	accessToken, expiresIn, err := u.tokenService.Generate(row.UserID, row.Username, row.Role)
	if err != nil {
		return model.LoginResponse{}, err
	}

	return model.LoginResponse{
		AccessToken: accessToken,
		TokenType:   constant.TokenTypeBearer,
		ExpiresIn:   expiresIn,
	}, nil
}

func (u *usecase) Logout(ctx context.Context, claims token.Claims) error {
	return u.blacklist.Revoke(ctx, claims)
}

func (u *usecase) Me(ctx context.Context, userID string) (model.AuthUserResponse, error) {
	row, err := u.repo.FindByID(ctx, userID)
	if err != nil {
		return model.AuthUserResponse{}, err
	}

	return toAuthUserResponse(row), nil
}

func toAuthUserResponse(row entity.User) model.AuthUserResponse {
	return model.AuthUserResponse{
		UserID:   row.UserID,
		FullName: row.FullName,
		Username: row.Username,
		Email:    row.Email,
		Role:     row.Role,
	}
}

func (u *usecase) hashVerificationCode(email string, code string) string {
	secret := u.cfg.JWT.Secret
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(strings.ToLower(strings.TrimSpace(email))))
	mac.Write([]byte(":"))
	mac.Write([]byte(code))
	return hex.EncodeToString(mac.Sum(nil))
}

func (u *usecase) hashPasswordResetCode(email string, code string) string {
	secret := u.cfg.JWT.Secret
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte("password_reset:"))
	mac.Write([]byte(strings.ToLower(strings.TrimSpace(email))))
	mac.Write([]byte(":"))
	mac.Write([]byte(code))
	return hex.EncodeToString(mac.Sum(nil))
}

func (u *usecase) enforceVerificationEmailLimit(ctx context.Context, email string) error {
	since := time.Now().Add(-emailCodeLimitWindow)
	count, err := u.repo.CountVerificationsByEmailSince(ctx, email, since)
	if err != nil {
		return err
	}
	if count >= maxEmailCodeRequests {
		return customerror.TooManyRequests(constant.MessageEmailRequestLimitExceeded)
	}

	latest, err := u.repo.FindLatestVerificationByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	if time.Since(latest.CreatedAt) < emailCodeCooldown {
		return customerror.TooManyRequests(constant.MessageEmailRequestTooFrequent)
	}

	return nil
}

func (u *usecase) enforcePasswordResetEmailLimit(ctx context.Context, email string) error {
	since := time.Now().Add(-emailCodeLimitWindow)
	count, err := u.repo.CountPasswordResetsByEmailSince(ctx, email, since)
	if err != nil {
		return err
	}
	if count >= maxEmailCodeRequests {
		return customerror.TooManyRequests(constant.MessageEmailRequestLimitExceeded)
	}

	latest, err := u.repo.FindLatestPasswordResetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	if time.Since(latest.CreatedAt) < emailCodeCooldown {
		return customerror.TooManyRequests(constant.MessageEmailRequestTooFrequent)
	}

	return nil
}

func generateVerificationCode() (string, error) {
	var b [3]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}

	value := (int(b[0])<<16 | int(b[1])<<8 | int(b[2])) % 1000000
	return fmt.Sprintf("%06d", value), nil
}
