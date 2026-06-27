package health

import (
	"net/http"

	"github.com/labstack/echo/v4"
	healthUsecase "github.com/rizkiar00/homework/internal/usecase/health"
)

type Controller struct {
	uc healthUsecase.Usecase
}

func NewController(uc healthUsecase.Usecase) *Controller {
	return &Controller{uc: uc}
}

func (c *Controller) RegisterRoutes(e *echo.Echo) {
	e.GET("/health", c.Health)
	e.GET("/readiness", c.Readiness)
}

func (c *Controller) Health(ctx echo.Context) error {
	res, err := c.uc.Health(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, res)
}

func (c *Controller) Readiness(ctx echo.Context) error {
	res, err := c.uc.Readiness(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	status := http.StatusOK
	if res.Status != "ready" {
		status = http.StatusServiceUnavailable
	}

	return ctx.JSON(status, res)
}
