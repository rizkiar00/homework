package health

import (
	"context"

	"github.com/rizkiar00/homework/internal/model"
)

type Usecase interface {
	Health(ctx context.Context) (model.HealthResponse, error)
	Readiness(ctx context.Context) (model.ReadinessResponse, error)
}
