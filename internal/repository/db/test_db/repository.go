package test_db

import (
	"context"

	"github.com/rizkiar00/homework/internal/entity"
)

type Repository interface {
	Create(ctx context.Context, data entity.TestTable) (entity.TestTable, error)
	FindAll(ctx context.Context) ([]entity.TestTable, error)
	FindByID(ctx context.Context, id string) (entity.TestTable, error)
	Update(ctx context.Context, data entity.TestTable) (entity.TestTable, error)
	Delete(ctx context.Context, id string) error
}
