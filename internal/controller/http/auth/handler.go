package auth

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rizkiar00/homework/internal/controller/http/middleware"
	"github.com/rizkiar00/homework/internal/model"
	authUsecase "github.com/rizkiar00/homework/internal/usecase/auth"
	"github.com/rizkiar00/homework/pkg/constant"
	"github.com/rizkiar00/homework/pkg/token"
	"gorm.io/gorm"
)

type Controller struct {
	uc           authUsecase.Usecase
	tokenService *token.Service
}

func NewController(uc authUsecase.Usecase, tokenService *token.Service) *Controller {
	return &Controller{
		uc:           uc,
		tokenService: tokenService,
	}
}

func (c *Controller) RegisterRoutes(e *echo.Echo) {
	e.POST("/auth/register", c.Register)
	e.POST("/auth/login", c.Login)
	e.GET("/auth/me", c.Me, middleware.JWT(c.tokenService))
}

func (c *Controller) Register(ctx echo.Context) error {
	var request model.RegisterRequest
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(constant.MessageInvalidRequestBody))
	}

	response, err := c.uc.Register(ctx.Request().Context(), request)
	if err != nil {
		return writeError(ctx, err)
	}

	return ctx.JSON(http.StatusCreated, response)
}

func (c *Controller) Login(ctx echo.Context) error {
	var request model.LoginRequest
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse(constant.MessageInvalidRequestBody))
	}

	response, err := c.uc.Login(ctx.Request().Context(), request)
	if err != nil {
		return writeError(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

func (c *Controller) Me(ctx echo.Context) error {
	claims, ok := ctx.Get(middleware.UserClaimsKey).(token.Claims)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, errorResponse(constant.MessageUnauthorized))
	}

	response, err := c.uc.Me(ctx.Request().Context(), claims.UserID)
	if err != nil {
		return writeError(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

func writeError(ctx echo.Context, err error) error {
	if errors.Is(err, model.ErrInvalidCredential) {
		return ctx.JSON(http.StatusUnauthorized, errorResponse(err.Error()))
	}
	if errors.Is(err, model.ErrUsernameAlreadyExists) {
		return ctx.JSON(http.StatusConflict, errorResponse(err.Error()))
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ctx.JSON(http.StatusNotFound, errorResponse(constant.MessageDataNotFound))
	}

	return ctx.JSON(http.StatusInternalServerError, errorResponse(err.Error()))
}

func errorResponse(message string) map[string]string {
	return map[string]string{"error": message}
}
