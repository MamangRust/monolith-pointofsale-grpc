package recordmapper

import (
	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"

	"github.com/MamangRust/monolith-point-of-sale-shared/domain/record"
)

type productRecordMapper struct {
}

func NewProductRecordMapper() *productRecordMapper {
	return &productRecordMapper{}
}

func (s *productRecordMapper) ToProductRecord(product *db.Product) *record.ProductRecord {
	var deletedAt *string
	if product.DeletedAt.Valid {
		deletedAtStr := product.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.ProductRecord{
		ID:           int(product.ProductID),
		MerchantID:   int(product.MerchantID),
		CategoryID:   int(product.CategoryID),
		Name:         product.Name,
		Description:  product.Description.String,
		Price:        int(product.Price),
		CountInStock: int(product.CountInStock),
		Brand:        product.Brand.String,
		Weight:       int(product.Weight.Int32),
		SlugProduct:  product.SlugProduct.String,
		ImageProduct: product.ImageProduct.String,
		Barcode:      product.Barcode.String,
		CreatedAt:    product.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:    product.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:    deletedAt,
	}
}

func (s *productRecordMapper) ToProductsRecord(products []*db.Product) []*record.ProductRecord {
	var result []*record.ProductRecord

	for _, product := range products {
		result = append(result, s.ToProductRecord(product))
	}

	return result
}

func (s *productRecordMapper) ToProductRecordPagination(product *db.GetProductsRow) *record.ProductRecord {
	var deletedAt *string
	if product.DeletedAt.Valid {
		deletedAtStr := product.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.ProductRecord{
		ID:           int(product.ProductID),
		MerchantID:   int(product.MerchantID),
		CategoryID:   int(product.CategoryID),
		Name:         product.Name,
		Description:  product.Description.String,
		Price:        int(product.Price),
		CountInStock: int(product.CountInStock),
		Brand:        product.Brand.String,
		Weight:       int(product.Weight.Int32),
		SlugProduct:  product.SlugProduct.String,
		ImageProduct: product.ImageProduct.String,
		Barcode:      product.Barcode.String,
		CreatedAt:    product.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:    product.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:    deletedAt,
	}
}

func (s *productRecordMapper) ToProductsRecordPagination(products []*db.GetProductsRow) []*record.ProductRecord {
	var result []*record.ProductRecord

	for _, product := range products {
		result = append(result, s.ToProductRecordPagination(product))
	}

	return result
}

func (s *productRecordMapper) ToProductRecordActivePagination(product *db.GetProductsActiveRow) *record.ProductRecord {
	var deletedAt *string
	if product.DeletedAt.Valid {
		deletedAtStr := product.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.ProductRecord{
		ID:           int(product.ProductID),
		MerchantID:   int(product.MerchantID),
		CategoryID:   int(product.CategoryID),
		Name:         product.Name,
		Description:  product.Description.String,
		Price:        int(product.Price),
		CountInStock: int(product.CountInStock),
		Brand:        product.Brand.String,
		Weight:       int(product.Weight.Int32),
		SlugProduct:  product.SlugProduct.String,
		ImageProduct: product.ImageProduct.String,
		Barcode:      product.Barcode.String,
		CreatedAt:    product.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:    product.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:    deletedAt,
	}
}

func (s *productRecordMapper) ToProductsRecordActivePagination(products []*db.GetProductsActiveRow) []*record.ProductRecord {
	var result []*record.ProductRecord

	for _, product := range products {
		result = append(result, s.ToProductRecordActivePagination(product))
	}

	return result
}

func (s *productRecordMapper) ToProductRecordTrashedPagination(product *db.GetProductsTrashedRow) *record.ProductRecord {
	var deletedAt *string
	if product.DeletedAt.Valid {
		deletedAtStr := product.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.ProductRecord{
		ID:           int(product.ProductID),
		MerchantID:   int(product.MerchantID),
		CategoryID:   int(product.CategoryID),
		Name:         product.Name,
		Description:  product.Description.String,
		Price:        int(product.Price),
		CountInStock: int(product.CountInStock),
		Brand:        product.Brand.String,
		Weight:       int(product.Weight.Int32),
		SlugProduct:  product.SlugProduct.String,
		ImageProduct: product.ImageProduct.String,
		Barcode:      product.Barcode.String,
		CreatedAt:    product.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:    product.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:    deletedAt,
	}
}

func (s *productRecordMapper) ToProductsRecordTrashedPagination(products []*db.GetProductsTrashedRow) []*record.ProductRecord {
	var result []*record.ProductRecord

	for _, product := range products {
		result = append(result, s.ToProductRecordTrashedPagination(product))
	}

	return result
}

func (s *productRecordMapper) ToProductRecordMerchantPagination(product *db.GetProductsByMerchantRow) *record.ProductRecord {
	return &record.ProductRecord{
		ID:           int(product.ProductID),
		Name:         product.Name,
		Description:  product.Description.String,
		Price:        int(product.Price),
		CountInStock: int(product.CountInStock),
		Brand:        product.Brand.String,
		ImageProduct: product.ImageProduct.String,
		CreatedAt:    product.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
	}
}

func (s *productRecordMapper) ToProductsRecordMerchantPagination(products []*db.GetProductsByMerchantRow) []*record.ProductRecord {
	var result []*record.ProductRecord

	for _, product := range products {
		result = append(result, s.ToProductRecordMerchantPagination(product))
	}

	return result
}

func (s *productRecordMapper) ToProductRecordCategoryPagination(product *db.GetProductsByCategoryNameRow) *record.ProductRecord {
	var deletedAt *string
	if product.DeletedAt.Valid {
		deletedAtStr := product.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.ProductRecord{
		ID:           int(product.ProductID),
		MerchantID:   int(product.MerchantID),
		CategoryID:   int(product.CategoryID),
		Name:         product.Name,
		Description:  product.Description.String,
		Price:        int(product.Price),
		CountInStock: int(product.CountInStock),
		Brand:        product.Brand.String,
		Weight:       int(product.Weight.Int32),
		SlugProduct:  product.SlugProduct.String,
		ImageProduct: product.ImageProduct.String,
		Barcode:      product.Barcode.String,
		CreatedAt:    product.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:    product.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:    deletedAt,
	}
}

func (s *productRecordMapper) ToProductsRecordCategoryPagination(products []*db.GetProductsByCategoryNameRow) []*record.ProductRecord {
	var result []*record.ProductRecord

	for _, product := range products {
		result = append(result, s.ToProductRecordCategoryPagination(product))
	}

	return result
}
