package service

import (
	"context"

	mencache "github.com/MamangRust/monolith-point-of-sale-order/internal/redis"
	"github.com/MamangRust/monolith-point-of-sale-order/internal/repository"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
)

type Service struct {
	OrderQuery           OrderQueryService
	OrderCommand         OrderCommandService
	OrderStats           OrderStatsService
	OrderStatsByMerchant OrderStatByMerchantService
}

type Deps struct {
	Ctx           context.Context
	Mencache      mencache.Mencache
	Repositories  *repository.Repositories
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

func NewService(deps *Deps) *Service {
	return &Service{
		OrderQuery: NewOrderQueryService(&orderQueryDeps{
			Cache:                deps.Mencache,
			OrderQueryRepository: deps.Repositories.OrderQuery,
			Logger:               deps.Logger,
			Observability:        deps.Observability,
		}),
		OrderCommand: NewOrderCommandService(&orderCommandDeps{
			Cache:                      deps.Mencache,
			CashierQueryRepository:     deps.Repositories.CashierQuery,
			OrderQueryRepository:       deps.Repositories.OrderQuery,
			OrderCommandRepository:     deps.Repositories.OrderCommand,
			OrderItemQueryRepository:   deps.Repositories.OrderItemQuery,
			OrderItemCommandRepository: deps.Repositories.OrderItemCommand,
			MerchantQueryRepository:    deps.Repositories.MerchantQuery,
			ProductQueryRepository:     deps.Repositories.ProductQuery,
			ProductCommandRepository:   deps.Repositories.ProductCommand,
			Logger:                     deps.Logger,
			Observability:              deps.Observability,
		}),
		OrderStats: NewOrderStatsService(&orderStatsDeps{
			Cache:                deps.Mencache,
			OrderStatsRepository: deps.Repositories.OrderStats,
			Logger:               deps.Logger,
			Observability:        deps.Observability,
		}),
		OrderStatsByMerchant: NewOrderStatsByMerchantService(&orderStatsByMerchantDeps{
			Cache:                          deps.Mencache,
			OrderStatsByMerchantRepository: deps.Repositories.OrderStatsByMerchant,
			Logger:                         deps.Logger,
			Observability:                  deps.Observability,
		}),
	}
}
