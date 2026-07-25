package auth

import (
	userRepo "github.com/rizkiar00/homework/internal/repository/db/user"
	emailRepo "github.com/rizkiar00/homework/internal/repository/email"
	"github.com/rizkiar00/homework/internal/usecase/auth/impl"
	"github.com/rizkiar00/homework/pkg/config"
	"github.com/rizkiar00/homework/pkg/token"
	"go.uber.org/dig"
)

func Register(container *dig.Container) error {
	return container.Provide(func(cfg config.Config, repo userRepo.Repository, emailRepo emailRepo.Repository, tokenService *token.Service, blacklist *token.Blacklist) Usecase {
		return impl.New(cfg, repo, emailRepo, tokenService, blacklist)
	})
}
