package service

import (
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
)

type ProductQueryService interface {
	FindAll(req *requests.FindAllProducts) ([]*response.ProductResponse, *int, *response.ErrorResponse)
	FindByMerchant(req *requests.ProductByMerchantRequest) ([]*response.ProductResponse, *int, *response.ErrorResponse)
	FindByCategory(req *requests.ProductByCategoryRequest) ([]*response.ProductResponse, *int, *response.ErrorResponse)
	FindById(productID int) (*response.ProductResponse, *response.ErrorResponse)
	FindByActive(req *requests.FindAllProducts) ([]*response.ProductResponseDeleteAt, *int, *response.ErrorResponse)
	FindByTrashed(req *requests.FindAllProducts) ([]*response.ProductResponseDeleteAt, *int, *response.ErrorResponse)
}

type ProductCommandService interface {
	CreateProduct(req *requests.CreateProductRequest) (*response.ProductResponse, *response.ErrorResponse)
	UpdateProduct(req *requests.UpdateProductRequest) (*response.ProductResponse, *response.ErrorResponse)
	TrashProduct(productID int) (*response.ProductResponseDeleteAt, *response.ErrorResponse)
	RestoreProduct(productID int) (*response.ProductResponseDeleteAt, *response.ErrorResponse)
	DeleteProductPermanent(productID int) (bool, *response.ErrorResponse)
	RestoreAllProducts() (bool, *response.ErrorResponse)
	DeleteAllProductsPermanent() (bool, *response.ErrorResponse)
}
