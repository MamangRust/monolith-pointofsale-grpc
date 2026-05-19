package repository

import (
	"context"

	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/order_errors"
)

type orderCommandRepository struct {
	db *db.Queries
}

func NewOrderCommandRepository(db *db.Queries) OrderCommandRepository {
	return &orderCommandRepository{
		db: db,
	}
}

func (r *orderCommandRepository) CreateOrder(ctx context.Context, request *requests.CreateOrderRecordRequest) (*db.Order, error) {
	req := db.CreateOrderParams{
		MerchantID: int32(request.MerchantID),
		CashierID:  int32(request.CashierID),
		TotalPrice: int64(request.TotalPrice),
	}

	user, err := r.db.CreateOrder(ctx, req)
	if err != nil {
		return nil, order_errors.ErrCreateOrder
	}

	return user, nil
}

func (r *orderCommandRepository) UpdateOrder(ctx context.Context, request *requests.UpdateOrderRecordRequest) (*db.Order, error) {
	req := db.UpdateOrderParams{
		OrderID:    int32(request.OrderID),
		TotalPrice: int64(request.TotalPrice),
	}

	res, err := r.db.UpdateOrder(ctx, req)
	if err != nil {
		return nil, order_errors.ErrUpdateOrder
	}

	return res, nil
}

func (r *orderCommandRepository) TrashedOrder(ctx context.Context, orderID int) (*db.Order, error) {
	res, err := r.db.TrashedOrder(ctx, int32(orderID))
	if err != nil {
		return nil, order_errors.ErrTrashedOrder
	}

	return res, nil
}

func (r *orderCommandRepository) RestoreOrder(ctx context.Context, orderID int) (*db.Order, error) {
	res, err := r.db.RestoreOrder(ctx, int32(orderID))
	if err != nil {
		return nil, order_errors.ErrRestoreOrder
	}

	return res, nil
}

func (r *orderCommandRepository) DeleteOrderPermanent(ctx context.Context, orderID int) (bool, error) {
	err := r.db.DeleteOrderPermanently(ctx, int32(orderID))
	if err != nil {
		return false, order_errors.ErrDeleteOrderPermanent
	}

	return true, nil
}

func (r *orderCommandRepository) RestoreAllOrder(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllOrders(ctx)
	if err != nil {
		return false, order_errors.ErrRestoreAllOrder
	}
	return true, nil
}

func (r *orderCommandRepository) DeleteAllOrderPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentOrders(ctx)
	if err != nil {
		return false, order_errors.ErrDeleteAllOrderPermanent
	}
	return true, nil
}
