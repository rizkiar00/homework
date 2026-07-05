package http

import (
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime/types"
	"github.com/rizkiar00/homework/internal/controller/http/access"
	"github.com/rizkiar00/homework/internal/controller/http/auth"
	"github.com/rizkiar00/homework/internal/controller/http/health"
	"github.com/rizkiar00/homework/internal/controller/http/test_db"
)

type APIHandler struct {
	AuthController   *auth.Controller
	AccessController *access.Controller
	HealthController *health.Controller
	TestDBController *test_db.Controller
}

func (h *APIHandler) Login(ctx echo.Context) error {
	return h.AuthController.Login(ctx)
}

func (h *APIHandler) GetMe(ctx echo.Context) error {
	return h.AuthController.Me(ctx)
}

func (h *APIHandler) Register(ctx echo.Context) error {
	return h.AuthController.Register(ctx)
}

func (h *APIHandler) GetHealth(ctx echo.Context) error {
	return h.HealthController.Health(ctx)
}

func (h *APIHandler) GetReadiness(ctx echo.Context) error {
	return h.HealthController.Readiness(ctx)
}

func (h *APIHandler) GetTestDBList(ctx echo.Context, params GetTestDBListParams) error {
	return h.TestDBController.FindAll(ctx)
}

func (h *APIHandler) CreateTestDB(ctx echo.Context) error {
	return h.TestDBController.Create(ctx)
}

func (h *APIHandler) DeleteTestDB(ctx echo.Context, idTest types.UUID) error {
	return h.TestDBController.Delete(ctx)
}

func (h *APIHandler) GetTestDBByID(ctx echo.Context, idTest types.UUID) error {
	return h.TestDBController.FindByID(ctx)
}

func (h *APIHandler) UpdateTestDB(ctx echo.Context, idTest types.UUID) error {
	return h.TestDBController.Update(ctx)
}

func (h *APIHandler) GetActions(ctx echo.Context) error {
	return h.AccessController.FindAllActions(ctx)
}

func (h *APIHandler) CreateRole(ctx echo.Context) error {
	return h.AccessController.CreateRole(ctx)
}

func (h *APIHandler) UpdateRole(ctx echo.Context, roleID int64) error {
	return h.AccessController.UpdateRole(ctx, roleID)
}

func (h *APIHandler) SetRoleActions(ctx echo.Context, roleID int64) error {
	return h.AccessController.SetRoleActions(ctx, roleID)
}

func (h *APIHandler) AssignUserRole(ctx echo.Context, userID types.UUID) error {
	return h.AccessController.AssignUserRole(ctx)
}
