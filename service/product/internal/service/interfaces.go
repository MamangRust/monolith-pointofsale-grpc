package service

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
)

type ProductQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllProducts) ([]*response.ProductResponse, *int, *response.ErrorResponse)
	FindByMerchant(ctx context.Context, req *requests.ProductByMerchantRequest) ([]*response.ProductResponse, *int, *response.ErrorResponse)
	FindByCategory(ctx context.Context, req *requests.ProductByCategoryRequest) ([]*response.ProductResponse, *int, *response.ErrorResponse)
	FindById(ctx context.Context, productID int) (*response.ProductResponse, *response.ErrorResponse)
	FindByActive(ctx context.Context, req *requests.FindAllProducts) ([]*response.ProductResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(ctx context.Context, req *requests.FindAllProducts) ([]*response.ProductResponseDeleteAt, *int, *response.ErrorResponse)
}

type ProductCommandService interface {
	CreateProduct(ctx context.Context, req *requests.CreateProductRequest) (*response.ProductResponse, *response.ErrorResponse)
	UpdateProduct(ctx context.Context, req *requests.UpdateProductRequest) (*response.ProductResponse, *response.ErrorResponse)
	TrashProduct(ctx context.Context, productID int) (*response.ProductResponseDeleteAt, *response.ErrorResponse)
	RestoreProduct(ctx context.Context, productID int) (*response.ProductResponseDeleteAt, *response.ErrorResponse)
	DeleteProductPermanent(ctx context.Context, productID int) (bool, *response.ErrorResponse)
	RestoreAllProducts(ctx context.Context) (bool, *response.ErrorResponse)
	DeleteAllProductsPermanent(ctx context.Context) (bool, *response.ErrorResponse)
}
