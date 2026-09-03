package repository

import (
	"context"

	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-shared/convert"
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

func (r *orderItemQueryRepository) FindOrderItemByOrder(ctx context.Context, order_id int) ([]*db.OrderItem, error) {
	resp, err := r.client.FindOrderItemByOrder(ctx, &pb.FindByIdOrderItemRequest{
		Id: int32(order_id),
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
			CreatedAt:   convert.PgTimestamp(item.CreatedAt),
			UpdatedAt:   convert.PgTimestamp(item.UpdatedAt),
		})
	}

	return res, nil
}
