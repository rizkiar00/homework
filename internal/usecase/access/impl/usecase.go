package impl

import (
	"context"
	"strings"
	"time"

	"github.com/rizkiar00/homework/internal/entity"
	"github.com/rizkiar00/homework/internal/model"
	accessRepo "github.com/rizkiar00/homework/internal/repository/db/access"
	"github.com/rizkiar00/homework/pkg/customerror"
)

type usecase struct {
	repo accessRepo.Repository
}

func New(repo accessRepo.Repository) *usecase {
	return &usecase{repo: repo}
}

func (u *usecase) FindAllActions(ctx context.Context) ([]model.ActionResponse, error) {
	rows, err := u.repo.FindAllActions(ctx)
	if err != nil {
		return nil, err
	}

	return toActionResponses(rows), nil
}

func (u *usecase) CreateRole(ctx context.Context, request model.CreateRoleRequest) (model.RoleResponse, error) {
	request.RoleDesc = strings.TrimSpace(request.RoleDesc)
	if request.RoleDesc == "" {
		return model.RoleResponse{}, customerror.BadRequest("role_desc is required")
	}

	row, actions, err := u.repo.CreateRole(ctx, entity.Role{
		RoleDesc:  request.RoleDesc,
		IsActive:  true,
		CreatedAt: time.Now(),
	}, uniqueActionIDs(request.ActionIDs))
	if err != nil {
		return model.RoleResponse{}, err
	}

	return toRoleResponse(row, actions), nil
}

func (u *usecase) UpdateRole(ctx context.Context, roleID int64, request model.UpdateRoleRequest) (model.RoleResponse, error) {
	if roleID <= 0 {
		return model.RoleResponse{}, customerror.BadRequest("role_id must be greater than 0")
	}

	updateDesc := request.RoleDesc != nil
	roleDesc := ""
	if request.RoleDesc != nil {
		roleDesc = strings.TrimSpace(*request.RoleDesc)
		if roleDesc == "" {
			return model.RoleResponse{}, customerror.BadRequest("role_desc is required")
		}
	}

	row, actions, err := u.repo.UpdateRole(ctx, entity.Role{
		RoleID:   roleID,
		RoleDesc: roleDesc,
		IsActive: request.IsActive != nil && *request.IsActive,
	}, updateDesc, request.IsActive != nil, uniqueActionIDs(request.ActionIDs), request.ActionIDs != nil)
	if err != nil {
		return model.RoleResponse{}, err
	}

	return toRoleResponse(row, actions), nil
}

func (u *usecase) SetRoleActions(ctx context.Context, roleID int64, request model.SetRoleActionsRequest) (model.RoleResponse, error) {
	if roleID <= 0 {
		return model.RoleResponse{}, customerror.BadRequest("role_id must be greater than 0")
	}

	actions, err := u.repo.ReplaceRoleActions(ctx, roleID, uniqueActionIDs(request.ActionIDs))
	if err != nil {
		return model.RoleResponse{}, err
	}

	return model.RoleResponse{
		RoleID:    roleID,
		ActionIDs: actionIDs(actions),
		Actions:   toActionResponses(actions),
	}, nil
}

func (u *usecase) AssignUserRole(ctx context.Context, userID string, request model.AssignUserRoleRequest) error {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return customerror.BadRequest("user_id is required")
	}
	if request.RoleID <= 0 {
		return customerror.BadRequest("role_id must be greater than 0")
	}

	return u.repo.AssignUserRole(ctx, userID, request.RoleID)
}

func toActionResponses(rows []entity.Action) []model.ActionResponse {
	responses := make([]model.ActionResponse, 0, len(rows))
	for _, row := range rows {
		responses = append(responses, model.ActionResponse{
			ActionID:   row.ActionID,
			ActionDesc: row.ActionDesc,
			ActionType: row.ActionType,
			Endpoint:   row.Endpoint,
		})
	}
	return responses
}

func toRoleResponse(row entity.Role, actions []entity.Action) model.RoleResponse {
	return model.RoleResponse{
		RoleID:    row.RoleID,
		RoleDesc:  row.RoleDesc,
		IsActive:  row.IsActive,
		ActionIDs: actionIDs(actions),
		Actions:   toActionResponses(actions),
	}
}

func actionIDs(actions []entity.Action) []int64 {
	ids := make([]int64, 0, len(actions))
	for _, action := range actions {
		ids = append(ids, action.ActionID)
	}
	return ids
}

func uniqueActionIDs(actionIDs []int64) []int64 {
	seen := make(map[int64]struct{}, len(actionIDs))
	ids := make([]int64, 0, len(actionIDs))
	for _, actionID := range actionIDs {
		if actionID <= 0 {
			continue
		}
		if _, ok := seen[actionID]; ok {
			continue
		}
		seen[actionID] = struct{}{}
		ids = append(ids, actionID)
	}
	return ids
}
