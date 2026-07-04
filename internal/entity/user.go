package entity

import "time"

type User struct {
	IDUser       string    `gorm:"column:id_user;primaryKey;type:uuid"`
	Username     string    `gorm:"column:username"`
	PasswordHash string    `gorm:"column:password_hash"`
	Role         string    `gorm:"column:role"`
	CreatedAt    time.Time `gorm:"column:created_at"`
}

func (User) TableName() string {
	return "public.users"
}
