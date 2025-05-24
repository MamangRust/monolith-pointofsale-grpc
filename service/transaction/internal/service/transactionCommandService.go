package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	traceunic "github.com/MamangRust/monolith-point-of-sale-pkg/trace_unic"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/cashier_errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/merchant_errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/order_errors"
	orderitem_errors "github.com/MamangRust/monolith-point-of-sale-shared/errors/order_item_errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/transaction_errors"
	response_service "github.com/MamangRust/monolith-point-of-sale-shared/mapper/response/service"
	"github.com/MamangRust/monolith-point-of-sale-transacton/internal/repository"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type transactionCommandService struct {
	ctx                          context.Context
	trace                        trace.Tracer
	cashierQueryRepository       repository.CashierQueryRepository
	merchantQueryRepository      repository.MerchantQueryRepository
	transactionQueryRepository   repository.TransactionQueryRepository
	transactionCommandRepository repository.TransactionCommandRepository
	orderQueryRepository         repository.OrderQueryRepository
	orderItemQueryRepository     repository.OrderItemQueryRepository
	mapping                      response_service.TransactionResponseMapper
	logger                       logger.LoggerInterface
	requestCounter               *prometheus.CounterVec
	requestDuration              *prometheus.HistogramVec
}

func NewTransactionCommandService(
	ctx context.Context,
	cashierQueryRepository repository.CashierQueryRepository,
	merchantQueryRepository repository.MerchantQueryRepository,
	transactionQueryRepository repository.TransactionQueryRepository,
	transactionCommandRepository repository.TransactionCommandRepository,
	orderQueryRepository repository.OrderQueryRepository,
	orderItemQueryRepository repository.OrderItemQueryRepository,
	mapping response_service.TransactionResponseMapper,
	logger logger.LoggerInterface,
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
			Name:    "transaction_command_service_request_duration",
			Help:    "Histogram of request durations for the TransactionCommandService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &transactionCommandService{
		ctx:                          ctx,
		trace:                        otel.Tracer("transaction-command-service"),
		cashierQueryRepository:       cashierQueryRepository,
		merchantQueryRepository:      merchantQueryRepository,
		transactionQueryRepository:   transactionQueryRepository,
		transactionCommandRepository: transactionCommandRepository,
		orderQueryRepository:         orderQueryRepository,
		orderItemQueryRepository:     orderItemQueryRepository,
		mapping:                      mapping,
		logger:                       logger,
		requestCounter:               requestCounter,
		requestDuration:              requestDuration,
	}
}

