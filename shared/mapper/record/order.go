package recordmapper

import (
	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"

	"github.com/MamangRust/monolith-point-of-sale-shared/domain/record"
)

type orderRecordMapper struct {
}

func NewOrderRecordMapper() *orderRecordMapper {
	return &orderRecordMapper{}
}

func (s *orderRecordMapper) ToOrderRecord(order *db.Order) *record.OrderRecord {
	var deletedAt *string
	if order.DeletedAt.Valid {
		deletedAtStr := order.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.OrderRecord{
		ID:         int(order.OrderID),
		MerchantID: int(order.MerchantID),
		CashierID:  int(order.CashierID),
		TotalPrice: int(order.TotalPrice),
		CreatedAt:  order.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:  order.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:  deletedAt,
	}
}

func (s *orderRecordMapper) ToOrdersRecord(orders []*db.Order) []*record.OrderRecord {
	var result []*record.OrderRecord

	for _, order := range orders {
		result = append(result, s.ToOrderRecord(order))
	}

	return result
}

func (s *orderRecordMapper) ToOrderRecordPagination(order *db.GetOrdersRow) *record.OrderRecord {
	var deletedAt *string
	if order.DeletedAt.Valid {
		deletedAtStr := order.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.OrderRecord{
		ID:         int(order.OrderID),
		MerchantID: int(order.MerchantID),
		CashierID:  int(order.CashierID),
		TotalPrice: int(order.TotalPrice),
		CreatedAt:  order.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:  order.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:  deletedAt,
	}
}

func (s *orderRecordMapper) ToOrdersRecordPagination(orders []*db.GetOrdersRow) []*record.OrderRecord {
	var result []*record.OrderRecord

	for _, order := range orders {
		result = append(result, s.ToOrderRecordPagination(order))
	}

	return result
}

func (s *orderRecordMapper) ToOrderRecordActivePagination(order *db.GetOrdersActiveRow) *record.OrderRecord {
	var deletedAt *string
	if order.DeletedAt.Valid {
		deletedAtStr := order.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.OrderRecord{
		ID:         int(order.OrderID),
		MerchantID: int(order.MerchantID),
		CashierID:  int(order.CashierID),
		TotalPrice: int(order.TotalPrice),
		CreatedAt:  order.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:  order.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:  deletedAt,
	}
}

func (s *orderRecordMapper) ToOrdersRecordActivePagination(orders []*db.GetOrdersActiveRow) []*record.OrderRecord {
	var result []*record.OrderRecord

	for _, order := range orders {
		result = append(result, s.ToOrderRecordActivePagination(order))
	}

	return result
}

func (s *orderRecordMapper) ToOrderRecordTrashedPagination(order *db.GetOrdersTrashedRow) *record.OrderRecord {
	var deletedAt *string
	if order.DeletedAt.Valid {
		deletedAtStr := order.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.OrderRecord{
		ID:         int(order.OrderID),
		MerchantID: int(order.MerchantID),
		CashierID:  int(order.CashierID),
		TotalPrice: int(order.TotalPrice),
		CreatedAt:  order.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:  order.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:  deletedAt,
	}
}

func (s *orderRecordMapper) ToOrdersRecordTrashedPagination(orders []*db.GetOrdersTrashedRow) []*record.OrderRecord {
	var result []*record.OrderRecord

	for _, order := range orders {
		result = append(result, s.ToOrderRecordTrashedPagination(order))
	}

	return result
}

func (s *orderRecordMapper) ToOrderRecordByMerchantPagination(order *db.GetOrdersByMerchantRow) *record.OrderRecord {
	var deletedAt *string
	if order.DeletedAt.Valid {
		deletedAtStr := order.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.OrderRecord{
		ID:         int(order.OrderID),
		MerchantID: int(order.MerchantID),
		CashierID:  int(order.CashierID),
		TotalPrice: int(order.TotalPrice),
		CreatedAt:  order.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:  order.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:  deletedAt,
	}
}

func (s *orderRecordMapper) ToOrdersRecordByMerchantPagination(orders []*db.GetOrdersByMerchantRow) []*record.OrderRecord {
	var result []*record.OrderRecord

	for _, order := range orders {
		result = append(result, s.ToOrderRecordByMerchantPagination(order))
	}

	return result
}

func (s *orderRecordMapper) ToOrderMonthlyPrice(order *db.GetMonthlyOrderRow) *record.OrderMonthlyRecord {
	return &record.OrderMonthlyRecord{
		Month:          order.Month,
		OrderCount:     int(order.OrderCount),
		TotalRevenue:   int(order.TotalRevenue),
		TotalItemsSold: int(order.TotalItemsSold),
	}
}

func (s *orderRecordMapper) ToOrderMonthlyPrices(c []*db.GetMonthlyOrderRow) []*record.OrderMonthlyRecord {
	var orderRecords []*record.OrderMonthlyRecord

	for _, order := range c {
		orderRecords = append(orderRecords, s.ToOrderMonthlyPrice(order))
	}

	return orderRecords
}

func (s *orderRecordMapper) ToOrderYearlyPrice(order *db.GetYearlyOrderRow) *record.OrderYearlyRecord {
	return &record.OrderYearlyRecord{
		Year:               order.Year,
		OrderCount:         int(order.OrderCount),
		TotalRevenue:       int(order.TotalRevenue),
		TotalItemsSold:     int(order.TotalItemsSold),
		ActiveCashiers:     int(order.ActiveCashiers),
		UniqueProductsSold: int(order.UniqueProductsSold),
	}
}

func (s *orderRecordMapper) ToOrderYearlyPrices(c []*db.GetYearlyOrderRow) []*record.OrderYearlyRecord {
	var orderRecords []*record.OrderYearlyRecord

	for _, order := range c {
		orderRecords = append(orderRecords, s.ToOrderYearlyPrice(order))
	}

	return orderRecords
}

func (s *orderRecordMapper) ToOrderMonthlyPriceByMerchant(order *db.GetMonthlyOrderByMerchantRow) *record.OrderMonthlyRecord {

	return &record.OrderMonthlyRecord{
		Month:          order.Month,
		OrderCount:     int(order.OrderCount),
		TotalRevenue:   int(order.TotalRevenue),
		TotalItemsSold: int(order.TotalItemsSold),
	}
}

func (s *orderRecordMapper) ToOrderMonthlyPricesByMerchant(c []*db.GetMonthlyOrderByMerchantRow) []*record.OrderMonthlyRecord {
	var orderRecords []*record.OrderMonthlyRecord

	for _, order := range c {
		orderRecords = append(orderRecords, s.ToOrderMonthlyPriceByMerchant(order))
	}

	return orderRecords
}

func (s *orderRecordMapper) ToOrderYearlyPriceByMerchant(order *db.GetYearlyOrderByMerchantRow) *record.OrderYearlyRecord {
	return &record.OrderYearlyRecord{
		Year:               order.Year,
		OrderCount:         int(order.OrderCount),
		TotalRevenue:       int(order.TotalRevenue),
		TotalItemsSold:     int(order.TotalItemsSold),
		ActiveCashiers:     int(order.ActiveCashiers),
		UniqueProductsSold: int(order.UniqueProductsSold),
	}
}

func (s *orderRecordMapper) ToOrderYearlyPricesByMerchant(c []*db.GetYearlyOrderByMerchantRow) []*record.OrderYearlyRecord {
	var orderRecords []*record.OrderYearlyRecord

	for _, order := range c {
		orderRecords = append(orderRecords, s.ToOrderYearlyPriceByMerchant(order))
	}

	return orderRecords
}

func (s *orderRecordMapper) ToOrderMonthlyTotalRevenue(c *db.GetMonthlyTotalRevenueRow) *record.OrderMonthlyTotalRevenueRecord {
	return &record.OrderMonthlyTotalRevenueRecord{
		Year:         c.Year,
		Month:        c.Month,
		TotalRevenue: int(c.TotalRevenue),
	}
}

func (s *orderRecordMapper) ToOrderMonthlyTotalRevenues(c []*db.GetMonthlyTotalRevenueRow) []*record.OrderMonthlyTotalRevenueRecord {
	var orderRecords []*record.OrderMonthlyTotalRevenueRecord

	for _, row := range c {
		orderRecords = append(orderRecords, s.ToOrderMonthlyTotalRevenue(row))
	}

	return orderRecords
}

func (s *orderRecordMapper) ToOrderYearlyTotalRevenue(c *db.GetYearlyTotalRevenueRow) *record.OrderYearlyTotalRevenueRecord {
	return &record.OrderYearlyTotalRevenueRecord{
		Year:         c.Year,
		TotalRevenue: int(c.TotalRevenue),
	}
}

func (s *orderRecordMapper) ToOrderYearlyTotalRevenues(c []*db.GetYearlyTotalRevenueRow) []*record.OrderYearlyTotalRevenueRecord {
	var orderRecords []*record.OrderYearlyTotalRevenueRecord

	for _, row := range c {
		orderRecords = append(orderRecords, s.ToOrderYearlyTotalRevenue(row))
	}

	return orderRecords
}

func (s *orderRecordMapper) ToOrderMonthlyTotalRevenueById(c *db.GetMonthlyTotalRevenueByIdRow) *record.OrderMonthlyTotalRevenueRecord {
	return &record.OrderMonthlyTotalRevenueRecord{
		Year:         c.Year,
		Month:        c.Month,
		TotalRevenue: int(c.TotalRevenue),
	}
}

func (s *orderRecordMapper) ToOrderMonthlyTotalRevenuesById(c []*db.GetMonthlyTotalRevenueByIdRow) []*record.OrderMonthlyTotalRevenueRecord {
	var orderRecords []*record.OrderMonthlyTotalRevenueRecord

	for _, row := range c {
		orderRecords = append(orderRecords, s.ToOrderMonthlyTotalRevenueById(row))
	}

	return orderRecords
}

func (s *orderRecordMapper) ToOrderYearlyTotalRevenueById(c *db.GetYearlyTotalRevenueByIdRow) *record.OrderYearlyTotalRevenueRecord {
	return &record.OrderYearlyTotalRevenueRecord{
		Year:         c.Year,
		TotalRevenue: int(c.TotalRevenue),
	}
}

func (s *orderRecordMapper) ToOrderYearlyTotalRevenuesById(c []*db.GetYearlyTotalRevenueByIdRow) []*record.OrderYearlyTotalRevenueRecord {
	var orderRecords []*record.OrderYearlyTotalRevenueRecord

	for _, row := range c {
		orderRecords = append(orderRecords, s.ToOrderYearlyTotalRevenueById(row))
	}

	return orderRecords
}

func (s *orderRecordMapper) ToOrderMonthlyTotalRevenueByMerchant(c *db.GetMonthlyTotalRevenueByMerchantRow) *record.OrderMonthlyTotalRevenueRecord {
	return &record.OrderMonthlyTotalRevenueRecord{
		Year:         c.Year,
		Month:        c.Month,
		TotalRevenue: int(c.TotalRevenue),
	}
}

func (s *orderRecordMapper) ToOrderMonthlyTotalRevenuesByMerchant(c []*db.GetMonthlyTotalRevenueByMerchantRow) []*record.OrderMonthlyTotalRevenueRecord {
	var orderRecords []*record.OrderMonthlyTotalRevenueRecord

	for _, row := range c {
		orderRecords = append(orderRecords, s.ToOrderMonthlyTotalRevenueByMerchant(row))
	}

	return orderRecords
}

func (s *orderRecordMapper) ToOrderYearlyTotalRevenueByMerchant(c *db.GetYearlyTotalRevenueByMerchantRow) *record.OrderYearlyTotalRevenueRecord {
	return &record.OrderYearlyTotalRevenueRecord{
		Year:         c.Year,
		TotalRevenue: int(c.TotalRevenue),
	}
}

func (s *orderRecordMapper) ToOrderYearlyTotalRevenuesByMerchant(c []*db.GetYearlyTotalRevenueByMerchantRow) []*record.OrderYearlyTotalRevenueRecord {
	var orderRecords []*record.OrderYearlyTotalRevenueRecord

	for _, row := range c {
		orderRecords = append(orderRecords, s.ToOrderYearlyTotalRevenueByMerchant(row))
	}

	return orderRecords
}
