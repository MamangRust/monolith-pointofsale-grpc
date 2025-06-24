package mencache

import (
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
)

type OrderStatsCache interface {
	GetMonthlyTotalRevenueCache(req *requests.MonthTotalRevenue) ([]*response.OrderMonthlyTotalRevenueResponse, bool)
	SetMonthlyTotalRevenueCache(req *requests.MonthTotalRevenue, res []*response.OrderMonthlyTotalRevenueResponse)

	GetYearlyTotalRevenueCache(year int) ([]*response.OrderYearlyTotalRevenueResponse, bool)
	SetYearlyTotalRevenueCache(year int, res []*response.OrderYearlyTotalRevenueResponse)

	GetMonthlyOrderCache(year int) ([]*response.OrderMonthlyResponse, bool)
	SetMonthlyOrderCache(year int, res []*response.OrderMonthlyResponse)

	GetYearlyOrderCache(year int) ([]*response.OrderYearlyResponse, bool)
	SetYearlyOrderCache(year int, res []*response.OrderYearlyResponse)
}

type OrderStatsByMerchantCache interface {
	GetMonthlyTotalRevenueByMerchantCache(req *requests.MonthTotalRevenueMerchant) ([]*response.OrderMonthlyTotalRevenueResponse, bool)
	SetMonthlyTotalRevenueByMerchantCache(req *requests.MonthTotalRevenueMerchant, res []*response.OrderMonthlyTotalRevenueResponse)

	GetYearlyTotalRevenueByMerchantCache(req *requests.YearTotalRevenueMerchant) ([]*response.OrderYearlyTotalRevenueResponse, bool)
	SetYearlyTotalRevenueByMerchantCache(req *requests.YearTotalRevenueMerchant, res []*response.OrderYearlyTotalRevenueResponse)

	GetMonthlyOrderByMerchantCache(req *requests.MonthOrderMerchant) ([]*response.OrderMonthlyResponse, bool)
	SetMonthlyOrderByMerchantCache(req *requests.MonthOrderMerchant, res []*response.OrderMonthlyResponse)

	GetYearlyOrderByMerchantCache(req *requests.YearOrderMerchant) ([]*response.OrderYearlyResponse, bool)
	SetYearlyOrderByMerchantCache(req *requests.YearOrderMerchant, res []*response.OrderYearlyResponse)
}

type OrderQueryCache interface {
	GetOrderAllCache(req *requests.FindAllOrders) ([]*response.OrderResponse, *int, bool)
	SetOrderAllCache(req *requests.FindAllOrders, data []*response.OrderResponse, total *int)

	GetCachedOrderCache(order_id int) (*response.OrderResponse, bool)
	SetCachedOrderCache(data *response.OrderResponse)

	GetCachedOrderMerchant(req *requests.FindAllOrderMerchant) ([]*response.OrderResponse, *int, bool)
	SetCachedOrderMerchant(req *requests.FindAllOrderMerchant, res []*response.OrderResponse, total *int)

	GetOrderActiveCache(req *requests.FindAllOrders) ([]*response.OrderResponseDeleteAt, *int, bool)
	SetOrderActiveCache(req *requests.FindAllOrders, data []*response.OrderResponseDeleteAt, total *int)

	GetOrderTrashedCache(req *requests.FindAllOrders) ([]*response.OrderResponseDeleteAt, *int, bool)
	SetOrderTrashedCache(req *requests.FindAllOrders, data []*response.OrderResponseDeleteAt, total *int)
}

type OrderCommandCache interface {
	DeleteOrderCache(id int)
}
