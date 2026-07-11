package http

import (
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/rizkiar00/homework/internal/controller/http/access"
	"github.com/rizkiar00/homework/internal/controller/http/auth"
	"github.com/rizkiar00/homework/internal/controller/http/health"
	"github.com/rizkiar00/homework/internal/controller/http/middleware"
	"github.com/rizkiar00/homework/internal/controller/http/test_db"
	accessRepo "github.com/rizkiar00/homework/internal/repository/db/access"
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
	AccessRepository accessRepo.Repository
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
	e.Use(echoMiddleware.BodyLimitWithConfig(middleware.BodyLimit(params.Config.HTTP)))
	e.Use(middleware.RateLimit(params.Config.HTTP.RateLimitRequestsPerMinute, params.Config.HTTP.RateLimitBurst))
	e.Use(echoMiddleware.ContextTimeoutWithConfig(middleware.Timeout(params.Config.HTTP)))

	handler := &APIHandler{
		AuthController:   params.AuthController,
		AccessController: params.AccessController,
		HealthController: params.HealthController,
		TestDBController: params.TestDBController,
	}
	RegisterHandlersWithOptions(e, handler, RegisterHandlersOptions{
		OperationMiddlewares: map[string][]echo.MiddlewareFunc{
			"Register":       authRateLimitMiddleware(params.Config.HTTP),
			"Login":          authRateLimitMiddleware(params.Config.HTTP),
			"GetMe":          privateMiddlewares(params.TokenService, params.AccessRepository),
			"GetTestDBList":  privateMiddlewares(params.TokenService, params.AccessRepository),
			"CreateTestDB":   privateMiddlewares(params.TokenService, params.AccessRepository),
			"GetTestDBByID":  privateMiddlewares(params.TokenService, params.AccessRepository),
			"UpdateTestDB":   privateMiddlewares(params.TokenService, params.AccessRepository),
			"DeleteTestDB":   privateMiddlewares(params.TokenService, params.AccessRepository),
			"GetActions":     privateMiddlewares(params.TokenService, params.AccessRepository),
			"CreateRole":     privateMiddlewares(params.TokenService, params.AccessRepository),
			"UpdateRole":     privateMiddlewares(params.TokenService, params.AccessRepository),
			"SetRoleActions": privateMiddlewares(params.TokenService, params.AccessRepository),
			"AssignUserRole": privateMiddlewares(params.TokenService, params.AccessRepository),
		},
	})
	registerSwaggerRoutes(e)

	return e
}

func authRateLimitMiddleware(cfg config.HTTPConfig) []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		middleware.RateLimit(cfg.AuthRateLimitRequestsPerMinute, cfg.AuthRateLimitBurst),
	}
}

func privateMiddlewares(tokenService *token.Service, repo accessRepo.Repository) []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		middleware.JWT(tokenService),
		middleware.ActionAccess(repo),
	}
}
