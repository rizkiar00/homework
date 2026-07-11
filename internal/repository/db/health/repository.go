package health

import "context"

type Repository interface {
	Check(ctx context.Context) error
	CheckRedis(ctx context.Context) error
}
