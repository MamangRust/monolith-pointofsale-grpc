package repository

import (
	"context"

	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/record"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/cashier_errors"
	recordmapper "github.com/MamangRust/monolith-point-of-sale-shared/mapper/record"
)

type cashierCommandRepository struct {
	db      *db.Queries
	mapping recordmapper.CashierRecordMapping
}

func NewCashierCommandRepository(db *db.Queries, mapping recordmapper.CashierRecordMapping) *cashierCommandRepository {
	return &cashierCommandRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *cashierCommandRepository) CreateCashier(ctx context.Context, request *requests.CreateCashierRequest) (*record.CashierRecord, error) {
	req := db.CreateCashierParams{
		MerchantID: int32(request.MerchantID),
		UserID:     int32(request.UserID),
		Name:       request.Name,
	}

	cashier, err := r.db.CreateCashier(ctx, req)

	if err != nil {
		return nil, cashier_errors.ErrCreateCashier
	}

	return r.mapping.ToCashierRecord(cashier), nil
}

func (r *cashierCommandRepository) UpdateCashier(ctx context.Context, request *requests.UpdateCashierRequest) (*record.CashierRecord, error) {
	req := db.UpdateCashierParams{
		CashierID: int32(*request.CashierID),
		Name:      request.Name,
	}

	res, err := r.db.UpdateCashier(ctx, req)

	if err != nil {
		return nil, cashier_errors.ErrUpdateCashier
	}

	return r.mapping.ToCashierRecord(res), nil
}

func (r *cashierCommandRepository) TrashedCashier(ctx context.Context, cashier_id int) (*record.CashierRecord, error) {
	res, err := r.db.TrashCashier(ctx, int32(cashier_id))

	if err != nil {
		return nil, cashier_errors.ErrTrashedCashier
	}

	return r.mapping.ToCashierRecord(res), nil
}

func (r *cashierCommandRepository) RestoreCashier(ctx context.Context, cashier_id int) (*record.CashierRecord, error) {
	res, err := r.db.RestoreCashier(ctx, int32(cashier_id))

	if err != nil {
		return nil, cashier_errors.ErrRestoreCashier
	}

	return r.mapping.ToCashierRecord(res), nil
}

func (r *cashierCommandRepository) DeleteCashierPermanent(ctx context.Context, cashier_id int) (bool, error) {
	err := r.db.DeleteCashierPermanently(ctx, int32(cashier_id))

	if err != nil {
		return false, cashier_errors.ErrDeleteCashierPermanent
	}

	return true, nil
}

func (r *cashierCommandRepository) RestoreAllCashier(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllCashiers(ctx)

	if err != nil {
		return false, cashier_errors.ErrRestoreAllCashiers
	}

	return true, nil
}

func (r *cashierCommandRepository) DeleteAllCashierPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentCashiers(ctx)

	if err != nil {
		return false, cashier_errors.ErrDeleteAllCashiersPermanent
	}

	return true, nil
}
