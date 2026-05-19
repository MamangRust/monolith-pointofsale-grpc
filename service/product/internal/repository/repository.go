package repository

import (
	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
)

type Repositories struct {
	ProductQuery   ProductQueryRepository
	ProductCommand ProductCommandRepository
	CategoryQuery  CategoryQueryRepository
	MerchantQuery  MerchantQueryRepository
}

func NewRepositories(DB *db.Queries) *Repositories {
	return &Repositories{
		ProductQuery:   NewProductQueryRepository(DB),
		ProductCommand: NewProductCommandRepository(DB),
		CategoryQuery:  NewCategoryQueryRepository(DB),
		MerchantQuery:  NewMerchantQueryRepository(DB),
	}
}
