package test_db

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rizkiar00/homework/internal/model"
	testDBUsecase "github.com/rizkiar00/homework/internal/usecase/test_db"
	"gorm.io/gorm"
)

type Controller struct {
	uc testDBUsecase.Usecase
}

func NewController(uc testDBUsecase.Usecase) *Controller {
	return &Controller{uc: uc}
}

func (c *Controller) RegisterRoutes(e *echo.Echo) {
	e.POST("/test_db", c.Create)
	e.GET("/test_db", c.FindAll)
	e.GET("/test_db/:id_test", c.FindByID)
	e.PUT("/test_db/:id_test", c.Update)
	e.DELETE("/test_db/:id_test", c.Delete)
}

func (c *Controller) Create(ctx echo.Context) error {
	var request model.CreateTestDBRequest
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse("invalid request body"))
	}

	response, err := c.uc.Create(ctx.Request().Context(), request)
	if err != nil {
		return writeError(ctx, err)
	}

	return ctx.JSON(http.StatusCreated, response)
}

func (c *Controller) FindAll(ctx echo.Context) error {
	response, err := c.uc.FindAll(ctx.Request().Context())
	if err != nil {
		return writeError(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *Controller) FindByID(ctx echo.Context) error {
	response, err := c.uc.FindByID(ctx.Request().Context(), ctx.Param("id_test"))
	if err != nil {
		return writeError(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *Controller) Update(ctx echo.Context) error {
	var request model.UpdateTestDBRequest
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse("invalid request body"))
	}

	response, err := c.uc.Update(ctx.Request().Context(), ctx.Param("id_test"), request)
	if err != nil {
		return writeError(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *Controller) Delete(ctx echo.Context) error {
	if err := c.uc.Delete(ctx.Request().Context(), ctx.Param("id_test")); err != nil {
		return writeError(ctx, err)
	}

	return ctx.NoContent(http.StatusNoContent)
}

func writeError(ctx echo.Context, err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ctx.JSON(http.StatusNotFound, errorResponse("data not found"))
	}

	return ctx.JSON(http.StatusInternalServerError, errorResponse(err.Error()))
}

func errorResponse(message string) map[string]string {
	return map[string]string{"error": message}
}
