package repository

import (
	"context"

	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	orderitem_errors "github.com/MamangRust/monolith-point-of-sale-shared/errors/order_item_errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/pb"
)

type orderItemQueryRepository struct {
	client pb.OrderItemServiceClient
}

func NewOrderItemQueryRepository(client pb.OrderItemServiceClient) OrderItemQueryRepository {
	return &orderItemQueryRepository{
		client: client,
	}
}

func (r *orderItemQueryRepository) CalculateTotalPrice(ctx context.Context, orderID int) (*int32, error) {
	items, err := r.FindOrderItemByOrder(ctx, orderID)
	if err != nil {
		return nil, orderitem_errors.ErrCalculateTotalPrice
	}

	var total int32 = 0
	for _, item := range items {
		if item != nil {
			total += item.Quantity * item.Price
		}
	}

	return &total, nil
}

func (r *orderItemQueryRepository) FindOrderItemByOrder(ctx context.Context, orderID int) ([]*db.OrderItem, error) {
	resp, err := r.client.FindOrderItemByOrder(ctx, &pb.FindByIdOrderItemRequest{
		Id: int32(orderID),
	})
	if err != nil {
		return nil, orderitem_errors.ErrFindOrderItemByOrder
	}

	if resp == nil || resp.Data == nil {
		return nil, orderitem_errors.ErrFindOrderItemByOrder
	}

	var res []*db.OrderItem
	for _, item := range resp.Data {
		if item == nil {
			continue
		}
		res = append(res, &db.OrderItem{
			OrderItemID: item.Id,
			OrderID:     item.OrderId,
			ProductID:   item.ProductId,
			Quantity:    item.Quantity,
			Price:       item.Price,
			CreatedAt:    parsePgTimestamp(item.CreatedAt),
			UpdatedAt:    parsePgTimestamp(item.UpdatedAt),
		})
	}

	return res, nil
}
