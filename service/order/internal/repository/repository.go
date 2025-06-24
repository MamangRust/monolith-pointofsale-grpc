package repository

import (
	"context"

	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	recordmapper "github.com/MamangRust/monolith-point-of-sale-shared/mapper/record"
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

type Deps struct {
	DB  *db.Queries
	Ctx context.Context
}

func NewRepositories(deps *Deps) *Repositories {
	mapperCashier := recordmapper.NewCashierRecordMapper()
	mapperMerchant := recordmapper.NewMerchantRecordMapper()
	mapperProduct := recordmapper.NewProductRecordMapper()
	mapperOrder := recordmapper.NewOrderRecordMapper()
	mapperOrderItem := recordmapper.NewOrderItemRecordMapper()

	return &Repositories{
		CashierQuery:         NewCashierQueryRepository(deps.DB, deps.Ctx, mapperCashier),
		MerchantQuery:        NewMerchantQueryRepository(deps.DB, deps.Ctx, mapperMerchant),
		ProductQuery:         NewProductQueryRepository(deps.DB, deps.Ctx, mapperProduct),
		ProductCommand:       NewProductCommandRepository(deps.DB, deps.Ctx, mapperProduct),
		OrderQuery:           NewOrderQueryRepository(deps.DB, deps.Ctx, mapperOrder),
		OrderCommand:         NewOrderCommandRepository(deps.DB, deps.Ctx, mapperOrder),
		OrderItemQuery:       NewOrderItemQueryRepository(deps.DB, deps.Ctx, mapperOrderItem),
		OrderItemCommand:     NewOrderItemCommandRepository(deps.DB, deps.Ctx, mapperOrderItem),
		OrderStats:           NewOrderStatsRepository(deps.DB, deps.Ctx, mapperOrder),
		OrderStatsByMerchant: NewOrderStatsByMerchantRepository(deps.DB, deps.Ctx, mapperOrder),
	}
}
