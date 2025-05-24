package service

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-product/internal/repository"
	response_service "github.com/MamangRust/monolith-point-of-sale-shared/mapper/response/service"
)

type Service struct {
	ProductQuery   ProductQueryService
	ProductCommand ProductCommandService
}

type Deps struct {
	Ctx          context.Context
	Repositories *repository.Repositories
	Logger       logger.LoggerInterface
}

func NewService(deps Deps) *Service {
	mapper := response_service.NewProductResponseMapper()

	return &Service{
		ProductQuery:   NewProductQueryService(deps.Ctx, deps.Repositories.ProductQuery, mapper, deps.Logger),
		ProductCommand: NewProductCommandService(deps.Ctx, deps.Repositories.CategoryQuery, deps.Repositories.MerchantQuery, deps.Repositories.ProductQuery, deps.Repositories.ProductCommand, mapper, deps.Logger),
	}
}
