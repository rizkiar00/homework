package repository

// Repository defines persistence methods used by usecases.
type Repository interface {
	Save(value string) error
}

// InMemoryRepo is a trivial in-memory implementation used for examples.
type InMemoryRepo struct{}

func NewInMemoryRepo() *InMemoryRepo { return &InMemoryRepo{} }

func (r *InMemoryRepo) Save(value string) error {
	// no-op for example
	return nil
}
