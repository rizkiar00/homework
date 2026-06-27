package model

type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Env     string `json:"env"`
}

type ReadinessResponse struct {
	Status string                    `json:"status"`
	Checks map[string]ReadinessCheck `json:"checks"`
}

type ReadinessCheck struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}
