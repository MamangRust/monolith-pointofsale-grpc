package mencache

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
)

type RoleCommandCache interface {
	DeleteCachedRole(ctx context.Context, id int)
}

type RoleQueryCache interface {
	SetCachedRoles(ctx context.Context, req *requests.FindAllRoles, data []*response.RoleResponse, total *int)
	SetCachedRoleById(ctx context.Context, data *response.RoleResponse)
	SetCachedRoleByUserId(ctx context.Context, userId int, data []*response.RoleResponse)
	SetCachedRoleActive(ctx context.Context, req *requests.FindAllRoles, data []*response.RoleResponseDeleteAt, total *int)
	SetCachedRoleTrashed(ctx context.Context, req *requests.FindAllRoles, data []*response.RoleResponseDeleteAt, total *int)

	GetCachedRoles(ctx context.Context, req *requests.FindAllRoles) ([]*response.RoleResponse, *int, bool)
	GetCachedRoleByUserId(ctx context.Context, userId int) ([]*response.RoleResponse, bool)
	GetCachedRoleById(ctx context.Context, id int) (*response.RoleResponse, bool)
	GetCachedRoleActive(ctx context.Context, req *requests.FindAllRoles) ([]*response.RoleResponseDeleteAt, *int, bool)
	GetCachedRoleTrashed(ctx context.Context, req *requests.FindAllRoles) ([]*response.RoleResponseDeleteAt, *int, bool)
}
