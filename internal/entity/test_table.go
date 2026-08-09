package entity

import (
	"time"

	"github.com/rizkiar00/homework/pkg/constant"
)

type TestTable struct {
	TestID    string     `gorm:"column:test_id;primaryKey;type:uuid"`
	DescTest  *string    `gorm:"column:desc_test"`
	IsActive  bool       `gorm:"column:is_active"`
	CreatedBy *string    `gorm:"column:created_by;type:uuid"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedBy *string    `gorm:"column:updated_by;type:uuid"`
	UpdatedAt *time.Time `gorm:"column:updated_at;autoUpdateTime:false"`

	CreatedByUsername *string `gorm:"column:created_by_username;->"`
	UpdatedByUsername *string `gorm:"column:updated_by_username;->"`
}

func (TestTable) TableName() string {
	return constant.TableTestTable
}
