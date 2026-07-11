package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rizkiar00/homework/pkg/constant"
	httpresponse "github.com/rizkiar00/homework/pkg/response"
	"github.com/rizkiar00/homework/pkg/token"
)

const UserClaimsKey = constant.ContextUserClaimsKey

func JWT(tokenService *token.Service, blacklist *token.Blacklist) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			header := ctx.Request().Header.Get(echo.HeaderAuthorization)
			if header == "" {
				return httpresponse.Error(ctx, http.StatusUnauthorized, constant.CodeUnauthorized, constant.MessageUnauthorized)
			}

			value := strings.TrimPrefix(header, constant.AuthorizationBearerPrefix)
			if value == header || value == "" {
				return httpresponse.Error(ctx, http.StatusUnauthorized, constant.CodeUnauthorized, constant.MessageUnauthorized)
			}

			claims, err := tokenService.Parse(value)
			if err != nil {
				return httpresponse.Error(ctx, http.StatusUnauthorized, constant.CodeUnauthorized, constant.MessageUnauthorized)
			}
			if blacklist != nil && blacklist.IsRevoked(ctx.Request().Context(), claims) {
				return httpresponse.Error(ctx, http.StatusUnauthorized, constant.CodeUnauthorized, constant.MessageUnauthorized)
			}

			ctx.Set(UserClaimsKey, claims)
			return next(ctx)
		}
	}
}
