package repository

import (
	dbrepo "github.com/rizkiar00/homework/internal/repository/db"
	"go.uber.org/dig"
)

func Register(container *dig.Container) error {
	return dbrepo.Register(container)
}
