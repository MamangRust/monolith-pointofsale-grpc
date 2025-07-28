package service

import (
	"github.com/MamangRust/monolith-point-of-sale-auth/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-point-of-sale-auth/internal/redis"
	"github.com/MamangRust/monolith-point-of-sale-auth/internal/repository"
	"github.com/MamangRust/monolith-point-of-sale-pkg/auth"
	"github.com/MamangRust/monolith-point-of-sale-pkg/hash"
	"github.com/MamangRust/monolith-point-of-sale-pkg/kafka"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	response_service "github.com/MamangRust/monolith-point-of-sale-shared/mapper/response/service"
)

type Service struct {
	Login         LoginService
	Register      RegistrationService
	PasswordReset PasswordResetService
	Identify      IdentifyService
}

type Deps struct {
	ErrorHandler *errorhandler.ErrorHandler
	Mencache     *mencache.Mencache
	Repositories *repository.Repositories
	Token        auth.TokenManager
	Hash         hash.HashPassword
	Logger       logger.LoggerInterface
	Kafka        *kafka.Kafka
	Mapper       response_service.UserResponseMapper
}

func NewService(deps *Deps) *Service {
	tokenService := NewTokenService(deps.Repositories.RefreshToken, deps.Token, deps.Logger)

	mapper := response_service.NewUserResponseMapper()

	return &Service{
		Login:         NewLoginService(deps.ErrorHandler.PasswordError, deps.ErrorHandler.TokenError, deps.ErrorHandler.LoginError, deps.Mencache.LoginCache, deps.Logger, deps.Hash, deps.Repositories.User, deps.Repositories.RefreshToken, deps.Token, *tokenService),
		Register:      NewRegisterService(deps.ErrorHandler.RegisterError, deps.ErrorHandler.PasswordError, deps.ErrorHandler.RandomString, deps.ErrorHandler.MarshalError, deps.ErrorHandler.KafkaError, deps.Mencache.RegisterCache, deps.Repositories.User, deps.Repositories.Role, deps.Repositories.UserRole, deps.Hash, deps.Kafka, deps.Logger, mapper),
		PasswordReset: NewPasswordResetService(deps.ErrorHandler.PasswordResetError, deps.ErrorHandler.RandomString, deps.ErrorHandler.MarshalError, deps.ErrorHandler.PasswordError, deps.ErrorHandler.KafkaError, deps.Mencache.PasswordResetCache, deps.Kafka, deps.Logger, deps.Repositories.User, deps.Repositories.ResetToken),
		Identify:      NewIdentityService(deps.ErrorHandler.IdentityError, deps.ErrorHandler.TokenError, deps.Mencache.IdentityCache, deps.Token, deps.Repositories.RefreshToken, deps.Repositories.User, deps.Logger, mapper, *tokenService),
	}
}
