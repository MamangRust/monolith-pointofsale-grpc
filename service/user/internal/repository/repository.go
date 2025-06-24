package repository

import (
	"context"

	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-point-of-sale-shared/mapper/record"
)

type Repositories struct {
	UserCommand UserCommandRepository
	UserQuery   UserQueryRepository
	Role        RoleQueryRepository
}

type Deps struct {
	DB           *db.Queries
	Ctx          context.Context
	MapperRecord *recordmapper.RecordMapper
}

func NewRepositories(deps *Deps) *Repositories {
	return &Repositories{
		UserCommand: NewUserCommandRepository(deps.DB, deps.Ctx, deps.MapperRecord.UserRecordMapper),
		UserQuery:   NewUserQueryRepository(deps.DB, deps.Ctx, deps.MapperRecord.UserRecordMapper),
		Role:        NewRoleRepository(deps.DB, deps.Ctx, deps.MapperRecord.RoleRecordMapper),
	}
}
