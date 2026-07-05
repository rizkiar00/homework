package model

type ActionResponse struct {
	ActionID   int64  `json:"action_id"`
	ActionDesc string `json:"action_desc"`
	ActionType string `json:"action_type"`
	Endpoint   string `json:"endpoint"`
}

type CreateRoleRequest struct {
	RoleDesc  string  `json:"role_desc"`
	ActionIDs []int64 `json:"action_ids,omitempty"`
}

type UpdateRoleRequest struct {
	RoleDesc  *string `json:"role_desc,omitempty"`
	IsActive  *bool   `json:"is_active,omitempty"`
	ActionIDs []int64 `json:"action_ids,omitempty"`
}

type SetRoleActionsRequest struct {
	ActionIDs []int64 `json:"action_ids"`
}

type AssignUserRoleRequest struct {
	RoleID int64 `json:"role_id"`
}

type RoleResponse struct {
	RoleID    int64            `json:"role_id"`
	RoleDesc  string           `json:"role_desc"`
	IsActive  bool             `json:"is_active"`
	ActionIDs []int64          `json:"action_ids,omitempty"`
	Actions   []ActionResponse `json:"actions,omitempty"`
}
