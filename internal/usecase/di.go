package usecase

import (
	"github.com/rizkiar00/homework/internal/usecase/health"
	"go.uber.org/dig"
)

func Register(container *dig.Container) error {
	return health.Register(container)
}
