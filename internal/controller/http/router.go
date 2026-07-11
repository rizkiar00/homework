package http

import (
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/rizkiar00/homework/internal/controller/http/access"
	"github.com/rizkiar00/homework/internal/controller/http/auth"
	"github.com/rizkiar00/homework/internal/controller/http/health"
	"github.com/rizkiar00/homework/internal/controller/http/middleware"
	"github.com/rizkiar00/homework/internal/controller/http/test_db"
	"github.com/rizkiar00/homework/internal/model"
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
	TokenBlacklist   *token.Blacklist
	Logger           *logrus.Logger
	Config           config.Config
	Redis            model.Redis
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
	e.Use(middleware.RateLimit(params.Config.HTTP.RateLimitRequestsPerMinute, params.Config.HTTP.RateLimitBurst, params.Redis, "global"))
	e.Use(echoMiddleware.ContextTimeoutWithConfig(middleware.Timeout(params.Config.HTTP)))

	handler := &APIHandler{
		AuthController:   params.AuthController,
		AccessController: params.AccessController,
		HealthController: params.HealthController,
		TestDBController: params.TestDBController,
	}
	RegisterHandlersWithOptions(e, handler, RegisterHandlersOptions{
		OperationMiddlewares: map[string][]echo.MiddlewareFunc{
			"Register":       authRateLimitMiddleware(params.Config.HTTP, params.Redis),
			"Login":          authRateLimitMiddleware(params.Config.HTTP, params.Redis),
			"Logout":         privateMiddlewares(params.TokenService, params.TokenBlacklist, params.AccessRepository, params.Redis),
			"GetMe":          privateMiddlewares(params.TokenService, params.TokenBlacklist, params.AccessRepository, params.Redis),
			"GetTestDBList":  privateMiddlewares(params.TokenService, params.TokenBlacklist, params.AccessRepository, params.Redis),
			"CreateTestDB":   privateMiddlewares(params.TokenService, params.TokenBlacklist, params.AccessRepository, params.Redis),
			"GetTestDBByID":  privateMiddlewares(params.TokenService, params.TokenBlacklist, params.AccessRepository, params.Redis),
			"UpdateTestDB":   privateMiddlewares(params.TokenService, params.TokenBlacklist, params.AccessRepository, params.Redis),
			"DeleteTestDB":   privateMiddlewares(params.TokenService, params.TokenBlacklist, params.AccessRepository, params.Redis),
			"GetActions":     privateMiddlewares(params.TokenService, params.TokenBlacklist, params.AccessRepository, params.Redis),
			"CreateRole":     privateMiddlewares(params.TokenService, params.TokenBlacklist, params.AccessRepository, params.Redis),
			"UpdateRole":     privateMiddlewares(params.TokenService, params.TokenBlacklist, params.AccessRepository, params.Redis),
			"SetRoleActions": privateMiddlewares(params.TokenService, params.TokenBlacklist, params.AccessRepository, params.Redis),
			"AssignUserRole": privateMiddlewares(params.TokenService, params.TokenBlacklist, params.AccessRepository, params.Redis),
		},
	})
	registerSwaggerRoutes(e)

	return e
}

func authRateLimitMiddleware(cfg config.HTTPConfig, redis model.Redis) []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		middleware.RateLimit(cfg.AuthRateLimitRequestsPerMinute, cfg.AuthRateLimitBurst, redis, "auth"),
	}
}

func privateMiddlewares(tokenService *token.Service, blacklist *token.Blacklist, repo accessRepo.Repository, redis model.Redis) []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		middleware.JWT(tokenService, blacklist),
		middleware.ActionAccess(repo, redis),
	}
}
