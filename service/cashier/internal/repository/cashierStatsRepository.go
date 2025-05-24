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
	ctx     context.Context
	mapping recordmapper.CashierRecordMapping
}

func NewCashierStatsRepository(db *db.Queries, ctx context.Context, mapping recordmapper.CashierRecordMapping) *cashierStatsRepository {
	return &cashierStatsRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *cashierStatsRepository) GetMonthlyTotalSales(req *requests.MonthTotalSales) ([]*record.CashierRecordMonthTotalSales, error) {
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

	res, err := r.db.GetMonthlyTotalSalesCashier(r.ctx, params)

	if err != nil {
		return nil, cashier_errors.ErrGetMonthlyTotalSales
	}

	return r.mapping.ToCashierMonthlyTotalSales(res), nil
}

func (r *cashierStatsRepository) GetYearlyTotalSales(year int) ([]*record.CashierRecordYearTotalSales, error) {
	res, err := r.db.GetYearlyTotalSalesCashier(r.ctx, int32(year))

	if err != nil {
		return nil, cashier_errors.ErrGetYearlyTotalSales
	}

	so := r.mapping.ToCashierYearlyTotalSales(res)

	return so, nil
}

func (r *cashierStatsRepository) GetMonthyCashier(year int) ([]*record.CashierRecordMonthSales, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyCashier(r.ctx, yearStart)

	if err != nil {
		return nil, cashier_errors.ErrGetMonthlyCashier
	}

	return r.mapping.ToCashierMonthlySales(res), nil

}

func (r *cashierStatsRepository) GetYearlyCashier(year int) ([]*record.CashierRecordYearSales, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetYearlyCashier(r.ctx, yearStart)

	if err != nil {
		return nil, cashier_errors.ErrGetYearlyCashier
	}

	return r.mapping.ToCashierYearlySales(res), nil
}
