package repository

import (
	"context"

	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/record"
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

func (r *cashierQueryRepository) FindById(cashier_id int) (*record.CashierRecord, error) {
	res, err := r.db.GetCashierById(r.ctx, int32(cashier_id))

	if err != nil {
		return nil, cashier_errors.ErrFindCashierById
	}

	return r.mapping.ToCashierRecord(res), nil
}
