package impl

import (
	"context"
	"errors"
	"time"

	"github.com/rizkiar00/homework/internal/entity"
	"github.com/rizkiar00/homework/internal/model"
	"github.com/rizkiar00/homework/pkg/constant"
	"gorm.io/gorm"
)

type repository struct {
	db model.Database
}

func New(db model.Database) *repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, data entity.User) (entity.User, error) {
	if r.db == nil {
		return entity.User{}, errors.New(constant.MessageDatabaseNotConfigured)
	}
	if err := r.db.WithContext(ctx).Create(&data).Error; err != nil {
		return entity.User{}, err
	}

	return data, nil
}

func (r *repository) FindByUsername(ctx context.Context, username string) (entity.User, error) {
	if r.db == nil {
		return entity.User{}, errors.New(constant.MessageDatabaseNotConfigured)
	}

	var row entity.User
	if err := r.db.WithContext(ctx).Where(constant.ColumnUsername+" = ?", username).First(&row).Error; err != nil {
		return entity.User{}, err
	}

	return row, nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (entity.User, error) {
	if r.db == nil {
		return entity.User{}, errors.New(constant.MessageDatabaseNotConfigured)
	}

	var row entity.User
	if err := r.db.WithContext(ctx).Where(constant.ColumnEmail+" = ?", email).First(&row).Error; err != nil {
		return entity.User{}, err
	}

	return row, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (entity.User, error) {
	if r.db == nil {
		return entity.User{}, errors.New(constant.MessageDatabaseNotConfigured)
	}

	var row entity.User
	if err := r.db.WithContext(ctx).Where(constant.ColumnUserID+" = ? AND is_active = ?", id, true).First(&row).Error; err != nil {
		return entity.User{}, err
	}

	return row, nil
}

func (r *repository) CreatePendingRegistration(ctx context.Context, user entity.User, verification entity.UserEmailVerification) (entity.User, error) {
	if r.db == nil {
		return entity.User{}, errors.New(constant.MessageDatabaseNotConfigured)
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		return replaceActiveVerification(ctx, tx, verification)
	})
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (r *repository) UpdatePendingRegistration(ctx context.Context, userID string, data entity.User, verification entity.UserEmailVerification) (entity.User, error) {
	if r.db == nil {
		return entity.User{}, errors.New(constant.MessageDatabaseNotConfigured)
	}

	var row entity.User
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		updates := map[string]interface{}{
			"full_name":      data.FullName,
			"username":       data.Username,
			"password_hash":  data.PasswordHash,
			"role":           data.Role,
			"role_id":        data.RoleID,
			"is_active":      false,
			"email_verified": false,
			"updated_at":     time.Now(),
		}

		result := tx.Model(&entity.User{}).Where("user_id = ? AND email_verified = ?", userID, false).Updates(updates)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		if err := replaceActiveVerification(ctx, tx, verification); err != nil {
			return err
		}

		return tx.Where("user_id = ?", userID).First(&row).Error
	})
	if err != nil {
		return entity.User{}, err
	}

	return row, nil
}

func (r *repository) FindActiveVerificationByEmail(ctx context.Context, email string) (entity.UserEmailVerification, error) {
	if r.db == nil {
		return entity.UserEmailVerification{}, errors.New(constant.MessageDatabaseNotConfigured)
	}

	var row entity.UserEmailVerification
	err := r.db.WithContext(ctx).
		Where("email = ? AND used_at IS NULL AND expires_at > ?", email, time.Now()).
		Order("created_at desc").
		First(&row).Error
	if err != nil {
		return entity.UserEmailVerification{}, err
	}

	return row, nil
}

func (r *repository) FindLatestVerificationByEmail(ctx context.Context, email string) (entity.UserEmailVerification, error) {
	if r.db == nil {
		return entity.UserEmailVerification{}, errors.New(constant.MessageDatabaseNotConfigured)
	}

	var row entity.UserEmailVerification
	err := r.db.WithContext(ctx).
		Where("email = ?", email).
		Order("created_at desc").
		First(&row).Error
	if err != nil {
		return entity.UserEmailVerification{}, err
	}

	return row, nil
}

