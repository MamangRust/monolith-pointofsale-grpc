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

type cashierStatsService struct {
	ctx             context.Context
	trace           trace.Tracer
	cashierStats    repository.CashierStatsRepository
	logger          logger.LoggerInterface
	mapping         response_service.CashierResponseMapper
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewCashierStatsService(ctx context.Context, cashierStats repository.CashierStatsRepository,
	logger logger.LoggerInterface, mapping response_service.CashierResponseMapper,
) *cashierStatsService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cashier_stats_service_requests_total",
			Help: "Total number of requests to the CashierStatsService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cashier_stats_service_request_duration_seconds",
			Help:    "Histogram of request durations for the CashierStatsService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &cashierStatsService{
		ctx:             ctx,
		trace:           otel.Tracer("cashier-stats-service"),
		cashierStats:    cashierStats,
		logger:          logger,
		mapping:         mapping,
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}
}
func (s *cashierStatsService) FindMonthlyTotalSales(req *requests.MonthTotalSales) ([]*response.CashierResponseMonthTotalSales, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("FindMonthlyTotalSales", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyTotalSales")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month),
	)

	res, err := s.cashierStats.GetMonthlyTotalSales(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_MONTHLY_TOTAL_SALES")
		status = "failed_monthly_total_sales"

		s.logger.Error("failed to get monthly total sales",
			zap.Int("year", req.Year),
			zap.Int("month", req.Month),
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "monthly_sales_query_failed"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to get monthly total sales")

		return nil, cashier_errors.ErrFailedFindMonthlyTotalSales
	}

	span.SetAttributes(
		attribute.Int("result.count", len(res)),
	)

	return s.mapping.ToCashierMonthlyTotalSales(res), nil
}

func (s *cashierStatsService) FindYearlyTotalSales(year int) ([]*response.CashierResponseYearTotalSales, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("FindYearlyTotalSales", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyTotalSales")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	res, err := s.cashierStats.GetYearlyTotalSales(year)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_YEARLY_TOTAL_SALES")
		status = "failed_yearly_total_sales"

		s.logger.Error("failed to get yearly total sales",
			zap.Int("year", year),
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "yearly_sales_query_failed"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to get yearly total sales")

		return nil, cashier_errors.ErrFailedFindYearlyTotalSales
	}

	span.SetAttributes(
		attribute.Int("result.count", len(res)),
	)

	return s.mapping.ToCashierYearlyTotalSales(res), nil
}

func (s *cashierStatsService) FindMonthlySales(year int) ([]*response.CashierResponseMonthSales, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("FindMonthlySales", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlySales")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	res, err := s.cashierStats.GetMonthyCashier(year)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_MONTHLY_SALES")
		status = "failed_monthly_sales"

		s.logger.Error("failed to get monthly cashier sales",
			zap.Int("year", year),
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "monthly_cashier_query_failed"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to get monthly sales")

		return nil, cashier_errors.ErrFailedFindMonthlySales
	}

	span.SetAttributes(
		attribute.Int("result.count", len(res)),
	)

	return s.mapping.ToCashierMonthlySales(res), nil
}

func (s *cashierStatsService) FindYearlySales(year int) ([]*response.CashierResponseYearSales, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("FindYearlySales", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlySales")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
	)

	res, err := s.cashierStats.GetYearlyCashier(year)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_YEARLY_SALES")
		status = "failed_yearly_sales"

		s.logger.Error("failed to get yearly cashier sales",
			zap.Int("year", year),
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "yearly_cashier_query_failed"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to get yearly sales")

		return nil, cashier_errors.ErrFailedFindYearlySales
	}

	span.SetAttributes(
		attribute.Int("result.count", len(res)),
	)

	return s.mapping.ToCashierYearlySales(res), nil
}

func (s *cashierStatsService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
