package repository

import (
	"context"

	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/record"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/cashier_errors"
	recordmapper "github.com/MamangRust/monolith-point-of-sale-shared/mapper/record"
)

type cashierQueryRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.CashierRecordMapping
}

func NewCashierQueryRepository(db *db.Queries, ctx context.Context, mapping recordmapper.CashierRecordMapping) *cashierQueryRepository {
	return &cashierQueryRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *cashierQueryRepository) FindAllCashiers(req *requests.FindAllCashiers) ([]*record.CashierRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetCashiersParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetCashiers(r.ctx, reqDb)

	if err != nil {
		return nil, nil, cashier_errors.ErrFindAllCashiers
	}

	var totalCount int

	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToCashiersRecordPagination(res), &totalCount, nil
}

func (r *cashierQueryRepository) FindByActive(req *requests.FindAllCashiers) ([]*record.CashierRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetCashiersActiveParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetCashiersActive(r.ctx, reqDb)

	if err != nil {
		return nil, nil, cashier_errors.ErrFindActiveCashiers
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToCashiersRecordActivePagination(res), &totalCount, nil
}

func (r *cashierQueryRepository) FindByTrashed(req *requests.FindAllCashiers) ([]*record.CashierRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetCashiersTrashedParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetCashiersTrashed(r.ctx, reqDb)

	if err != nil {
		return nil, nil, cashier_errors.ErrFindTrashedCashiers
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToCashiersRecordTrashedPagination(res), &totalCount, nil
}

func (r *cashierQueryRepository) FindByMerchant(req *requests.FindAllCashierMerchant) ([]*record.CashierRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetCashiersByMerchantParams{
		MerchantID: int32(req.MerchantID),
		Column2:    req.Search,
		Limit:      int32(req.PageSize),
		Offset:     int32(offset),
	}

	res, err := r.db.GetCashiersByMerchant(r.ctx, reqDb)

	if err != nil {
		return nil, nil, cashier_errors.ErrFindCashiersByMerchant
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToCashiersMerchantRecordPagination(res), &totalCount, nil
}

func (r *cashierQueryRepository) FindById(cashier_id int) (*record.CashierRecord, error) {
	res, err := r.db.GetCashierById(r.ctx, int32(cashier_id))

	if err != nil {
		return nil, cashier_errors.ErrFindCashierById
	}

	return r.mapping.ToCashierRecord(res), nil
}
