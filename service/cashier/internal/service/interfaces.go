package service

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
)

type CashierStatsService interface {
	FindMonthlyTotalSales(ctx context.Context, req *requests.MonthTotalSales) ([]*response.CashierResponseMonthTotalSales, *response.ErrorResponse)
	FindYearlyTotalSales(ctx context.Context, year int) ([]*response.CashierResponseYearTotalSales, *response.ErrorResponse)

	FindMonthlySales(ctx context.Context, year int) ([]*response.CashierResponseMonthSales, *response.ErrorResponse)
	FindYearlySales(ctx context.Context, year int) ([]*response.CashierResponseYearSales, *response.ErrorResponse)
}

type CashierStatsByIdService interface {
	FindMonthlyTotalSalesById(ctx context.Context, req *requests.MonthTotalSalesCashier) ([]*response.CashierResponseMonthTotalSales, *response.ErrorResponse)
	FindYearlyTotalSalesById(ctx context.Context, req *requests.YearTotalSalesCashier) ([]*response.CashierResponseYearTotalSales, *response.ErrorResponse)
	FindMonthlyCashierById(ctx context.Context, req *requests.MonthCashierId) ([]*response.CashierResponseMonthSales, *response.ErrorResponse)
	FindYearlyCashierById(ctx context.Context, req *requests.YearCashierId) ([]*response.CashierResponseYearSales, *response.ErrorResponse)
}

type CashierStatsByMerchant interface {
	FindMonthlyTotalSalesByMerchant(ctx context.Context, req *requests.MonthTotalSalesMerchant) ([]*response.CashierResponseMonthTotalSales, *response.ErrorResponse)
	FindYearlyTotalSalesByMerchant(ctx context.Context, req *requests.YearTotalSalesMerchant) ([]*response.CashierResponseYearTotalSales, *response.ErrorResponse)

	FindMonthlyCashierByMerchant(ctx context.Context, req *requests.MonthCashierMerchant) ([]*response.CashierResponseMonthSales, *response.ErrorResponse)
	FindYearlyCashierByMerchant(ctx context.Context, req *requests.YearCashierMerchant) ([]*response.CashierResponseYearSales, *response.ErrorResponse)
}

type CashierQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllCashiers) ([]*response.CashierResponse, *int, *response.ErrorResponse)
	FindById(ctx context.Context, cashierID int) (*response.CashierResponse, *response.ErrorResponse)
	FindByActive(ctx context.Context, req *requests.FindAllCashiers) ([]*response.CashierResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(ctx context.Context, req *requests.FindAllCashiers) ([]*response.CashierResponseDeleteAt, *int, *response.ErrorResponse)
	FindByMerchant(ctx context.Context, req *requests.FindAllCashierMerchant) ([]*response.CashierResponse, *int, *response.ErrorResponse)
}

type CashierCommandService interface {
	CreateCashier(ctx context.Context, req *requests.CreateCashierRequest) (*response.CashierResponse, *response.ErrorResponse)
	UpdateCashier(ctx context.Context, req *requests.UpdateCashierRequest) (*response.CashierResponse, *response.ErrorResponse)
	TrashedCashier(ctx context.Context, cashierID int) (*response.CashierResponseDeleteAt, *response.ErrorResponse)
	RestoreCashier(ctx context.Context, cashierID int) (*response.CashierResponseDeleteAt, *response.ErrorResponse)
	DeleteCashierPermanent(ctx context.Context, cashierID int) (bool, *response.ErrorResponse)
	RestoreAllCashier(ctx context.Context) (bool, *response.ErrorResponse)
	DeleteAllCashierPermanent(ctx context.Context) (bool, *response.ErrorResponse)
}
