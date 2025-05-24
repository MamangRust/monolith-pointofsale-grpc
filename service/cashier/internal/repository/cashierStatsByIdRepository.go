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

type cashierStatsByIdRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.CashierRecordMapping
}

func NewCashierStatsByIdRepository(db *db.Queries, ctx context.Context, mapping recordmapper.CashierRecordMapping) *cashierStatsByIdRepository {
	return &cashierStatsByIdRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *cashierStatsByIdRepository) GetMonthlyTotalSalesById(req *requests.MonthTotalSalesCashier) ([]*record.CashierRecordMonthTotalSales, error) {
	currentMonthStart := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	currentMonthEnd := currentMonthStart.AddDate(0, 1, -1)

	prevMonthStart := currentMonthStart.AddDate(0, -1, 0)
	prevMonthEnd := prevMonthStart.AddDate(0, 1, -1)

	res, err := r.db.GetMonthlyTotalSalesById(r.ctx, db.GetMonthlyTotalSalesByIdParams{
		Extract:     currentMonthStart,
		CreatedAt:   sql.NullTime{Time: currentMonthEnd, Valid: true},
		CreatedAt_2: sql.NullTime{Time: prevMonthStart, Valid: true},
		CreatedAt_3: sql.NullTime{Time: prevMonthEnd, Valid: true},
		CashierID:   int32(req.CashierID),
	})

	if err != nil {
		return nil, cashier_errors.ErrGetMonthlyTotalSalesById
	}

	so := r.mapping.ToCashierMonthlyTotalSalesById(res)

	return so, nil
}

func (r *cashierStatsByIdRepository) GetYearlyTotalSalesById(req *requests.YearTotalSalesCashier) ([]*record.CashierRecordYearTotalSales, error) {
	res, err := r.db.GetYearlyTotalSalesById(r.ctx, db.GetYearlyTotalSalesByIdParams{
		Column1:   int32(req.Year),
		CashierID: int32(req.CashierID),
	})

	if err != nil {
		return nil, cashier_errors.ErrGetYearlyTotalSalesById
	}

	so := r.mapping.ToCashierYearlyTotalSalesById(res)

	return so, nil
}

func (r *cashierStatsByIdRepository) GetMonthlyCashierById(req *requests.MonthCashierId) ([]*record.CashierRecordMonthSales, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetMonthlyCashierByCashierId(r.ctx, db.GetMonthlyCashierByCashierIdParams{
		Column1:   yearStart,
		CashierID: int32(req.Year),
	})

	if err != nil {
		return nil, cashier_errors.ErrGetMonthlyCashierById
	}

	return r.mapping.ToCashierMonthlySalesById(res), nil
}

func (r *cashierStatsByIdRepository) GetYearlyCashierById(req *requests.YearCashierId) ([]*record.CashierRecordYearSales, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	res, err := r.db.GetYearlyCashierByCashierId(r.ctx, db.GetYearlyCashierByCashierIdParams{
		Column1:   yearStart,
		CashierID: int32(req.CashierID),
	})

	if err != nil {
		return nil, cashier_errors.ErrGetYearlyCashierById
	}

	return r.mapping.ToCashierYearlySalesById(res), nil
}
