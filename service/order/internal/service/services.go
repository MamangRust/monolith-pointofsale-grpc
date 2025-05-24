package service

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-order/internal/repository"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	response_service "github.com/MamangRust/monolith-point-of-sale-shared/mapper/response/service"
)

type Service struct {
	OrderQuery           OrderQueryService
	OrderCommand         OrderCommandService
	OrderStats           OrderStatsService
	OrderStatsByMerchant OrderStatByMerchantService
}

type Deps struct {
	Ctx          context.Context
	Repositories *repository.Repositories
	Logger       logger.LoggerInterface
}

func NewService(deps Deps) *Service {
	mapper := response_service.NewOrderResponseMapper()
	return &Service{
		OrderQuery:           NewOrderQueryService(deps.Ctx, deps.Repositories.OrderQuery, deps.Logger, mapper),
		OrderCommand:         NewOrderCommandService(deps.Ctx, deps.Repositories.CashierQuery, deps.Repositories.OrderItemQuery, deps.Repositories.OrderItemCommand, deps.Repositories.OrderQuery, deps.Repositories.OrderCommand, deps.Repositories.ProductQuery, deps.Repositories.ProductCommand, deps.Repositories.MerchantQuery, deps.Logger, mapper),
		OrderStats:           NewOrderStatsService(deps.Ctx, deps.Repositories.OrderStats, deps.Logger, mapper),
		OrderStatsByMerchant: NewOrderStatsByMerchantService(deps.Ctx, deps.Repositories.OrderStatsByMerchant, deps.Logger, mapper),
	}
}
