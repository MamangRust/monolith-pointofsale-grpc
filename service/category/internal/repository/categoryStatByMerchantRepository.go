package repository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"github.com/jackc/pgx/v5/pgtype"
)

type categoryStatsByMerchantRepository struct {
	db *db.Queries
}

func NewCategoryStatsByMerchantRepository(db *db.Queries) CategoryStatsByMerchantRepository {
	return &categoryStatsByMerchantRepository{
		db: db,
	}
}

func (r *categoryStatsByMerchantRepository) GetMonthlyTotalPriceByMerchant(ctx context.Context, req *requests.MonthTotalPriceMerchant) ([]*db.GetMonthlyTotalPriceByMerchantRow, error) {
	currentMonthStart := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	currentMonthEnd := currentMonthStart.AddDate(0, 1, -1)
	prevMonthStart := currentMonthStart.AddDate(0, -1, 0)
	prevMonthEnd := prevMonthStart.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyTotalPriceByMerchant(ctx, db.GetMonthlyTotalPriceByMerchantParams{
		Extract:     pgtype.Date{Time: currentMonthStart, Valid: true},
		CreatedAt:   pgtype.Timestamp{Time: currentMonthEnd, Valid: true},
		CreatedAt_2: pgtype.Timestamp{Time: prevMonthStart, Valid: true},
		CreatedAt_3: pgtype.Timestamp{Time: prevMonthEnd, Valid: true},
		MerchantID:  int32(req.MerchantID),
	})
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return res, nil
}

func (r *categoryStatsByMerchantRepository) GetYearlyTotalPricesByMerchant(ctx context.Context, req *requests.YearTotalPriceMerchant) ([]*db.GetYearlyTotalPriceByMerchantRow, error) {
	res, err := r.db.GetYearlyTotalPriceByMerchant(ctx, db.GetYearlyTotalPriceByMerchantParams{
		Column1:    int32(req.Year),
		MerchantID: int32(req.MerchantID),
	})
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return res, nil
}

func (r *categoryStatsByMerchantRepository) GetMonthPriceByMerchant(ctx context.Context, req *requests.MonthPriceMerchant) ([]*db.GetMonthlyCategoryByMerchantRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyCategoryByMerchant(ctx, db.GetMonthlyCategoryByMerchantParams{
		Column1:    yearStart,
		MerchantID: int32(req.MerchantID),
	})
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return res, nil
}

func (r *categoryStatsByMerchantRepository) GetYearPriceByMerchant(ctx context.Context, req *requests.YearPriceMerchant) ([]*db.GetYearlyCategoryByMerchantRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetYearlyCategoryByMerchant(ctx, db.GetYearlyCategoryByMerchantParams{
		Column1:    yearStart,
		MerchantID: int32(req.MerchantID),
	})
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return res, nil
}
