package repository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/cashier_errors"
	"github.com/jackc/pgx/v5/pgtype"
)

type cashierStatsByIdRepository struct {
	db *db.Queries
}

func NewCashierStatsByIdRepository(db *db.Queries) CashierStatByIdRepository {
	return &cashierStatsByIdRepository{
		db: db,
	}
}

func (r *cashierStatsByIdRepository) GetMonthlyTotalSalesById(ctx context.Context, req *requests.MonthTotalSalesCashier) ([]*db.GetMonthlyTotalSalesByIdRow, error) {
	currentMonthStart := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	currentMonthEnd := currentMonthStart.AddDate(0, 1, -1)

	prevMonthStart := currentMonthStart.AddDate(0, -1, 0)
	prevMonthEnd := prevMonthStart.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyTotalSalesById(ctx, db.GetMonthlyTotalSalesByIdParams{
		Extract:     pgtype.Date{Time: currentMonthStart, Valid: true},
		CreatedAt:   pgtype.Timestamp{Time: currentMonthEnd, Valid: true},
		CreatedAt_2: pgtype.Timestamp{Time: prevMonthStart, Valid: true},
		CreatedAt_3: pgtype.Timestamp{Time: prevMonthEnd, Valid: true},
		CashierID:   int32(req.CashierID),
	})
	if err != nil {
		return nil, cashier_errors.ErrGetMonthlyTotalSalesById
	}

	return res, nil
}

func (r *cashierStatsByIdRepository) GetYearlyTotalSalesById(ctx context.Context, req *requests.YearTotalSalesCashier) ([]*db.GetYearlyTotalSalesByIdRow, error) {
	res, err := r.db.GetYearlyTotalSalesById(ctx, db.GetYearlyTotalSalesByIdParams{
		Column1:   int32(req.Year),
		CashierID: int32(req.CashierID),
	})
	if err != nil {
		return nil, cashier_errors.ErrGetYearlyTotalSalesById
	}

	return res, nil
}

func (r *cashierStatsByIdRepository) GetMonthlyCashierById(ctx context.Context, req *requests.MonthCashierId) ([]*db.GetMonthlyCashierByCashierIdRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyCashierByCashierId(ctx, db.GetMonthlyCashierByCashierIdParams{
		Column1:   yearStart,
		CashierID: int32(req.CashierID),
	})
	if err != nil {
		return nil, cashier_errors.ErrGetMonthlyCashierById
	}

	return res, nil
}

func (r *cashierStatsByIdRepository) GetYearlyCashierById(ctx context.Context, req *requests.YearCashierId) ([]*db.GetYearlyCashierByCashierIdRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetYearlyCashierByCashierId(ctx, db.GetYearlyCashierByCashierIdParams{
		Column1:   yearStart,
		CashierID: int32(req.CashierID),
	})
	if err != nil {
		return nil, cashier_errors.ErrGetYearlyCashierById
	}

	return res, nil
}
