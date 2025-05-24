package repository

import (
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/record"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type OrderItemQueryRepository interface {
	FindAllOrderItems(req *requests.FindAllOrderItems) ([]*record.OrderItemRecord, *int, error)
	FindByActive(req *requests.FindAllOrderItems) ([]*record.OrderItemRecord, *int, error)
	FindByTrashed(req *requests.FindAllOrderItems) ([]*record.OrderItemRecord, *int, error)
	FindOrderItemByOrder(order_id int) ([]*record.OrderItemRecord, error)
}
