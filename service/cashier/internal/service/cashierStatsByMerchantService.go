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

type cashierStatsByMerchantService struct {
	ctx             context.Context
	trace           trace.Tracer
	cashierStats    repository.CashierStatByMerchantRepository
	logger          logger.LoggerInterface
	mapping         response_service.CashierResponseMapper
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewCashierStatsByMerchantService(ctx context.Context, cashierStats repository.CashierStatByMerchantRepository,
	logger logger.LoggerInterface, mapping response_service.CashierResponseMapper,
) *cashierStatsByMerchantService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cashier_stats_by_merchant_service_requests_total",
			Help: "Total number of requests to the CashierStatsByMerchantService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cashier_stats_by_merchant_service_request_duration_seconds",
			Help:    "Histogram of request durations for the CashierStatsByMerchantService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &cashierStatsByMerchantService{
		ctx:             ctx,
		trace:           otel.Tracer("cashier-stats-by-merchant-service"),
		cashierStats:    cashierStats,
		logger:          logger,
		mapping:         mapping,
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}
}

func (s *cashierStatsByMerchantService) FindMonthlyTotalSalesByMerchant(req *requests.MonthTotalSalesMerchant) ([]*response.CashierResponseMonthTotalSales, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyTotalSalesByMerchant", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyTotalSalesByMerchant")
	defer span.End()

	month := req.Month
	year := req.Year

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("month", month),
	)

	s.logger.Debug("FindMonthlyTotalSalesByMerchant",
		zap.Int("year", year),
		zap.Int("month", month))

	res, err := s.cashierStats.GetMonthlyTotalSalesByMerchant(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTHLY_TOTAL_SALES_BY_MERCHANT")

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

		status = "failed_find_monthly_total_sales_by_merchant"

		return nil, cashier_errors.ErrFailedFindMonthlyTotalSalesByMerchant
	}

	return s.mapping.ToCashierMonthlyTotalSales(res), nil
}

func (s *cashierStatsByMerchantService) FindYearlyTotalSalesByMerchant(req *requests.YearTotalSalesMerchant) ([]*response.CashierResponseYearTotalSales, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyTotalSalesByMerchant", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyTotalSalesByMerchant")
	defer span.End()

	year := req.Year
	merchant_id := req.MerchantID

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("merchant_id", merchant_id),
	)

	res, err := s.cashierStats.GetYearlyTotalSalesByMerchant(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEARLY_TOTAL_SALES_BY_merchant")

		s.logger.Error("failed to get yearly total sales",
			zap.Int("year", year),
			zap.Int("merchant_id", merchant_id),
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get yearly total sales")

		status = "failed_find_yearly_total_sales_by_merchant"

		return nil, cashier_errors.ErrFailedFindYearlyTotalSalesByMerchant
	}

	return s.mapping.ToCashierYearlyTotalSales(res), nil
}

func (s *cashierStatsByMerchantService) FindMonthlyCashierByMerchant(req *requests.MonthCashierMerchant) ([]*response.CashierResponseMonthSales, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindMonthlyCashierByMerchant", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindMonthlyCashierByMerchant")
	defer span.End()

	year := req.Year
	merchant_id := req.MerchantID

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("merchant_id", merchant_id),
	)

	res, err := s.cashierStats.GetMonthlyCashierByMerchant(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MONTHLY_CASHIER_BY_MERCHANT")

		s.logger.Error("failed to get monthly cashier sales by Merchant",
			zap.Int("year", year),
			zap.Int("merchant_id", merchant_id),
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get monthly cashier sales by Merchant")

		status = "failed_find_monthly_cashier_by_merchant"

		return nil, cashier_errors.ErrFailedFindMonthlyCashierByMerchant
	}

	return s.mapping.ToCashierMonthlySales(res), nil
}

func (s *cashierStatsByMerchantService) FindYearlyCashierByMerchant(req *requests.YearCashierMerchant) ([]*response.CashierResponseYearSales, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindYearlyCashierByMerchant", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindYearlyCashierByMerchant")
	defer span.End()

	year := req.Year
	merchant_id := req.MerchantID

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.Int("merchant_id", merchant_id),
	)

	s.logger.Debug("find yearly cashier sales by Merchant",
		zap.Int("year", year),
		zap.Int("merchant_id", merchant_id))

	res, err := s.cashierStats.GetYearlyCashierByMerchant(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_YEARLY_CASHIER_BY_MERCHANT")

		s.logger.Error("failed to get yearly cashier sales by Merchant",
			zap.Int("year", year),
			zap.Int("merchant_id", merchant_id),
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get yearly cashier sales by ID")

		status = "failed_find_yearly_cashier_by_merchant"

		return nil, cashier_errors.ErrFailedFindYearlyCashierByMerchant
	}

	return s.mapping.ToCashierYearlySales(res), nil
}

func (s *cashierStatsByMerchantService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
