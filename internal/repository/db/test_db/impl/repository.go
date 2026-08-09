package impl

import (
	"context"
	"errors"
	"fmt"
	"time"

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

func (r *repository) Create(ctx context.Context, data entity.TestTable) (entity.TestTable, error) {
	if r.db == nil {
		return entity.TestTable{}, errors.New(constant.MessageDatabaseNotConfigured)
	}

	if err := r.db.WithContext(ctx).Create(&data).Error; err != nil {
		return entity.TestTable{}, err
	}

	return data, nil
}

func (r *repository) FindAll(ctx context.Context, option model.TestDBFindAllOption) ([]entity.TestTable, int64, error) {
	if r.db == nil {
		return nil, 0, errors.New(constant.MessageDatabaseNotConfigured)
	}

	query := r.db.WithContext(ctx).Model(&entity.TestTable{}).Where("is_active = ?", true)
	query = applyOwnerScope(query, option.UserID, option.Role)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []entity.TestTable
	order := fmt.Sprintf("%s %s", option.OrderBy, option.OrderDir)
	if err := query.Order(order).Limit(option.Limit).Offset(option.Offset).Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

func (r *repository) FindByID(ctx context.Context, id string, userID string, role string) (entity.TestTable, error) {
	if r.db == nil {
		return entity.TestTable{}, errors.New(constant.MessageDatabaseNotConfigured)
	}

	var row entity.TestTable
	query := r.db.WithContext(ctx).Where(constant.ColumnTestID+" = ? AND is_active = ?", id, true)
	query = applyOwnerScope(query, userID, role)
	if err := query.First(&row).Error; err != nil {
		return entity.TestTable{}, err
	}

	return row, nil
}

func (r *repository) Update(ctx context.Context, data entity.TestTable, userID string, role string) (entity.TestTable, error) {
	if r.db == nil {
		return entity.TestTable{}, errors.New(constant.MessageDatabaseNotConfigured)
	}

	now := time.Now()
	query := r.db.WithContext(ctx).
		Model(&entity.TestTable{}).
		Where(constant.ColumnTestID+" = ? AND is_active = ?", data.TestID, true)
	query = applyOwnerScope(query, userID, role)
	result := query.Updates(map[string]interface{}{
		constant.ColumnDescTest: data.DescTest,
		"updated_by":            userID,
		"updated_at":            now,
	})
	if result.Error != nil {
		return entity.TestTable{}, result.Error
	}
	if result.RowsAffected == 0 {
		return entity.TestTable{}, gorm.ErrRecordNotFound
	}

	return r.FindByID(ctx, data.TestID, userID, role)
}

func (r *repository) Delete(ctx context.Context, id string, userID string, role string) error {
	if r.db == nil {
		return errors.New(constant.MessageDatabaseNotConfigured)
	}

	now := time.Now()
	query := r.db.WithContext(ctx).
		Model(&entity.TestTable{}).
		Where(constant.ColumnTestID+" = ? AND is_active = ?", id, true)
	query = applyOwnerScope(query, userID, role)
	result := query.Updates(map[string]interface{}{
		"is_active":  false,
		"updated_by": userID,
		"updated_at": now,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func applyOwnerScope(query *gorm.DB, userID string, role string) *gorm.DB {
	if role == constant.RoleAdmin {
		return query
	}

	return query.Where("created_by = ?", userID)
}
