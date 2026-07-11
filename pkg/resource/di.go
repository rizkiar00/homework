package resource

import "go.uber.org/dig"

func Register(container *dig.Container) error {
	if err := container.Provide(NewLogger); err != nil {
		return err
	}

	if err := container.Provide(NewDatabase); err != nil {
		return err
	}

	return container.Provide(NewRedis)
}
