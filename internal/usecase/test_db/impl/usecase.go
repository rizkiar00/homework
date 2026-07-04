package impl

import (
	"context"
	"fmt"
	"math"
	"strings"

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

func (u *usecase) FindAll(ctx context.Context, request model.TestDBListRequest) (model.TestDBListResponse, error) {
	request = normalizeListRequest(request)
	orderBy, err := normalizeOrderBy(request.OrderBy)
	if err != nil {
		return model.TestDBListResponse{}, err
	}

	orderDir, err := normalizeOrderDir(request.OrderDir)
	if err != nil {
		return model.TestDBListResponse{}, err
	}

	rows, total, err := u.repo.FindAll(ctx, model.TestDBFindAllOption{
		Limit:    request.Limit,
		Offset:   (request.Page - 1) * request.Limit,
		OrderBy:  orderBy,
		OrderDir: orderDir,
	})
	if err != nil {
		return model.TestDBListResponse{}, err
	}

	responses := make([]model.TestDBResponse, 0, len(rows))
	for _, row := range rows {
		responses = append(responses, toResponse(row))
	}

	return model.TestDBListResponse{
		Data: responses,
		Meta: model.PaginationMeta{
			Page:       request.Page,
			Limit:      request.Limit,
			Total:      total,
			TotalPages: int(math.Ceil(float64(total) / float64(request.Limit))),
			OrderBy:    request.OrderBy,
			OrderDir:   orderDir,
		},
	}, nil
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

func normalizeListRequest(request model.TestDBListRequest) model.TestDBListRequest {
	if request.Page < 1 {
		request.Page = 1
	}
	if request.Limit < 1 {
		request.Limit = 10
	}
	if request.Limit > 100 {
		request.Limit = 100
	}
	if request.OrderBy == "" {
		request.OrderBy = "id_test"
	}
	if request.OrderDir == "" {
		request.OrderDir = "asc"
	}

	return request
}

func normalizeOrderBy(orderBy string) (string, error) {
	switch orderBy {
	case "id_test":
		return "id_test", nil
	case "desc_test":
		return "desc_test", nil
	default:
		return "", fmt.Errorf("order_by must be one of: id_test, desc_test")
	}
}

func normalizeOrderDir(orderDir string) (string, error) {
	switch strings.ToLower(orderDir) {
	case "asc":
		return "asc", nil
	case "desc":
		return "desc", nil
	default:
		return "", fmt.Errorf("order_dir must be one of: asc, desc")
	}
}
