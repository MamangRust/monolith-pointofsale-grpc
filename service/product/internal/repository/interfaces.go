package repository

import (
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/record"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type CategoryQueryRepository interface {
	FindById(category_id int) (*record.CategoriesRecord, error)
}

type MerchantQueryRepository interface {
	FindById(merchant_id int) (*record.MerchantRecord, error)
}

type ProductQueryRepository interface {
	FindAllProducts(req *requests.FindAllProducts) ([]*record.ProductRecord, *int, error)
	FindByActive(req *requests.FindAllProducts) ([]*record.ProductRecord, *int, error)
	FindByTrashed(req *requests.FindAllProducts) ([]*record.ProductRecord, *int, error)
	FindByMerchant(req *requests.ProductByMerchantRequest) ([]*record.ProductRecord, *int, error)
	FindByCategory(req *requests.ProductByCategoryRequest) ([]*record.ProductRecord, *int, error)
	FindById(user_id int) (*record.ProductRecord, error)
	FindByIdTrashed(id int) (*record.ProductRecord, error)
}

type ProductCommandRepository interface {
	CreateProduct(request *requests.CreateProductRequest) (*record.ProductRecord, error)
	UpdateProduct(request *requests.UpdateProductRequest) (*record.ProductRecord, error)
	UpdateProductCountStock(product_id int, stock int) (*record.ProductRecord, error)
	TrashedProduct(user_id int) (*record.ProductRecord, error)
	RestoreProduct(user_id int) (*record.ProductRecord, error)
	DeleteProductPermanent(user_id int) (bool, error)
	RestoreAllProducts() (bool, error)
	DeleteAllProductPermanent() (bool, error)
}
