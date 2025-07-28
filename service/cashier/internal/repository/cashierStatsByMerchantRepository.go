package repository

import (
	"context"
	"database/sql"
	"time"

	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/record"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/cashier_errors"
	recordmapper "github.com/MamangRust/monolith-point-of-sale-shared/mapper/record"
)

type cashierStatsByMerchantRepository struct {
	db      *db.Queries
	mapping recordmapper.CashierRecordMapping
}

func NewCashierStatsByMerchantRepository(db *db.Queries, mapping recordmapper.CashierRecordMapping) *cashierStatsByMerchantRepository {
	return &cashierStatsByMerchantRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *cashierStatsByMerchantRepository) GetMonthlyTotalSalesByMerchant(ctx context.Context, req *requests.MonthTotalSalesMerchant) ([]*record.CashierRecordMonthTotalSales, error) {
	currentMonthStart := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	currentMonthEnd := currentMonthStart.AddDate(0, 1, -1)
	prevMonthStart := currentMonthStart.AddDate(0, -1, 0)
	prevMonthEnd := prevMonthStart.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyTotalSalesByMerchant(ctx, db.GetMonthlyTotalSalesByMerchantParams{
		Extract:     currentMonthStart,
		CreatedAt:   sql.NullTime{Time: currentMonthEnd, Valid: true},
		CreatedAt_2: sql.NullTime{Time: prevMonthStart, Valid: true},
		CreatedAt_3: sql.NullTime{Time: prevMonthEnd, Valid: true},
		MerchantID:  int32(req.MerchantID),
	})

	if err != nil {
		return nil, cashier_errors.ErrGetMonthlyTotalSalesByMerchant
	}

	so := r.mapping.ToCashierMonthlyTotalSalesByMerchant(res)

	return so, nil
}

func (r *cashierStatsByMerchantRepository) GetYearlyTotalSalesByMerchant(ctx context.Context, req *requests.YearTotalSalesMerchant) ([]*record.CashierRecordYearTotalSales, error) {
	res, err := r.db.GetYearlyTotalSalesByMerchant(ctx, db.GetYearlyTotalSalesByMerchantParams{
		Column1:    int32(req.Year),
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, cashier_errors.ErrGetYearlyTotalSalesByMerchant
	}

	so := r.mapping.ToCashierYearlyTotalSalesByMerchant(res)

	return so, nil
}

func (r *cashierStatsByMerchantRepository) GetMonthlyCashierByMerchant(ctx context.Context, req *requests.MonthCashierMerchant) ([]*record.CashierRecordMonthSales, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyCashierByMerchant(ctx, db.GetMonthlyCashierByMerchantParams{
		Column1:    yearStart,
		MerchantID: int32(req.MerchantID),
	})

	if err != nil {
		return nil, cashier_errors.ErrGetMonthlyCashierByMerchant
	}

	return r.mapping.ToCashierMonthlySalesByMerchant(res), nil

}

func (r *cashierStatsByMerchantRepository) GetYearlyCashierByMerchant(ctx context.Context, req *requests.YearCashierMerchant) ([]*record.CashierRecordYearSales, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetYearlyCashierByMerchant(ctx, db.GetYearlyCashierByMerchantParams{
		Column1:    yearStart,
		MerchantID: int32(req.Year),
	})

	if err != nil {
		return nil, cashier_errors.ErrGetYearlyCashierByMerchant
	}

	return r.mapping.ToCashierYearlySalesByMerchant(res), nil
}
