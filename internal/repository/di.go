package repository

import (
	dbrepo "github.com/rizkiar00/homework/internal/repository/db"
	emailrepo "github.com/rizkiar00/homework/internal/repository/email"
	"github.com/rizkiar00/homework/internal/repository/email/resend"
	"go.uber.org/dig"
)

func Register(container *dig.Container) error {
	if err := dbrepo.Register(container); err != nil {
		return err
	}

	if err := container.Provide(resend.New); err != nil {
		return err
	}

	return emailrepo.Register(container)
}
