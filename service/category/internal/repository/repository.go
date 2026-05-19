package repository

import (
	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
)

type Repositories struct {
	CategoryQuery           CategoryQueryRepository
	CategoryCommand         CategoryCommandRepository
	CategoryStats           CategoryStatsRepository
	CategoryStatsById       CategoryStatsByIdRepository
	CategoryStatsByMerchant CategoryStatsByMerchantRepository
}

func NewRepositories(DB *db.Queries) *Repositories {
	return &Repositories{
		CategoryQuery:           NewCategoryQueryRepository(DB),
		CategoryCommand:         NewCategoryCommandRepository(DB),
		CategoryStats:           NewCategoryStatsRepository(DB),
		CategoryStatsById:       NewCategoryStatsByIdRepository(DB),
		CategoryStatsByMerchant: NewCategoryStatsByMerchantRepository(DB),
	}
}
