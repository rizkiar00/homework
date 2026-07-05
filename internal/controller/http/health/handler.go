package health

import (
	"net/http"

	"github.com/labstack/echo/v4"
	healthUsecase "github.com/rizkiar00/homework/internal/usecase/health"
	"github.com/rizkiar00/homework/pkg/constant"
	httpresponse "github.com/rizkiar00/homework/pkg/response"
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
		return httpresponse.Error(ctx, http.StatusInternalServerError, constant.CodeInternalServer, err.Error())
	}

	return httpresponse.JSON(ctx, http.StatusOK, constant.CodeSuccess, constant.MessageSuccess, res)
}

func (c *Controller) Readiness(ctx echo.Context) error {
	res, err := c.uc.Readiness(ctx.Request().Context())
	if err != nil {
		return httpresponse.Error(ctx, http.StatusInternalServerError, constant.CodeInternalServer, err.Error())
	}

	status := http.StatusOK
	code := constant.CodeSuccess
	message := constant.MessageSuccess
	if res.Status != constant.StatusReady {
		status = http.StatusServiceUnavailable
		code = constant.CodeServiceUnavailable
		message = constant.MessageServiceUnavailable
	}

	return httpresponse.JSON(ctx, status, code, message, res)
}
