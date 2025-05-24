package service

import (
	"context"

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
	Repositoriees *repository.Repositories
	Logger        logger.LoggerInterface
}

func NewService(deps Deps) *Service {
	mapper := response_service.NewCashierResponseMapper()

	return &Service{
		CashierQuery:           NewCashierQueryService(deps.Ctx, deps.Repositoriees.CashierQuery, deps.Logger, mapper),
		CashierCommand:         NewCashierCommandService(deps.Ctx, deps.Repositoriees.MerchantQuery, deps.Repositoriees.UserQuery, deps.Repositoriees.CashierCommand, mapper, deps.Logger),
		CashierStats:           NewCashierStatsService(deps.Ctx, deps.Repositoriees.CashierStats, deps.Logger, mapper),
		CashierStatsById:       NewCashierStatsByIdService(deps.Ctx, deps.Repositoriees.CashierStatsById, deps.Logger, mapper),
		CashierStatsByMerchant: NewCashierStatsByMerchantService(deps.Ctx, deps.Repositoriees.CashierStatsByMerchant, deps.Logger, mapper),
	}
}
