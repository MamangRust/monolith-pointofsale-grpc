package service

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-category/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-point-of-sale-category/internal/redis"
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
	ErrorHandler *errorhandler.ErrorHandler
	Mencache     *mencache.Mencache
	Repositories *repository.Repositories
	Logger       logger.LoggerInterface
}

func NewService(deps *Deps) *Service {
	categoryMapper := response_service.NewCategoryResponseMapper()

	return &Service{
		CategoryQuery:           NewCategoryQueryService(deps.ErrorHandler.CategoryQueryError, deps.Mencache.CategoryQueryCache, deps.Repositories.CategoryQuery, deps.Logger, categoryMapper),
		CategoryCommand:         NewCategoryCommandService(deps.Mencache.CategoryCommandCache, deps.ErrorHandler.CategoryCommandError, deps.Repositories.CategoryCommand, deps.Repositories.CategoryQuery, deps.Logger, categoryMapper),
		CategoryStats:           NewCategoryStatsService(deps.Mencache.CategoryStatsCache, deps.ErrorHandler.CategoryStatsByIdError, deps.Repositories.CategoryStats, deps.Logger, categoryMapper),
		CategoryStatsById:       NewCategoryStatsByIdService(deps.Mencache.CategoryStatsByIdCache, deps.ErrorHandler.CategoryStatsByIdError, deps.Repositories.CategoryStatsById, deps.Logger, categoryMapper),
		CategoryStatsByMerchant: NewCategoryStatsByMerchantService(deps.Mencache.CategoryStatsByMerchantCache, deps.ErrorHandler.CategoryStatsByMerchantError, deps.Repositories.CategoryStatsByMerchant, deps.Logger, categoryMapper),
	}
}
