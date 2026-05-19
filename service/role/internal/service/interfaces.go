package service

import (
	"context"

	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type RoleQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllRoles) ([]*db.GetRolesRow, *int, error)
	FindByActiveRole(ctx context.Context, req *requests.FindAllRoles) ([]*db.GetActiveRolesRow, *int, error)
	FindByTrashedRole(ctx context.Context, req *requests.FindAllRoles) ([]*db.GetTrashedRolesRow, *int, error)
	FindById(ctx context.Context, roleID int) (*db.Role, error)
	FindByUserId(ctx context.Context, id int) ([]*db.Role, error)
}

type RoleCommandService interface {
	CreateRole(ctx context.Context, request *requests.CreateRoleRequest) (*db.Role, error)
	UpdateRole(ctx context.Context, request *requests.UpdateRoleRequest) (*db.Role, error)
	TrashedRole(ctx context.Context, roleID int) (*db.Role, error)
	RestoreRole(ctx context.Context, roleID int) (*db.Role, error)
	DeleteRolePermanent(ctx context.Context, roleID int) (bool, error)
	RestoreAllRole(ctx context.Context) (bool, error)
	DeleteAllRolePermanent(ctx context.Context) (bool, error)
}
