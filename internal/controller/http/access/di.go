package access

import (
	accessUsecase "github.com/rizkiar00/homework/internal/usecase/access"
	"go.uber.org/dig"
)

func Register(container *dig.Container) error {
	return container.Provide(func(uc accessUsecase.Usecase) *Controller {
		return NewController(uc)
	})
}
