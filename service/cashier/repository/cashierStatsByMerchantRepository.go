package repository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/cashier_errors"
	"github.com/jackc/pgx/v5/pgtype"
)

type cashierStatsByMerchantRepository struct {
	db *db.Queries
}

func NewCashierStatsByMerchantRepository(db *db.Queries) CashierStatByMerchantRepository {
	return &cashierStatsByMerchantRepository{
		db: db,
	}
}

func (r *cashierStatsByMerchantRepository) GetMonthlyTotalSalesByMerchant(ctx context.Context, req *requests.MonthTotalSalesMerchant) ([]*db.GetMonthlyTotalSalesByMerchantRow, error) {
	currentMonthStart := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	currentMonthEnd := currentMonthStart.AddDate(0, 1, -1)
	prevMonthStart := currentMonthStart.AddDate(0, -1, 0)
	prevMonthEnd := prevMonthStart.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyTotalSalesByMerchant(ctx, db.GetMonthlyTotalSalesByMerchantParams{
		Extract:     pgtype.Date{Time: currentMonthStart, Valid: true},
		CreatedAt:   pgtype.Timestamp{Time: currentMonthEnd, Valid: true},
		CreatedAt_2: pgtype.Timestamp{Time: prevMonthStart, Valid: true},
		CreatedAt_3: pgtype.Timestamp{Time: prevMonthEnd, Valid: true},
		MerchantID:  int32(req.MerchantID),
	})
	if err != nil {
		return nil, cashier_errors.ErrGetMonthlyTotalSalesByMerchant
	}

	return res, nil
}

func (r *cashierStatsByMerchantRepository) GetYearlyTotalSalesByMerchant(ctx context.Context, req *requests.YearTotalSalesMerchant) ([]*db.GetYearlyTotalSalesByMerchantRow, error) {
	res, err := r.db.GetYearlyTotalSalesByMerchant(ctx, db.GetYearlyTotalSalesByMerchantParams{
		Column1:    int32(req.Year),
		MerchantID: int32(req.MerchantID),
	})
	if err != nil {
		return nil, cashier_errors.ErrGetYearlyTotalSalesByMerchant
	}

	return res, nil
}

func (r *cashierStatsByMerchantRepository) GetMonthlyCashierByMerchant(ctx context.Context, req *requests.MonthCashierMerchant) ([]*db.GetMonthlyCashierByMerchantRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyCashierByMerchant(ctx, db.GetMonthlyCashierByMerchantParams{
		Column1:    yearStart,
		MerchantID: int32(req.MerchantID),
	})
	if err != nil {
		return nil, cashier_errors.ErrGetMonthlyCashierByMerchant
	}

	return res, nil
}

func (r *cashierStatsByMerchantRepository) GetYearlyCashierByMerchant(ctx context.Context, req *requests.YearCashierMerchant) ([]*db.GetYearlyCashierByMerchantRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetYearlyCashierByMerchant(ctx, db.GetYearlyCashierByMerchantParams{
		Column1:    yearStart,
		MerchantID: int32(req.Year),
	})
	if err != nil {
		return nil, cashier_errors.ErrGetYearlyCashierByMerchant
	}

	return res, nil
}
