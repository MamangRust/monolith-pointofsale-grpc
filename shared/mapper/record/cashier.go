package recordmapper

import (
	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"

	"github.com/MamangRust/monolith-point-of-sale-shared/domain/record"
)

type cashierRecordMapper struct {
}

func NewCashierRecordMapper() *cashierRecordMapper {
	return &cashierRecordMapper{}
}

func (s *cashierRecordMapper) ToCashierRecord(cashier *db.Cashier) *record.CashierRecord {
	var deletedAt *string
	if cashier.DeletedAt.Valid {
		deletedAtStr := cashier.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.CashierRecord{
		ID:         int(cashier.CashierID),
		MerchantID: int(cashier.MerchantID),
		UserID:     int(cashier.UserID),
		Name:       cashier.Name,
		CreatedAt:  cashier.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:  cashier.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:  deletedAt,
	}
}

func (s *cashierRecordMapper) ToCashierRecordPagination(cashier *db.GetCashiersRow) *record.CashierRecord {
	var deletedAt *string
	if cashier.DeletedAt.Valid {
		deletedAtStr := cashier.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.CashierRecord{
		ID:         int(cashier.CashierID),
		MerchantID: int(cashier.MerchantID),
		UserID:     int(cashier.UserID),
		Name:       cashier.Name,
		CreatedAt:  cashier.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:  cashier.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:  deletedAt,
	}
}

func (s *cashierRecordMapper) ToCashiersRecordPagination(cashiers []*db.GetCashiersRow) []*record.CashierRecord {
	var result []*record.CashierRecord

	for _, cashier := range cashiers {
		result = append(result, s.ToCashierRecordPagination(cashier))
	}

	return result
}

func (s *cashierRecordMapper) ToCashierRecordActivePagination(cashier *db.GetCashiersActiveRow) *record.CashierRecord {
	var deletedAt *string
	if cashier.DeletedAt.Valid {
		deletedAtStr := cashier.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.CashierRecord{
		ID:         int(cashier.CashierID),
		MerchantID: int(cashier.MerchantID),
		UserID:     int(cashier.UserID),
		Name:       cashier.Name,
		CreatedAt:  cashier.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:  cashier.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:  deletedAt,
	}
}

func (s *cashierRecordMapper) ToCashiersRecordActivePagination(cashiers []*db.GetCashiersActiveRow) []*record.CashierRecord {
	var result []*record.CashierRecord

	for _, cashier := range cashiers {
		result = append(result, s.ToCashierRecordActivePagination(cashier))
	}

	return result
}

func (s *cashierRecordMapper) ToCashierRecordTrashedPagination(cashier *db.GetCashiersTrashedRow) *record.CashierRecord {
	var deletedAt *string
	if cashier.DeletedAt.Valid {
		deletedAtStr := cashier.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.CashierRecord{
		ID:         int(cashier.CashierID),
		MerchantID: int(cashier.MerchantID),
		UserID:     int(cashier.UserID),
		Name:       cashier.Name,
		CreatedAt:  cashier.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:  cashier.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:  deletedAt,
	}
}

func (s *cashierRecordMapper) ToCashiersRecordTrashedPagination(cashiers []*db.GetCashiersTrashedRow) []*record.CashierRecord {
	var result []*record.CashierRecord

	for _, cashier := range cashiers {
		result = append(result, s.ToCashierRecordTrashedPagination(cashier))
	}

	return result
}

func (s *cashierRecordMapper) ToCashierMerchantRecordPagination(cashier *db.GetCashiersByMerchantRow) *record.CashierRecord {
	var deletedAt *string
	if cashier.DeletedAt.Valid {
		deletedAtStr := cashier.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.CashierRecord{
		ID:         int(cashier.CashierID),
		MerchantID: int(cashier.MerchantID),
		UserID:     int(cashier.UserID),
		Name:       cashier.Name,
		CreatedAt:  cashier.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:  cashier.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:  deletedAt,
	}
}

func (s *cashierRecordMapper) ToCashiersMerchantRecordPagination(cashiers []*db.GetCashiersByMerchantRow) []*record.CashierRecord {
	var result []*record.CashierRecord

	for _, cashier := range cashiers {
		result = append(result, s.ToCashierMerchantRecordPagination(cashier))
	}

	return result
}

func (s *cashierRecordMapper) ToCashierMonthlySale(cashier *db.GetMonthlyCashierRow) *record.CashierRecordMonthSales {
	return &record.CashierRecordMonthSales{
		Month:       cashier.Month,
		CashierID:   int(cashier.CashierID),
		CashierName: cashier.CashierName,
		OrderCount:  int(cashier.OrderCount),
		TotalSales:  int(cashier.TotalSales),
	}
}

func (s *cashierRecordMapper) ToCashierMonthlySales(c []*db.GetMonthlyCashierRow) []*record.CashierRecordMonthSales {
	var cashierRecords []*record.CashierRecordMonthSales

	for _, cashier := range c {
		cashierRecords = append(cashierRecords, s.ToCashierMonthlySale(cashier))
	}

	return cashierRecords
}

func (s *cashierRecordMapper) ToCashierYearlySale(cashier *db.GetYearlyCashierRow) *record.CashierRecordYearSales {
	return &record.CashierRecordYearSales{
		Year:        cashier.Year,
		CashierID:   int(cashier.CashierID),
		CashierName: cashier.CashierName,
		OrderCount:  int(cashier.OrderCount),
		TotalSales:  int(cashier.TotalSales),
	}
}

func (s *cashierRecordMapper) ToCashierYearlySales(c []*db.GetYearlyCashierRow) []*record.CashierRecordYearSales {
	var cashierRecords []*record.CashierRecordYearSales

	for _, cashier := range c {
		cashierRecords = append(cashierRecords, s.ToCashierYearlySale(cashier))
	}

	return cashierRecords
}

func (s *cashierRecordMapper) ToCashierMonthlySaleById(cashier *db.GetMonthlyCashierByCashierIdRow) *record.CashierRecordMonthSales {
	return &record.CashierRecordMonthSales{
		Month:       cashier.Month,
		CashierID:   int(cashier.CashierID),
		CashierName: cashier.CashierName,
		OrderCount:  int(cashier.OrderCount),
		TotalSales:  int(cashier.TotalSales),
	}
}

func (s *cashierRecordMapper) ToCashierMonthlySalesById(c []*db.GetMonthlyCashierByCashierIdRow) []*record.CashierRecordMonthSales {
	var cashierRecords []*record.CashierRecordMonthSales

	for _, cashier := range c {
		cashierRecords = append(cashierRecords, s.ToCashierMonthlySaleById(cashier))
	}

	return cashierRecords
}

func (s *cashierRecordMapper) ToCashierYearlySaleById(cashier *db.GetYearlyCashierByCashierIdRow) *record.CashierRecordYearSales {
	return &record.CashierRecordYearSales{
		Year:        cashier.Year,
		CashierID:   int(cashier.CashierID),
		CashierName: cashier.CashierName,
		OrderCount:  int(cashier.OrderCount),
		TotalSales:  int(cashier.TotalSales),
	}
}

func (s *cashierRecordMapper) ToCashierYearlySalesById(c []*db.GetYearlyCashierByCashierIdRow) []*record.CashierRecordYearSales {
	var cashierRecords []*record.CashierRecordYearSales

	for _, cashier := range c {
		cashierRecords = append(cashierRecords, s.ToCashierYearlySaleById(cashier))
	}

	return cashierRecords
}

func (s *cashierRecordMapper) ToCashierMonthlySaleByMerchant(cashier *db.GetMonthlyCashierByMerchantRow) *record.CashierRecordMonthSales {
	return &record.CashierRecordMonthSales{
		Month:       cashier.Month,
		CashierID:   int(cashier.CashierID),
		CashierName: cashier.CashierName,
		OrderCount:  int(cashier.OrderCount),
		TotalSales:  int(cashier.TotalSales),
	}
}

func (s *cashierRecordMapper) ToCashierMonthlySalesByMerchant(c []*db.GetMonthlyCashierByMerchantRow) []*record.CashierRecordMonthSales {
	var cashierRecords []*record.CashierRecordMonthSales

	for _, cashier := range c {
		cashierRecords = append(cashierRecords, s.ToCashierMonthlySaleByMerchant(cashier))
	}

	return cashierRecords
}

func (s *cashierRecordMapper) ToCashierYearlySaleByMerchant(cashier *db.GetYearlyCashierByMerchantRow) *record.CashierRecordYearSales {
	return &record.CashierRecordYearSales{
		Year:        cashier.Year,
		CashierID:   int(cashier.CashierID),
		CashierName: cashier.CashierName,
		OrderCount:  int(cashier.OrderCount),
		TotalSales:  int(cashier.TotalSales),
	}
}

func (s *cashierRecordMapper) ToCashierYearlySalesByMerchant(c []*db.GetYearlyCashierByMerchantRow) []*record.CashierRecordYearSales {
	var cashierRecords []*record.CashierRecordYearSales

	for _, cashier := range c {
		cashierRecords = append(cashierRecords, s.ToCashierYearlySaleByMerchant(cashier))
	}

	return cashierRecords
}

func (s *cashierRecordMapper) ToCashierMonthlyTotalSale(c *db.GetMonthlyTotalSalesCashierRow) *record.CashierRecordMonthTotalSales {
	return &record.CashierRecordMonthTotalSales{
		Year:       c.Year,
		Month:      c.Month,
		TotalSales: int(c.TotalSales),
	}
}

func (s *cashierRecordMapper) ToCashierMonthlyTotalSales(c []*db.GetMonthlyTotalSalesCashierRow) []*record.CashierRecordMonthTotalSales {
	var cashierRecords []*record.CashierRecordMonthTotalSales

	for _, cashier := range c {
		cashierRecords = append(cashierRecords, s.ToCashierMonthlyTotalSale(cashier))
	}

	return cashierRecords
}

func (s *cashierRecordMapper) ToCashierYearlyTotalSale(c *db.GetYearlyTotalSalesCashierRow) *record.CashierRecordYearTotalSales {
	return &record.CashierRecordYearTotalSales{
		Year:       c.Year,
		TotalSales: int(c.TotalSales),
	}
}

func (s *cashierRecordMapper) ToCashierYearlyTotalSales(c []*db.GetYearlyTotalSalesCashierRow) []*record.CashierRecordYearTotalSales {
	var cashierRecords []*record.CashierRecordYearTotalSales

	for _, cashier := range c {
		cashierRecords = append(cashierRecords, s.ToCashierYearlyTotalSale(cashier))
	}

	return cashierRecords
}

func (s *cashierRecordMapper) ToCashierMonthlyTotalSaleById(c *db.GetMonthlyTotalSalesByIdRow) *record.CashierRecordMonthTotalSales {
	return &record.CashierRecordMonthTotalSales{
		Year:       c.Year,
		Month:      c.Month,
		TotalSales: int(c.TotalSales),
	}
}

func (s *cashierRecordMapper) ToCashierMonthlyTotalSalesById(c []*db.GetMonthlyTotalSalesByIdRow) []*record.CashierRecordMonthTotalSales {
	var cashierRecords []*record.CashierRecordMonthTotalSales

	for _, cashier := range c {
		cashierRecords = append(cashierRecords, s.ToCashierMonthlyTotalSaleById(cashier))
	}

	return cashierRecords
}

func (s *cashierRecordMapper) ToCashierYearlyTotalSaleById(c *db.GetYearlyTotalSalesByIdRow) *record.CashierRecordYearTotalSales {
	return &record.CashierRecordYearTotalSales{
		Year:       c.Year,
		TotalSales: int(c.TotalSales),
	}
}

func (s *cashierRecordMapper) ToCashierYearlyTotalSalesById(c []*db.GetYearlyTotalSalesByIdRow) []*record.CashierRecordYearTotalSales {
	var cashierRecords []*record.CashierRecordYearTotalSales

	for _, cashier := range c {
		cashierRecords = append(cashierRecords, s.ToCashierYearlyTotalSaleById(cashier))
	}

	return cashierRecords
}

func (s *cashierRecordMapper) ToCashierMonthlyTotalSaleByMerchant(c *db.GetMonthlyTotalSalesByMerchantRow) *record.CashierRecordMonthTotalSales {
	return &record.CashierRecordMonthTotalSales{
		Year:       c.Year,
		Month:      c.Month,
		TotalSales: int(c.TotalSales),
	}
}

func (s *cashierRecordMapper) ToCashierMonthlyTotalSalesByMerchant(c []*db.GetMonthlyTotalSalesByMerchantRow) []*record.CashierRecordMonthTotalSales {
	var cashierRecords []*record.CashierRecordMonthTotalSales

	for _, cashier := range c {
		cashierRecords = append(cashierRecords, s.ToCashierMonthlyTotalSaleByMerchant(cashier))
	}

	return cashierRecords
}

func (s *cashierRecordMapper) ToCashierYearlyTotalSaleByMerchant(c *db.GetYearlyTotalSalesByMerchantRow) *record.CashierRecordYearTotalSales {
	return &record.CashierRecordYearTotalSales{
		Year:       c.Year,
		TotalSales: int(c.TotalSales),
	}
}

func (s *cashierRecordMapper) ToCashierYearlyTotalSalesByMerchant(c []*db.GetYearlyTotalSalesByMerchantRow) []*record.CashierRecordYearTotalSales {
	var cashierRecords []*record.CashierRecordYearTotalSales

	for _, cashier := range c {
		cashierRecords = append(cashierRecords, s.ToCashierYearlyTotalSaleByMerchant(cashier))
	}

	return cashierRecords
}
