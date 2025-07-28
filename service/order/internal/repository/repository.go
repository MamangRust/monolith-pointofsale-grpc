package repository

import (
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

func NewRepositories(DB *db.Queries) *Repositories {
	mapperCashier := recordmapper.NewCashierRecordMapper()
	mapperMerchant := recordmapper.NewMerchantRecordMapper()
	mapperProduct := recordmapper.NewProductRecordMapper()
	mapperOrder := recordmapper.NewOrderRecordMapper()
	mapperOrderItem := recordmapper.NewOrderItemRecordMapper()

	return &Repositories{
		CashierQuery:         NewCashierQueryRepository(DB, mapperCashier),
		MerchantQuery:        NewMerchantQueryRepository(DB, mapperMerchant),
		ProductQuery:         NewProductQueryRepository(DB, mapperProduct),
		ProductCommand:       NewProductCommandRepository(DB, mapperProduct),
		OrderQuery:           NewOrderQueryRepository(DB, mapperOrder),
		OrderCommand:         NewOrderCommandRepository(DB, mapperOrder),
		OrderItemQuery:       NewOrderItemQueryRepository(DB, mapperOrderItem),
		OrderItemCommand:     NewOrderItemCommandRepository(DB, mapperOrderItem),
		OrderStats:           NewOrderStatsRepository(DB, mapperOrder),
		OrderStatsByMerchant: NewOrderStatsByMerchantRepository(DB, mapperOrder),
	}
}
