package repository

import (
	"context"

	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-shared/convert"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/cashier_errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/pb"
)

type cashierQueryRepository struct {
	client pb.CashierServiceClient
}

func NewCashierQueryRepository(client pb.CashierServiceClient) CashierQueryRepository {
	return &cashierQueryRepository{
		client: client,
	}
}

func (r *cashierQueryRepository) FindById(ctx context.Context, cashier_id int) (*db.Cashier, error) {
	resp, err := r.client.FindById(ctx, &pb.FindByIdCashierRequest{
		Id: int32(cashier_id),
	})
	if err != nil {
		return nil, cashier_errors.ErrFindCashierById
	}

	if resp == nil || resp.Data == nil {
		return nil, cashier_errors.ErrFindCashierById
	}

	c := resp.Data
	res := &db.Cashier{
		CashierID:  c.Id,
		MerchantID: c.MerchantId,
		Name:       c.Name,
		CreatedAt:  convert.PgTimestamp(c.CreatedAt),
		UpdatedAt:  convert.PgTimestamp(c.UpdatedAt),
	}

	return res, nil
}
