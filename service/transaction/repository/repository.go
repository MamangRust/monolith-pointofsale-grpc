package repository

import (
	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-shared/pb"
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

func NewRepositories(
	DB *db.Queries,
	cashierClient pb.CashierServiceClient,
	merchantClient pb.MerchantServiceClient,
	orderClient pb.OrderServiceClient,
	orderItemClient pb.OrderItemServiceClient,
) *Repositories {
	return &Repositories{
		CashierQuery:                 NewCashierQueryRepository(cashierClient),
		MerchantQuery:                NewMerchantQueryRepository(merchantClient),
		OrderQuery:                   NewOrderQueryRepository(orderClient),
		OrderItemQuery:               NewOrderItemQueryRepository(orderItemClient),
		TransactionCommandRepository: NewTransactionCommandRepository(DB),
		TransactionQueryRepository:   NewTransactionQueryRepository(DB),
		TransactionStatsRepository:   NewTransactionStatsRepository(DB),
		TransactionStatsByMerchant:   NewTransactionStatsByMerchantRepository(DB),
	}
}