func (s *transactionCommandService) CreateTransaction(req *requests.CreateTransactionRequest) (*response.TransactionResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("CreateTransaction", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "CreateTransaction")
	defer span.End()

	span.SetAttributes(
		attribute.Int("order.id", req.OrderID),
		attribute.Int("cashier.id", req.CashierID),
	)

	s.logger.Debug("Creating new transaction",
		zap.Int("orderID", req.OrderID))

	cashier, err := s.cashierQueryRepository.FindById(req.CashierID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_CASHIER")
		status = "failed_find_cashier"

		s.logger.Error("Cashier not found",
			zap.Int("cashierId", req.CashierID),
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "cashier_not_found"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Cashier not found")

		return nil, cashier_errors.ErrFailedFindCashierById
	}

	_, err = s.merchantQueryRepository.FindById(cashier.MerchantID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MERCHANT")
		status = "failed_find_merchant"

		s.logger.Error("Merchant not found",
			zap.Int("merchantId", cashier.MerchantID),
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "merchant_not_found"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Merchant not found")

		return nil, merchant_errors.ErrFailedFindMerchantById
	}
	req.MerchantID = cashier.MerchantID

	_, err = s.orderQueryRepository.FindById(req.OrderID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ORDER")
		status = "failed_find_order"

		s.logger.Error("Order not found",
			zap.Int("orderID", req.OrderID),
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "order_not_found"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Order not found")

		return nil, order_errors.ErrFailedFindOrderById
	}

	orderItems, err := s.orderItemQueryRepository.FindOrderItemByOrder(req.OrderID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ORDER_ITEMS")
		status = "failed_find_order_items"

		s.logger.Error("Failed to retrieve order items",
			zap.Int("orderID", req.OrderID),
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "order_items_not_found"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find order items")

		return nil, orderitem_errors.ErrFailedFindOrderItemByOrder
	}

	if len(orderItems) == 0 {
		traceID := traceunic.GenerateTraceID("EMPTY_ORDER_ITEMS")
		status = "empty_order_items"

		s.logger.Error("Order items empty",
			zap.Int("orderID", req.OrderID),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "empty_order_items"),
		)
		span.SetStatus(codes.Error, "Empty order items")

		return nil, orderitem_errors.ErrFailedOrderItemEmpty
	}

	var totalAmount int
	for _, item := range orderItems {
		if item.Quantity <= 0 {
			traceID := traceunic.GenerateTraceID("INVALID_ORDER_ITEM_QUANTITY")
			status = "invalid_order_item_quantity"

			s.logger.Error("Invalid order item quantity",
				zap.Int("orderID", req.OrderID),
				zap.Int("quantity", item.Quantity),
				zap.String("traceID", traceID))

			span.SetAttributes(
				attribute.String("error.trace_id", traceID),
				attribute.String("error.type", "invalid_quantity"),
			)
			span.SetStatus(codes.Error, "Invalid order item quantity")

			return nil, orderitem_errors.ErrFailedFindOrderItemByOrder
		}
		totalAmount += item.Price * item.Quantity
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
		traceID := traceunic.GenerateTraceID("INSUFFICIENT_PAYMENT")
		status = "insufficient_payment"

		s.logger.Error("Insufficient payment amount",
			zap.Int("amount.paid", req.Amount),
			zap.Int("amount.required", totalAmountWithTax),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "insufficient_payment"),
			attribute.Int("amount.paid", req.Amount),
			attribute.Int("amount.required", totalAmountWithTax),
		)
		span.SetStatus(codes.Error, "Insufficient payment")

		return nil, transaction_errors.ErrFailedPaymentInsufficientBalance
	}

	req.Amount = totalAmountWithTax
	req.PaymentStatus = &paymentStatus

	transaction, err := s.transactionCommandRepository.CreateTransaction(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_CREATE_TRANSACTION")
		status = "failed_create_transaction"

		s.logger.Error("Failed to create transaction record",
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "create_transaction_failed"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to create transaction")

		return nil, transaction_errors.ErrFailedCreateTransaction
	}

	span.SetAttributes(
		attribute.Int("transaction.id", transaction.ID),
		attribute.String("transaction.status", paymentStatus),
		attribute.Int("transaction.amount", transaction.Amount),
	)

	s.logger.Debug("Transaction created successfully",
		zap.Int("transactionID", transaction.ID),
		zap.String("status", paymentStatus),
		zap.Int("amount", transaction.Amount))

	return s.mapping.ToTransactionResponse(transaction), nil
}

func (s *transactionCommandService) UpdateTransaction(req *requests.UpdateTransactionRequest) (*response.TransactionResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("UpdateTransaction", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "UpdateTransaction")
	defer span.End()

	span.SetAttributes(
		attribute.Int("transaction.id", *req.TransactionID),
		attribute.Int("cashier.id", req.CashierID),
		attribute.Int("order.id", req.OrderID),
	)

	s.logger.Debug("Updating transaction",
		zap.Int("transactionID", *req.TransactionID))

	cashier, err := s.cashierQueryRepository.FindById(req.CashierID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_CASHIER")
		status = "failed_find_cashier"

		s.logger.Error("Cashier not found",
			zap.Int("cashierId", req.CashierID),
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "cashier_not_found"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Cashier not found")

		return nil, cashier_errors.ErrFailedFindCashierById
	}

	existingTx, err := s.transactionQueryRepository.FindById(*req.TransactionID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_TRANSACTION")
		status = "failed_find_transaction"

		s.logger.Error("Transaction not found",
			zap.Int("transactionID", *req.TransactionID),
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "transaction_not_found"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Transaction not found")

		return nil, transaction_errors.ErrFailedFindTransactionById
	}

	if existingTx.PaymentStatus == "paid" || existingTx.PaymentStatus == "refunded" {
		traceID := traceunic.GenerateTraceID("INVALID_PAYMENT_STATUS")
		status = "invalid_payment_status"

		s.logger.Error("Payment status cannot be modified",
			zap.String("current_status", existingTx.PaymentStatus),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "immutable_payment_status"),
			attribute.String("current.status", existingTx.PaymentStatus),
		)
		span.SetStatus(codes.Error, "Payment status cannot be modified")

		return nil, transaction_errors.ErrFailedPaymentStatusCannotBeModified
	}

	_, err = s.merchantQueryRepository.FindById(cashier.MerchantID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MERCHANT")
		status = "failed_find_merchant"

		s.logger.Error("Merchant not found",
			zap.Int("merchantId", cashier.MerchantID),
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "merchant_not_found"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Merchant not found")

		return nil, merchant_errors.ErrFailedFindMerchantById
	}
	req.MerchantID = cashier.MerchantID

	_, err = s.orderQueryRepository.FindById(req.OrderID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ORDER")
		status = "failed_find_order"

		s.logger.Error("Order not found",
			zap.Int("orderID", req.OrderID),
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "order_not_found"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Order not found")

		return nil, order_errors.ErrFailedFindOrderById
	}

	orderItems, err := s.orderItemQueryRepository.FindOrderItemByOrder(req.OrderID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ORDER_ITEMS")
		status = "failed_find_order_items"

		s.logger.Error("Failed to retrieve order items",
			zap.Int("orderID", req.OrderID),
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "order_items_not_found"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find order items")

		return nil, orderitem_errors.ErrFailedFindOrderItemByOrder
	}

	var totalAmount int
	for _, item := range orderItems {
		totalAmount += item.Price * item.Quantity
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
		traceID := traceunic.GenerateTraceID("INSUFFICIENT_PAYMENT")
		status = "insufficient_payment"

		s.logger.Error("Insufficient payment amount",
			zap.Int("amount.paid", req.Amount),
			zap.Int("amount.required", totalAmountWithTax),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "insufficient_payment"),
			attribute.Int("amount.paid", req.Amount),
			attribute.Int("amount.required", totalAmountWithTax),
		)
		span.SetStatus(codes.Error, "Insufficient payment")

		return nil, transaction_errors.ErrFailedPaymentInsufficientBalance
	}

	req.Amount = totalAmountWithTax
	req.PaymentStatus = &paymentStatus

	transaction, err := s.transactionCommandRepository.UpdateTransaction(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_TRANSACTION")
		status = "failed_update_transaction"

		s.logger.Error("Failed to update transaction",
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "update_transaction_failed"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update transaction")

		return nil, transaction_errors.ErrFailedUpdateTransaction
	}

	span.SetAttributes(
		attribute.String("transaction.status", *req.PaymentStatus),
		attribute.Int("transaction.amount", transaction.Amount),
		attribute.String("transaction.previous_status", existingTx.PaymentStatus),
	)

	s.logger.Debug("Transaction updated successfully",
		zap.Int("transactionID", transaction.ID),
		zap.String("status", *req.PaymentStatus),
		zap.Int("amount", transaction.Amount))

	return s.mapping.ToTransactionResponse(transaction), nil
}

func (s *transactionCommandService) TrashedTransaction(transactionID int) (*response.TransactionResponseDeleteAt, *response.ErrorResponse) {
	start := time.Now()
	status := "success"
	defer func() { s.recordMetrics("TrashedTransaction", status, start) }()

	_, span := s.trace.Start(s.ctx, "TrashedTransaction")
	defer span.End()

	span.SetAttributes(attribute.Int("transaction.id", transactionID))
	s.logger.Debug("Trashing transaction", zap.Int("transaction_id", transactionID))

	res, err := s.transactionCommandRepository.TrashTransaction(transactionID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_TRASH_TRANSACTION")
		s.logger.Error("Failed to move transaction to trash",
			zap.Int("transaction_id", transactionID), zap.String("trace_id", traceID), zap.Error(err))

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to trash transaction")
		status = "failed_trash"
		return nil, transaction_errors.ErrFailedTrashedTransaction
	}

	span.SetStatus(codes.Ok, "Transaction trashed")
	s.logger.Debug("Successfully trashed transaction", zap.Int("transaction_id", transactionID))
	return s.mapping.ToTransactionResponseDeleteAt(res), nil
}

func (s *transactionCommandService) RestoreTransaction(transactionID int) (*response.TransactionResponseDeleteAt, *response.ErrorResponse) {
	start := time.Now()
	status := "success"
	defer func() { s.recordMetrics("RestoreTransaction", status, start) }()

	_, span := s.trace.Start(s.ctx, "RestoreTransaction")
	defer span.End()

	span.SetAttributes(attribute.Int("transaction.id", transactionID))
	s.logger.Debug("Restoring transaction", zap.Int("transaction_id", transactionID))

	res, err := s.transactionCommandRepository.RestoreTransaction(transactionID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_TRANSACTION")
		s.logger.Error("Failed to restore transaction",
			zap.Int("transaction_id", transactionID), zap.String("trace_id", traceID), zap.Error(err))

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore transaction")
		status = "failed_restore"
		return nil, transaction_errors.ErrFailedRestoreTransaction
	}

	span.SetStatus(codes.Ok, "Transaction restored")
	s.logger.Debug("Successfully restored transaction", zap.Int("transaction_id", transactionID))
	return s.mapping.ToTransactionResponseDeleteAt(res), nil
}

func (s *transactionCommandService) DeleteTransactionPermanently(transactionID int) (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"
	defer func() { s.recordMetrics("DeleteTransactionPermanently", status, start) }()

	_, span := s.trace.Start(s.ctx, "DeleteTransactionPermanently")
	defer span.End()

	span.SetAttributes(attribute.Int("transaction.id", transactionID))
	s.logger.Debug("Permanently deleting transaction", zap.Int("transactionID", transactionID))

	success, err := s.transactionCommandRepository.DeleteTransactionPermanently(transactionID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_TRANSACTION")
		s.logger.Error("Failed to permanently delete transaction",
			zap.Int("transaction_id", transactionID), zap.String("trace_id", traceID), zap.Error(err))

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to delete transaction")
		status = "failed_delete"
		return false, transaction_errors.ErrFailedDeleteTransactionPermanently
	}

	span.SetStatus(codes.Ok, "Transaction permanently deleted")
	return success, nil
}

func (s *transactionCommandService) RestoreAllTransactions() (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"
	defer func() { s.recordMetrics("RestoreAllTransactions", status, start) }()

	_, span := s.trace.Start(s.ctx, "RestoreAllTransactions")
	defer span.End()

	s.logger.Debug("Restoring all trashed transactions")

	success, err := s.transactionCommandRepository.RestoreAllTransactions()
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_ALL")
		s.logger.Error("Failed to restore all trashed transactions", zap.String("trace_id", traceID), zap.Error(err))

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore all transactions")
		status = "failed_restore_all"
		return false, transaction_errors.ErrFailedRestoreAllTransactions
	}

	span.SetStatus(codes.Ok, "All transactions restored")
	return success, nil
}

func (s *transactionCommandService) DeleteAllTransactionPermanent() (bool, *response.ErrorResponse) {
	start := time.Now()
	status := "success"
	defer func() { s.recordMetrics("DeleteAllTransactionPermanent", status, start) }()

	_, span := s.trace.Start(s.ctx, "DeleteAllTransactionPermanent")
	defer span.End()

	s.logger.Debug("Permanently deleting all transactions")

	success, err := s.transactionCommandRepository.DeleteAllTransactionPermanent()
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_ALL")
		s.logger.Error("Failed to permanently delete all transactions", zap.String("trace_id", traceID), zap.Error(err))

		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to delete all transactions")
		status = "failed_delete_all"
		return false, transaction_errors.ErrFailedDeleteAllTransactionPermanent
	}

	span.SetStatus(codes.Ok, "All transactions permanently deleted")
	return success, nil
}

func (s *transactionCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
