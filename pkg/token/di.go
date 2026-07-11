package token

import "go.uber.org/dig"

func Register(container *dig.Container) error {
	if err := container.Provide(New); err != nil {
		return err
	}

	return container.Provide(NewBlacklist)
}
