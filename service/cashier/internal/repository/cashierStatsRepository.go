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

type cashierStatsRepository struct {
	db      *db.Queries
	mapping recordmapper.CashierRecordMapping
}

func NewCashierStatsRepository(db *db.Queries, mapping recordmapper.CashierRecordMapping) *cashierStatsRepository {
	return &cashierStatsRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *cashierStatsRepository) GetMonthlyTotalSales(ctx context.Context, req *requests.MonthTotalSales) ([]*record.CashierRecordMonthTotalSales, error) {
	currentMonthStart := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	currentMonthEnd := currentMonthStart.AddDate(0, 1, -1)

	prevMonthStart := currentMonthStart.AddDate(0, -1, 0)
	prevMonthEnd := prevMonthStart.AddDate(0, 1, -1)

	params := db.GetMonthlyTotalSalesCashierParams{
		Extract:     currentMonthStart,
		CreatedAt:   sql.NullTime{Time: currentMonthEnd, Valid: true},
		CreatedAt_2: sql.NullTime{Time: prevMonthStart, Valid: true},
		CreatedAt_3: sql.NullTime{Time: prevMonthEnd, Valid: true},
	}

	res, err := r.db.GetMonthlyTotalSalesCashier(ctx, params)

	if err != nil {
		return nil, cashier_errors.ErrGetMonthlyTotalSales
	}

	return r.mapping.ToCashierMonthlyTotalSales(res), nil
}

func (r *cashierStatsRepository) GetYearlyTotalSales(ctx context.Context, year int) ([]*record.CashierRecordYearTotalSales, error) {
	res, err := r.db.GetYearlyTotalSalesCashier(ctx, int32(year))

	if err != nil {
		return nil, cashier_errors.ErrGetYearlyTotalSales
	}

	so := r.mapping.ToCashierYearlyTotalSales(res)

	return so, nil
}

func (r *cashierStatsRepository) GetMonthyCashier(ctx context.Context, year int) ([]*record.CashierRecordMonthSales, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyCashier(ctx, yearStart)

	if err != nil {
		return nil, cashier_errors.ErrGetMonthlyCashier
	}

	return r.mapping.ToCashierMonthlySales(res), nil

}

func (r *cashierStatsRepository) GetYearlyCashier(ctx context.Context, year int) ([]*record.CashierRecordYearSales, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetYearlyCashier(ctx, yearStart)

	if err != nil {
		return nil, cashier_errors.ErrGetYearlyCashier
	}

	return r.mapping.ToCashierYearlySales(res), nil
}