func (r *repository) CountVerificationsByEmailSince(ctx context.Context, email string, since time.Time) (int64, error) {
	if r.db == nil {
		return 0, errors.New(constant.MessageDatabaseNotConfigured)
	}

	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.UserEmailVerification{}).
		Where("email = ? AND created_at >= ?", email, since).
		Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *repository) IncrementVerificationAttempt(ctx context.Context, verificationID string) error {
	if r.db == nil {
		return errors.New(constant.MessageDatabaseNotConfigured)
	}

	return r.db.WithContext(ctx).
		Model(&entity.UserEmailVerification{}).
		Where("verification_id = ?", verificationID).
		UpdateColumn("attempt_count", gorm.Expr("attempt_count + 1")).Error
}

func (r *repository) VerifyEmail(ctx context.Context, userID string, verificationID string) error {
	if r.db == nil {
		return errors.New(constant.MessageDatabaseNotConfigured)
	}

	now := time.Now()
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&entity.UserEmailVerification{}).
			Where("verification_id = ? AND user_id = ? AND used_at IS NULL", verificationID, userID).
			Update("used_at", now)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		result = tx.Model(&entity.User{}).
			Where("user_id = ?", userID).
			Updates(map[string]interface{}{
				"is_active":      true,
				"email_verified": true,
				"verified_at":    now,
				"updated_at":     now,
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return nil
	})
}

func (r *repository) CreatePasswordReset(ctx context.Context, reset entity.UserPasswordReset) error {
	if r.db == nil {
		return errors.New(constant.MessageDatabaseNotConfigured)
	}

	now := time.Now()
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&entity.UserPasswordReset{}).
			Where("user_id = ? AND used_at IS NULL", reset.UserID).
			Update("used_at", now).Error; err != nil {
			return err
		}

		return tx.Create(&reset).Error
	})
}

func (r *repository) FindActivePasswordResetByEmail(ctx context.Context, email string) (entity.UserPasswordReset, error) {
	if r.db == nil {
		return entity.UserPasswordReset{}, errors.New(constant.MessageDatabaseNotConfigured)
	}

	var row entity.UserPasswordReset
	err := r.db.WithContext(ctx).
		Where("email = ? AND used_at IS NULL AND expires_at > ?", email, time.Now()).
		Order("created_at desc").
		First(&row).Error
	if err != nil {
		return entity.UserPasswordReset{}, err
	}

	return row, nil
}

func (r *repository) FindLatestPasswordResetByEmail(ctx context.Context, email string) (entity.UserPasswordReset, error) {
	if r.db == nil {
		return entity.UserPasswordReset{}, errors.New(constant.MessageDatabaseNotConfigured)
	}

	var row entity.UserPasswordReset
	err := r.db.WithContext(ctx).
		Where("email = ?", email).
		Order("created_at desc").
		First(&row).Error
	if err != nil {
		return entity.UserPasswordReset{}, err
	}

	return row, nil
}

func (r *repository) CountPasswordResetsByEmailSince(ctx context.Context, email string, since time.Time) (int64, error) {
	if r.db == nil {
		return 0, errors.New(constant.MessageDatabaseNotConfigured)
	}

	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.UserPasswordReset{}).
		Where("email = ? AND created_at >= ?", email, since).
		Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *repository) IncrementPasswordResetAttempt(ctx context.Context, resetID string) error {
	if r.db == nil {
		return errors.New(constant.MessageDatabaseNotConfigured)
	}

	return r.db.WithContext(ctx).
		Model(&entity.UserPasswordReset{}).
		Where("reset_id = ?", resetID).
		UpdateColumn("attempt_count", gorm.Expr("attempt_count + 1")).Error
}

func (r *repository) ResetPassword(ctx context.Context, userID string, resetID string, passwordHash string) error {
	if r.db == nil {
		return errors.New(constant.MessageDatabaseNotConfigured)
	}

	now := time.Now()
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&entity.UserPasswordReset{}).
			Where("reset_id = ? AND user_id = ? AND used_at IS NULL", resetID, userID).
			Update("used_at", now)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		result = tx.Model(&entity.User{}).
			Where("user_id = ? AND is_active = ? AND email_verified = ?", userID, true, true).
			Updates(map[string]interface{}{
				"password_hash": passwordHash,
				"updated_at":    now,
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return nil
	})
}

func replaceActiveVerification(ctx context.Context, tx *gorm.DB, verification entity.UserEmailVerification) error {
	now := time.Now()
	if err := tx.WithContext(ctx).
		Model(&entity.UserEmailVerification{}).
		Where("user_id = ? AND used_at IS NULL", verification.UserID).
		Update("used_at", now).Error; err != nil {
		return err
	}

	return tx.WithContext(ctx).Create(&verification).Error
}
