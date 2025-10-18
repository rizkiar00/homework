package usecase

import "github.com/rizkiar00/homework/internal/repository"

// EchoUsecase provides business logic for echoing values.
type EchoUsecase struct{
	repo repository.Repository
}

func NewEchoUsecase(r repository.Repository) *EchoUsecase {
	return &EchoUsecase{repo: r}
}

func (e *EchoUsecase) Echo(value string) (string, error) {
	// business rules could go here; persist via repo if needed
	if err := e.repo.Save(value); err != nil {
		return "", err
	}
	return value, nil
}
