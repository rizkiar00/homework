package access

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rizkiar00/homework/internal/model"
	accessUsecase "github.com/rizkiar00/homework/internal/usecase/access"
	"github.com/rizkiar00/homework/pkg/constant"
	httpresponse "github.com/rizkiar00/homework/pkg/response"
)

type Controller struct {
	uc accessUsecase.Usecase
}

func NewController(uc accessUsecase.Usecase) *Controller {
	return &Controller{uc: uc}
}

func (c *Controller) FindAllActions(ctx echo.Context) error {
	response, err := c.uc.FindAllActions(ctx.Request().Context())
	if err != nil {
		return httpresponse.CustomError(ctx, err)
	}

	return httpresponse.JSON(ctx, http.StatusOK, constant.CodeSuccess, constant.MessageSuccess, response)
}

func (c *Controller) CreateRole(ctx echo.Context) error {
	var request model.CreateRoleRequest
	if err := ctx.Bind(&request); err != nil {
		return httpresponse.Error(ctx, http.StatusBadRequest, constant.CodeBadRequest, constant.MessageInvalidRequestBody)
	}

	response, err := c.uc.CreateRole(ctx.Request().Context(), request)
	if err != nil {
		return httpresponse.CustomError(ctx, err)
	}

	return httpresponse.JSON(ctx, http.StatusCreated, constant.CodeSuccess, constant.MessageCreated, response)
}

func (c *Controller) UpdateRole(ctx echo.Context, roleID int64) error {
	var request model.UpdateRoleRequest
	if err := ctx.Bind(&request); err != nil {
		return httpresponse.Error(ctx, http.StatusBadRequest, constant.CodeBadRequest, constant.MessageInvalidRequestBody)
	}

	response, err := c.uc.UpdateRole(ctx.Request().Context(), roleID, request)
	if err != nil {
		return httpresponse.CustomError(ctx, err)
	}

	return httpresponse.JSON(ctx, http.StatusOK, constant.CodeSuccess, constant.MessageSuccess, response)
}

func (c *Controller) SetRoleActions(ctx echo.Context, roleID int64) error {
	var request model.SetRoleActionsRequest
	if err := ctx.Bind(&request); err != nil {
		return httpresponse.Error(ctx, http.StatusBadRequest, constant.CodeBadRequest, constant.MessageInvalidRequestBody)
	}

	response, err := c.uc.SetRoleActions(ctx.Request().Context(), roleID, request)
	if err != nil {
		return httpresponse.CustomError(ctx, err)
	}

	return httpresponse.JSON(ctx, http.StatusOK, constant.CodeSuccess, constant.MessageSuccess, response)
}

func (c *Controller) AssignUserRole(ctx echo.Context) error {
	var request model.AssignUserRoleRequest
	if err := ctx.Bind(&request); err != nil {
		return httpresponse.Error(ctx, http.StatusBadRequest, constant.CodeBadRequest, constant.MessageInvalidRequestBody)
	}

	if err := c.uc.AssignUserRole(ctx.Request().Context(), ctx.Param("user_id"), request); err != nil {
		return httpresponse.CustomError(ctx, err)
	}

	return httpresponse.JSON(ctx, http.StatusOK, constant.CodeSuccess, constant.MessageSuccess, nil)
}

func ParseInt64Param(value string, name string) (int64, error) {
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusBadRequest, name+" must be a number")
	}
	return parsed, nil
}
