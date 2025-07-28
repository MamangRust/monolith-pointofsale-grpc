package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/record"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	recordmapper "github.com/MamangRust/monolith-point-of-sale-shared/mapper/record"
)

type roleQueryRepository struct {
	db      *db.Queries
	mapping recordmapper.RoleRecordMapping
}

func NewRoleQueryRepository(db *db.Queries, mapping recordmapper.RoleRecordMapping) *roleQueryRepository {
	return &roleQueryRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *roleQueryRepository) FindAllRoles(ctx context.Context, req *requests.FindAllRoles) ([]*record.RoleRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetRolesParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetRoles(ctx, reqDb)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch roles: invalid pagination (page %d, size %d) or search query '%s'", req.Page, req.PageSize, req.Search)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToRolesRecordAll(res), &totalCount, nil
}

func (r *roleQueryRepository) FindById(ctx context.Context, id int) (*record.RoleRecord, error) {
	res, err := r.db.GetRole(ctx, int32(id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("role not found with ID: %d", id)
		}
		return nil, fmt.Errorf("failed to retrieve role with ID %d: %w", id, err)
	}

	return r.mapping.ToRoleRecord(res), nil
}

func (r *roleQueryRepository) FindByName(ctx context.Context, name string) (*record.RoleRecord, error) {
	res, err := r.db.GetRoleByName(ctx, name)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("role not found with name: '%s'", name)
		}
		return nil, fmt.Errorf("failed to retrieve role with name '%s': %w", name, err)
	}

	return r.mapping.ToRoleRecord(res), nil
}

func (r *roleQueryRepository) FindByUserId(ctx context.Context, user_id int) ([]*record.RoleRecord, error) {
	res, err := r.db.GetUserRoles(ctx, int32(user_id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no roles found for user ID: %d", user_id)
		}
		return nil, fmt.Errorf("failed to retrieve roles for user ID %d: %w", user_id, err)
	}

	return r.mapping.ToRolesRecord(res), nil
}

func (r *roleQueryRepository) FindByActiveRole(ctx context.Context, req *requests.FindAllRoles) ([]*record.RoleRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveRolesParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveRoles(ctx, reqDb)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch active roles: invalid parameters (page %d, size %d, search '%s')", req.Page, req.PageSize, req.Search)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToRolesRecordActive(res), &totalCount, nil
}

func (r *roleQueryRepository) FindByTrashedRole(ctx context.Context, req *requests.FindAllRoles) ([]*record.RoleRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedRolesParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedRoles(ctx, reqDb)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch trashed roles: invalid parameters (page %d, size %d, search '%s')", req.Page, req.PageSize, req.Search)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToRolesRecordTrashed(res), &totalCount, nil
}
