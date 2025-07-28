package mencache

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
)

type CashierQueryCache interface {
	GetCachedCashiersCache(ctx context.Context, req *requests.FindAllCashiers) ([]*response.CashierResponse, *int, bool)
	SetCachedCashiersCache(ctx context.Context, req *requests.FindAllCashiers, res []*response.CashierResponse, total *int)

	GetCachedCashier(ctx context.Context, cashierID int) (*response.CashierResponse, bool)
	SetCachedCashier(ctx context.Context, res *response.CashierResponse)

	GetCachedCashiersActive(ctx context.Context, req *requests.FindAllCashiers) ([]*response.CashierResponseDeleteAt, *int, bool)
	SetCachedCashiersActive(ctx context.Context, req *requests.FindAllCashiers, res []*response.CashierResponseDeleteAt, total *int)

	GetCachedCashiersTrashed(ctx context.Context, req *requests.FindAllCashiers) ([]*response.CashierResponseDeleteAt, *int, bool)
	SetCachedCashiersTrashed(ctx context.Context, req *requests.FindAllCashiers, res []*response.CashierResponseDeleteAt, total *int)

	GetCachedCashiersByMerchant(ctx context.Context, req *requests.FindAllCashierMerchant) ([]*response.CashierResponse, *int, bool)
	SetCachedCashiersByMerchant(ctx context.Context, req *requests.FindAllCashierMerchant, res []*response.CashierResponse, total *int)
}

type CashierCommandCache interface {
	DeleteCashierCache(ctx context.Context, id int)
}

type CashierStatsCache interface {
	GetMonthlyTotalSalesCache(ctx context.Context, req *requests.MonthTotalSales) ([]*response.CashierResponseMonthTotalSales, bool)
	SetMonthlyTotalSalesCache(ctx context.Context, req *requests.MonthTotalSales, res []*response.CashierResponseMonthTotalSales)

	GetYearlyTotalSalesCache(ctx context.Context, year int) ([]*response.CashierResponseYearTotalSales, bool)
	SetYearlyTotalSalesCache(ctx context.Context, year int, res []*response.CashierResponseYearTotalSales)

	GetMonthlySalesCache(ctx context.Context, year int) ([]*response.CashierResponseMonthSales, bool)
	SetMonthlySalesCache(ctx context.Context, year int, res []*response.CashierResponseMonthSales)

	GetYearlySalesCache(ctx context.Context, year int) ([]*response.CashierResponseYearSales, bool)
	SetYearlySalesCache(ctx context.Context, year int, res []*response.CashierResponseYearSales)
}

type CashierStatsByIdCache interface {
	GetMonthlyTotalSalesByIdCache(ctx context.Context, req *requests.MonthTotalSalesCashier) ([]*response.CashierResponseMonthTotalSales, bool)
	SetMonthlyTotalSalesByIdCache(ctx context.Context, req *requests.MonthTotalSalesCashier, res []*response.CashierResponseMonthTotalSales)

	GetYearlyTotalSalesByIdCache(ctx context.Context, req *requests.YearTotalSalesCashier) ([]*response.CashierResponseYearTotalSales, bool)
	SetYearlyTotalSalesByIdCache(ctx context.Context, req *requests.YearTotalSalesCashier, res []*response.CashierResponseYearTotalSales)

	GetMonthlyCashierByIdCache(ctx context.Context, req *requests.MonthCashierId) ([]*response.CashierResponseMonthSales, bool)
	SetMonthlyCashierByIdCache(ctx context.Context, req *requests.MonthCashierId, res []*response.CashierResponseMonthSales)

	GetYearlyCashierByIdCache(ctx context.Context, req *requests.YearCashierId) ([]*response.CashierResponseYearSales, bool)
	SetYearlyCashierByIdCache(ctx context.Context, req *requests.YearCashierId, res []*response.CashierResponseYearSales)
}

type CashierStatsByMerchantCache interface {
	GetMonthlyTotalSalesByMerchantCache(ctx context.Context, req *requests.MonthTotalSalesMerchant) ([]*response.CashierResponseMonthTotalSales, bool)
	SetMonthlyTotalSalesByMerchantCache(ctx context.Context, req *requests.MonthTotalSalesMerchant, res []*response.CashierResponseMonthTotalSales)

	GetYearlyTotalSalesByMerchantCache(ctx context.Context, req *requests.YearTotalSalesMerchant) ([]*response.CashierResponseYearTotalSales, bool)
	SetYearlyTotalSalesByMerchantCache(ctx context.Context, req *requests.YearTotalSalesMerchant, res []*response.CashierResponseYearTotalSales)

	GetMonthlyCashierByMerchantCache(ctx context.Context, req *requests.MonthCashierMerchant) ([]*response.CashierResponseMonthSales, bool)
	SetMonthlyCashierByMerchantCache(ctx context.Context, req *requests.MonthCashierMerchant, res []*response.CashierResponseMonthSales)

	GetYearlyCashierByMerchantCache(ctx context.Context, req *requests.YearCashierMerchant) ([]*response.CashierResponseYearSales, bool)
	SetYearlyCashierByMerchantCache(ctx context.Context, req *requests.YearCashierMerchant, res []*response.CashierResponseYearSales)
}
