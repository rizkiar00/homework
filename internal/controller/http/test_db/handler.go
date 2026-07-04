package test_db

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rizkiar00/homework/internal/controller/http/middleware"
	"github.com/rizkiar00/homework/internal/model"
	testDBUsecase "github.com/rizkiar00/homework/internal/usecase/test_db"
	"github.com/rizkiar00/homework/pkg/token"
	"gorm.io/gorm"
)

type Controller struct {
	uc           testDBUsecase.Usecase
	tokenService *token.Service
}

func NewController(uc testDBUsecase.Usecase, tokenService *token.Service) *Controller {
	return &Controller{
		uc:           uc,
		tokenService: tokenService,
	}
}

func (c *Controller) RegisterRoutes(e *echo.Echo) {
	group := e.Group("/test_db", middleware.JWT(c.tokenService))
	group.POST("", c.Create)
	group.GET("", c.FindAll)
	group.GET("/:id_test", c.FindByID)
	group.PUT("/:id_test", c.Update)
	group.DELETE("/:id_test", c.Delete)
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
	request, err := parseListRequest(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err.Error()))
	}

	response, err := c.uc.FindAll(ctx.Request().Context(), request)
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
	if strings.Contains(err.Error(), "order_by must be") || strings.Contains(err.Error(), "order_dir must be") {
		return ctx.JSON(http.StatusBadRequest, errorResponse(err.Error()))
	}

	return ctx.JSON(http.StatusInternalServerError, errorResponse(err.Error()))
}

func errorResponse(message string) map[string]string {
	return map[string]string{"error": message}
}

func parseListRequest(ctx echo.Context) (model.TestDBListRequest, error) {
	page, err := parseIntQuery(ctx, "page", 1)
	if err != nil {
		return model.TestDBListRequest{}, err
	}

	limit, err := parseIntQuery(ctx, "limit", 10)
	if err != nil {
		return model.TestDBListRequest{}, err
	}

	return model.TestDBListRequest{
		Page:     page,
		Limit:    limit,
		OrderBy:  ctx.QueryParam("order_by"),
		OrderDir: ctx.QueryParam("order_dir"),
	}, nil
}

func parseIntQuery(ctx echo.Context, name string, defaultValue int) (int, error) {
	value := ctx.QueryParam(name)
	if value == "" {
		return defaultValue, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, errors.New(name + " must be a number")
	}

	return parsed, nil
}
