package model

type CreateTestDBRequest struct {
	DescTest *string `json:"desc_test"`
}

type UpdateTestDBRequest struct {
	DescTest *string `json:"desc_test"`
}

type TestDBResponse struct {
	IDTest   string  `json:"id_test"`
	DescTest *string `json:"desc_test"`
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
