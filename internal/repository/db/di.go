package db

import (
	"github.com/rizkiar00/homework/internal/repository/db/health"
	"github.com/rizkiar00/homework/internal/repository/db/test_db"
	"go.uber.org/dig"
)

func Register(container *dig.Container) error {
	if err := health.Register(container); err != nil {
		return err
	}

	return test_db.Register(container)
}
