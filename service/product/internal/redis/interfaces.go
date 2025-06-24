package mencache

import (
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
)

type ProductQueryCache interface {
	GetCachedProducts(req *requests.FindAllProducts) ([]*response.ProductResponse, *int, bool)
	SetCachedProducts(req *requests.FindAllProducts, data []*response.ProductResponse, total *int)

	GetCachedProductsByMerchant(req *requests.ProductByMerchantRequest) ([]*response.ProductResponse, *int, bool)
	SetCachedProductsByMerchant(req *requests.ProductByMerchantRequest, data []*response.ProductResponse, total *int)

	GetCachedProductsByCategory(req *requests.ProductByCategoryRequest) ([]*response.ProductResponse, *int, bool)
	SetCachedProductsByCategory(req *requests.ProductByCategoryRequest, data []*response.ProductResponse, total *int)

	GetCachedProductActive(req *requests.FindAllProducts) ([]*response.ProductResponseDeleteAt, *int, bool)
	SetCachedProductActive(req *requests.FindAllProducts, data []*response.ProductResponseDeleteAt, total *int)

	GetCachedProductTrashed(req *requests.FindAllProducts) ([]*response.ProductResponseDeleteAt, *int, bool)
	SetCachedProductTrashed(req *requests.FindAllProducts, data []*response.ProductResponseDeleteAt, total *int)

	GetCachedProduct(productID int) (*response.ProductResponse, bool)
	SetCachedProduct(data *response.ProductResponse)
}

type ProductCommandCache interface {
	DeleteCachedProduct(productID int)
}
