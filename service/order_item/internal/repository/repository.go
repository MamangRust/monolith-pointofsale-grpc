package repository

import (
	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-point-of-sale-shared/mapper/record"
)

type Repositories struct {
	OrderItemQuery OrderItemQueryRepository
}

func NewRepositories(DB *db.Queries) *Repositories {
	mapper := recordmapper.NewOrderItemRecordMapper()

	return &Repositories{
		OrderItemQuery: NewOrderItemQueryRepository(DB, mapper),
	}
}
