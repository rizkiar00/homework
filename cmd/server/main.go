package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/rizkiar00/homework/pkg/config"
	"github.com/rizkiar00/homework/pkg/di"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
)

func main() {
	container, err := di.NewContainer()
	if err != nil {
		log.Fatal(err)
	}

	if err := container.Invoke(run); err != nil {
		log.Fatal(err)
	}
}

type appParams struct {
	dig.In

	Config config.Config
	Logger *logrus.Logger
	Router *echo.Echo
}

func run(params appParams) error {
	params.Logger.WithFields(logrus.Fields{
		"service": params.Config.AppConfig.Name,
		"env":     params.Config.AppConfig.Env,
		"address": params.Config.AppConfig.Address(),
	}).Info("starting server")

	return params.Router.Start(params.Config.AppConfig.Address())
}
