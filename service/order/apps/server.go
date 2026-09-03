package apps

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-order/handler"
	mencache "github.com/MamangRust/monolith-point-of-sale-order/cache"
	"github.com/MamangRust/monolith-point-of-sale-order/repository"
	"github.com/MamangRust/monolith-point-of-sale-order/service"
	"github.com/MamangRust/monolith-point-of-sale-pkg/server"
	"github.com/MamangRust/monolith-point-of-sale-shared/cache"
	"github.com/MamangRust/monolith-point-of-sale-shared/convert"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	"github.com/MamangRust/monolith-point-of-sale-shared/pb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	cashierAddr := convert.EnvOr("GRPC_CASHIER_ADDR", "localhost:50055")
	merchantAddr := convert.EnvOr("GRPC_MERCHANT_ADDR", "localhost:50056")
	productAddr := convert.EnvOr("GRPC_PRODUCT_ADDR", "localhost:50059")
	orderItemAddr := convert.EnvOr("GRPC_ORDERITEM_ADDR", "localhost:50057")

	srv.Logger.Info("Connecting to gRPC microservices from Order microservice",
		zap.String("cashier", cashierAddr),
		zap.String("merchant", merchantAddr),
		zap.String("product", productAddr),
		zap.String("order_item", orderItemAddr),
	)

	cashierConn, err := server.NewGRPCClient(cashierAddr)
	if err != nil {
		return nil, err
	}

	merchantConn, err := server.NewGRPCClient(merchantAddr)
	if err != nil {
		cashierConn.Close()
		return nil, err
	}

	productConn, err := server.NewGRPCClient(productAddr)
	if err != nil {
		cashierConn.Close()
		merchantConn.Close()
		return nil, err
	}

	orderItemConn, err := server.NewGRPCClient(orderItemAddr)
	if err != nil {
		cashierConn.Close()
		merchantConn.Close()
		productConn.Close()
		return nil, err
	}

	go func() {
		<-srv.Ctx.Done()
		srv.Logger.Info("Closing gRPC client connections in Order microservice")
		cashierConn.Close()
		merchantConn.Close()
		productConn.Close()
		orderItemConn.Close()
	}()

	cashierClient := pb.NewCashierServiceClient(cashierConn)
	merchantClient := pb.NewMerchantServiceClient(merchantConn)
	productClient := pb.NewProductServiceClient(productConn)
	orderItemClient := pb.NewOrderItemServiceClient(orderItemConn)

	repos := repository.NewRepositories(srv.DB, cashierClient, merchantClient, productClient, orderItemClient)
	traceLoggerObservability := observability.NewTraceLoggerObservability(srv.Logger)

	cacheMetrics, err := observability.NewCacheMetrics("order")
	if err != nil {
		return nil, err
	}

	cacheStore := cache.NewCacheStore(srv.Redis, srv.Logger, cacheMetrics)
	mencacheObj := mencache.NewMencache(cacheStore)

	services := service.NewService(&service.Deps{
		Mencache:      mencacheObj,
		Ctx:           context.Background(),
		Repositories:  repos,
		Logger:        srv.Logger,
		Observability: traceLoggerObservability,
	})

	handlers := handler.NewHandler(&handler.Deps{
		Service: services,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pb.RegisterOrderServiceServer(gs, handlers.Order)
	}

	return srv, nil
}
