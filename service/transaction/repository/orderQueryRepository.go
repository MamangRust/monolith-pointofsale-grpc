package repository

import (
	"context"

	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-shared/convert"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/order_errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/pb"
	"github.com/jackc/pgx/v5/pgtype"
)

type orderQueryRepository struct {
	client pb.OrderServiceClient
}

func NewOrderQueryRepository(client pb.OrderServiceClient) OrderQueryRepository {
	return &orderQueryRepository{
		client: client,
	}
}

func (r *orderQueryRepository) FindById(ctx context.Context, order_id int) (*db.Order, error) {
	resp, err := r.client.FindById(ctx, &pb.FindByIdOrderRequest{
		Id: int32(order_id),
	})
	if err != nil {
		return nil, order_errors.ErrFindById
	}

	if resp == nil || resp.Data == nil {
		return nil, order_errors.ErrFindById
	}

	o := resp.Data
	res := &db.Order{
		OrderID:    o.Id,
		MerchantID: o.MerchantId,
		CashierID:  o.CashierId,
		TotalPrice: int64(o.TotalPrice),
		CreatedAt:  convert.PgTimestamp(o.CreatedAt),
		UpdatedAt:  convert.PgTimestamp(o.UpdatedAt),
		DeletedAt:  pgtype.Timestamp{},
	}

	return res, nil
}
