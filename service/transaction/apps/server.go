package apps

import (
	"os"

	"github.com/MamangRust/monolith-point-of-sale-pkg/kafka"
	"github.com/MamangRust/monolith-point-of-sale-pkg/outbox"
	"github.com/MamangRust/monolith-point-of-sale-pkg/server"
	"github.com/MamangRust/monolith-point-of-sale-shared/convert"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	"github.com/MamangRust/monolith-point-of-sale-shared/pb"
	mencache "github.com/MamangRust/monolith-point-of-sale-transacton/cache"
	"github.com/MamangRust/monolith-point-of-sale-transacton/handler"
	"github.com/MamangRust/monolith-point-of-sale-transacton/repository"
	"github.com/MamangRust/monolith-point-of-sale-transacton/service"
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
	orderAddr := convert.EnvOr("GRPC_ORDER_ADDR", "localhost:50058")
	orderItemAddr := convert.EnvOr("GRPC_ORDERITEM_ADDR", "localhost:50057")

	srv.Logger.Info("Connecting to gRPC microservices from Transaction microservice",
		zap.String("cashier", cashierAddr),
		zap.String("merchant", merchantAddr),
		zap.String("order", orderAddr),
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

	orderConn, err := server.NewGRPCClient(orderAddr)
	if err != nil {
		cashierConn.Close()
		merchantConn.Close()
		return nil, err
	}

	orderItemConn, err := server.NewGRPCClient(orderItemAddr)
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

	var myKafka *kafka.Kafka
	if brokers := os.Getenv("KAFKA_BROKERS"); brokers != "" {
		myKafka = kafka.NewKafka(srv.Logger, []string{brokers})
	}

	mencacheObj := mencache.NewMencache(srv.CacheStore)

	traceLoggerObservability := observability.NewTraceLoggerObservability(srv.Logger)

	outboxService := outbox.NewOutboxService(srv.DB, myKafka, srv.Logger)

	services := service.NewService(&service.Deps{
		Mencache:      mencacheObj,
		Repositories:  repos,
		Logger:        srv.Logger,
		Kafka:         myKafka,
		Pool:          srv.DBPool,
		Outbox:        outboxService,
		Observability: traceLoggerObservability,
	})

	handlers := handler.NewHandler(&handler.Deps{
		Service: services,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pb.RegisterTransactionServiceServer(gs, handlers.Transaction)
	}

	// Start the outbox relay so events committed with the business writes are
	// published to Kafka with durable retry and dead-letter semantics.
	go outboxService.Start(srv.Ctx, outbox.OutboxRelayInterval, outbox.OutboxRelayBatchSize)

	return srv, nil
}
