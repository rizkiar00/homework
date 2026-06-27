package impl

import (
	"context"
	"errors"

	"github.com/rizkiar00/homework/pkg/config"
)

type repository struct {
	cfg config.Config
}

func New(cfg config.Config) *repository {
	return &repository{cfg: cfg}
}

func (r *repository) Check(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if !r.cfg.Database.IsConfigured() {
		return errors.New("database is not configured")
	}

	return nil
}
