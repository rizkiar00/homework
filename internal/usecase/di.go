package usecase

import (
	"github.com/rizkiar00/homework/internal/usecase/auth"
	"github.com/rizkiar00/homework/internal/usecase/health"
	"github.com/rizkiar00/homework/internal/usecase/test_db"
	"go.uber.org/dig"
)

func Register(container *dig.Container) error {
	if err := health.Register(container); err != nil {
		return err
	}
	if err := auth.Register(container); err != nil {
		return err
	}

	return test_db.Register(container)
}
