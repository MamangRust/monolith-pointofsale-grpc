package service

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/kafka"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	response_service "github.com/MamangRust/monolith-point-of-sale-shared/mapper/response/service"
	"github.com/MamangRust/monolith-point-of-sale-transacton/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-point-of-sale-transacton/internal/redis"
	"github.com/MamangRust/monolith-point-of-sale-transacton/internal/repository"
)

type Service struct {
	TransactionQuery           TransactionQueryService
	TransactionCommand         TransactionCommandService
	TransactionStats           TransactionStatsService
	TransactionStatsByMerchant TransactionStatsByMerchantService
}

type Deps struct {
	Ctx          context.Context
	Kafka        *kafka.Kafka
	ErrorHandler *errorhandler.ErrorHandler
	Mencache     *mencache.Mencache
	Repositories *repository.Repositories
	Logger       logger.LoggerInterface
}

func NewService(deps *Deps) *Service {
	mapper := response_service.NewTransactionResponseMapper()

	return &Service{
		TransactionQuery:           NewTransactionQueryService(deps.Mencache.TransactionQueryCache, deps.ErrorHandler.TransactionQueryError, deps.Repositories.TransactionQueryRepository, mapper, deps.Logger),
		TransactionCommand:         NewTransactionCommandService(deps.Mencache.TransactionCommandCache, deps.ErrorHandler.TransactionCommandError, deps.Repositories.CashierQuery, deps.Repositories.MerchantQuery, deps.Repositories.TransactionQueryRepository, deps.Repositories.TransactionCommandRepository, deps.Repositories.OrderQuery, deps.Repositories.OrderItemQuery, mapper, deps.Logger),
		TransactionStats:           NewTransactionStatsService(deps.ErrorHandler.TransactionStatsError, deps.Mencache.TransactionStatsCache, deps.Repositories.TransactionStatsRepository, mapper, deps.Logger),
		TransactionStatsByMerchant: NewTransactionStatsByMerchantService(deps.ErrorHandler.TransactonStatsByMerchantError, deps.Mencache.TransactionStatsByMerchant, deps.Repositories.TransactionStatsByMerchant, mapper, deps.Logger),
	}
}
