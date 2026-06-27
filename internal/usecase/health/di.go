package health

import (
	healthRepo "github.com/rizkiar00/homework/internal/repository/db/health"
	"github.com/rizkiar00/homework/internal/usecase/health/impl"
	"github.com/rizkiar00/homework/pkg/config"
	"go.uber.org/dig"
)

func Register(container *dig.Container) error {
	return container.Provide(func(cfg config.Config, repo healthRepo.Repository) Usecase {
		return impl.New(cfg, repo)
	})
}
