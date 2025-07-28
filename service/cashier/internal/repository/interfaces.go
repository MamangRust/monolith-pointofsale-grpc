package repository

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-shared/domain/record"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type MerchantQueryRepository interface {
	FindById(ctx context.Context, id int) (*record.MerchantRecord, error)
}

type UserQueryRepository interface {
	FindById(ctx context.Context, id int) (*record.UserRecord, error)
}

type CashierStatsRepository interface {
	GetMonthlyTotalSales(ctx context.Context, req *requests.MonthTotalSales) ([]*record.CashierRecordMonthTotalSales, error)
	GetYearlyTotalSales(ctx context.Context, year int) ([]*record.CashierRecordYearTotalSales, error)

	GetMonthyCashier(ctx context.Context, year int) ([]*record.CashierRecordMonthSales, error)
	GetYearlyCashier(ctx context.Context, year int) ([]*record.CashierRecordYearSales, error)
}

type CashierStatByIdRepository interface {
	GetMonthlyTotalSalesById(ctx context.Context, req *requests.MonthTotalSalesCashier) ([]*record.CashierRecordMonthTotalSales, error)
	GetYearlyTotalSalesById(ctx context.Context, req *requests.YearTotalSalesCashier) ([]*record.CashierRecordYearTotalSales, error)

	GetMonthlyCashierById(ctx context.Context, req *requests.MonthCashierId) ([]*record.CashierRecordMonthSales, error)
	GetYearlyCashierById(ctx context.Context, req *requests.YearCashierId) ([]*record.CashierRecordYearSales, error)
}

type CashierStatByMerchantRepository interface {
	GetMonthlyTotalSalesByMerchant(ctx context.Context, req *requests.MonthTotalSalesMerchant) ([]*record.CashierRecordMonthTotalSales, error)
	GetYearlyTotalSalesByMerchant(ctx context.Context, req *requests.YearTotalSalesMerchant) ([]*record.CashierRecordYearTotalSales, error)

	GetMonthlyCashierByMerchant(ctx context.Context, req *requests.MonthCashierMerchant) ([]*record.CashierRecordMonthSales, error)
	GetYearlyCashierByMerchant(ctx context.Context, req *requests.YearCashierMerchant) ([]*record.CashierRecordYearSales, error)
}

type CashierQueryRepository interface {
	FindAllCashiers(ctx context.Context, req *requests.FindAllCashiers) ([]*record.CashierRecord, *int, error)
	FindById(ctx context.Context, cashier_id int) (*record.CashierRecord, error)
	FindByActive(ctx context.Context, req *requests.FindAllCashiers) ([]*record.CashierRecord, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllCashiers) ([]*record.CashierRecord, *int, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllCashierMerchant) ([]*record.CashierRecord, *int, error)
}

type CashierCommandRepository interface {
	CreateCashier(ctx context.Context, request *requests.CreateCashierRequest) (*record.CashierRecord, error)
	UpdateCashier(ctx context.Context, request *requests.UpdateCashierRequest) (*record.CashierRecord, error)
	TrashedCashier(ctx context.Context, cashier_id int) (*record.CashierRecord, error)
	RestoreCashier(ctx context.Context, cashier_id int) (*record.CashierRecord, error)
	DeleteCashierPermanent(ctx context.Context, cashier_id int) (bool, error)
	RestoreAllCashier(ctx context.Context) (bool, error)
	DeleteAllCashierPermanent(ctx context.Context) (bool, error)
}
