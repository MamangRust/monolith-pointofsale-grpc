package apps

import (
	"github.com/MamangRust/monolith-point-of-sale-auth/internal/handler"
	mencache "github.com/MamangRust/monolith-point-of-sale-auth/internal/redis"
	"github.com/MamangRust/monolith-point-of-sale-auth/internal/repository"
	"github.com/MamangRust/monolith-point-of-sale-auth/internal/service"
	"github.com/MamangRust/monolith-point-of-sale-pkg/auth"
	"github.com/MamangRust/monolith-point-of-sale-pkg/hash"
	"github.com/MamangRust/monolith-point-of-sale-pkg/kafka"
	"github.com/MamangRust/monolith-point-of-sale-pkg/server"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	"github.com/MamangRust/monolith-point-of-sale-shared/pb"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	tokenManager, err := auth.NewManager(viper.GetString("SECRET_KEY"))
	if err != nil {
		return nil, err
	}

	repos := repository.NewRepositories(srv.DB)
	hash := hash.NewHashingPassword()
	myKafka := kafka.NewKafka(srv.Logger, []string{viper.GetString("KAFKA_BROKERS")})

	mencacheObj := mencache.NewMencache(srv.CacheStore)

	services := service.NewService(&service.Deps{
		Mencache:      mencacheObj,
		Repositories:  repos,
		Hash:          hash,
		Token:         tokenManager,
		Logger:        srv.Logger,
		Kafka:         myKafka,
		Observability: observability.NewTraceLoggerObservability(srv.Logger),
	})

	handlers := handler.NewHandler(&handler.Deps{
		Service: services,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pb.RegisterAuthServiceServer(gs, handlers.Auth)
	}

	return srv, nil
}
