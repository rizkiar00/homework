package entity

import (
	"time"

	"github.com/rizkiar00/homework/pkg/constant"
)

type Role struct {
	RoleID    int64      `gorm:"column:role_id;primaryKey;autoIncrement"`
	RoleDesc  string     `gorm:"column:role_desc"`
	IsActive  bool       `gorm:"column:is_active"`
	CreatedBy *string    `gorm:"column:created_by;type:uuid"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedBy *string    `gorm:"column:updated_by;type:uuid"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
}

func (Role) TableName() string {
	return constant.TableRoles
}
