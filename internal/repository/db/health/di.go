package health

import (
	"github.com/rizkiar00/homework/internal/repository/db/health/impl"
	"github.com/rizkiar00/homework/pkg/config"
	"go.uber.org/dig"
)

func Register(container *dig.Container) error {
	return container.Provide(func(cfg config.Config) Repository {
		return impl.New(cfg)
	})
}
