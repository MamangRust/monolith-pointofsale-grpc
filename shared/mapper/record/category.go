package recordmapper

import (
	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"

	"github.com/MamangRust/monolith-point-of-sale-shared/domain/record"
)

type categoryRecordMapper struct {
}

func NewCategoryRecordMapper() *categoryRecordMapper {
	return &categoryRecordMapper{}
}

func (s *categoryRecordMapper) ToCategoryRecord(category *db.Category) *record.CategoriesRecord {
	var deletedAt *string
	if category.DeletedAt.Valid {
		deletedAtStr := category.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.CategoriesRecord{
		ID:           int(category.CategoryID),
		Name:         category.Name,
		Description:  category.Description.String,
		SlugCategory: category.SlugCategory.String,
		CreatedAt:    category.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:    category.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:    deletedAt,
	}
}

func (s *categoryRecordMapper) ToCategoryRecordPagination(category *db.GetCategoriesRow) *record.CategoriesRecord {
	var deletedAt *string
	if category.DeletedAt.Valid {
		deletedAtStr := category.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.CategoriesRecord{
		ID:           int(category.CategoryID),
		Name:         category.Name,
		Description:  category.Description.String,
		SlugCategory: category.SlugCategory.String,
		CreatedAt:    category.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:    category.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:    deletedAt,
	}
}

func (s *categoryRecordMapper) ToCategoriesRecordPagination(categories []*db.GetCategoriesRow) []*record.CategoriesRecord {
	var result []*record.CategoriesRecord

	for _, category := range categories {
		result = append(result, s.ToCategoryRecordPagination(category))
	}

	return result
}

func (s *categoryRecordMapper) ToCategoryRecordActivePagination(category *db.GetCategoriesActiveRow) *record.CategoriesRecord {
	var deletedAt *string
	if category.DeletedAt.Valid {
		formatted := category.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &formatted
	} else {
		deletedAt = nil
	}

	return &record.CategoriesRecord{
		ID:           int(category.CategoryID),
		Name:         category.Name,
		Description:  category.Description.String,
		SlugCategory: category.SlugCategory.String,
		CreatedAt:    category.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:    category.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:    deletedAt,
	}
}

func (s *categoryRecordMapper) ToCategoriesRecordActivePagination(categories []*db.GetCategoriesActiveRow) []*record.CategoriesRecord {
	var result []*record.CategoriesRecord

	for _, category := range categories {
		result = append(result, s.ToCategoryRecordActivePagination(category))
	}

	return result
}

func (s *categoryRecordMapper) ToCategoryRecordTrashedPagination(category *db.GetCategoriesTrashedRow) *record.CategoriesRecord {
	var deletedAt *string
	if category.DeletedAt.Valid {
		deletedAtStr := category.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.CategoriesRecord{
		ID:           int(category.CategoryID),
		Name:         category.Name,
		Description:  category.Description.String,
		SlugCategory: category.SlugCategory.String,
		CreatedAt:    category.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:    category.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:    deletedAt,
	}
}

func (s *categoryRecordMapper) ToCategoriesRecordTrashedPagination(categories []*db.GetCategoriesTrashedRow) []*record.CategoriesRecord {
	var result []*record.CategoriesRecord

	for _, category := range categories {
		result = append(result, s.ToCategoryRecordTrashedPagination(category))
	}

	return result
}

func (s *categoryRecordMapper) ToCategoryMonthlyPrice(category *db.GetMonthlyCategoryRow) *record.CategoriesMonthPriceRecord {
	return &record.CategoriesMonthPriceRecord{
		Month:        category.Month,
		CategoryID:   int(category.CategoryID),
		CategoryName: category.CategoryName,
		OrderCount:   int(category.OrderCount),
		ItemsSold:    int(category.ItemsSold),
		TotalRevenue: int(category.TotalRevenue),
	}
}

func (s *categoryRecordMapper) ToCategoryMonthlyPrices(c []*db.GetMonthlyCategoryRow) []*record.CategoriesMonthPriceRecord {
	var categoryRecords []*record.CategoriesMonthPriceRecord

	for _, category := range c {
		categoryRecords = append(categoryRecords, s.ToCategoryMonthlyPrice(category))
	}

	return categoryRecords
}

func (s *categoryRecordMapper) ToCategoryYearlyPrice(category *db.GetYearlyCategoryRow) *record.CategoriesYearPriceRecord {
	return &record.CategoriesYearPriceRecord{
		Year:               category.Year,
		CategoryID:         int(category.CategoryID),
		CategoryName:       category.CategoryName,
		OrderCount:         int(category.OrderCount),
		ItemsSold:          int(category.ItemsSold),
		TotalRevenue:       int(category.TotalRevenue),
		UniqueProductsSold: int(category.UniqueProductsSold),
	}
}

func (s *categoryRecordMapper) ToCategoryYearlyPrices(c []*db.GetYearlyCategoryRow) []*record.CategoriesYearPriceRecord {
	var categoryRecords []*record.CategoriesYearPriceRecord

	for _, category := range c {
		categoryRecords = append(categoryRecords, s.ToCategoryYearlyPrice(category))
	}

	return categoryRecords
}

func (s *categoryRecordMapper) ToCategoryMonthlyPriceById(category *db.GetMonthlyCategoryByIdRow) *record.CategoriesMonthPriceRecord {
	return &record.CategoriesMonthPriceRecord{
		Month:        category.Month,
		CategoryID:   int(category.CategoryID),
		CategoryName: category.CategoryName,
		OrderCount:   int(category.OrderCount),
		ItemsSold:    int(category.ItemsSold),
		TotalRevenue: int(category.TotalRevenue),
	}
}

func (s *categoryRecordMapper) ToCategoryMonthlyPricesById(c []*db.GetMonthlyCategoryByIdRow) []*record.CategoriesMonthPriceRecord {
	var categoryRecords []*record.CategoriesMonthPriceRecord

	for _, category := range c {
		categoryRecords = append(categoryRecords, s.ToCategoryMonthlyPriceById(category))
	}

	return categoryRecords
}

func (s *categoryRecordMapper) ToCategoryYearlyPriceById(category *db.GetYearlyCategoryByIdRow) *record.CategoriesYearPriceRecord {
	return &record.CategoriesYearPriceRecord{
		Year:               category.Year,
		CategoryID:         int(category.CategoryID),
		CategoryName:       category.CategoryName,
		OrderCount:         int(category.OrderCount),
		ItemsSold:          int(category.ItemsSold),
		TotalRevenue:       int(category.TotalRevenue),
		UniqueProductsSold: int(category.UniqueProductsSold),
	}
}

func (s *categoryRecordMapper) ToCategoryYearlyPricesById(c []*db.GetYearlyCategoryByIdRow) []*record.CategoriesYearPriceRecord {
	var categoryRecords []*record.CategoriesYearPriceRecord

	for _, category := range c {
		categoryRecords = append(categoryRecords, s.ToCategoryYearlyPriceById(category))
	}

	return categoryRecords
}

func (s *categoryRecordMapper) ToCategoryMonthlyPriceByMerchant(category *db.GetMonthlyCategoryByMerchantRow) *record.CategoriesMonthPriceRecord {
	return &record.CategoriesMonthPriceRecord{
		Month:        category.Month,
		CategoryID:   int(category.CategoryID),
		CategoryName: category.CategoryName,
		OrderCount:   int(category.OrderCount),
		ItemsSold:    int(category.ItemsSold),
		TotalRevenue: int(category.TotalRevenue),
	}
}

func (s *categoryRecordMapper) ToCategoryMonthlyPricesByMerchant(c []*db.GetMonthlyCategoryByMerchantRow) []*record.CategoriesMonthPriceRecord {
	var categoryRecords []*record.CategoriesMonthPriceRecord

	for _, category := range c {
		categoryRecords = append(categoryRecords, s.ToCategoryMonthlyPriceByMerchant(category))
	}

	return categoryRecords
}

func (s *categoryRecordMapper) ToCategoryYearlyPriceByMerchant(category *db.GetYearlyCategoryByMerchantRow) *record.CategoriesYearPriceRecord {
	return &record.CategoriesYearPriceRecord{
		Year:               category.Year,
		CategoryID:         int(category.CategoryID),
		CategoryName:       category.CategoryName,
		OrderCount:         int(category.OrderCount),
		ItemsSold:          int(category.ItemsSold),
		TotalRevenue:       int(category.TotalRevenue),
		UniqueProductsSold: int(category.UniqueProductsSold),
	}
}

func (s *categoryRecordMapper) ToCategoryYearlyPricesByMerchant(c []*db.GetYearlyCategoryByMerchantRow) []*record.CategoriesYearPriceRecord {
	var categoryRecords []*record.CategoriesYearPriceRecord

	for _, category := range c {
		categoryRecords = append(categoryRecords, s.ToCategoryYearlyPriceByMerchant(category))
	}

	return categoryRecords
}

func (s *categoryRecordMapper) ToCategoryMonthlyTotalPrice(c *db.GetMonthlyTotalPriceRow) *record.CategoriesMonthlyTotalPriceRecord {
	return &record.CategoriesMonthlyTotalPriceRecord{
		Year:         c.Year,
		Month:        c.Month,
		TotalRevenue: int(c.TotalRevenue),
	}
}

func (s *categoryRecordMapper) ToCategoryMonthlyTotalPrices(c []*db.GetMonthlyTotalPriceRow) []*record.CategoriesMonthlyTotalPriceRecord {
	var categoryRecords []*record.CategoriesMonthlyTotalPriceRecord

	for _, category := range c {
		categoryRecords = append(categoryRecords, s.ToCategoryMonthlyTotalPrice(category))
	}

	return categoryRecords
}

func (s *categoryRecordMapper) ToCategoryYearlyTotalPrice(c *db.GetYearlyTotalPriceRow) *record.CategoriesYearlyTotalPriceRecord {
	return &record.CategoriesYearlyTotalPriceRecord{
		Year:         c.Year,
		TotalRevenue: int(c.TotalRevenue),
	}
}

func (s *categoryRecordMapper) ToCategoryYearlyTotalPrices(c []*db.GetYearlyTotalPriceRow) []*record.CategoriesYearlyTotalPriceRecord {
	var categoryRecords []*record.CategoriesYearlyTotalPriceRecord

	for _, category := range c {
		categoryRecords = append(categoryRecords, s.ToCategoryYearlyTotalPrice(category))
	}

	return categoryRecords
}

func (s *categoryRecordMapper) ToCategoryMonthlyTotalPriceById(c *db.GetMonthlyTotalPriceByIdRow) *record.CategoriesMonthlyTotalPriceRecord {
	return &record.CategoriesMonthlyTotalPriceRecord{
		Year:         c.Year,
		Month:        c.Month,
		TotalRevenue: int(c.TotalRevenue),
	}
}

func (s *categoryRecordMapper) ToCategoryMonthlyTotalPricesById(c []*db.GetMonthlyTotalPriceByIdRow) []*record.CategoriesMonthlyTotalPriceRecord {
	var categoryRecords []*record.CategoriesMonthlyTotalPriceRecord

	for _, category := range c {
		categoryRecords = append(categoryRecords, s.ToCategoryMonthlyTotalPriceById(category))
	}

	return categoryRecords
}

func (s *categoryRecordMapper) ToCategoryYearlyTotalPriceById(c *db.GetYearlyTotalPriceByIdRow) *record.CategoriesYearlyTotalPriceRecord {
	return &record.CategoriesYearlyTotalPriceRecord{
		Year:         c.Year,
		TotalRevenue: int(c.TotalRevenue),
	}
}

func (s *categoryRecordMapper) ToCategoryYearlyTotalPricesById(c []*db.GetYearlyTotalPriceByIdRow) []*record.CategoriesYearlyTotalPriceRecord {
	var categoryRecords []*record.CategoriesYearlyTotalPriceRecord

	for _, category := range c {
		categoryRecords = append(categoryRecords, s.ToCategoryYearlyTotalPriceById(category))
	}

	return categoryRecords
}

func (s *categoryRecordMapper) ToCategoryMonthlyTotalPriceByMerchant(c *db.GetMonthlyTotalPriceByMerchantRow) *record.CategoriesMonthlyTotalPriceRecord {
	return &record.CategoriesMonthlyTotalPriceRecord{
		Year:         c.Year,
		Month:        c.Month,
		TotalRevenue: int(c.TotalRevenue),
	}
}

func (s *categoryRecordMapper) ToCategoryMonthlyTotalPricesByMerchant(c []*db.GetMonthlyTotalPriceByMerchantRow) []*record.CategoriesMonthlyTotalPriceRecord {
	var categoryRecords []*record.CategoriesMonthlyTotalPriceRecord

	for _, category := range c {
		categoryRecords = append(categoryRecords, s.ToCategoryMonthlyTotalPriceByMerchant(category))
	}

	return categoryRecords
}

func (s *categoryRecordMapper) ToCategoryYearlyTotalPriceByMerchant(c *db.GetYearlyTotalPriceByMerchantRow) *record.CategoriesYearlyTotalPriceRecord {
	return &record.CategoriesYearlyTotalPriceRecord{
		Year:         c.Year,
		TotalRevenue: int(c.TotalRevenue),
	}
}

func (s *categoryRecordMapper) ToCategoryYearlyTotalPricesByMerchant(c []*db.GetYearlyTotalPriceByMerchantRow) []*record.CategoriesYearlyTotalPriceRecord {
	var categoryRecords []*record.CategoriesYearlyTotalPriceRecord

	for _, category := range c {
		categoryRecords = append(categoryRecords, s.ToCategoryYearlyTotalPriceByMerchant(category))
	}

	return categoryRecords
}
