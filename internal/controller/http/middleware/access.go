package middleware

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rizkiar00/homework/internal/model"
	accessRepo "github.com/rizkiar00/homework/internal/repository/db/access"
	"github.com/rizkiar00/homework/pkg/constant"
	httpresponse "github.com/rizkiar00/homework/pkg/response"
	"github.com/rizkiar00/homework/pkg/token"
)

func ActionAccess(repo accessRepo.Repository, redis model.Redis) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			claims, ok := ctx.Get(UserClaimsKey).(token.Claims)
			if !ok {
				return httpresponse.Error(ctx, http.StatusUnauthorized, constant.CodeUnauthorized, constant.MessageUnauthorized)
			}

			method := ctx.Request().Method
			endpoint := normalizeEndpoint(ctx.Path())

			allowed, ok := cachedActionAccess(ctx, redis, claims.UserID, method, endpoint)
			if ok && !allowed {
				return httpresponse.Error(ctx, http.StatusForbidden, constant.CodeForbidden, constant.MessageForbidden)
			}
			if ok && allowed {
				return next(ctx)
			}

			allowed, err := repo.HasAccess(ctx.Request().Context(), claims.UserID, method, endpoint)
			if err != nil {
				return httpresponse.CustomError(ctx, err)
			}
			cacheActionAccess(ctx, redis, claims.UserID, method, endpoint, allowed)
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

func cachedActionAccess(ctx echo.Context, redis model.Redis, userID string, method string, endpoint string) (bool, bool) {
	if redis == nil {
		return false, false
	}

	value, err := redis.Get(ctx.Request().Context(), actionAccessKey(userID, method, endpoint)).Result()
	if err != nil {
		return false, false
	}

	var allowed bool
	if err := json.Unmarshal([]byte(value), &allowed); err != nil {
		return false, false
	}

	return allowed, true
}

func cacheActionAccess(ctx echo.Context, redis model.Redis, userID string, method string, endpoint string, allowed bool) {
	if redis == nil {
		return
	}

	value, err := json.Marshal(allowed)
	if err != nil {
		return
	}

	_ = redis.Set(ctx.Request().Context(), actionAccessKey(userID, method, endpoint), value, 5*time.Minute).Err()
}

func actionAccessKey(userID string, method string, endpoint string) string {
	return "access:user:" + userID + ":" + method + ":" + endpoint
}
