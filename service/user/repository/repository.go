package repository

import (
	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
)

type Repositories struct {
	UserCommand UserCommandRepository
	UserQuery   UserQueryRepository
	Role        RoleQueryRepository
}

func NewRepositories(DB *db.Queries) *Repositories {
	return &Repositories{
		UserCommand: NewUserCommandRepository(DB),
		UserQuery:   NewUserQueryRepository(DB),
		Role:        NewRoleRepository(DB),
	}
}
