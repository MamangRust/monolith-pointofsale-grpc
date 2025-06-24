package repository

import (
	"context"

	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-point-of-sale-shared/mapper/record"
)

type Repositories struct {
	CashierQuery                 CashierQueryRepository
	MerchantQuery                MerchantQueryRepository
	OrderQuery                   OrderQueryRepository
	OrderItemQuery               OrderItemQueryRepository
	TransactionCommandRepository TransactionCommandRepository
	TransactionQueryRepository   TransactionQueryRepository
	TransactionStatsRepository   TransactionStatsRepository
	TransactionStatsByMerchant   TransactionStatsByMerchantRepository
}

type Deps struct {
	DB  *db.Queries
	Ctx context.Context
}

func NewRepositories(deps *Deps) *Repositories {
	mapperOrderItem := recordmapper.NewOrderItemRecordMapper()
	mapperOrder := recordmapper.NewOrderRecordMapper()
	mapperTransaction := recordmapper.NewTransactionRecordMapper()
	mapperMerchant := recordmapper.NewMerchantRecordMapper()
	mapperCashier := recordmapper.NewCashierRecordMapper()

	return &Repositories{
		CashierQuery:                 NewCashierQueryRepository(deps.DB, deps.Ctx, mapperCashier),
		MerchantQuery:                NewMerchantQueryRepository(deps.DB, deps.Ctx, mapperMerchant),
		OrderQuery:                   NewOrderQueryRepository(deps.DB, deps.Ctx, mapperOrder),
		OrderItemQuery:               NewOrderItemQueryRepository(deps.DB, deps.Ctx, mapperOrderItem),
		TransactionCommandRepository: NewTransactionCommandRepository(deps.DB, deps.Ctx, mapperTransaction),
		TransactionQueryRepository:   NewTransactionQueryRepository(deps.DB, deps.Ctx, mapperTransaction),
		TransactionStatsRepository:   NewTransactionStatsRepository(deps.DB, deps.Ctx, mapperTransaction),
		TransactionStatsByMerchant:   NewTransactionStatsByMerchantRepository(deps.DB, deps.Ctx, mapperTransaction),
	}
}
