package access

import (
	accessRepo "github.com/rizkiar00/homework/internal/repository/db/access"
	"github.com/rizkiar00/homework/internal/usecase/access/impl"
	"go.uber.org/dig"
)

func Register(container *dig.Container) error {
	return container.Provide(func(repo accessRepo.Repository) Usecase {
		return impl.New(repo)
	})
}
