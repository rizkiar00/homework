package email

import (
	"github.com/rizkiar00/homework/internal/repository/email/resend"
	"go.uber.org/dig"
)

func Register(container *dig.Container) error {
	return container.Provide(func(repo *resend.Repository) Repository {
		return repo
	})
}
