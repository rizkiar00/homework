package db

import (
	"github.com/rizkiar00/homework/internal/repository/db/health"
	"go.uber.org/dig"
)

func Register(container *dig.Container) error {
	return health.Register(container)
}
