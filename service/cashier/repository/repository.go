package repository

import (
	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-shared/pb"
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

func NewRepositories(
	DB *db.Queries,
	userClient pb.UserServiceClient,
	merchantClient pb.MerchantServiceClient,
) *Repositories {
	return &Repositories{
		UserQuery:              NewUserQueryRepository(userClient),
		MerchantQuery:          NewMerchantQueryRepository(merchantClient),
		CashierQuery:           NewCashierQueryRepository(DB),
		CashierCommand:         NewCashierCommandRepository(DB),
		CashierStats:           NewCashierStatsRepository(DB),
		CashierStatsByMerchant: NewCashierStatsByMerchantRepository(DB),
		CashierStatsById:       NewCashierStatsByIdRepository(DB),
	}
}
