package apps

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-category/handler"
	mencache "github.com/MamangRust/monolith-point-of-sale-category/cache"
	"github.com/MamangRust/monolith-point-of-sale-category/repository"
	"github.com/MamangRust/monolith-point-of-sale-category/service"
	"github.com/MamangRust/monolith-point-of-sale-pkg/server"
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
	traceLoggerObservability := observability.NewTraceLoggerObservability(srv.Logger)

	mencacheObj := mencache.NewMencache(srv.CacheStore)

	services := service.NewService(&service.Deps{
		Ctx:           context.Background(),
		Mencache:      mencacheObj,
		Repositories:  repos,
		Logger:        srv.Logger,
		Observability: traceLoggerObservability,
	})

	handlers := handler.NewHandler(&handler.Deps{
		Service: services,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pb.RegisterCategoryServiceServer(gs, handlers.Category)
	}

	return srv, nil
}
