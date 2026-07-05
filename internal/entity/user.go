package entity

import "time"

import "github.com/rizkiar00/homework/pkg/constant"

type User struct {
	UserID       string     `gorm:"column:user_id;primaryKey;type:uuid"`
	Username     string     `gorm:"column:username"`
	PasswordHash string     `gorm:"column:password_hash"`
	Role         string     `gorm:"column:role"`
	RoleID       *int64     `gorm:"column:role_id"`
	IsActive     bool       `gorm:"column:is_active"`
	CreatedBy    *string    `gorm:"column:created_by;type:uuid"`
	CreatedAt    time.Time  `gorm:"column:created_at"`
	UpdatedBy    *string    `gorm:"column:updated_by;type:uuid"`
	UpdatedAt    *time.Time `gorm:"column:updated_at"`
}

func (User) TableName() string {
	return constant.TableUsers
}
