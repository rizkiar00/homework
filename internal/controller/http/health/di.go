package health

import (
	healthUsecase "github.com/rizkiar00/homework/internal/usecase/health"
	"go.uber.org/dig"
)

func Register(container *dig.Container) error {
	return container.Provide(func(uc healthUsecase.Usecase) *Controller {
		return NewController(uc)
	})
}
