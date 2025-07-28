package repository

import (
	"context"

	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/record"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/merchant_errors"
	recordmapper "github.com/MamangRust/monolith-point-of-sale-shared/mapper/record"
)

type merchantQueryRepository struct {
	db      *db.Queries
	mapping recordmapper.MerchantRecordMapping
}

func NewMerchantQueryRepository(db *db.Queries, mapping recordmapper.MerchantRecordMapping) *merchantQueryRepository {
	return &merchantQueryRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *merchantQueryRepository) FindById(ctx context.Context, user_id int) (*record.MerchantRecord, error) {
	res, err := r.db.GetMerchantByID(ctx, int32(user_id))

	if err != nil {
		return nil, merchant_errors.ErrFindById
	}

	return r.mapping.ToMerchantRecord(res), nil
}
