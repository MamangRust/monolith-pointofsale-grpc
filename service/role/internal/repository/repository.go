package repository

import (
	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
)

type Repositories struct {
	RoleCommand RoleCommandRepository
	RoleQuery   RoleQueryRepository
}

func NewRepositories(DB *db.Queries) *Repositories {
	return &Repositories{
		RoleCommand: NewRoleCommandRepository(DB),
		RoleQuery:   NewRoleQueryRepository(DB),
	}
}
