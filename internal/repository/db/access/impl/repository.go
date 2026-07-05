package impl

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rizkiar00/homework/internal/entity"
	"github.com/rizkiar00/homework/internal/model"
	"github.com/rizkiar00/homework/pkg/constant"
	"gorm.io/gorm"
)

type repository struct {
	db model.Database
}

func New(db model.Database) *repository {
	return &repository{db: db}
}

func (r *repository) FindAllActions(ctx context.Context) ([]entity.Action, error) {
	if r.db == nil {
		return nil, errors.New(constant.MessageDatabaseNotConfigured)
	}

	var rows []entity.Action
	if err := r.db.WithContext(ctx).Order("action_id asc").Find(&rows).Error; err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *repository) CreateRole(ctx context.Context, data entity.Role, actionIDs []int64) (entity.Role, []entity.Action, error) {
	if r.db == nil {
		return entity.Role{}, nil, errors.New(constant.MessageDatabaseNotConfigured)
	}

	var actions []entity.Action
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&data).Error; err != nil {
			return err
		}

		var err error
		actions, err = replaceRoleActions(ctx, tx, data.RoleID, actionIDs)
		return err
	})
	if err != nil {
		return entity.Role{}, nil, err
	}

	return data, actions, nil
}

func (r *repository) UpdateRole(ctx context.Context, data entity.Role, updateDesc bool, updateActive bool, actionIDs []int64, updateActions bool) (entity.Role, []entity.Action, error) {
	if r.db == nil {
		return entity.Role{}, nil, errors.New(constant.MessageDatabaseNotConfigured)
	}

	var row entity.Role
	var actions []entity.Action
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		updates := map[string]interface{}{"updated_at": time.Now()}
		if updateDesc {
			updates["role_desc"] = data.RoleDesc
		}
		if updateActive {
			updates["is_active"] = data.IsActive
		}

		result := tx.Model(&entity.Role{}).
			Where("role_id = ?", data.RoleID).
			Updates(updates)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		if err := tx.Where("role_id = ?", data.RoleID).First(&row).Error; err != nil {
			return err
		}

		if updateActions {
			var err error
			actions, err = replaceRoleActions(ctx, tx, data.RoleID, actionIDs)
			return err
		}

		var err error
		actions, err = findRoleActions(ctx, tx, data.RoleID)
		return err
	})
	if err != nil {
		return entity.Role{}, nil, err
	}

	return row, actions, nil
}

func (r *repository) ReplaceRoleActions(ctx context.Context, roleID int64, actionIDs []int64) ([]entity.Action, error) {
	if r.db == nil {
		return nil, errors.New(constant.MessageDatabaseNotConfigured)
	}

	var actions []entity.Action
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var role entity.Role
		if err := tx.Where("role_id = ?", roleID).First(&role).Error; err != nil {
			return err
		}

		var err error
		actions, err = replaceRoleActions(ctx, tx, roleID, actionIDs)
		return err
	})
	if err != nil {
		return nil, err
	}

	return actions, nil
}

func (r *repository) AssignUserRole(ctx context.Context, userID string, roleID int64) error {
	if r.db == nil {
		return errors.New(constant.MessageDatabaseNotConfigured)
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var role entity.Role
		if err := tx.Where("role_id = ? AND is_active = ?", roleID, true).First(&role).Error; err != nil {
			return err
		}

		result := tx.Model(&entity.User{}).
			Where("user_id = ? AND is_active = ?", userID, true).
			Updates(map[string]interface{}{
				"role_id":    roleID,
				"role":       role.RoleDesc,
				"updated_at": time.Now(),
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return nil
	})
}

func replaceRoleActions(ctx context.Context, tx *gorm.DB, roleID int64, actionIDs []int64) ([]entity.Action, error) {
	actions, err := findActionsByID(ctx, tx, actionIDs)
	if err != nil {
		return nil, err
	}

	if err := tx.WithContext(ctx).Where("role_id = ?", roleID).Delete(&entity.RoleAccess{}).Error; err != nil {
		return nil, err
	}

	now := time.Now()
	rows := make([]entity.RoleAccess, 0, len(actionIDs))
	for _, actionID := range actionIDs {
		rows = append(rows, entity.RoleAccess{
			RoleAccessID: uuid.NewString(),
			RoleID:       roleID,
			ActionID:     actionID,
			CreatedAt:    now,
		})
	}

	if len(rows) > 0 {
		if err := tx.WithContext(ctx).Create(&rows).Error; err != nil {
			return nil, err
		}
	}

	return actions, nil
}

func findActionsByID(ctx context.Context, tx *gorm.DB, actionIDs []int64) ([]entity.Action, error) {
	if len(actionIDs) == 0 {
		return nil, nil
	}

	var actions []entity.Action
	if err := tx.WithContext(ctx).Where("action_id IN ?", actionIDs).Order("action_id asc").Find(&actions).Error; err != nil {
		return nil, err
	}
	if len(actions) != len(uniqueInt64(actionIDs)) {
		return nil, gorm.ErrRecordNotFound
	}

	return actions, nil
}

func findRoleActions(ctx context.Context, tx *gorm.DB, roleID int64) ([]entity.Action, error) {
	var rows []entity.Action
	err := tx.WithContext(ctx).
		Table(constant.TableActions).
		Select("actions.action_id, actions.action_desc, actions.action_type, actions.endpoint").
		Joins("JOIN public.role_accesses ON role_accesses.action_id = actions.action_id").
		Where("role_accesses.role_id = ?", roleID).
		Order("actions.action_id asc").
		Scan(&rows).Error
	return rows, err
}

func uniqueInt64(values []int64) map[int64]struct{} {
	result := make(map[int64]struct{}, len(values))
	for _, value := range values {
		result[value] = struct{}{}
	}
	return result
}
