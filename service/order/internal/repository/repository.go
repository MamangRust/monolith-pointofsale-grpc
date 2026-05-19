package repository

import (
	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-shared/pb"
)

type Repositories struct {
	CashierQuery         CashierQueryRepository
	MerchantQuery        MerchantQueryRepository
	ProductQuery         ProductQueryRepository
	ProductCommand       ProductCommandRepository
	OrderQuery           OrderQueryRepository
	OrderCommand         OrderCommandRepository
	OrderItemQuery       OrderItemQueryRepository
	OrderItemCommand     OrderItemCommandRepository
	OrderStats           OrderStatsRepository
	OrderStatsByMerchant OrderStatByMerchantRepository
}

func NewRepositories(
	DB *db.Queries,
	cashierClient pb.CashierServiceClient,
	merchantClient pb.MerchantServiceClient,
	productClient pb.ProductServiceClient,
	orderItemClient pb.OrderItemServiceClient,
) *Repositories {
	return &Repositories{
		CashierQuery:         NewCashierQueryRepository(cashierClient),
		MerchantQuery:        NewMerchantQueryRepository(merchantClient),
		ProductQuery:         NewProductQueryRepository(productClient),
		ProductCommand:       NewProductCommandRepository(productClient),
		OrderQuery:           NewOrderQueryRepository(DB),
		OrderCommand:         NewOrderCommandRepository(DB),
		OrderItemQuery:       NewOrderItemQueryRepository(orderItemClient),
		OrderItemCommand:     NewOrderItemCommandRepository(DB),
		OrderStats:           NewOrderStatsRepository(DB),
		OrderStatsByMerchant: NewOrderStatsByMerchantRepository(DB),
	}
}
