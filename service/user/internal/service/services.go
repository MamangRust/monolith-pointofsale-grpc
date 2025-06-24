package service

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/hash"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	response_service "github.com/MamangRust/monolith-point-of-sale-shared/mapper/response/service"
	"github.com/MamangRust/monolith-point-of-sale-user/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-point-of-sale-user/internal/redis"
	"github.com/MamangRust/monolith-point-of-sale-user/internal/repository"
)

type Service struct {
	UserQuery   UserQueryService
	UserCommand UserCommandService
}

type Deps struct {
	Ctx          context.Context
	ErrorHandler *errorhandler.ErrorHandler
	Mencache     *mencache.Mencache
	Repositories *repository.Repositories
	Hash         hash.HashPassword
	Logger       logger.LoggerInterface
}

func NewService(deps *Deps) *Service {
	userMapper := response_service.NewUserResponseMapper()

	return &Service{
		UserQuery:   NewUserQueryService(deps.Ctx, deps.ErrorHandler.UserQueryError, deps.Mencache.UserQueryCache, deps.Repositories.UserQuery, deps.Logger, userMapper),
		UserCommand: NewUserCommandService(deps.Ctx, deps.ErrorHandler.UserCommandError, deps.Mencache.UserCommandCache, deps.Repositories.UserQuery, deps.Repositories.UserCommand, deps.Repositories.Role, deps.Logger, userMapper, deps.Hash),
	}
}
