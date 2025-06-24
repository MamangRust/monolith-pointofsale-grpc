package repository

import (
	"context"

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

type Deps struct {
	DB  *db.Queries
	Ctx context.Context
}

func NewRepositories(deps *Deps) *Repositories {
	mapperUser := recordmapper.NewUserRecordMapper()
	mapperMerchant := recordmapper.NewMerchantRecordMapper()
	mapperCashier := recordmapper.NewCashierRecordMapper()

	return &Repositories{
		UserQuery:              NewUserQueryRepository(deps.DB, deps.Ctx, mapperUser),
		MerchantQuery:          NewMerchantQueryRepository(deps.DB, deps.Ctx, mapperMerchant),
		CashierQuery:           NewCashierQueryRepository(deps.DB, deps.Ctx, mapperCashier),
		CashierCommand:         NewCashierCommandRepository(deps.DB, deps.Ctx, mapperCashier),
		CashierStats:           NewCashierStatsRepository(deps.DB, deps.Ctx, mapperCashier),
		CashierStatsByMerchant: NewCashierStatsByMerchantRepository(deps.DB, deps.Ctx, mapperCashier),
		CashierStatsById:       NewCashierStatsByIdRepository(deps.DB, deps.Ctx, mapperCashier),
	}
}
