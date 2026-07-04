package test_db

import (
	"context"

	"github.com/rizkiar00/homework/internal/model"
)

type Usecase interface {
	Create(ctx context.Context, request model.CreateTestDBRequest) (model.TestDBResponse, error)
	FindAll(ctx context.Context, request model.TestDBListRequest) (model.TestDBListResponse, error)
	FindByID(ctx context.Context, id string) (model.TestDBResponse, error)
	Update(ctx context.Context, id string, request model.UpdateTestDBRequest) (model.TestDBResponse, error)
	Delete(ctx context.Context, id string) error
}
