package health

import "context"

type Repository interface {
	Check(ctx context.Context) error
}
