package test_db

import (
	"github.com/rizkiar00/homework/internal/model"
	"github.com/rizkiar00/homework/internal/repository/db/test_db/impl"
	"go.uber.org/dig"
)

func Register(container *dig.Container) error {
	return container.Provide(func(db model.Database) Repository {
		return impl.New(db)
	})
}
