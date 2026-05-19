package service

import (
	"context"
	"errors"

	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	sharederrorhandler "github.com/MamangRust/monolith-point-of-sale-shared/errorhandler"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/cashier_errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/merchant_errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/order_errors"
	orderitem_errors "github.com/MamangRust/monolith-point-of-sale-shared/errors/order_item_errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/transaction_errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	mencache "github.com/MamangRust/monolith-point-of-sale-transacton/internal/redis"
	"github.com/MamangRust/monolith-point-of-sale-transacton/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type transactionCommandService struct {
	mencache                     mencache.TransactionCommandCache
	trace                        trace.Tracer
	cashierQueryRepository       repository.CashierQueryRepository
	merchantQueryRepository      repository.MerchantQueryRepository
	transactionQueryRepository   repository.TransactionQueryRepository
	transactionCommandRepository repository.TransactionCommandRepository
	orderQueryRepository         repository.OrderQueryRepository
	orderItemQueryRepository     repository.OrderItemQueryRepository
	logger                       logger.LoggerInterface
	observability                observability.TraceLoggerObservability
	requestCounter               *prometheus.CounterVec
	requestDuration              *prometheus.HistogramVec
}

func NewTransactionCommandService(
	mencache mencache.TransactionCommandCache,
	cashierQueryRepository repository.CashierQueryRepository,
	merchantQueryRepository repository.MerchantQueryRepository,
	transactionQueryRepository repository.TransactionQueryRepository,
	transactionCommandRepository repository.TransactionCommandRepository,
	orderQueryRepository repository.OrderQueryRepository,
	orderItemQueryRepository repository.OrderItemQueryRepository,
	logger logger.LoggerInterface,
	obs observability.TraceLoggerObservability,
) *transactionCommandService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transaction_command_service_request_total",
			Help: "Total number of requests to the TransactionCommandService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transaction_command_service_request_duration_seconds",
			Help:    "Histogram of request durations for the TransactionCommandService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &transactionCommandService{
		mencache:                     mencache,
		trace:                        otel.Tracer("transaction-command-service"),
		cashierQueryRepository:       cashierQueryRepository,
		merchantQueryRepository:      merchantQueryRepository,
		transactionQueryRepository:   transactionQueryRepository,
		transactionCommandRepository: transactionCommandRepository,
		orderQueryRepository:         orderQueryRepository,
		orderItemQueryRepository:     orderItemQueryRepository,
		logger:                       logger,
		observability:                obs,
		requestCounter:               requestCounter,
		requestDuration:              requestDuration,
	}
}

