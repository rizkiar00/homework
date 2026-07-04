package user

import (
	"context"

	"github.com/rizkiar00/homework/internal/entity"
)

type Repository interface {
	Create(ctx context.Context, data entity.User) (entity.User, error)
	FindByUsername(ctx context.Context, username string) (entity.User, error)
	FindByID(ctx context.Context, id string) (entity.User, error)
}
