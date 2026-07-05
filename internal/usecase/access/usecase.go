package access

import (
	"context"

	"github.com/rizkiar00/homework/internal/model"
)

type Usecase interface {
	FindAllActions(ctx context.Context) ([]model.ActionResponse, error)
	CreateRole(ctx context.Context, request model.CreateRoleRequest) (model.RoleResponse, error)
	UpdateRole(ctx context.Context, roleID int64, request model.UpdateRoleRequest) (model.RoleResponse, error)
	SetRoleActions(ctx context.Context, roleID int64, request model.SetRoleActionsRequest) (model.RoleResponse, error)
	AssignUserRole(ctx context.Context, userID string, request model.AssignUserRoleRequest) error
}
