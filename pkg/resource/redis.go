package resource

import (
	"github.com/redis/go-redis/v9"
	"github.com/rizkiar00/homework/internal/model"
	"github.com/rizkiar00/homework/pkg/config"
)

func NewRedis(cfg config.Config) model.Redis {
	if !cfg.Redis.Enabled {
		return nil
	}

	return redis.NewClient(&redis.Options{
		Addr:        cfg.Redis.Address(),
		Username:    cfg.Redis.Username,
		Password:    cfg.Redis.Password,
		DB:          cfg.Redis.DB,
		DialTimeout: cfg.Redis.DialTimeout(),
	})
}
