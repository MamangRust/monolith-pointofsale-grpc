package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-cashier/internal/repository"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	traceunic "github.com/MamangRust/monolith-point-of-sale-pkg/trace_unic"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/cashier_errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/merchant_errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/user_errors"
	response_service "github.com/MamangRust/monolith-point-of-sale-shared/mapper/response/service"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type cashierCommandService struct {
	ctx             context.Context
	trace           trace.Tracer
	merchantQuery   repository.MerchantQueryRepository
	userQuery       repository.UserQueryRepository
	cashierCommand  repository.CashierCommandRepository
	mapping         response_service.CashierResponseMapper
	logger          logger.LoggerInterface
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewCashierCommandService(ctx context.Context, merchantQuery repository.MerchantQueryRepository,
	userQuery repository.UserQueryRepository, cashierCommand repository.CashierCommandRepository,
	mapping response_service.CashierResponseMapper, logger logger.LoggerInterface,
) *cashierCommandService {

	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cashier_command_service_requests_total",
			Help: "Total number of requests to the CashierCommandService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cashier_command_service_request_duration_seconds",
			Help:    "Histogram of request durations for the CashierCommandService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &cashierCommandService{
		ctx:             ctx,
		trace:           otel.Tracer("cashier-command-service"),
		merchantQuery:   merchantQuery,
		userQuery:       userQuery,
		cashierCommand:  cashierCommand,
		mapping:         mapping,
		logger:          logger,
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}
}

