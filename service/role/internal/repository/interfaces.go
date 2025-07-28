package repository

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-shared/domain/record"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type RoleQueryRepository interface {
	FindAllRoles(ctx context.Context, req *requests.FindAllRoles) ([]*record.RoleRecord, *int, error)
	FindById(ctx context.Context, role_id int) (*record.RoleRecord, error)
	FindByName(ctx context.Context, name string) (*record.RoleRecord, error)
	FindByUserId(ctx context.Context, user_id int) ([]*record.RoleRecord, error)
	FindByActiveRole(ctx context.Context, req *requests.FindAllRoles) ([]*record.RoleRecord, *int, error)
	FindByTrashedRole(ctx context.Context, req *requests.FindAllRoles) ([]*record.RoleRecord, *int, error)
}

type RoleCommandRepository interface {
	CreateRole(ctx context.Context, request *requests.CreateRoleRequest) (*record.RoleRecord, error)
	UpdateRole(ctx context.Context, request *requests.UpdateRoleRequest) (*record.RoleRecord, error)
	TrashedRole(ctx context.Context, role_id int) (*record.RoleRecord, error)

	RestoreRole(ctx context.Context, role_id int) (*record.RoleRecord, error)
	DeleteRolePermanent(ctx context.Context, role_id int) (bool, error)
	RestoreAllRole(ctx context.Context) (bool, error)
	DeleteAllRolePermanent(ctx context.Context) (bool, error)
}
