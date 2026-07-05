package access

import (
	"github.com/rizkiar00/homework/internal/model"
	"github.com/rizkiar00/homework/internal/repository/db/access/impl"
	"go.uber.org/dig"
)

func Register(container *dig.Container) error {
	return container.Provide(func(db model.Database) Repository {
		return impl.New(db)
	})
}
