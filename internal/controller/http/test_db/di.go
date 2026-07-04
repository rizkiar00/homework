package test_db

import (
	testDBUsecase "github.com/rizkiar00/homework/internal/usecase/test_db"
	"github.com/rizkiar00/homework/pkg/token"
	"go.uber.org/dig"
)

func Register(container *dig.Container) error {
	return container.Provide(func(uc testDBUsecase.Usecase, tokenService *token.Service) *Controller {
		return NewController(uc, tokenService)
	})
}
