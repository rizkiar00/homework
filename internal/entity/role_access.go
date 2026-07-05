package entity

import (
	"time"

	"github.com/rizkiar00/homework/pkg/constant"
)

type RoleAccess struct {
	RoleAccessID string     `gorm:"column:role_access_id;primaryKey;type:uuid"`
	RoleID       int64      `gorm:"column:role_id"`
	ActionID     int64      `gorm:"column:action_id"`
	CreatedBy    *string    `gorm:"column:created_by;type:uuid"`
	CreatedAt    time.Time  `gorm:"column:created_at"`
	UpdatedBy    *string    `gorm:"column:updated_by;type:uuid"`
	UpdatedAt    *time.Time `gorm:"column:updated_at"`
}

func (RoleAccess) TableName() string {
	return constant.TableRoleAccesses
}
