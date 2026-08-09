package test_db

import (
	"context"

	"github.com/rizkiar00/homework/internal/model"
	"github.com/rizkiar00/homework/pkg/token"
)

type Usecase interface {
	Create(ctx context.Context, createdBy string, request model.CreateTestDBRequest) (model.TestDBResponse, error)
	FindAll(ctx context.Context, claims token.Claims, request model.TestDBListRequest) (model.TestDBListResponse, error)
	FindByID(ctx context.Context, claims token.Claims, id string) (model.TestDBResponse, error)
	Update(ctx context.Context, claims token.Claims, id string, request model.UpdateTestDBRequest) (model.TestDBResponse, error)
	Delete(ctx context.Context, claims token.Claims, id string) error
}
