package impl

import (
	"context"
	"errors"

	"github.com/rizkiar00/homework/internal/entity"
	"github.com/rizkiar00/homework/internal/model"
	"github.com/rizkiar00/homework/pkg/constant"
)

type repository struct {
	db model.Database
}

func New(db model.Database) *repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, data entity.User) (entity.User, error) {
	if r.db == nil {
		return entity.User{}, errors.New(constant.MessageDatabaseNotConfigured)
	}
	if err := r.db.WithContext(ctx).Create(&data).Error; err != nil {
		return entity.User{}, err
	}

	return data, nil
}

func (r *repository) FindByUsername(ctx context.Context, username string) (entity.User, error) {
	if r.db == nil {
		return entity.User{}, errors.New(constant.MessageDatabaseNotConfigured)
	}

	var row entity.User
	if err := r.db.WithContext(ctx).Where(constant.ColumnUsername+" = ? AND is_active = ?", username, true).First(&row).Error; err != nil {
		return entity.User{}, err
	}

	return row, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (entity.User, error) {
	if r.db == nil {
		return entity.User{}, errors.New(constant.MessageDatabaseNotConfigured)
	}

	var row entity.User
	if err := r.db.WithContext(ctx).Where(constant.ColumnUserID+" = ? AND is_active = ?", id, true).First(&row).Error; err != nil {
		return entity.User{}, err
	}

	return row, nil
}
