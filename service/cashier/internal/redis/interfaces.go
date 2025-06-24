package mencache

import (
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
)

type CashierQueryCache interface {
	GetCachedCashiersCache(req *requests.FindAllCashiers) ([]*response.CashierResponse, *int, bool)
	SetCachedCashiersCache(req *requests.FindAllCashiers, res []*response.CashierResponse, total *int)

	GetCachedCashier(cashierID int) (*response.CashierResponse, bool)
	SetCachedCashier(res *response.CashierResponse)

	GetCachedCashiersActive(req *requests.FindAllCashiers) ([]*response.CashierResponseDeleteAt, *int, bool)
	SetCachedCashiersActive(req *requests.FindAllCashiers, res []*response.CashierResponseDeleteAt, total *int)

	GetCachedCashiersTrashed(req *requests.FindAllCashiers) ([]*response.CashierResponseDeleteAt, *int, bool)
	SetCachedCashiersTrashed(req *requests.FindAllCashiers, res []*response.CashierResponseDeleteAt, total *int)

	GetCachedCashiersByMerchant(req *requests.FindAllCashierMerchant) ([]*response.CashierResponse, *int, bool)
	SetCachedCashiersByMerchant(req *requests.FindAllCashierMerchant, res []*response.CashierResponse, total *int)
}

type CashierCommandCache interface {
	DeleteCashierCache(id int)
}

type CashierStatsCache interface {
	GetMonthlyTotalSalesCache(req *requests.MonthTotalSales) ([]*response.CashierResponseMonthTotalSales, bool)
	SetMonthlyTotalSalesCache(req *requests.MonthTotalSales, res []*response.CashierResponseMonthTotalSales)

	GetYearlyTotalSalesCache(year int) ([]*response.CashierResponseYearTotalSales, bool)
	SetYearlyTotalSalesCache(year int, res []*response.CashierResponseYearTotalSales)

	GetMonthlySalesCache(year int) ([]*response.CashierResponseMonthSales, bool)
	SetMonthlySalesCache(year int, res []*response.CashierResponseMonthSales)

	GetYearlySalesCache(year int) ([]*response.CashierResponseYearSales, bool)
	SetYearlySalesCache(year int, res []*response.CashierResponseYearSales)
}

type CashierStatsByIdCache interface {
	GetMonthlyTotalSalesByIdCache(req *requests.MonthTotalSalesCashier) ([]*response.CashierResponseMonthTotalSales, bool)
	SetMonthlyTotalSalesByIdCache(req *requests.MonthTotalSalesCashier, res []*response.CashierResponseMonthTotalSales)

	GetYearlyTotalSalesByIdCache(req *requests.YearTotalSalesCashier) ([]*response.CashierResponseYearTotalSales, bool)
	SetYearlyTotalSalesByIdCache(req *requests.YearTotalSalesCashier, res []*response.CashierResponseYearTotalSales)

	GetMonthlyCashierByIdCache(req *requests.MonthCashierId) ([]*response.CashierResponseMonthSales, bool)
	SetMonthlyCashierByIdCache(req *requests.MonthCashierId, res []*response.CashierResponseMonthSales)

	GetYearlyCashierByIdCache(req *requests.YearCashierId) ([]*response.CashierResponseYearSales, bool)
	SetYearlyCashierByIdCache(req *requests.YearCashierId, res []*response.CashierResponseYearSales)
}

type CashierStatsByMerchantCache interface {
	GetMonthlyTotalSalesByMerchantCache(req *requests.MonthTotalSalesMerchant) ([]*response.CashierResponseMonthTotalSales, bool)
	SetMonthlyTotalSalesByMerchantCache(req *requests.MonthTotalSalesMerchant, res []*response.CashierResponseMonthTotalSales)

	GetYearlyTotalSalesByMerchantCache(req *requests.YearTotalSalesMerchant) ([]*response.CashierResponseYearTotalSales, bool)
	SetYearlyTotalSalesByMerchantCache(req *requests.YearTotalSalesMerchant, res []*response.CashierResponseYearTotalSales)

	GetMonthlyCashierByMerchantCache(req *requests.MonthCashierMerchant) ([]*response.CashierResponseMonthSales, bool)
	SetMonthlyCashierByMerchantCache(req *requests.MonthCashierMerchant, res []*response.CashierResponseMonthSales)

	GetYearlyCashierByMerchantCache(req *requests.YearCashierMerchant) ([]*response.CashierResponseYearSales, bool)
	SetYearlyCashierByMerchantCache(req *requests.YearCashierMerchant, res []*response.CashierResponseYearSales)
}
