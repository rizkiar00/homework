package impl

import (
	"context"

	"github.com/google/uuid"
	"github.com/rizkiar00/homework/internal/entity"
	"github.com/rizkiar00/homework/internal/model"
	testDBRepo "github.com/rizkiar00/homework/internal/repository/db/test_db"
)

type usecase struct {
	repo testDBRepo.Repository
}

func New(repo testDBRepo.Repository) *usecase {
	return &usecase{repo: repo}
}

func (u *usecase) Create(ctx context.Context, request model.CreateTestDBRequest) (model.TestDBResponse, error) {
	row, err := u.repo.Create(ctx, entity.TestTable{
		IDTest:   uuid.NewString(),
		DescTest: request.DescTest,
	})
	if err != nil {
		return model.TestDBResponse{}, err
	}

	return toResponse(row), nil
}

func (u *usecase) FindAll(ctx context.Context) ([]model.TestDBResponse, error) {
	rows, err := u.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]model.TestDBResponse, 0, len(rows))
	for _, row := range rows {
		responses = append(responses, toResponse(row))
	}

	return responses, nil
}

func (u *usecase) FindByID(ctx context.Context, id string) (model.TestDBResponse, error) {
	row, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return model.TestDBResponse{}, err
	}

	return toResponse(row), nil
}

func (u *usecase) Update(ctx context.Context, id string, request model.UpdateTestDBRequest) (model.TestDBResponse, error) {
	row, err := u.repo.Update(ctx, entity.TestTable{
		IDTest:   id,
		DescTest: request.DescTest,
	})
	if err != nil {
		return model.TestDBResponse{}, err
	}

	return toResponse(row), nil
}

func (u *usecase) Delete(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}

func toResponse(row entity.TestTable) model.TestDBResponse {
	return model.TestDBResponse{
		IDTest:   row.IDTest,
		DescTest: row.DescTest,
	}
}
