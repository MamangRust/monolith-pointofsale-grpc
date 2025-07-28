package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-cashier/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-point-of-sale-cashier/internal/redis"
	"github.com/MamangRust/monolith-point-of-sale-cashier/internal/repository"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
	response_service "github.com/MamangRust/monolith-point-of-sale-shared/mapper/response/service"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type cashierStatsByMerchantService struct {
	mencache        mencache.CashierStatsByMerchantCache
	errorhandler    errorhandler.CashierStatsByMerchantError
	trace           trace.Tracer
	cashierStats    repository.CashierStatByMerchantRepository
	logger          logger.LoggerInterface
	mapping         response_service.CashierResponseMapper
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewCashierStatsByMerchantService(
	mencache mencache.CashierStatsByMerchantCache,
	errorhandler errorhandler.CashierStatsByMerchantError,
	cashierStats repository.CashierStatByMerchantRepository,
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
		trace:           otel.Tracer("cashier-stats-by-merchant-service"),
		cashierStats:    cashierStats,
		logger:          logger,
		mapping:         mapping,
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}
}

func (s *cashierStatsByMerchantService) FindMonthlyTotalSalesByMerchant(ctx context.Context, req *requests.MonthTotalSalesMerchant) ([]*response.CashierResponseMonthTotalSales, *response.ErrorResponse) {
	const method = "FindMonthlyTotalSalesByMerchant"

	month := req.Month
	year := req.Year

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.Int("month", month))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthlyTotalSalesByMerchantCache(ctx, req); found {
		logSuccess("Successfully fetched monthly total sales by ID from cache", zap.Int("year", year), zap.Int("month", month))

		return data, nil
	}

	res, err := s.cashierStats.GetMonthlyTotalSalesByMerchant(ctx, req)

	if err != nil {
		return s.errorhandler.HandleMonthlyTotalSalesByMerchantError(err, method, "FAILED_FIND_MONTHLY_TOTAL_SALES_BY_MERCHANT", span, &status, zap.Error(err))
	}

	so := s.mapping.ToCashierMonthlyTotalSales(res)

	s.mencache.SetMonthlyTotalSalesByMerchantCache(ctx, req, so)

	logSuccess("Successfully fetched monthly total sales by ID", zap.Int("year", year), zap.Int("month", month))

	return so, nil
}

func (s *cashierStatsByMerchantService) FindYearlyTotalSalesByMerchant(ctx context.Context, req *requests.YearTotalSalesMerchant) ([]*response.CashierResponseYearTotalSales, *response.ErrorResponse) {
	const method = "FindMonthlyTotalSalesByMerchant"

	year := req.Year
	merchant_id := req.MerchantID

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.Int("merchant_id", merchant_id))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearlyTotalSalesByMerchantCache(ctx, req); found {
		logSuccess("Successfully fetched yearly total sales by merchant id from cache", zap.Int("year", year), zap.Int("merchant.id", merchant_id))

		return data, nil
	}

	res, err := s.cashierStats.GetYearlyTotalSalesByMerchant(ctx, req)

	if err != nil {
		return s.errorhandler.HandleYearlyTotalSalesByMerchantError(err, method, "FAILED_FIND_YEARLY_TOTAL_SALES_BY_MERCHANT", span, &status, zap.Error(err))
	}

	so := s.mapping.ToCashierYearlyTotalSales(res)

	s.mencache.SetYearlyTotalSalesByMerchantCache(ctx, req, so)

	logSuccess("Successfully fetched yearly total sales by merchant id", zap.Int("year", year), zap.Int("merchant.id", merchant_id))

	return so, nil
}

func (s *cashierStatsByMerchantService) FindMonthlyCashierByMerchant(ctx context.Context, req *requests.MonthCashierMerchant) ([]*response.CashierResponseMonthSales, *response.ErrorResponse) {
	const method = "FindMonthlyCashierByMerchant"

	year := req.Year
	merchant_id := req.MerchantID

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.Int("merchant.id", merchant_id))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthlyCashierByMerchantCache(ctx, req); found {
		logSuccess("Successfully fetched monthly cashier sales by ID from cache", zap.Int("year", year), zap.Int("merchant.id", merchant_id))

		return data, nil
	}

	res, err := s.cashierStats.GetMonthlyCashierByMerchant(ctx, req)

	if err != nil {
		return s.errorhandler.HandleMonthlySalesByMerchantError(err, method, "FAILED_FIND_MONTHLY_CASHIER_BY_MERCHANT", span, &status, zap.Error(err))
	}

	so := s.mapping.ToCashierMonthlySales(res)

	s.mencache.SetMonthlyCashierByMerchantCache(ctx, req, so)

	logSuccess("Successfully fetched monthly cashier sales by ID", zap.Int("year", year), zap.Int("merchant.id", merchant_id))

	return so, nil
}

func (s *cashierStatsByMerchantService) FindYearlyCashierByMerchant(ctx context.Context, req *requests.YearCashierMerchant) ([]*response.CashierResponseYearSales, *response.ErrorResponse) {
	const method = "FindMonthlyCashierByMerchant"

	year := req.Year
	merchant_id := req.MerchantID

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.Int("merchant.id", merchant_id))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearlyCashierByMerchantCache(ctx, req); found {
		logSuccess("Successfully fetched yearly cashier sales by Merchant ID from cache", zap.Int("year", year), zap.Int("merchant.id", merchant_id))

		return data, nil
	}

	res, err := s.cashierStats.GetYearlyCashierByMerchant(ctx, req)

	if err != nil {
		return s.errorhandler.HandleYearlySalesByMerchantError(err, method, "FAILED_FIND_YEARLY_CASHIER_BY_MERCHANT", span, &status, zap.Error(err))
	}

	so := s.mapping.ToCashierYearlySales(res)

	s.mencache.SetYearlyCashierByMerchantCache(ctx, req, so)

	logSuccess("Successfully fetched yearly cashier sales by Merchant ID", zap.Int("year", year), zap.Int("merchant.id", merchant_id))

	return so, nil
}

func (s *cashierStatsByMerchantService) startTracingAndLogging(ctx context.Context, method string, attrs ...attribute.KeyValue) (
	context.Context,
	trace.Span,
	func(string),
	string,
	func(string, ...zap.Field),
) {
	start := time.Now()
	status := "success"

	ctx, span := s.trace.Start(ctx, method)

	if len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}

	span.AddEvent("Start: " + method)

	s.logger.Debug("Start: " + method)

	end := func(status string) {
		s.recordMetrics(method, status, start)
		code := codes.Ok
		if status != "success" {
			code = codes.Error
		}
		span.SetStatus(code, status)
		span.End()
	}

	logSuccess := func(msg string, fields ...zap.Field) {
		span.AddEvent(msg)
		s.logger.Debug(msg, fields...)
	}

	return ctx, span, end, status, logSuccess
}

func (s *cashierStatsByMerchantService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
