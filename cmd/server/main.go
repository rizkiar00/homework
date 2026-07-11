package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rizkiar00/homework/internal/model"
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
	DB     model.Database
	Redis  model.Redis
}

func run(params appParams) error {
	params.Logger.WithFields(logrus.Fields{
		"service": params.Config.AppConfig.Name,
		"env":     params.Config.AppConfig.Env,
		"address": params.Config.AppConfig.Address(),
	}).Info("starting server")

	serverErr := make(chan error, 1)
	go func() {
		if err := params.Router.Start(params.Config.AppConfig.Address()); err != nil && err != http.ErrServerClosed {
			serverErr <- err
			return
		}
		serverErr <- nil
	}()

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(shutdownSignal)

	select {
	case err := <-serverErr:
		return err
	case sig := <-shutdownSignal:
		params.Logger.WithField("signal", sig.String()).Info("shutdown signal received")
	}

	timeout := time.Duration(params.Config.AppConfig.ShutdownTimeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := params.Router.Shutdown(ctx); err != nil {
		params.Logger.WithError(err).Error("server shutdown failed")
		return err
	}

	if err := closeDatabase(params.DB); err != nil {
		params.Logger.WithError(err).Error("database close failed")
		return err
	}

	if err := closeRedis(params.Redis); err != nil {
		params.Logger.WithError(err).Error("redis close failed")
		return err
	}

	params.Logger.Info("server stopped gracefully")
	return nil
}

func closeDatabase(db model.Database) error {
	if db == nil {
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

func closeRedis(redis model.Redis) error {
	if redis == nil {
		return nil
	}

	return redis.Close()
}
