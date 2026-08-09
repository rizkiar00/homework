package model

import "time"

type CreateTestDBRequest struct {
	DescTest *string `json:"desc_test"`
}

type UpdateTestDBRequest struct {
	DescTest *string `json:"desc_test"`
}

type TestDBResponse struct {
	TestID    string     `json:"test_id"`
	DescTest  *string    `json:"desc_test"`
	CreatedBy *string    `json:"created_by"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedBy *string    `json:"updated_by"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type TestDBListRequest struct {
	Page     int
	Limit    int
	OrderBy  string
	OrderDir string
}

type TestDBFindAllOption struct {
	Limit    int
	Offset   int
	OrderBy  string
	OrderDir string
	UserID   string
	Role     string
}

type TestDBListResponse struct {
	Data []TestDBResponse `json:"data"`
	Meta PaginationMeta   `json:"meta"`
}

type PaginationMeta struct {
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	Total      int64  `json:"total"`
	TotalPages int    `json:"total_pages"`
	OrderBy    string `json:"order_by"`
	OrderDir   string `json:"order_dir"`
}
