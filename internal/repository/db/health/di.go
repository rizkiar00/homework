package health

import (
	"github.com/rizkiar00/homework/internal/model"
	"github.com/rizkiar00/homework/internal/repository/db/health/impl"
	"go.uber.org/dig"
)

func Register(container *dig.Container) error {
	return container.Provide(func(db model.Database, redis model.Redis) Repository {
		return impl.New(db, redis)
	})
}
