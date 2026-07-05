package impl

import (
	"context"
	"errors"

	"github.com/rizkiar00/homework/internal/model"
	"github.com/rizkiar00/homework/pkg/constant"
)

type repository struct {
	db model.Database
}

func New(db model.Database) *repository {
	return &repository{db: db}
}

func (r *repository) Check(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if r.db == nil {
		return errors.New(constant.MessageDatabaseNotConfigured)
	}

	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.PingContext(ctx)
}
