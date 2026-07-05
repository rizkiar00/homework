package entity

import "github.com/rizkiar00/homework/pkg/constant"

type TestTable struct {
	IDTest   string  `gorm:"column:id_test;primaryKey;type:uuid"`
	DescTest *string `gorm:"column:desc_test"`
}

func (TestTable) TableName() string {
	return constant.TableTestTable
}
