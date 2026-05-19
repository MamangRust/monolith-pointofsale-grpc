package apps

import (
	"context"
	"os"

	"github.com/MamangRust/monolith-point-of-sale-order/internal/handler"
	mencache "github.com/MamangRust/monolith-point-of-sale-order/internal/redis"
	"github.com/MamangRust/monolith-point-of-sale-order/internal/repository"
	"github.com/MamangRust/monolith-point-of-sale-order/internal/service"
	"github.com/MamangRust/monolith-point-of-sale-pkg/server"
	"github.com/MamangRust/monolith-point-of-sale-shared/cache"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	"github.com/MamangRust/monolith-point-of-sale-shared/pb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	cashierAddr := getEnv("GRPC_CASHIER_ADDR", "localhost:50054")
	merchantAddr := getEnv("GRPC_MERCHANT_ADDR", "localhost:50055")
	productAddr := getEnv("GRPC_PRODUCT_ADDR", "localhost:50058")
	orderItemAddr := getEnv("GRPC_ORDER_ITEM_ADDR", "localhost:50056")

	srv.Logger.Info("Connecting to gRPC microservices from Order microservice",
		zap.String("cashier", cashierAddr),
		zap.String("merchant", merchantAddr),
		zap.String("product", productAddr),
		zap.String("order_item", orderItemAddr),
	)

	cashierConn, err := grpc.NewClient(cashierAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	merchantConn, err := grpc.NewClient(merchantAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		cashierConn.Close()
		return nil, err
	}

	productConn, err := grpc.NewClient(productAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		cashierConn.Close()
		merchantConn.Close()
		return nil, err
	}

	orderItemConn, err := grpc.NewClient(orderItemAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
