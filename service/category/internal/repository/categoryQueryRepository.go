package repository

import (
	"context"

	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/record"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/category_errors"
	recordmapper "github.com/MamangRust/monolith-point-of-sale-shared/mapper/record"
)

type categoryQueryRepository struct {
	db      *db.Queries
	mapping recordmapper.CategoryRecordMapper
}

func NewCategoryQueryRepository(db *db.Queries, mapping recordmapper.CategoryRecordMapper) *categoryQueryRepository {
	return &categoryQueryRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *categoryQueryRepository) FindAllCategory(ctx context.Context, req *requests.FindAllCategory) ([]*record.CategoriesRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetCategoriesParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetCategories(ctx, reqDb)

	if err != nil {
		return nil, nil, category_errors.ErrFindAllCategory
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToCategoriesRecordPagination(res), &totalCount, nil
}

func (r *categoryQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllCategory) ([]*record.CategoriesRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetCategoriesActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetCategoriesActive(ctx, reqDb)

	if err != nil {
		return nil, nil, category_errors.ErrFindByActive
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToCategoriesRecordActivePagination(res), &totalCount, nil
}

func (r *categoryQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllCategory) ([]*record.CategoriesRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetCategoriesTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetCategoriesTrashed(ctx, reqDb)

	if err != nil {
		return nil, nil, category_errors.ErrFindByTrashed
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToCategoriesRecordTrashedPagination(res), &totalCount, nil
}

func (r *categoryQueryRepository) FindById(ctx context.Context, category_id int) (*record.CategoriesRecord, error) {
	res, err := r.db.GetCategoryByID(ctx, int32(category_id))

	if err != nil {
		return nil, category_errors.ErrFindById
	}

	return r.mapping.ToCategoryRecord(res), nil
}

func (r *categoryQueryRepository) FindByIdTrashed(ctx context.Context, category_id int) (*record.CategoriesRecord, error) {
	res, err := r.db.GetCategoryByIDTrashed(ctx, int32(category_id))

	if err != nil {
		return nil, category_errors.ErrFindByTrashed
	}

	return r.mapping.ToCategoryRecord(res), nil
}

func (r *categoryQueryRepository) FindByName(ctx context.Context, name string) (*record.CategoriesRecord, error) {
	res, err := r.db.GetCategoryByName(ctx, name)

	if err != nil {
		return nil, category_errors.ErrFindByName
	}

	return r.mapping.ToCategoryRecord(res), nil
}

func (r *categoryQueryRepository) FindByNameAndId(ctx context.Context, req *requests.CategoryNameAndId) (*record.CategoriesRecord, error) {
	res, err := r.db.GetCategoryByNameAndId(ctx, db.GetCategoryByNameAndIdParams{
		Name:       req.Name,
		CategoryID: int32(req.CategoryID),
	})

	if err != nil {
		return nil, category_errors.ErrFindByNameAndId
	}

	return r.mapping.ToCategoryRecord(res), nil
}
