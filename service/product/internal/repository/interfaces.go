package repository

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-shared/domain/record"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type CategoryQueryRepository interface {
	FindById(ctx context.Context, category_id int) (*record.CategoriesRecord, error)
}

type MerchantQueryRepository interface {
	FindById(ctx context.Context, merchant_id int) (*record.MerchantRecord, error)
}

type ProductQueryRepository interface {
	FindAllProducts(ctx context.Context, req *requests.FindAllProducts) ([]*record.ProductRecord, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllProducts) ([]*record.ProductRecord, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllProducts) ([]*record.ProductRecord, *int, error)
	FindByMerchant(ctx context.Context, req *requests.ProductByMerchantRequest) ([]*record.ProductRecord, *int, error)
	FindByCategory(ctx context.Context, req *requests.ProductByCategoryRequest) ([]*record.ProductRecord, *int, error)
	FindById(ctx context.Context, user_id int) (*record.ProductRecord, error)
	FindByIdTrashed(ctx context.Context, id int) (*record.ProductRecord, error)
}

type ProductCommandRepository interface {
	CreateProduct(ctx context.Context, request *requests.CreateProductRequest) (*record.ProductRecord, error)
	UpdateProduct(ctx context.Context, request *requests.UpdateProductRequest) (*record.ProductRecord, error)
	UpdateProductCountStock(ctx context.Context, product_id int, stock int) (*record.ProductRecord, error)
	TrashedProduct(ctx context.Context, user_id int) (*record.ProductRecord, error)
	RestoreProduct(ctx context.Context, user_id int) (*record.ProductRecord, error)
	DeleteProductPermanent(ctx context.Context, user_id int) (bool, error)
	RestoreAllProducts(ctx context.Context) (bool, error)
	DeleteAllProductPermanent(ctx context.Context) (bool, error)
}
