package service

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-order-item/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-point-of-sale-order-item/internal/redis"
	"github.com/MamangRust/monolith-point-of-sale-order-item/internal/repository"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	response_service "github.com/MamangRust/monolith-point-of-sale-shared/mapper/response/service"
)

type Service struct {
	OrderItemQuery OrderItemQueryService
}

type Deps struct {
	Ctx          context.Context
	ErrorHandler *errorhandler.ErrorHandler
	Mencache     *mencache.Mencache
	Repositories *repository.Repositories
	Logger       logger.LoggerInterface
}

func NewService(deps *Deps) *Service {
	mapper := response_service.NewOrderItemResponseMapper()

	return &Service{
		OrderItemQuery: NewOrderItemQueryService(deps.Ctx, deps.ErrorHandler.OrderItemQueryError, deps.Mencache.OrderItemQueryCache, deps.Repositories.OrderItemQuery, deps.Logger, mapper),
	}
}
