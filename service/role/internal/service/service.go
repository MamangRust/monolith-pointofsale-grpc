package service

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-role/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-point-of-sale-role/internal/redis"
	"github.com/MamangRust/monolith-point-of-sale-role/internal/repository"
	response_service "github.com/MamangRust/monolith-point-of-sale-shared/mapper/response/service"
)

type Service struct {
	RoleQuery   RoleQueryService
	RoleCommand RoleCommandService
}

type Deps struct {
	Ctx          context.Context
	ErrorHandler *errorhandler.ErrorHandler
	Mencache     *mencache.Mencache
	Repositories *repository.Repositories
	Logger       logger.LoggerInterface
}

func NewService(deps *Deps) *Service {
	roleMapper := response_service.NewRoleResponseMapper()

	return &Service{
		RoleQuery:   NewRoleQueryService(deps.ErrorHandler.RoleQueryError, deps.Mencache.RoleQueryCache, deps.Repositories.RoleQuery, deps.Logger, roleMapper),
		RoleCommand: NewRoleCommandService(deps.ErrorHandler.RoleCommandError, deps.Mencache.RoleCommandCache, deps.Repositories.RoleCommand, deps.Logger, roleMapper),
	}
}
