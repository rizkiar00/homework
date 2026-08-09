package test_db

import (
	"context"

	"github.com/rizkiar00/homework/internal/entity"
	"github.com/rizkiar00/homework/internal/model"
)

type Repository interface {
	Create(ctx context.Context, data entity.TestTable) (entity.TestTable, error)
	FindAll(ctx context.Context, option model.TestDBFindAllOption) ([]entity.TestTable, int64, error)
	FindByID(ctx context.Context, id string, userID string, role string) (entity.TestTable, error)
	Update(ctx context.Context, data entity.TestTable, userID string, role string) (entity.TestTable, error)
	Delete(ctx context.Context, id string, userID string, role string) error
}
