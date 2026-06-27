package config

import "go.uber.org/dig"

func Register(container *dig.Container) error {
	return container.Provide(New)
}
