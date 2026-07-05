package impl

import (
	"context"

	"github.com/rizkiar00/homework/internal/model"
	healthRepo "github.com/rizkiar00/homework/internal/repository/db/health"
	"github.com/rizkiar00/homework/pkg/config"
	"github.com/rizkiar00/homework/pkg/constant"
)

type usecase struct {
	cfg  config.Config
	repo healthRepo.Repository
}

func New(cfg config.Config, repo healthRepo.Repository) *usecase {
	return &usecase{
		cfg:  cfg,
		repo: repo,
	}
}

func (u *usecase) Health(ctx context.Context) (model.HealthResponse, error) {
	if err := ctx.Err(); err != nil {
		return model.HealthResponse{}, err
	}

	return model.HealthResponse{
		Status:  constant.StatusOK,
		Service: u.cfg.AppConfig.Name,
		Env:     u.cfg.AppConfig.Env,
	}, nil
}

func (u *usecase) Readiness(ctx context.Context) (model.ReadinessResponse, error) {
	checks := map[string]model.ReadinessCheck{
		"app": {
			Status:  constant.StatusOK,
			Message: constant.MessageApplicationReady,
		},
	}

	status := constant.StatusReady
	if err := u.repo.Check(ctx); err != nil {
		checks["database"] = model.ReadinessCheck{
			Status:  "skipped",
			Message: err.Error(),
		}
	} else {
		checks["database"] = model.ReadinessCheck{
			Status:  constant.StatusOK,
			Message: constant.MessageDatabaseConfigIsAvailable,
		}
	}

	return model.ReadinessResponse{
		Status: status,
		Checks: checks,
	}, nil
}
