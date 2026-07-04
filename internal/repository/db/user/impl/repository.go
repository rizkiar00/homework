package impl

import (
	"context"
	"errors"

	"github.com/rizkiar00/homework/internal/entity"
	"github.com/rizkiar00/homework/internal/model"
)

type repository struct {
	db model.Database
}

func New(db model.Database) *repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, data entity.User) (entity.User, error) {
	if r.db == nil {
		return entity.User{}, errors.New("database is not configured")
	}
	if err := r.db.WithContext(ctx).Create(&data).Error; err != nil {
		return entity.User{}, err
	}

	return data, nil
}

func (r *repository) FindByUsername(ctx context.Context, username string) (entity.User, error) {
	if r.db == nil {
		return entity.User{}, errors.New("database is not configured")
	}

	var row entity.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&row).Error; err != nil {
		return entity.User{}, err
	}

	return row, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (entity.User, error) {
	if r.db == nil {
		return entity.User{}, errors.New("database is not configured")
	}

	var row entity.User
	if err := r.db.WithContext(ctx).Where("id_user = ?", id).First(&row).Error; err != nil {
		return entity.User{}, err
	}

	return row, nil
}
