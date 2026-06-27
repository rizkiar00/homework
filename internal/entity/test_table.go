package entity

type TestTable struct {
	IDTest   string  `gorm:"column:id_test;primaryKey;type:uuid"`
	DescTest *string `gorm:"column:desc_test"`
}

func (TestTable) TableName() string {
	return "public.test_table"
}
