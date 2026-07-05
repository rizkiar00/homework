package db

import (
	"github.com/rizkiar00/homework/internal/repository/db/access"
	"github.com/rizkiar00/homework/internal/repository/db/health"
	"github.com/rizkiar00/homework/internal/repository/db/test_db"
	"github.com/rizkiar00/homework/internal/repository/db/user"
	"go.uber.org/dig"
)

func Register(container *dig.Container) error {
	if err := health.Register(container); err != nil {
		return err
	}
	if err := user.Register(container); err != nil {
		return err
	}
	if err := access.Register(container); err != nil {
		return err
	}

	return test_db.Register(container)
}
