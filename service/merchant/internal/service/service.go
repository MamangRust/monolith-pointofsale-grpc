package service

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-merchant/internal/repository"
	"github.com/MamangRust/monolith-point-of-sale-pkg/kafka"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	response_service "github.com/MamangRust/monolith-point-of-sale-shared/mapper/response/service"
)

type Service struct {
	MerchantQuery           MerchantQueryService
	MerchantCommand         MerchantCommandService
	MerchantDocumentCommand MerchantDocumentCommandService
	MerchantDocumentQuery   MerchantDocumentQueryService
}

type Deps struct {
	Kafka        kafka.Kafka
	Ctx          context.Context
	Repositories *repository.Repositories
	Logger       logger.LoggerInterface
}

func NewService(deps Deps) *Service {
	merchantMapper := response_service.NewMerchantResponseMapper()
	merchantDocument := response_service.NewMerchantDocumentResponseMapper()

	return &Service{
		MerchantQuery:           NewMerchantQueryService(deps.Ctx, deps.Repositories.MerchantQuery, deps.Logger, merchantMapper),
		MerchantCommand:         NewMerchantCommandService(deps.Kafka, deps.Ctx, deps.Repositories.UserQuery, deps.Repositories.MerchantQuery, deps.Repositories.MerchantCommand, deps.Logger, merchantMapper),
		MerchantDocumentCommand: NewMerchantDocumentCommandService(deps.Kafka, deps.Ctx, deps.Repositories.MerchantDocumentCommand, deps.Repositories.MerchantQuery, deps.Repositories.UserQuery, deps.Logger, merchantDocument),
		MerchantDocumentQuery:   NewMerchantDocumentQueryService(deps.Ctx, deps.Repositories.MerchantDocumentQuery, deps.Logger, merchantDocument),
	}
}
