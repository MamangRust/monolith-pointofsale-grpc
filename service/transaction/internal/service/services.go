package service

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/kafka"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	mencache "github.com/MamangRust/monolith-point-of-sale-transacton/internal/redis"
	"github.com/MamangRust/monolith-point-of-sale-transacton/internal/repository"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
)

type Service struct {
	TransactionQuery           TransactionQueryService
	TransactionCommand         TransactionCommandService
	TransactionStats           TransactionStatsService
	TransactionStatsByMerchant TransactionStatsByMerchantService
}

type Deps struct {
	Ctx           context.Context
	Kafka         *kafka.Kafka
	Mencache      mencache.Mencache
	Repositories  *repository.Repositories
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

func NewService(deps *Deps) *Service {
	return &Service{
		TransactionQuery:           NewTransactionQueryService(deps.Mencache, deps.Repositories.TransactionQueryRepository, deps.Logger, deps.Observability),
		TransactionCommand:         NewTransactionCommandService(deps.Mencache, deps.Repositories.CashierQuery, deps.Repositories.MerchantQuery, deps.Repositories.TransactionQueryRepository, deps.Repositories.TransactionCommandRepository, deps.Repositories.OrderQuery, deps.Repositories.OrderItemQuery, deps.Logger, deps.Observability),
		TransactionStats:           NewTransactionStatsService(deps.Mencache, deps.Repositories.TransactionStatsRepository, deps.Logger, deps.Observability),
		TransactionStatsByMerchant: NewTransactionStatsByMerchantService(deps.Mencache, deps.Repositories.TransactionStatsByMerchant, deps.Logger, deps.Observability),
	}
}
