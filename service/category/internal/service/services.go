package service

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-category/internal/repository"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	response_service "github.com/MamangRust/monolith-point-of-sale-shared/mapper/response/service"
)

type Service struct {
	CategoryQuery           CategoryQueryService
	CategoryCommand         CategoryCommandService
	CategoryStats           CategoryStatsService
	CategoryStatsById       CategoryStatsByIdService
	CategoryStatsByMerchant CategoryStatsByMerchantService
}

type Deps struct {
	Ctx          context.Context
	Repositories *repository.Repositories
	Logger       logger.LoggerInterface
}

func NewService(deps Deps) *Service {
	categoryMapper := response_service.NewCategoryResponseMapper()

	return &Service{
		CategoryQuery:           NewCategoryQueryService(deps.Ctx, deps.Repositories.CategoryQuery, deps.Logger, categoryMapper),
		CategoryCommand:         NewCategoryCommandService(deps.Ctx, deps.Repositories.CategoryCommand, deps.Repositories.CategoryQuery, deps.Logger, categoryMapper),
		CategoryStats:           NewCategoryStatsService(deps.Ctx, deps.Repositories.CategoryStats, deps.Logger, categoryMapper),
		CategoryStatsById:       NewCategoryStatsByIdService(deps.Ctx, deps.Repositories.CategoryStatsById, deps.Logger, categoryMapper),
		CategoryStatsByMerchant: NewCategoryStatsByMerchantService(deps.Ctx, deps.Repositories.CategoryStatsByMerchant, deps.Logger, categoryMapper),
	}
}
