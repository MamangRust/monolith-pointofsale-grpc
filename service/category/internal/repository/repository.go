package repository

import (
	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-point-of-sale-shared/mapper/record"
)

type Repositories struct {
	CategoryQuery           CategoryQueryRepository
	CategoryCommand         CategoryCommandRepository
	CategoryStats           CategoryStatsRepository
	CategoryStatsById       CategoryStatsByIdRepository
	CategoryStatsByMerchant CategoryStatsByMerchantRepository
}

func NewRepositories(DB *db.Queries) *Repositories {
	categoryMapper := recordmapper.NewCategoryRecordMapper()

	return &Repositories{
		CategoryQuery:           NewCategoryQueryRepository(DB, categoryMapper),
		CategoryCommand:         NewCategoryCommandRepository(DB, categoryMapper),
		CategoryStats:           NewCategoryStatsRepository(DB, categoryMapper),
		CategoryStatsById:       NewCategoryStatsByIdRepository(DB, categoryMapper),
		CategoryStatsByMerchant: NewCategoryStatsByMerchantRepository(DB, categoryMapper),
	}
}
