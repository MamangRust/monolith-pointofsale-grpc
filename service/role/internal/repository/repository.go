package repository

import (
	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-point-of-sale-shared/mapper/record"
)

type Repositories struct {
	RoleCommand RoleCommandRepository
	RoleQuery   RoleQueryRepository
}

func NewRepositories(DB *db.Queries) *Repositories {
	roleMapper := recordmapper.NewRoleRecordMapper()

	return &Repositories{
		RoleCommand: NewRoleCommandRepository(DB, roleMapper),
		RoleQuery:   NewRoleQueryRepository(DB, roleMapper),
	}
}
