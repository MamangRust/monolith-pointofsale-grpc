package apps

import (
	"os"

	"github.com/MamangRust/monolith-point-of-sale-pkg/kafka"
	"github.com/MamangRust/monolith-point-of-sale-pkg/server"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	"github.com/MamangRust/monolith-point-of-sale-shared/pb"
	"github.com/MamangRust/monolith-point-of-sale-transacton/internal/handler"
	mencache "github.com/MamangRust/monolith-point-of-sale-transacton/internal/redis"
	"github.com/MamangRust/monolith-point-of-sale-transacton/internal/repository"
	"github.com/MamangRust/monolith-point-of-sale-transacton/internal/service"
	"github.com/spf13/viper"
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
	orderAddr := getEnv("GRPC_ORDER_ADDR", "localhost:50057")
	orderItemAddr := getEnv("GRPC_ORDER_ITEM_ADDR", "localhost:50056")

	srv.Logger.Info("Connecting to gRPC microservices from Transaction microservice",
		zap.String("cashier", cashierAddr),
		zap.String("merchant", merchantAddr),
		zap.String("order", orderAddr),
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

	orderConn, err := grpc.NewClient(orderAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		cashierConn.Close()
		merchantConn.Close()
		return nil, err
	}

	orderItemConn, err := grpc.NewClient(orderItemAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		cashierConn.Close()
		merchantConn.Close()
		orderConn.Close()
		return nil, err
	}

	go func() {
		<-srv.Ctx.Done()
		srv.Logger.Info("Closing gRPC client connections in Transaction microservice")
		cashierConn.Close()
		merchantConn.Close()
		orderConn.Close()
		orderItemConn.Close()
	}()

	cashierClient := pb.NewCashierServiceClient(cashierConn)
	merchantClient := pb.NewMerchantServiceClient(merchantConn)
	orderClient := pb.NewOrderServiceClient(orderConn)
	orderItemClient := pb.NewOrderItemServiceClient(orderItemConn)

	repos := repository.NewRepositories(srv.DB, cashierClient, merchantClient, orderClient, orderItemClient)
	myKafka := kafka.NewKafka(srv.Logger, []string{viper.GetString("KAFKA_BROKERS")})

	mencacheObj := mencache.NewMencache(srv.CacheStore)

	traceLoggerObservability := observability.NewTraceLoggerObservability(srv.Logger)

	services := service.NewService(&service.Deps{
		Mencache:      mencacheObj,
		Repositories:  repos,
		Logger:        srv.Logger,
		Kafka:         myKafka,
		Observability: traceLoggerObservability,
	})

	handlers := handler.NewHandler(&handler.Deps{
		Service: services,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pb.RegisterTransactionServiceServer(gs, handlers.Transaction)
	}

	return srv, nil
}
