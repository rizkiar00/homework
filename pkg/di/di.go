package di

import (
	"github.com/rizkiar00/homework/internal/controller/http"
	"github.com/rizkiar00/homework/internal/repository"
	"github.com/rizkiar00/homework/internal/usecase"
	"github.com/rizkiar00/homework/pkg/config"
	"github.com/rizkiar00/homework/pkg/resource"
	"go.uber.org/dig"
)

func NewContainer() (*dig.Container, error) {
	container := dig.New()

	registers := []func(*dig.Container) error{
		config.Register,
		resource.Register,
		repository.Register,
		usecase.Register,
		http.Register,
	}

	for _, register := range registers {
		if err := register(container); err != nil {
			return nil, err
		}
	}

	return container, nil
}
