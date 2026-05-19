package service

import (
	"context"

	mencache "github.com/MamangRust/monolith-point-of-sale-merchant/internal/redis"
	"github.com/MamangRust/monolith-point-of-sale-merchant/internal/repository"
	"github.com/MamangRust/monolith-point-of-sale-pkg/kafka"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
)

type Service struct {
	MerchantQuery           MerchantQueryService
	MerchantCommand         MerchantCommandService
	MerchantDocumentCommand MerchantDocumentCommandService
	MerchantDocumentQuery   MerchantDocumentQueryService
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
		MerchantQuery: NewMerchantQueryService(&merchantQueryDeps{
			Cache:         deps.Mencache,
			MerchantQuery: deps.Repositories.MerchantQuery,
			Logger:        deps.Logger,
			Observability: deps.Observability,
		}),
		MerchantCommand: NewMerchantCommandService(&merchantCommandDeps{
			Kafka:           deps.Kafka,
			Cache:           deps.Mencache,
			UserQuery:       deps.Repositories.UserQuery,
			MerchantQuery:   deps.Repositories.MerchantQuery,
			MerchantCommand: deps.Repositories.MerchantCommand,
			Logger:          deps.Logger,
			Observability:   deps.Observability,
		}),
		MerchantDocumentCommand: NewMerchantDocumentCommandService(&merchantDocumentCommandDeps{
			Kafka:                   deps.Kafka,
			Cache:                   deps.Mencache,
			MerchantQuery:           deps.Repositories.MerchantQuery,
			MerchantDocumentCommand: deps.Repositories.MerchantDocumentCommand,
			UserQuery:               deps.Repositories.UserQuery,
			Logger:                  deps.Logger,
			Observability:           deps.Observability,
		}),
		MerchantDocumentQuery: NewMerchantDocumentQueryService(&merchantDocumentQueryDeps{
			Cache:                 deps.Mencache,
			MerchantDocumentQuery: deps.Repositories.MerchantDocumentQuery,
			Logger:                deps.Logger,
			Observability:         deps.Observability,
		}),
	}
}
