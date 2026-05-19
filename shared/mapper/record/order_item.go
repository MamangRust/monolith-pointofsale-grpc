package recordmapper

import (
	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"

	"github.com/MamangRust/monolith-point-of-sale-shared/domain/record"
)

type orderItemRecordMapper struct {
}

func NewOrderItemRecordMapper() *orderItemRecordMapper {
	return &orderItemRecordMapper{}
}

func (s *orderItemRecordMapper) ToOrderItemRecord(orderItems *db.OrderItem) *record.OrderItemRecord {
	var deletedAt *string
	if orderItems.DeletedAt.Valid {
		deletedAtStr := orderItems.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.OrderItemRecord{
		ID:        int(orderItems.OrderItemID),
		OrderID:   int(orderItems.OrderID),
		ProductID: int(orderItems.ProductID),
		Quantity:  int(orderItems.Quantity),
		Price:     int(orderItems.Price),
		CreatedAt: orderItems.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: orderItems.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt: deletedAt,
	}
}

func (s *orderItemRecordMapper) ToOrderItemsRecord(orders []*db.OrderItem) []*record.OrderItemRecord {
	var result []*record.OrderItemRecord

	for _, Merchant := range orders {
		result = append(result, s.ToOrderItemRecord(Merchant))
	}

	return result
}

func (s *orderItemRecordMapper) ToOrderItemRecordPagination(orderItems *db.GetOrderItemsRow) *record.OrderItemRecord {
	var deletedAt *string
	if orderItems.DeletedAt.Valid {
		deletedAtStr := orderItems.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.OrderItemRecord{
		ID:        int(orderItems.OrderItemID),
		OrderID:   int(orderItems.OrderID),
		ProductID: int(orderItems.ProductID),
		Quantity:  int(orderItems.Quantity),
		Price:     int(orderItems.Price),
		CreatedAt: orderItems.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: orderItems.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt: deletedAt,
	}
}

func (s *orderItemRecordMapper) ToOrderItemsRecordPagination(orders []*db.GetOrderItemsRow) []*record.OrderItemRecord {
	var result []*record.OrderItemRecord

	for _, Merchant := range orders {
		result = append(result, s.ToOrderItemRecordPagination(Merchant))
	}

	return result
}

func (s *orderItemRecordMapper) ToOrderItemRecordActivePagination(orderItems *db.GetOrderItemsActiveRow) *record.OrderItemRecord {
	var deletedAt *string
	if orderItems.DeletedAt.Valid {
		deletedAtStr := orderItems.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.OrderItemRecord{
		ID:        int(orderItems.OrderItemID),
		OrderID:   int(orderItems.OrderID),
		ProductID: int(orderItems.ProductID),
		Quantity:  int(orderItems.Quantity),
		Price:     int(orderItems.Price),
		CreatedAt: orderItems.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: orderItems.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt: deletedAt,
	}
}

func (s *orderItemRecordMapper) ToOrderItemsRecordActivePagination(orders []*db.GetOrderItemsActiveRow) []*record.OrderItemRecord {
	var result []*record.OrderItemRecord

	for _, Merchant := range orders {
		result = append(result, s.ToOrderItemRecordActivePagination(Merchant))
	}

	return result
}

func (s *orderItemRecordMapper) ToOrderItemRecordTrashedPagination(orderItems *db.GetOrderItemsTrashedRow) *record.OrderItemRecord {
	var deletedAt *string
	if orderItems.DeletedAt.Valid {
		deletedAtStr := orderItems.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.OrderItemRecord{
		ID:        int(orderItems.OrderItemID),
		OrderID:   int(orderItems.OrderID),
		ProductID: int(orderItems.ProductID),
		Quantity:  int(orderItems.Quantity),
		Price:     int(orderItems.Price),
		CreatedAt: orderItems.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: orderItems.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt: deletedAt,
	}
}

func (s *orderItemRecordMapper) ToOrderItemsRecordTrashedPagination(orders []*db.GetOrderItemsTrashedRow) []*record.OrderItemRecord {
	var result []*record.OrderItemRecord

	for _, Merchant := range orders {
		result = append(result, s.ToOrderItemRecordTrashedPagination(Merchant))
	}

	return result
}
