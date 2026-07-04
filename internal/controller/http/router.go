package http

import (
	"github.com/labstack/echo/v4"
	"github.com/rizkiar00/homework/internal/controller/http/auth"
	"github.com/rizkiar00/homework/internal/controller/http/health"
	"github.com/rizkiar00/homework/internal/controller/http/middleware"
	"github.com/rizkiar00/homework/internal/controller/http/test_db"
	"github.com/rizkiar00/homework/pkg/token"
	"go.uber.org/dig"
)

type RouterParams struct {
	dig.In

	AuthController   *auth.Controller
	HealthController *health.Controller
	TestDBController *test_db.Controller
	TokenService     *token.Service
}

func Register(container *dig.Container) error {
	if err := health.Register(container); err != nil {
		return err
	}
	if err := auth.Register(container); err != nil {
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

	handler := &APIHandler{
		AuthController:   params.AuthController,
		HealthController: params.HealthController,
		TestDBController: params.TestDBController,
	}
	RegisterHandlersWithOptions(e, handler, RegisterHandlersOptions{
		OperationMiddlewares: map[string][]echo.MiddlewareFunc{
			"GetMe":         {middleware.JWT(params.TokenService)},
			"GetTestDBList": {middleware.JWT(params.TokenService)},
			"CreateTestDB":  {middleware.JWT(params.TokenService)},
			"GetTestDBByID": {middleware.JWT(params.TokenService)},
			"UpdateTestDB":  {middleware.JWT(params.TokenService)},
			"DeleteTestDB":  {middleware.JWT(params.TokenService)},
		},
	})
	registerSwaggerRoutes(e)

	return e
}