func (s *transactionCommandService) CreateTransaction(ctx context.Context, req *requests.CreateTransactionRequest) (*db.Transaction, error) {
	const method = "CreateTransaction"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("cashier.id", req.CashierID), attribute.Int("merchant.id", req.MerchantID), attribute.Int("order.id", req.OrderID))

	defer func() {
		end(status)
	}()

	cashier, err := s.cashierQueryRepository.FindById(ctx, req.CashierID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Transaction](
			s.logger,
			cashier_errors.ErrFailedFindCashierById.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	_, err = s.merchantQueryRepository.FindById(ctx, int(cashier.MerchantID))
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Transaction](
			s.logger,
			merchant_errors.ErrFailedFindMerchantById.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}
	req.MerchantID = int(cashier.MerchantID)

	_, err = s.orderQueryRepository.FindById(ctx, req.OrderID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Transaction](
			s.logger,
			order_errors.ErrFailedFindOrderById.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	orderItems, err := s.orderItemQueryRepository.FindOrderItemByOrder(ctx, req.OrderID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Transaction](
			s.logger,
			orderitem_errors.ErrFailedFindOrderItemByOrder.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	if len(orderItems) == 0 {
		status = "error"
		errEmpty := errors.New("order item is empty")
		return sharederrorhandler.HandleError[*db.Transaction](
			s.logger,
			transaction_errors.ErrFailedOrderItemEmpty.WithInternal(errEmpty),
			method,
			span,
			zap.Error(errEmpty),
		)
	}

	var totalAmount int
	for _, item := range orderItems {
		if item.Quantity <= 0 {
			status = "error"
			errQty := errors.New("invalid quantity for order item")
			return sharederrorhandler.HandleError[*db.Transaction](
				s.logger,
				transaction_errors.ErrFailedPaymentStatusInvalid.WithInternal(errQty),
				method,
				span,
				zap.Error(errQty),
			)
		}
		totalAmount += int(item.Price) * int(item.Quantity)
	}

	ppn := totalAmount * 11 / 100
	totalAmountWithTax := totalAmount + ppn

	span.SetAttributes(
		attribute.Int("amount.subtotal", totalAmount),
		attribute.Int("amount.tax", ppn),
		attribute.Int("amount.total", totalAmountWithTax),
	)

	var paymentStatus string
	if req.Amount >= totalAmountWithTax {
		paymentStatus = "success"
	} else {
		status = "error"
		errBalance := errors.New("insufficient balance")
		return sharederrorhandler.HandleError[*db.Transaction](
			s.logger,
			transaction_errors.ErrFailedPaymentInsufficientBalance.WithInternal(errBalance),
			method,
			span,
			zap.Error(errBalance),
		)
	}

	req.Amount = totalAmountWithTax
	req.PaymentStatus = &paymentStatus

	transaction, err := s.transactionCommandRepository.CreateTransaction(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Transaction](
			s.logger,
			transaction_errors.ErrFailedCreateTransaction.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	logSuccess("Successfully created transaction", zap.Bool("success", true))

	return transaction, nil
}

func (s *transactionCommandService) UpdateTransaction(ctx context.Context, req *requests.UpdateTransactionRequest) (*db.Transaction, error) {
	const method = "UpdateTransaction"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("transaction.id", *req.TransactionID), attribute.Int("merchant.id", req.MerchantID), attribute.Int("order.id", req.OrderID))

	defer func() {
		end(status)
	}()

	cashier, err := s.cashierQueryRepository.FindById(ctx, req.CashierID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Transaction](
			s.logger,
			cashier_errors.ErrFailedFindCashierById.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	existingTx, err := s.transactionQueryRepository.FindById(ctx, *req.TransactionID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Transaction](
			s.logger,
			transaction_errors.ErrFailedFindTransactionById.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	if existingTx.PaymentStatus == "paid" || existingTx.PaymentStatus == "refunded" {
		status = "error"
		errStatus := errors.New("transaction payment status cannot be modified")
		return sharederrorhandler.HandleError[*db.Transaction](
			s.logger,
			transaction_errors.ErrFailedPaymentStatusCannotBeModified.WithInternal(errStatus),
			method,
			span,
			zap.Error(errStatus),
		)
	}

	_, err = s.merchantQueryRepository.FindById(ctx, int(cashier.MerchantID))
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Transaction](
			s.logger,
			merchant_errors.ErrFailedFindMerchantById.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}
	req.MerchantID = int(cashier.MerchantID)

	_, err = s.orderQueryRepository.FindById(ctx, req.OrderID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Transaction](
			s.logger,
			order_errors.ErrFailedFindOrderById.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	orderItems, err := s.orderItemQueryRepository.FindOrderItemByOrder(ctx, req.OrderID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Transaction](
			s.logger,
			orderitem_errors.ErrFailedFindOrderItemByOrder.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	var totalAmount int
	for _, item := range orderItems {
		totalAmount += int(item.Price) * int(item.Quantity)
	}

	ppn := totalAmount * 11 / 100
	totalAmountWithTax := totalAmount + ppn

	var paymentStatus string
	if req.Amount >= totalAmountWithTax {
		paymentStatus = "success"
	} else {
		status = "error"
		errBalance := errors.New("insufficient balance")
		return sharederrorhandler.HandleError[*db.Transaction](
			s.logger,
			transaction_errors.ErrFailedPaymentInsufficientBalance.WithInternal(errBalance),
			method,
			span,
			zap.Error(errBalance),
		)
	}

	req.Amount = totalAmountWithTax
	req.PaymentStatus = &paymentStatus

	transaction, err := s.transactionCommandRepository.UpdateTransaction(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Transaction](
			s.logger,
			transaction_errors.ErrFailedUpdateTransaction.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.DeleteTransactionCache(ctx, *req.TransactionID)

	logSuccess("Successfully updated transaction", zap.Bool("success", true))

	return transaction, nil
}

func (s *transactionCommandService) TrashedTransaction(ctx context.Context, transactionID int) (*db.Transaction, error) {
	const method = "TrashedTransaction"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("transaction.id", transactionID))

	defer func() {
		end(status)
	}()

	res, err := s.transactionCommandRepository.TrashTransaction(ctx, transactionID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Transaction](
			s.logger,
			transaction_errors.ErrFailedTrashedTransaction.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.DeleteTransactionCache(ctx, transactionID)

	logSuccess("Successfully trashed transaction", zap.Int("transaction.id", transactionID), zap.Bool("success", true))

	return res, nil
}

func (s *transactionCommandService) RestoreTransaction(ctx context.Context, transactionID int) (*db.Transaction, error) {
	const method = "RestoreTransaction"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("transaction.id", transactionID))

	defer func() {
		end(status)
	}()

	res, err := s.transactionCommandRepository.RestoreTransaction(ctx, transactionID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Transaction](
			s.logger,
			transaction_errors.ErrFailedRestoreTransaction.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.DeleteTransactionCache(ctx, transactionID)

	logSuccess("Successfully restored transaction", zap.Int("transaction.id", transactionID), zap.Bool("success", true))

	return res, nil
}

func (s *transactionCommandService) DeleteTransactionPermanently(ctx context.Context, transactionID int) (bool, error) {
	const method = "DeleteTransactionPermanently"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("transaction.id", transactionID))

	defer func() {
		end(status)
	}()

	success, err := s.transactionCommandRepository.DeleteTransactionPermanently(ctx, transactionID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](
			s.logger,
			transaction_errors.ErrFailedDeleteTransactionPermanently.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.DeleteTransactionCache(ctx, transactionID)

	logSuccess("Successfully permanently deleted transaction", zap.Int("transaction.id", transactionID), zap.Bool("success", success))

	return success, nil
}

func (s *transactionCommandService) RestoreAllTransactions(ctx context.Context) (bool, error) {
	const method = "RestoreAllTransactions"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.transactionCommandRepository.RestoreAllTransactions(ctx)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](
			s.logger,
			transaction_errors.ErrFailedRestoreAllTransactions.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	logSuccess("All trashed transactions restored successfully", zap.Bool("success", success))

	return success, nil
}

func (s *transactionCommandService) DeleteAllTransactionPermanent(ctx context.Context) (bool, error) {
	const method = "DeleteAllTransactionPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	success, err := s.transactionCommandRepository.DeleteAllTransactionPermanent(ctx)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](
			s.logger,
			transaction_errors.ErrFailedDeleteAllTransactionPermanent.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	logSuccess("Successfully permanently deleted all trashed transactions", zap.Bool("success", success))

	return success, nil
}
