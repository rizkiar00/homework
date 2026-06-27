package test_db

import (
	testDBRepo "github.com/rizkiar00/homework/internal/repository/db/test_db"
	"github.com/rizkiar00/homework/internal/usecase/test_db/impl"
	"go.uber.org/dig"
)

func Register(container *dig.Container) error {
	return container.Provide(func(repo testDBRepo.Repository) Usecase {
		return impl.New(repo)
	})
}
