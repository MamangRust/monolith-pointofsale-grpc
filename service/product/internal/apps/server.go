package apps

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/server"
	"github.com/MamangRust/monolith-point-of-sale-product/internal/handler"
	mencache "github.com/MamangRust/monolith-point-of-sale-product/internal/redis"
	"github.com/MamangRust/monolith-point-of-sale-product/internal/repository"
	"github.com/MamangRust/monolith-point-of-sale-product/internal/service"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	"github.com/MamangRust/monolith-point-of-sale-shared/pb"
	"google.golang.org/grpc"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	repos := repository.NewRepositories(srv.DB)

	mencacheObj := mencache.NewMencache(srv.CacheStore)

	obs := observability.NewTraceLoggerObservability(srv.Logger)

	services := service.NewService(&service.Deps{
		Mencache:      mencacheObj,
		Ctx:           context.Background(),
		Repositories:  repos,
		Logger:        srv.Logger,
		Observability: obs,
	})

	handlers := handler.NewHandler(&handler.Deps{
		Service: services,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pb.RegisterProductServiceServer(gs, handlers.Product)
	}

	return srv, nil
}
