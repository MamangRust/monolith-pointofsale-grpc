package repository

import (
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/record"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type MerchantQueryRepository interface {
	FindById(id int) (*record.MerchantRecord, error)
}

type UserQueryRepository interface {
	FindById(id int) (*record.UserRecord, error)
}

type CashierStatsRepository interface {
	GetMonthlyTotalSales(req *requests.MonthTotalSales) ([]*record.CashierRecordMonthTotalSales, error)
	GetYearlyTotalSales(year int) ([]*record.CashierRecordYearTotalSales, error)

	GetMonthyCashier(year int) ([]*record.CashierRecordMonthSales, error)
	GetYearlyCashier(year int) ([]*record.CashierRecordYearSales, error)
}

type CashierStatByIdRepository interface {
	GetMonthlyTotalSalesById(req *requests.MonthTotalSalesCashier) ([]*record.CashierRecordMonthTotalSales, error)
	GetYearlyTotalSalesById(req *requests.YearTotalSalesCashier) ([]*record.CashierRecordYearTotalSales, error)

	GetMonthlyCashierById(req *requests.MonthCashierId) ([]*record.CashierRecordMonthSales, error)
	GetYearlyCashierById(req *requests.YearCashierId) ([]*record.CashierRecordYearSales, error)
}

type CashierStatByMerchantRepository interface {
	GetMonthlyTotalSalesByMerchant(req *requests.MonthTotalSalesMerchant) ([]*record.CashierRecordMonthTotalSales, error)
	GetYearlyTotalSalesByMerchant(req *requests.YearTotalSalesMerchant) ([]*record.CashierRecordYearTotalSales, error)

	GetMonthlyCashierByMerchant(req *requests.MonthCashierMerchant) ([]*record.CashierRecordMonthSales, error)
	GetYearlyCashierByMerchant(req *requests.YearCashierMerchant) ([]*record.CashierRecordYearSales, error)
}

type CashierQueryRepository interface {
	FindAllCashiers(req *requests.FindAllCashiers) ([]*record.CashierRecord, *int, error)
	FindById(cashier_id int) (*record.CashierRecord, error)
	FindByActive(req *requests.FindAllCashiers) ([]*record.CashierRecord, *int, error)
	FindByTrashed(req *requests.FindAllCashiers) ([]*record.CashierRecord, *int, error)
	FindByMerchant(req *requests.FindAllCashierMerchant) ([]*record.CashierRecord, *int, error)
}

type CashierCommandRepository interface {
	CreateCashier(request *requests.CreateCashierRequest) (*record.CashierRecord, error)
	UpdateCashier(request *requests.UpdateCashierRequest) (*record.CashierRecord, error)
	TrashedCashier(cashier_id int) (*record.CashierRecord, error)
	RestoreCashier(cashier_id int) (*record.CashierRecord, error)
	DeleteCashierPermanent(cashier_id int) (bool, error)
	RestoreAllCashier() (bool, error)
	DeleteAllCashierPermanent() (bool, error)
}
