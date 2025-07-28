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

type roleCommandRepository struct {
	db      *db.Queries
	mapping recordmapper.RoleRecordMapping
}

func NewRoleCommandRepository(db *db.Queries, mapping recordmapper.RoleRecordMapping) *roleCommandRepository {
	return &roleCommandRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *roleCommandRepository) CreateRole(ctx context.Context, req *requests.CreateRoleRequest) (*record.RoleRecord, error) {
	res, err := r.db.CreateRole(ctx, req.Name)

	if err != nil {
		return nil, fmt.Errorf("failed to create role: invalid name '%s' or duplicate role", req.Name)
	}

	return r.mapping.ToRoleRecord(res), nil
}

func (r *roleCommandRepository) UpdateRole(ctx context.Context, req *requests.UpdateRoleRequest) (*record.RoleRecord, error) {
	res, err := r.db.UpdateRole(ctx, db.UpdateRoleParams{
		RoleID:   int32(*req.ID),
		RoleName: req.Name,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to update role ID %d: role not found or invalid data", req.ID)
	}

	return r.mapping.ToRoleRecord(res), nil
}

func (r *roleCommandRepository) TrashedRole(ctx context.Context, id int) (*record.RoleRecord, error) {
	res, err := r.db.TrashRole(ctx, int32(id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("role ID %d not found or already trashed", id)
		}
		return nil, fmt.Errorf("failed to trash role ID %d: %w", id, err)
	}

	return r.mapping.ToRoleRecord(res), nil
}

func (r *roleCommandRepository) RestoreRole(ctx context.Context, id int) (*record.RoleRecord, error) {
	res, err := r.db.RestoreRole(ctx, int32(id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("role ID %d not found in trash", id)
		}
		return nil, fmt.Errorf("failed to restore role ID %d: %w", id, err)
	}

	return r.mapping.ToRoleRecord(res), nil
}

func (r *roleCommandRepository) DeleteRolePermanent(ctx context.Context, role_id int) (bool, error) {
	err := r.db.DeletePermanentRole(ctx, int32(role_id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("role ID %d not found or already deleted", role_id)
		}
		return false, fmt.Errorf("failed to permanently delete role ID %d: %w", role_id, err)
	}

	return true, nil
}

func (r *roleCommandRepository) RestoreAllRole(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllRoles(ctx)

	if err != nil {
		return false, fmt.Errorf("no trashed roles available to restore")
	}

	return true, nil
}

func (r *roleCommandRepository) DeleteAllRolePermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentRoles(ctx)

	if err != nil {
		return false, fmt.Errorf("cannot permanently delete all roles: operation disabled for system protection")
	}

	return true, nil
}
