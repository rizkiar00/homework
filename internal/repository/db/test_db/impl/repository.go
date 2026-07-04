package impl

import (
	"context"
	"errors"
	"fmt"

	"github.com/rizkiar00/homework/internal/entity"
	"github.com/rizkiar00/homework/internal/model"
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
		return entity.TestTable{}, errors.New("database is not configured")
	}

	if err := r.db.WithContext(ctx).Create(&data).Error; err != nil {
		return entity.TestTable{}, err
	}

	return data, nil
}

func (r *repository) FindAll(ctx context.Context, option model.TestDBFindAllOption) ([]entity.TestTable, int64, error) {
	if r.db == nil {
		return nil, 0, errors.New("database is not configured")
	}

	var total int64
	if err := r.db.WithContext(ctx).Model(&entity.TestTable{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []entity.TestTable
	order := fmt.Sprintf("%s %s", option.OrderBy, option.OrderDir)
	if err := r.db.WithContext(ctx).Order(order).Limit(option.Limit).Offset(option.Offset).Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (entity.TestTable, error) {
	if r.db == nil {
		return entity.TestTable{}, errors.New("database is not configured")
	}

	var row entity.TestTable
	if err := r.db.WithContext(ctx).Where("id_test = ?", id).First(&row).Error; err != nil {
		return entity.TestTable{}, err
	}

	return row, nil
}

func (r *repository) Update(ctx context.Context, data entity.TestTable) (entity.TestTable, error) {
	if r.db == nil {
		return entity.TestTable{}, errors.New("database is not configured")
	}

	result := r.db.WithContext(ctx).
		Model(&entity.TestTable{}).
		Where("id_test = ?", data.IDTest).
		Updates(map[string]interface{}{
			"desc_test": data.DescTest,
		})
	if result.Error != nil {
		return entity.TestTable{}, result.Error
	}
	if result.RowsAffected == 0 {
		return entity.TestTable{}, gorm.ErrRecordNotFound
	}

	return r.FindByID(ctx, data.IDTest)
}

func (r *repository) Delete(ctx context.Context, id string) error {
	if r.db == nil {
		return errors.New("database is not configured")
	}

	result := r.db.WithContext(ctx).Where("id_test = ?", id).Delete(&entity.TestTable{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
