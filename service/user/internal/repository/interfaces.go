package repository

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-shared/domain/record"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type UserQueryRepository interface {
	FindAllUsers(ctx context.Context, req *requests.FindAllUsers) ([]*record.UserRecord, *int, error)
	FindById(ctx context.Context, user_id int) (*record.UserRecord, error)
	FindByEmail(ctx context.Context, email string) (*record.UserRecord, error)
	FindByActive(ctx context.Context, req *requests.FindAllUsers) ([]*record.UserRecord, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*record.UserRecord, *int, error)
}

type UserCommandRepository interface {
	CreateUser(ctx context.Context, request *requests.CreateUserRequest) (*record.UserRecord, error)
	UpdateUser(ctx context.Context, request *requests.UpdateUserRequest) (*record.UserRecord, error)
	TrashedUser(ctx context.Context, user_id int) (*record.UserRecord, error)
	RestoreUser(ctx context.Context, user_id int) (*record.UserRecord, error)
	DeleteUserPermanent(ctx context.Context, user_id int) (bool, error)
	RestoreAllUser(ctx context.Context) (bool, error)
	DeleteAllUserPermanent(ctx context.Context) (bool, error)
}

type RoleQueryRepository interface {
	FindByName(ctx context.Context, name string) (*record.RoleRecord, error)
}
