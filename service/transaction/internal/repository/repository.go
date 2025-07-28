package repository

import (
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

func NewRepositories(DB *db.Queries) *Repositories {
	mapperOrderItem := recordmapper.NewOrderItemRecordMapper()
	mapperOrder := recordmapper.NewOrderRecordMapper()
	mapperTransaction := recordmapper.NewTransactionRecordMapper()
	mapperMerchant := recordmapper.NewMerchantRecordMapper()
	mapperCashier := recordmapper.NewCashierRecordMapper()

	return &Repositories{
		CashierQuery:                 NewCashierQueryRepository(DB, mapperCashier),
		MerchantQuery:                NewMerchantQueryRepository(DB, mapperMerchant),
		OrderQuery:                   NewOrderQueryRepository(DB, mapperOrder),
		OrderItemQuery:               NewOrderItemQueryRepository(DB, mapperOrderItem),
		TransactionCommandRepository: NewTransactionCommandRepository(DB, mapperTransaction),
		TransactionQueryRepository:   NewTransactionQueryRepository(DB, mapperTransaction),
		TransactionStatsRepository:   NewTransactionStatsRepository(DB, mapperTransaction),
		TransactionStatsByMerchant:   NewTransactionStatsByMerchantRepository(DB, mapperTransaction),
	}
}
