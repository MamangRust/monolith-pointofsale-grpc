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
	response_service "github.com/MamangRust/monolith-point-of-sale-shared/mapper/response/service"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type cashierStatsByIdService struct {
	ctx             context.Context
	trace           trace.Tracer
	cashierStats    repository.CashierStatByIdRepository
	logger          logger.LoggerInterface
	mapping         response_service.CashierResponseMapper
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewCashierStatsByIdService(ctx context.Context, cashierStats repository.CashierStatByIdRepository,
	logger logger.LoggerInterface, mapping response_service.CashierResponseMapper,
) *cashierStatsByIdService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cashier_stats_by_id_service_requests_total",
			Help: "Total number of requests to the CashierStatsByIdService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cashier_stats_by_id_service_request_duration_seconds",
			Help:    "Histogram of request durations for the CashierStatsByIdService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &cashierStatsByIdService{
		ctx:             ctx,
		trace:           otel.Tracer("cashier-stats-by-id-service"),
		cashierStats:    cashierStats,
		logger:          logger,
		mapping:         mapping,
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}
}

func (s *cashierStatsByIdService) FindMonthlyTotalSalesById(req *requests.MonthTotalSalesCashier) ([]*response.CashierResponseMonthTotalSales, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyTotalSalesById", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyTotalSalesById")
	defer span.End()

	month := req.Month
	year := req.Year

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("month", month),
	)

	s.logger.Debug("FindMonthlyTotalSalesById",
		zap.Int("year", year),
		zap.Int("month", month))

	res, err := s.cashierStats.GetMonthlyTotalSalesById(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTHLY_TOTAL_SALES_BY_ID")

		s.logger.Error("failed to get monthly total sales",
			zap.Int("year", year),
			zap.Int("month", month),
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get monthly total sales")

		status = "failed_find_monthly_total_sales_by_id"

		return nil, cashier_errors.ErrFailedFindMonthlyTotalSalesById
	}

	return s.mapping.ToCashierMonthlyTotalSales(res), nil
}

func (s *cashierStatsByIdService) FindYearlyTotalSalesById(req *requests.YearTotalSalesCashier) ([]*response.CashierResponseYearTotalSales, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyTotalSalesById", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyTotalSalesById")
	defer span.End()

	year := req.Year
	cashier_id := req.CashierID

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("cashier_id", cashier_id),
	)

	res, err := s.cashierStats.GetYearlyTotalSalesById(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEARLY_TOTAL_SALES_BY_ID")

		s.logger.Error("failed to get yearly total sales",
			zap.Int("year", year),
			zap.Int("cashier_id", cashier_id),
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get yearly total sales")

		status = "failed_find_yearly_total_sales_by_id"

		return nil, cashier_errors.ErrFailedFindYearlyTotalSalesById
	}

	return s.mapping.ToCashierYearlyTotalSales(res), nil
}

func (s *cashierStatsByIdService) FindMonthlyCashierById(req *requests.MonthCashierId) ([]*response.CashierResponseMonthSales, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyCashierById", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyCashierById")
	defer span.End()

	year := req.Year
	cashier_id := req.CashierID

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("cashier_id", cashier_id),
	)

	res, err := s.cashierStats.GetMonthlyCashierById(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTHLY_CASHIER_BY_ID")

		s.logger.Error("failed to get monthly cashier sales by ID",
			zap.Int("year", year),
			zap.Int("cashier_id", cashier_id),
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get monthly cashier sales by ID")

		status = "failed_find_monthly_cashier_by_id"

		return nil, cashier_errors.ErrFailedFindMonthlyCashierById
	}

	return s.mapping.ToCashierMonthlySales(res), nil
}

func (s *cashierStatsByIdService) FindYearlyCashierById(req *requests.YearCashierId) ([]*response.CashierResponseYearSales, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyCashierById", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyCashierById")
	defer span.End()

	year := req.Year
	cashier_id := req.CashierID

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("cashier_id", cashier_id),
	)

	s.logger.Debug("find yearly cashier sales by ID",
		zap.Int("year", year),
		zap.Int("cashier_id", cashier_id))

	res, err := s.cashierStats.GetYearlyCashierById(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEARLY_CASHIER_BY_ID")

		s.logger.Error("failed to get yearly cashier sales by ID",
			zap.Int("year", year),
			zap.Int("cashier_id", cashier_id),
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get yearly cashier sales by ID")

		status = "failed_find_yearly_cashier_by_id"

		return nil, cashier_errors.ErrFailedFindYearlyCashierById
	}

	return s.mapping.ToCashierYearlySales(res), nil
}

func (s *cashierStatsByIdService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
