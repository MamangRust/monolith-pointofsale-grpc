package mencache

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
)

type ProductQueryCache interface {
	GetCachedProducts(ctx context.Context, req *requests.FindAllProducts) ([]*response.ProductResponse, *int, bool)
	SetCachedProducts(ctx context.Context, req *requests.FindAllProducts, data []*response.ProductResponse, total *int)

	GetCachedProductsByMerchant(ctx context.Context, req *requests.ProductByMerchantRequest) ([]*response.ProductResponse, *int, bool)
	SetCachedProductsByMerchant(ctx context.Context, req *requests.ProductByMerchantRequest, data []*response.ProductResponse, total *int)

	GetCachedProductsByCategory(ctx context.Context, req *requests.ProductByCategoryRequest) ([]*response.ProductResponse, *int, bool)
	SetCachedProductsByCategory(ctx context.Context, req *requests.ProductByCategoryRequest, data []*response.ProductResponse, total *int)

	GetCachedProductActive(ctx context.Context, req *requests.FindAllProducts) ([]*response.ProductResponseDeleteAt, *int, bool)
	SetCachedProductActive(ctx context.Context, req *requests.FindAllProducts, data []*response.ProductResponseDeleteAt, total *int)

	GetCachedProductTrashed(ctx context.Context, req *requests.FindAllProducts) ([]*response.ProductResponseDeleteAt, *int, bool)
	SetCachedProductTrashed(ctx context.Context, req *requests.FindAllProducts, data []*response.ProductResponseDeleteAt, total *int)

	GetCachedProduct(ctx context.Context, productID int) (*response.ProductResponse, bool)
	SetCachedProduct(ctx context.Context, data *response.ProductResponse)
}

type ProductCommandCache interface {
	DeleteCachedProduct(ctx context.Context, productID int)
}
