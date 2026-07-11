package auth

import (
	userRepo "github.com/rizkiar00/homework/internal/repository/db/user"
	"github.com/rizkiar00/homework/internal/usecase/auth/impl"
	"github.com/rizkiar00/homework/pkg/token"
	"go.uber.org/dig"
)

func Register(container *dig.Container) error {
	return container.Provide(func(repo userRepo.Repository, tokenService *token.Service, blacklist *token.Blacklist) Usecase {
		return impl.New(repo, tokenService, blacklist)
	})
}
