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

type cashierStatsService struct {
	mencache        mencache.CashierStatsCache
	errorhandler    errorhandler.CashierStatsError
	trace           trace.Tracer
	cashierStats    repository.CashierStatsRepository
	logger          logger.LoggerInterface
	mapping         response_service.CashierResponseMapper
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewCashierStatsService(
	mencache mencache.CashierStatsCache,
	errorhandler errorhandler.CashierStatsError,
	cashierStats repository.CashierStatsRepository,
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
		mencache:        mencache,
		errorhandler:    errorhandler,
		trace:           otel.Tracer("cashier-stats-service"),
		cashierStats:    cashierStats,
		logger:          logger,
		mapping:         mapping,
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}
}
func (s *cashierStatsService) FindMonthlyTotalSales(ctx context.Context, req *requests.MonthTotalSales) ([]*response.CashierResponseMonthTotalSales, *response.ErrorResponse) {
	const method = "FindMonthlyTotalSales"

	month := req.Month
	year := req.Year

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("month", month), attribute.Int("year", year))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthlyTotalSalesCache(ctx, req); found {
		logSuccess("Fetched monthly total sales from cache", zap.Int("month", month), zap.Int("year", year))
		return data, nil
	}

	res, err := s.cashierStats.GetMonthlyTotalSales(ctx, req)
	if err != nil {
		return s.errorhandler.HandleMonthlyTotalSalesError(err, method, "FAILED_FIND_MONTHLY_TOTAL_SALES", span, &status, zap.Error(err))
	}

	mapped := s.mapping.ToCashierMonthlyTotalSales(res)
	s.mencache.SetMonthlyTotalSalesCache(ctx, req, mapped)

	logSuccess("Fetched monthly total sales from DB", zap.Int("month", month), zap.Int("year", year))
	return mapped, nil
}

func (s *cashierStatsService) FindYearlyTotalSales(ctx context.Context, year int) ([]*response.CashierResponseYearTotalSales, *response.ErrorResponse) {
	const method = "FindYearlyTotalSales"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("year", year))
	defer end(status)

	if data, found := s.mencache.GetYearlyTotalSalesCache(ctx, year); found {
		logSuccess("Fetched yearly total sales from cache", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cashierStats.GetYearlyTotalSales(ctx, year)
	if err != nil {
		return s.errorhandler.HandleYearlyTotalSalesError(err, method, "FAILED_FIND_YEARLY_TOTAL_SALES", span, &status, zap.Error(err))
	}

	mapped := s.mapping.ToCashierYearlyTotalSales(res)
	s.mencache.SetYearlyTotalSalesCache(ctx, year, mapped)

	logSuccess("Fetched yearly total sales from DB", zap.Int("year", year))
	return mapped, nil
}

func (s *cashierStatsService) FindMonthlySales(ctx context.Context, year int) ([]*response.CashierResponseMonthSales, *response.ErrorResponse) {
	const method = "FindMonthlySales"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("year", year))
	defer end(status)

	if data, found := s.mencache.GetMonthlySalesCache(ctx, year); found {
		logSuccess("Fetched monthly sales from cache", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cashierStats.GetMonthyCashier(ctx, year)

	if err != nil {
		return s.errorhandler.HandleMonthlySalesError(err, method, "FAILED_FIND_MONTHLY_SALES", span, &status, zap.Error(err))
	}

	mapped := s.mapping.ToCashierMonthlySales(res)
	s.mencache.SetMonthlySalesCache(ctx, year, mapped)

	logSuccess("Fetched monthly sales from DB", zap.Int("year", year))
	return mapped, nil
}

func (s *cashierStatsService) FindYearlySales(ctx context.Context, year int) ([]*response.CashierResponseYearSales, *response.ErrorResponse) {
	const method = "FindYearlySales"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("year", year))
	defer end(status)

	if data, found := s.mencache.GetYearlySalesCache(ctx, year); found {
		logSuccess("Fetched yearly sales from cache", zap.Int("year", year))
		return data, nil
	}

	res, err := s.cashierStats.GetYearlyCashier(ctx, year)
	if err != nil {
		return s.errorhandler.HandleYearlySalesError(err, method, "FAILED_FIND_YEARLY_SALES", span, &status, zap.Error(err))
	}

	mapped := s.mapping.ToCashierYearlySales(res)
	s.mencache.SetYearlySalesCache(ctx, year, mapped)

	logSuccess("Fetched yearly sales from DB", zap.Int("year", year))
	return mapped, nil
}

func (s *cashierStatsService) startTracingAndLogging(ctx context.Context, method string, attrs ...attribute.KeyValue) (
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

func (s *cashierStatsService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
