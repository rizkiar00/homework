package impl

import (
	"context"

	"github.com/rizkiar00/homework/internal/model"
	healthRepo "github.com/rizkiar00/homework/internal/repository/db/health"
	"github.com/rizkiar00/homework/pkg/config"
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
		Status:  "ok",
		Service: u.cfg.AppConfig.Name,
		Env:     u.cfg.AppConfig.Env,
	}, nil
}

func (u *usecase) Readiness(ctx context.Context) (model.ReadinessResponse, error) {
	checks := map[string]model.ReadinessCheck{
		"app": {
			Status:  "ok",
			Message: "application is ready",
		},
	}

	status := "ready"
	if err := u.repo.Check(ctx); err != nil {
		checks["database"] = model.ReadinessCheck{
			Status:  "skipped",
			Message: err.Error(),
		}
	} else {
		checks["database"] = model.ReadinessCheck{
			Status:  "ok",
			Message: "database configuration is present",
		}
	}

	return model.ReadinessResponse{
		Status: status,
		Checks: checks,
	}, nil
}
