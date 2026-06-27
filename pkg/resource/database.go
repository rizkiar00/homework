package resource

import (
	"github.com/rizkiar00/homework/internal/model"
	"github.com/rizkiar00/homework/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(cfg config.Config) (model.Database, error) {
	if !cfg.Database.IsConfigured() {
		return nil, nil
	}

	return gorm.Open(postgres.Open(cfg.Database.GenerateConnectionString()), &gorm.Config{})
}
