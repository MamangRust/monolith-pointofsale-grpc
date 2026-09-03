package repository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/cashier_errors"
	"github.com/jackc/pgx/v5/pgtype"
)

type cashierStatsRepository struct {
	db *db.Queries
}

func NewCashierStatsRepository(db *db.Queries) CashierStatsRepository {
	return &cashierStatsRepository{
		db: db,
	}
}

func (r *cashierStatsRepository) GetMonthlyTotalSales(ctx context.Context, req *requests.MonthTotalSales) ([]*db.GetMonthlyTotalSalesCashierRow, error) {
	currentMonthStart := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	currentMonthEnd := currentMonthStart.AddDate(0, 1, -1)

	prevMonthStart := currentMonthStart.AddDate(0, -1, 0)
	prevMonthEnd := prevMonthStart.AddDate(0, 1, -1)

	params := db.GetMonthlyTotalSalesCashierParams{
		Extract:     pgtype.Date{Time: currentMonthStart, Valid: true},
		CreatedAt:   pgtype.Timestamp{Time: currentMonthEnd, Valid: true},
		CreatedAt_2: pgtype.Timestamp{Time: prevMonthStart, Valid: true},
		CreatedAt_3: pgtype.Timestamp{Time: prevMonthEnd, Valid: true},
	}

	res, err := r.db.GetMonthlyTotalSalesCashier(ctx, params)
	if err != nil {
		return nil, cashier_errors.ErrGetMonthlyTotalSales
	}

	return res, nil
}

func (r *cashierStatsRepository) GetYearlyTotalSales(ctx context.Context, year int) ([]*db.GetYearlyTotalSalesCashierRow, error) {
	res, err := r.db.GetYearlyTotalSalesCashier(ctx, int32(year))
	if err != nil {
		return nil, cashier_errors.ErrGetYearlyTotalSales
	}

	return res, nil
}

func (r *cashierStatsRepository) GetMonthyCashier(ctx context.Context, year int) ([]*db.GetMonthlyCashierRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyCashier(ctx, yearStart)
	if err != nil {
		return nil, cashier_errors.ErrGetMonthlyCashier
	}

	return res, nil
}

func (r *cashierStatsRepository) GetYearlyCashier(ctx context.Context, year int) ([]*db.GetYearlyCashierRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetYearlyCashier(ctx, yearStart)
	if err != nil {
		return nil, cashier_errors.ErrGetYearlyCashier
	}

	return res, nil
}
