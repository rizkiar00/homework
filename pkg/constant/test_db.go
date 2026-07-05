package constant

const (
	DefaultPage  = 1
	DefaultLimit = 10
	MaxLimit     = 100
)

const (
	OrderByTestID   = ColumnTestID
	OrderByDescTest = ColumnDescTest
	OrderDirAsc     = "asc"
	OrderDirDesc    = "desc"
)

const (
	MessageInvalidOrderBy  = "order_by must be one of: test_id, desc_test"
	MessageInvalidOrderDir = "order_dir must be one of: asc, desc"
)
