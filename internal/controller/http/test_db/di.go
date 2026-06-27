package test_db

import (
	testDBUsecase "github.com/rizkiar00/homework/internal/usecase/test_db"
	"go.uber.org/dig"
)

func Register(container *dig.Container) error {
	return container.Provide(func(uc testDBUsecase.Usecase) *Controller {
		return NewController(uc)
	})
}