func (s *cashierCommandService) CreateCashier(req *requests.CreateCashierRequest) (*response.CashierResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("CreateCashier", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "CreateCashier")
	defer span.End()

	span.SetAttributes(
		attribute.Int("merchant.id", req.MerchantID),
		attribute.Int("user.id", req.UserID),
	)

	s.logger.Debug("Creating new cashier")

	_, err := s.merchantQuery.FindById(req.MerchantID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MERCHANT")
		status = "failed_find_merchant"

		s.logger.Error("Failed to retrieve merchant details",
			zap.Error(err),
			zap.Int("merchant_id", req.MerchantID),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "merchant_not_found"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Merchant not found")

		return nil, merchant_errors.ErrFailedFindMerchantById
	}

	_, err = s.userQuery.FindById(req.UserID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_USER")
		status = "failed_find_user"

		s.logger.Error("Failed to retrieve user details",
			zap.Error(err),
			zap.Int("user_id", req.UserID),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "user_not_found"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "User not found")

		return nil, user_errors.ErrFailedFindUserByID
	}

	cashier, err := s.cashierCommand.CreateCashier(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_CREATE_CASHIER")
		status = "failed_create_cashier"

		s.logger.Error("Failed to create new cashier",
			zap.Error(err),
			zap.Any("request", req),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "create_cashier_failed"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to create cashier")

		return nil, cashier_errors.ErrFailedCreateCashier
	}

	span.SetAttributes(
		attribute.Int("cashier.id", cashier.ID),
		attribute.String("cashier.name", cashier.Name),
	)

	return s.mapping.ToCashierResponse(cashier), nil
}

func (s *cashierCommandService) UpdateCashier(req *requests.UpdateCashierRequest) (*response.CashierResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("UpdateCashier", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "UpdateCashier")
	defer span.End()

	span.SetAttributes(
		attribute.Int("cashier.id", *req.CashierID),
	)

	s.logger.Debug("Updating cashier",
		zap.Int("cashier_id", *req.CashierID))

	cashier, err := s.cashierCommand.UpdateCashier(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_CASHIER")
		status = "failed_update_cashier"

		s.logger.Error("Failed to update cashier",
			zap.Error(err),
			zap.Int("cashier_id", *req.CashierID),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "update_cashier_failed"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update cashier")

		return nil, cashier_errors.ErrFailedUpdateCashier
	}

	span.SetAttributes(
		attribute.String("cashier.name", cashier.Name),
	)

	return s.mapping.ToCashierResponse(cashier), nil
}

func (s *cashierCommandService) TrashedCashier(cashierID int) (*response.CashierResponseDeleteAt, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("TrashedCashier", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "TrashedCashier")
	defer span.End()

	span.SetAttributes(
		attribute.Int("cashier.id", cashierID),
	)

	s.logger.Debug("Trashing cashier",
		zap.Int("cashierID", cashierID))

	cashier, err := s.cashierCommand.TrashedCashier(cashierID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_TRASH_CASHIER")
		status = "failed_trash_cashier"

		s.logger.Error("Failed to move cashier to trash",
			zap.Error(err),
			zap.Int("cashier_id", cashierID),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "trash_cashier_failed"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to trash cashier")

		return nil, cashier_errors.ErrFailedTrashedCashier
	}

	span.SetAttributes(
		attribute.String("cashier.deleted_at", *cashier.DeletedAt),
	)
	return s.mapping.ToCashierResponseDeleteAt(cashier), nil
}

func (s *cashierCommandService) RestoreCashier(cashierID int) (*response.CashierResponseDeleteAt, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("RestoreCashier", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreCashier")
	defer span.End()

	span.SetAttributes(
		attribute.Int("cashier.id", cashierID),
	)

	s.logger.Debug("Restoring cashier",
		zap.Int("cashierID", cashierID))

	cashier, err := s.cashierCommand.RestoreCashier(cashierID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_CASHIER")
		status = "failed_restore_cashier"

		s.logger.Error("Failed to restore cashier from trash",
			zap.Error(err),
			zap.Int("cashier_id", cashierID),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "restore_cashier_failed"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore cashier")

		return nil, cashier_errors.ErrFailedRestoreCashier
	}

	span.SetAttributes(
		attribute.Bool("cashier.restored", true),
	)
	return s.mapping.ToCashierResponseDeleteAt(cashier), nil
}

func (s *cashierCommandService) DeleteCashierPermanent(cashierID int) (bool, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("DeleteCashierPermanent", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteCashierPermanent")
	defer span.End()

	span.SetAttributes(
		attribute.Int("cashier.id", cashierID),
	)

	s.logger.Debug("Permanently deleting cashier",
		zap.Int("cashierID", cashierID))

	success, err := s.cashierCommand.DeleteCashierPermanent(cashierID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_CASHIER")
		status = "failed_delete_cashier"

		s.logger.Error("Failed to permanently delete cashier",
			zap.Error(err),
			zap.Int("cashier_id", cashierID),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "delete_cashier_failed"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to delete cashier")

		return false, cashier_errors.ErrFailedDeleteCashierPermanent
	}

	span.SetAttributes(
		attribute.Bool("operation.success", success),
	)
	return success, nil
}

func (s *cashierCommandService) RestoreAllCashier() (bool, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("RestoreAllCashier", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreAllCashier")
	defer span.End()

	s.logger.Debug("Restoring all trashed cashiers")

	success, err := s.cashierCommand.RestoreAllCashier()
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_ALL_CASHIERS")
		status = "failed_restore_all_cashiers"

		s.logger.Error("Failed to restore all trashed cashiers",
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "restore_all_cashiers_failed"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore all cashiers")

		return false, cashier_errors.ErrFailedRestoreAllCashiers
	}

	span.SetAttributes(
		attribute.Bool("operation.success", success),
	)
	return success, nil
}

func (s *cashierCommandService) DeleteAllCashierPermanent() (bool, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("DeleteAllCashierPermanent", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteAllCashierPermanent")
	defer span.End()

	s.logger.Debug("Permanently deleting all cashiers")

	success, err := s.cashierCommand.DeleteAllCashierPermanent()
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_ALL_CASHIERS")
		status = "failed_delete_all_cashiers"

		s.logger.Error("Failed to permanently delete all trashed cashiers",
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "delete_all_cashiers_failed"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to delete all cashiers")

		return false, cashier_errors.ErrFailedDeleteAllCashierPermanent
	}

	span.SetAttributes(
		attribute.Bool("operation.success", success),
	)
	return success, nil
}

func (s *cashierCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
