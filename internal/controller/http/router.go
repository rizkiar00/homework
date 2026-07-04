package http

import (
	"github.com/labstack/echo/v4"
	"github.com/rizkiar00/homework/internal/controller/http/health"
	"github.com/rizkiar00/homework/internal/controller/http/test_db"
	"go.uber.org/dig"
)

type RouterParams struct {
	dig.In

	HealthController *health.Controller
	TestDBController *test_db.Controller
}

func Register(container *dig.Container) error {
	if err := health.Register(container); err != nil {
		return err
	}

	if err := test_db.Register(container); err != nil {
		return err
	}

	return container.Provide(NewRouter)
}

func NewRouter(params RouterParams) *echo.Echo {
	e := echo.New()
	e.HideBanner = false

	params.HealthController.RegisterRoutes(e)
	params.TestDBController.RegisterRoutes(e)
	registerSwaggerRoutes(e)

	return e
}
