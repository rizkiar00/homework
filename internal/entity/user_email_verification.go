package entity

import (
	"time"

	"github.com/rizkiar00/homework/pkg/constant"
)

type UserEmailVerification struct {
	VerificationID string     `gorm:"column:verification_id;primaryKey;type:uuid"`
	UserID         string     `gorm:"column:user_id;type:uuid"`
	Email          string     `gorm:"column:email"`
	CodeHash       string     `gorm:"column:code_hash"`
	ExpiresAt      time.Time  `gorm:"column:expires_at"`
	AttemptCount   int        `gorm:"column:attempt_count"`
	UsedAt         *time.Time `gorm:"column:used_at"`
	CreatedAt      time.Time  `gorm:"column:created_at"`
}

func (UserEmailVerification) TableName() string {
	return constant.TableEmailVerifications
}
