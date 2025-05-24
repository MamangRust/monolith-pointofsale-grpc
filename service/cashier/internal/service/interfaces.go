package service

import (
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
)

type CashierStatsService interface {
	FindMonthlyTotalSales(req *requests.MonthTotalSales) ([]*response.CashierResponseMonthTotalSales, *response.ErrorResponse)
	FindYearlyTotalSales(year int) ([]*response.CashierResponseYearTotalSales, *response.ErrorResponse)

	FindMonthlySales(year int) ([]*response.CashierResponseMonthSales, *response.ErrorResponse)
	FindYearlySales(year int) ([]*response.CashierResponseYearSales, *response.ErrorResponse)
}

type CashierStatsByIdService interface {
	FindMonthlyTotalSalesById(req *requests.MonthTotalSalesCashier) ([]*response.CashierResponseMonthTotalSales, *response.ErrorResponse)
	FindYearlyTotalSalesById(req *requests.YearTotalSalesCashier) ([]*response.CashierResponseYearTotalSales, *response.ErrorResponse)
	FindMonthlyCashierById(req *requests.MonthCashierId) ([]*response.CashierResponseMonthSales, *response.ErrorResponse)
	FindYearlyCashierById(req *requests.YearCashierId) ([]*response.CashierResponseYearSales, *response.ErrorResponse)
}

type CashierStatsByMerchant interface {
	FindMonthlyTotalSalesByMerchant(req *requests.MonthTotalSalesMerchant) ([]*response.CashierResponseMonthTotalSales, *response.ErrorResponse)
	FindYearlyTotalSalesByMerchant(req *requests.YearTotalSalesMerchant) ([]*response.CashierResponseYearTotalSales, *response.ErrorResponse)

	FindMonthlyCashierByMerchant(req *requests.MonthCashierMerchant) ([]*response.CashierResponseMonthSales, *response.ErrorResponse)
	FindYearlyCashierByMerchant(req *requests.YearCashierMerchant) ([]*response.CashierResponseYearSales, *response.ErrorResponse)
}

type CashierQueryService interface {
	FindAll(req *requests.FindAllCashiers) ([]*response.CashierResponse, *int, *response.ErrorResponse)
	FindById(cashierID int) (*response.CashierResponse, *response.ErrorResponse)
	FindByActive(req *requests.FindAllCashiers) ([]*response.CashierResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllCashiers) ([]*response.CashierResponseDeleteAt, *int, *response.ErrorResponse)
	FindByMerchant(req *requests.FindAllCashierMerchant) ([]*response.CashierResponse, *int, *response.ErrorResponse)
}

type CashierCommandService interface {
	CreateCashier(req *requests.CreateCashierRequest) (*response.CashierResponse, *response.ErrorResponse)
	UpdateCashier(req *requests.UpdateCashierRequest) (*response.CashierResponse, *response.ErrorResponse)
	TrashedCashier(cashierID int) (*response.CashierResponseDeleteAt, *response.ErrorResponse)
	RestoreCashier(cashierID int) (*response.CashierResponseDeleteAt, *response.ErrorResponse)
	DeleteCashierPermanent(cashierID int) (bool, *response.ErrorResponse)
	RestoreAllCashier() (bool, *response.ErrorResponse)
	DeleteAllCashierPermanent() (bool, *response.ErrorResponse)
}
