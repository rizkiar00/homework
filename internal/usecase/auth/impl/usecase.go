package impl

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rizkiar00/homework/internal/entity"
	"github.com/rizkiar00/homework/internal/model"
	userRepo "github.com/rizkiar00/homework/internal/repository/db/user"
	"github.com/rizkiar00/homework/pkg/constant"
	"github.com/rizkiar00/homework/pkg/token"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type usecase struct {
	repo         userRepo.Repository
	tokenService *token.Service
}

func New(repo userRepo.Repository, tokenService *token.Service) *usecase {
	return &usecase{
		repo:         repo,
		tokenService: tokenService,
	}
}

func (u *usecase) Register(ctx context.Context, request model.RegisterRequest) (model.AuthUserResponse, error) {
	request.Username = strings.TrimSpace(request.Username)
	if request.Username == "" || len(request.Password) < 8 {
		return model.AuthUserResponse{}, errors.New("username is required and password minimum length is 8")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.AuthUserResponse{}, err
	}

	row, err := u.repo.Create(ctx, entity.User{
		IDUser:       uuid.NewString(),
		Username:     request.Username,
		PasswordHash: string(hash),
		Role:         constant.RoleUser,
		CreatedAt:    time.Now(),
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "users_username_key") {
			return model.AuthUserResponse{}, model.ErrUsernameAlreadyExists
		}
		return model.AuthUserResponse{}, err
	}

	return toAuthUserResponse(row), nil
}

func (u *usecase) Login(ctx context.Context, request model.LoginRequest) (model.LoginResponse, error) {
	row, err := u.repo.FindByUsername(ctx, strings.TrimSpace(request.Username))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.LoginResponse{}, model.ErrInvalidCredential
		}
		return model.LoginResponse{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(row.PasswordHash), []byte(request.Password)); err != nil {
		return model.LoginResponse{}, model.ErrInvalidCredential
	}

	accessToken, expiresIn, err := u.tokenService.Generate(row.IDUser, row.Username, row.Role)
	if err != nil {
		return model.LoginResponse{}, err
	}

	return model.LoginResponse{
		AccessToken: accessToken,
		TokenType:   constant.TokenTypeBearer,
		ExpiresIn:   expiresIn,
	}, nil
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
		IDUser:   row.IDUser,
		Username: row.Username,
		Role:     row.Role,
	}
}
