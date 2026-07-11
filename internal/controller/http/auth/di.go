package auth

import (
	authUsecase "github.com/rizkiar00/homework/internal/usecase/auth"
	"github.com/rizkiar00/homework/pkg/token"
	"go.uber.org/dig"
)

func Register(container *dig.Container) error {
	return container.Provide(func(uc authUsecase.Usecase, tokenService *token.Service, blacklist *token.Blacklist) *Controller {
		return NewController(uc, tokenService, blacklist)
	})
}
