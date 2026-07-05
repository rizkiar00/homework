package http

import (
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/rizkiar00/homework/internal/controller/http/access"
	"github.com/rizkiar00/homework/internal/controller/http/auth"
	"github.com/rizkiar00/homework/internal/controller/http/health"
	"github.com/rizkiar00/homework/internal/controller/http/middleware"
	"github.com/rizkiar00/homework/internal/controller/http/test_db"
	"github.com/rizkiar00/homework/pkg/config"
	"github.com/rizkiar00/homework/pkg/token"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
)

type RouterParams struct {
	dig.In

	AuthController   *auth.Controller
	AccessController *access.Controller
	HealthController *health.Controller
	TestDBController *test_db.Controller
	TokenService     *token.Service
	Logger           *logrus.Logger
	Config           config.Config
}

func Register(container *dig.Container) error {
	if err := health.Register(container); err != nil {
		return err
	}
	if err := auth.Register(container); err != nil {
		return err
	}
	if err := access.Register(container); err != nil {
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
	e.HTTPErrorHandler = middleware.HTTPErrorHandler(params.Logger)
	e.Use(middleware.RequestID())
	e.Use(middleware.RequestLogger(params.Logger))
	e.Use(middleware.Recover(params.Logger))
	e.Use(echoMiddleware.CORSWithConfig(middleware.CORS(params.Config.HTTP)))
	e.Use(echoMiddleware.SecureWithConfig(middleware.Secure()))
	e.Use(echoMiddleware.ContextTimeoutWithConfig(middleware.Timeout(params.Config.HTTP)))

	handler := &APIHandler{
		AuthController:   params.AuthController,
		AccessController: params.AccessController,
		HealthController: params.HealthController,
		TestDBController: params.TestDBController,
	}
	RegisterHandlersWithOptions(e, handler, RegisterHandlersOptions{
		OperationMiddlewares: map[string][]echo.MiddlewareFunc{
			"GetMe":          {middleware.JWT(params.TokenService)},
			"GetTestDBList":  {middleware.JWT(params.TokenService)},
			"CreateTestDB":   {middleware.JWT(params.TokenService)},
			"GetTestDBByID":  {middleware.JWT(params.TokenService)},
			"UpdateTestDB":   {middleware.JWT(params.TokenService)},
			"DeleteTestDB":   {middleware.JWT(params.TokenService)},
			"GetActions":     {middleware.JWT(params.TokenService)},
			"CreateRole":     {middleware.JWT(params.TokenService)},
			"UpdateRole":     {middleware.JWT(params.TokenService)},
			"SetRoleActions": {middleware.JWT(params.TokenService)},
			"AssignUserRole": {middleware.JWT(params.TokenService)},
		},
	})
	registerSwaggerRoutes(e)

	return e
}
