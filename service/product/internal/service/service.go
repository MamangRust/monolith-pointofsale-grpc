package service

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-product/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-point-of-sale-product/internal/redis"
	"github.com/MamangRust/monolith-point-of-sale-product/internal/repository"
	response_service "github.com/MamangRust/monolith-point-of-sale-shared/mapper/response/service"
)

type Service struct {
	ProductQuery   ProductQueryService
	ProductCommand ProductCommandService
}

type Deps struct {
	Ctx          context.Context
	ErrorHandler *errorhandler.ErrorHandler
	Mencache     *mencache.Mencache
	Repositories *repository.Repositories
	Logger       logger.LoggerInterface
}

func NewService(deps *Deps) *Service {
	mapper := response_service.NewProductResponseMapper()

	return &Service{
		ProductQuery:   NewProductQueryService(deps.ErrorHandler.ProductQueryError, deps.Mencache.ProductQuery, deps.Repositories.ProductQuery, mapper, deps.Logger),
		ProductCommand: NewProductCommandService(deps.ErrorHandler.ProductCommandError, deps.Mencache.ProductCommand, deps.Repositories.CategoryQuery, deps.Repositories.MerchantQuery, deps.Repositories.ProductQuery, deps.Repositories.ProductCommand, mapper, deps.Logger),
	}
}
