package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	accessRepo "github.com/rizkiar00/homework/internal/repository/db/access"
	"github.com/rizkiar00/homework/pkg/constant"
	httpresponse "github.com/rizkiar00/homework/pkg/response"
	"github.com/rizkiar00/homework/pkg/token"
)

func ActionAccess(repo accessRepo.Repository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			claims, ok := ctx.Get(UserClaimsKey).(token.Claims)
			if !ok {
				return httpresponse.Error(ctx, http.StatusUnauthorized, constant.CodeUnauthorized, constant.MessageUnauthorized)
			}

			allowed, err := repo.HasAccess(ctx.Request().Context(), claims.UserID, ctx.Request().Method, normalizeEndpoint(ctx.Path()))
			if err != nil {
				return httpresponse.CustomError(ctx, err)
			}
			if !allowed {
				return httpresponse.Error(ctx, http.StatusForbidden, constant.CodeForbidden, constant.MessageForbidden)
			}

			return next(ctx)
		}
	}
}

func normalizeEndpoint(path string) string {
	parts := strings.Split(path, "/")
	for index, part := range parts {
		if strings.HasPrefix(part, ":") {
			parts[index] = "{" + strings.TrimPrefix(part, ":") + "}"
		}
	}

	return strings.Join(parts, "/")
}
