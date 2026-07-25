package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rizkiar00/homework/internal/entity"
	"github.com/rizkiar00/homework/pkg/config"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const adminRoleID int64 = 1

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if len(os.Args) < 2 {
		printUsage()
		return nil
	}

	switch os.Args[1] {
	case "create":
		return createAdmin(os.Args[2:])
	default:
		printUsage()
		return fmt.Errorf("unknown admin command: %s", os.Args[1])
	}
}

func createAdmin(args []string) error {
	flags := flag.NewFlagSet("create", flag.ExitOnError)
	username := flags.String("username", "", "admin username")
	email := flags.String("email", "", "admin email")
	fullName := flags.String("full-name", "", "admin full name")
	password := flags.String("password", "", "admin password")
	if err := flags.Parse(args); err != nil {
		return err
	}

	*username = strings.TrimSpace(*username)
	if *username == "" {
		return errors.New("username is required")
	}
	*email = strings.ToLower(strings.TrimSpace(*email))
	if *email == "" {
		*email = *username + "@local.invalid"
	}
	*fullName = strings.TrimSpace(*fullName)
	if *fullName == "" {
		*fullName = *username
	}
	if len(*password) < 8 {
		return errors.New("password minimum length is 8")
	}

	cfg, err := config.New()
	if err != nil {
		return err
	}
	if !cfg.Database.IsConfigured() {
		return errors.New("database config is not complete")
	}

	db, err := gorm.Open(postgres.Open(cfg.Database.GenerateConnectionString()), &gorm.Config{})
	if err != nil {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	now := time.Now()
	data := entity.User{
		UserID:        uuid.NewString(),
		FullName:      *fullName,
		Username:      *username,
		Email:         *email,
		PasswordHash:  string(hash),
		Role:          "admin",
		RoleID:        int64Pointer(adminRoleID),
		IsActive:      true,
		EmailVerified: true,
		VerifiedAt:    &now,
		CreatedAt:     now,
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		var role entity.Role
		if err := tx.Where("role_id = ? AND is_active = ?", adminRoleID, true).First(&role).Error; err != nil {
			return err
		}

		var existing entity.User
		if err := tx.Where("username = ?", *username).First(&existing).Error; err == nil {
			return tx.Model(&entity.User{}).
				Where("user_id = ?", existing.UserID).
				Updates(map[string]interface{}{
					"password_hash":  data.PasswordHash,
					"full_name":      data.FullName,
					"email":          data.Email,
					"role":           data.Role,
					"role_id":        data.RoleID,
					"is_active":      true,
					"email_verified": true,
					"verified_at":    now,
					"updated_at":     now,
				}).Error
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		return tx.Create(&data).Error
	})
	if err != nil {
		return err
	}

	fmt.Printf("admin user %q is ready\n", *username)
	return nil
}

func int64Pointer(value int64) *int64 {
	return &value
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "usage:")
	fmt.Fprintln(os.Stderr, "  go run ./cmd/admin create -username=admin -email=admin@example.com -full-name=Admin -password=your_password")
}
