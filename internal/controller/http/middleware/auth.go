package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rizkiar00/homework/pkg/token"
)

const UserClaimsKey = "user_claims"

func JWT(tokenService *token.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			header := ctx.Request().Header.Get(echo.HeaderAuthorization)
			if header == "" {
				return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			}

			value := strings.TrimPrefix(header, "Bearer ")
			if value == header || value == "" {
				return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			}

			claims, err := tokenService.Parse(value)
			if err != nil {
				return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			}

			ctx.Set(UserClaimsKey, claims)
			return next(ctx)
		}
	}
}
