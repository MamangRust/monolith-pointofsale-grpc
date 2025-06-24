package service

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-cashier/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-point-of-sale-cashier/internal/redis"
	"github.com/MamangRust/monolith-point-of-sale-cashier/internal/repository"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	response_service "github.com/MamangRust/monolith-point-of-sale-shared/mapper/response/service"
)

type Service struct {
	CashierQuery           CashierQueryService
	CashierCommand         CashierCommandService
	CashierStats           CashierStatsService
	CashierStatsById       CashierStatsByIdService
	CashierStatsByMerchant CashierStatsByMerchant
}

type Deps struct {
	Ctx           context.Context
	ErrorHandler  *errorhandler.ErrorHandler
	Mencache      *mencache.Mencache
	Repositoriees *repository.Repositories
	Logger        logger.LoggerInterface
}

func NewService(deps *Deps) *Service {
	mapper := response_service.NewCashierResponseMapper()

	return &Service{
		CashierQuery:           NewCashierQueryService(deps.Ctx, deps.ErrorHandler.CashierQueryError, deps.Mencache.CashierQueryCache, deps.Repositoriees.CashierQuery, deps.Logger, mapper),
		CashierCommand:         NewCashierCommandService(deps.Ctx, deps.Mencache.CashierCommandCache, deps.ErrorHandler.CashierCommandError, deps.Repositoriees.MerchantQuery, deps.Repositoriees.UserQuery, deps.Repositoriees.CashierCommand, mapper, deps.Logger),
		CashierStats:           NewCashierStatsService(deps.Ctx, deps.Mencache.CashierStatsCache, deps.ErrorHandler.CashierStatsError, deps.Repositoriees.CashierStats, deps.Logger, mapper),
		CashierStatsById:       NewCashierStatsByIdService(deps.Ctx, deps.Mencache.CashierStatsByIdCache, deps.ErrorHandler.CashierStatsByIdError, deps.Repositoriees.CashierStatsById, deps.Logger, mapper),
		CashierStatsByMerchant: NewCashierStatsByMerchantService(deps.Ctx, deps.Mencache.CashierStatsByMerchantCache, deps.ErrorHandler.CashierStatsByMerchantError, deps.Repositoriees.CashierStatsByMerchant, deps.Logger, mapper),
	}
}
