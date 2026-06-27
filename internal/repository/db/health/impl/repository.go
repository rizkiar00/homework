package impl

import (
	"context"
	"errors"

	"github.com/rizkiar00/homework/internal/model"
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
		return errors.New("database is not configured")
	}

	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.PingContext(ctx)
}
