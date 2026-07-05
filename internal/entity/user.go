package entity

import "time"

import "github.com/rizkiar00/homework/pkg/constant"

type User struct {
	IDUser       string    `gorm:"column:id_user;primaryKey;type:uuid"`
	Username     string    `gorm:"column:username"`
	PasswordHash string    `gorm:"column:password_hash"`
	Role         string    `gorm:"column:role"`
	CreatedAt    time.Time `gorm:"column:created_at"`
}

func (User) TableName() string {
	return constant.TableUsers
}
