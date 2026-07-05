package constant

const (
	DefaultPage  = 1
	DefaultLimit = 10
	MaxLimit     = 100
)

const (
	OrderByIDTest   = ColumnIDTest
	OrderByDescTest = ColumnDescTest
	OrderDirAsc     = "asc"
	OrderDirDesc    = "desc"
)

const (
	MessageInvalidOrderBy  = "order_by must be one of: id_test, desc_test"
	MessageInvalidOrderDir = "order_dir must be one of: asc, desc"
)
