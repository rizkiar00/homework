package impl

import (
	"context"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rizkiar00/homework/internal/entity"
	"github.com/rizkiar00/homework/internal/model"
	testDBRepo "github.com/rizkiar00/homework/internal/repository/db/test_db"
	"github.com/rizkiar00/homework/pkg/constant"
	"github.com/rizkiar00/homework/pkg/customerror"
)

type usecase struct {
	repo testDBRepo.Repository
}

func New(repo testDBRepo.Repository) *usecase {
	return &usecase{repo: repo}
}

func (u *usecase) Create(ctx context.Context, createdBy string, request model.CreateTestDBRequest) (model.TestDBResponse, error) {
	row, err := u.repo.Create(ctx, entity.TestTable{
		TestID:    uuid.NewString(),
		DescTest:  request.DescTest,
		IsActive:  true,
		CreatedBy: stringPointer(createdBy),
		CreatedAt: time.Now(),
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
		TestID:   id,
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
		TestID:    row.TestID,
		DescTest:  row.DescTest,
		CreatedBy: row.CreatedBy,
		CreatedAt: row.CreatedAt,
	}
}

func stringPointer(value string) *string {
	return &value
}

func normalizeListRequest(request model.TestDBListRequest) model.TestDBListRequest {
	if request.Page < 1 {
		request.Page = constant.DefaultPage
	}
	if request.Limit < 1 {
		request.Limit = constant.DefaultLimit
	}
	if request.Limit > constant.MaxLimit {
		request.Limit = constant.MaxLimit
	}
	if request.OrderBy == "" {
		request.OrderBy = constant.OrderByTestID
	}
	if request.OrderDir == "" {
		request.OrderDir = constant.OrderDirAsc
	}

	return request
}

func normalizeOrderBy(orderBy string) (string, error) {
	switch orderBy {
	case constant.OrderByTestID:
		return constant.OrderByTestID, nil
	case constant.OrderByDescTest:
		return constant.OrderByDescTest, nil
	case constant.OrderByCreatedAt:
		return constant.OrderByCreatedAt, nil
	default:
		return "", customerror.BadRequest(constant.MessageInvalidOrderBy)
	}
}

func normalizeOrderDir(orderDir string) (string, error) {
	switch strings.ToLower(orderDir) {
	case constant.OrderDirAsc:
		return constant.OrderDirAsc, nil
	case constant.OrderDirDesc:
		return constant.OrderDirDesc, nil
	default:
		return "", customerror.BadRequest(constant.MessageInvalidOrderDir)
	}
}
