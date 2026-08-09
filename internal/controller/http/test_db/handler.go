package test_db

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rizkiar00/homework/internal/controller/http/middleware"
	"github.com/rizkiar00/homework/internal/model"
	testDBUsecase "github.com/rizkiar00/homework/internal/usecase/test_db"
	"github.com/rizkiar00/homework/pkg/constant"
	"github.com/rizkiar00/homework/pkg/customerror"
	httpresponse "github.com/rizkiar00/homework/pkg/response"
	"github.com/rizkiar00/homework/pkg/token"
)

type Controller struct {
	uc           testDBUsecase.Usecase
	tokenService *token.Service
	blacklist    *token.Blacklist
}

func NewController(uc testDBUsecase.Usecase, tokenService *token.Service, blacklist *token.Blacklist) *Controller {
	return &Controller{
		uc:           uc,
		tokenService: tokenService,
		blacklist:    blacklist,
	}
}

func (c *Controller) RegisterRoutes(e *echo.Echo) {
	group := e.Group("/test_db", middleware.JWT(c.tokenService, c.blacklist))
	group.POST("", c.Create)
	group.GET("", c.FindAll)
	group.GET("/:test_id", c.FindByID)
	group.PUT("/:test_id", c.Update)
	group.DELETE("/:test_id", c.Delete)
}

func (c *Controller) Create(ctx echo.Context) error {
	var request model.CreateTestDBRequest
	if err := ctx.Bind(&request); err != nil {
		return httpresponse.Error(ctx, http.StatusBadRequest, constant.CodeBadRequest, constant.MessageInvalidRequestBody)
	}

	claims, ok := ctx.Get(middleware.UserClaimsKey).(token.Claims)
	if !ok {
		return httpresponse.Error(ctx, http.StatusUnauthorized, constant.CodeUnauthorized, constant.MessageUnauthorized)
	}

	response, err := c.uc.Create(ctx.Request().Context(), claims.UserID, request)
	if err != nil {
		return httpresponse.CustomError(ctx, err)
	}

	return httpresponse.JSON(ctx, http.StatusCreated, constant.CodeSuccess, constant.MessageCreated, response)
}

func (c *Controller) FindAll(ctx echo.Context) error {
	request, err := parseListRequest(ctx)
	if err != nil {
		return httpresponse.CustomError(ctx, err)
	}

	response, err := c.uc.FindAll(ctx.Request().Context(), request)
	if err != nil {
		return httpresponse.CustomError(ctx, err)
	}

	return httpresponse.JSONWithMeta(ctx, http.StatusOK, constant.CodeSuccess, constant.MessageSuccess, response.Data, response.Meta)
}

func (c *Controller) FindByID(ctx echo.Context) error {
	response, err := c.uc.FindByID(ctx.Request().Context(), ctx.Param("test_id"))
	if err != nil {
		return httpresponse.CustomError(ctx, err)
	}

	return httpresponse.JSON(ctx, http.StatusOK, constant.CodeSuccess, constant.MessageSuccess, response)
}

func (c *Controller) Update(ctx echo.Context) error {
	var request model.UpdateTestDBRequest
	if err := ctx.Bind(&request); err != nil {
		return httpresponse.Error(ctx, http.StatusBadRequest, constant.CodeBadRequest, constant.MessageInvalidRequestBody)
	}

	response, err := c.uc.Update(ctx.Request().Context(), ctx.Param("test_id"), request)
	if err != nil {
		return httpresponse.CustomError(ctx, err)
	}

	return httpresponse.JSON(ctx, http.StatusOK, constant.CodeSuccess, constant.MessageSuccess, response)
}

func (c *Controller) Delete(ctx echo.Context) error {
	if err := c.uc.Delete(ctx.Request().Context(), ctx.Param("test_id")); err != nil {
		return httpresponse.CustomError(ctx, err)
	}

	return httpresponse.JSON(ctx, http.StatusOK, constant.CodeSuccess, constant.MessageDeleted, nil)
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
		return 0, customerror.BadRequest(name + " must be a number")
	}

	return parsed, nil
}
