package repository

import (
	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-point-of-sale-shared/mapper/record"
)

type Repositories struct {
	UserQuery              UserQueryRepository
	MerchantQuery          MerchantQueryRepository
	CashierQuery           CashierQueryRepository
	CashierCommand         CashierCommandRepository
	CashierStats           CashierStatsRepository
	CashierStatsByMerchant CashierStatByMerchantRepository
	CashierStatsById       CashierStatByIdRepository
}

func NewRepositories(DB *db.Queries) *Repositories {
	mapperUser := recordmapper.NewUserRecordMapper()
	mapperMerchant := recordmapper.NewMerchantRecordMapper()
	mapperCashier := recordmapper.NewCashierRecordMapper()

	return &Repositories{
		UserQuery:              NewUserQueryRepository(DB, mapperUser),
		MerchantQuery:          NewMerchantQueryRepository(DB, mapperMerchant),
		CashierQuery:           NewCashierQueryRepository(DB, mapperCashier),
		CashierCommand:         NewCashierCommandRepository(DB, mapperCashier),
		CashierStats:           NewCashierStatsRepository(DB, mapperCashier),
		CashierStatsByMerchant: NewCashierStatsByMerchantRepository(DB, mapperCashier),
		CashierStatsById:       NewCashierStatsByIdRepository(DB, mapperCashier),
	}
}
