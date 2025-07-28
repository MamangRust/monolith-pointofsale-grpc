package repository

import (
	"context"

	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/record"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/order_errors"
	recordmapper "github.com/MamangRust/monolith-point-of-sale-shared/mapper/record"
)

type orderQueryRepository struct {
	db      *db.Queries
	mapping recordmapper.OrderRecordMapping
}

func NewOrderQueryRepository(db *db.Queries, mapping recordmapper.OrderRecordMapping) *orderQueryRepository {
	return &orderQueryRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r orderQueryRepository) FindById(ctx context.Context, order_id int) (*record.OrderRecord, error) {
	res, err := r.db.GetOrderByID(ctx, int32(order_id))

	if err != nil {
		return nil, order_errors.ErrFindById
	}

	return r.mapping.ToOrderRecord(res), nil
}
