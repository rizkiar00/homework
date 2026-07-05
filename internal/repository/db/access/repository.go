package access

import (
	"context"

	"github.com/rizkiar00/homework/internal/entity"
)

type Repository interface {
	FindAllActions(ctx context.Context) ([]entity.Action, error)
	CreateRole(ctx context.Context, data entity.Role, actionIDs []int64) (entity.Role, []entity.Action, error)
	UpdateRole(ctx context.Context, data entity.Role, updateDesc bool, updateActive bool, actionIDs []int64, updateActions bool) (entity.Role, []entity.Action, error)
	ReplaceRoleActions(ctx context.Context, roleID int64, actionIDs []int64) ([]entity.Action, error)
	AssignUserRole(ctx context.Context, userID string, roleID int64) error
}
