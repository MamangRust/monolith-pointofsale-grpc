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

type cashierStatsByIdService struct {
	ctx             context.Context
	mencache        mencache.CashierStatsByIdCache
	errorhandler    errorhandler.CashierStatsByIdError
	trace           trace.Tracer
	cashierStats    repository.CashierStatByIdRepository
	logger          logger.LoggerInterface
	mapping         response_service.CashierResponseMapper
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewCashierStatsByIdService(ctx context.Context,
	mencache mencache.CashierStatsByIdCache,
	errorhandler errorhandler.CashierStatsByIdError,
	cashierStats repository.CashierStatByIdRepository,
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
		mencache:        mencache,
		errorhandler:    errorhandler,
		trace:           otel.Tracer("cashier-stats-by-id-service"),
		cashierStats:    cashierStats,
		logger:          logger,
		mapping:         mapping,
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}
}

func (s *cashierStatsByIdService) FindMonthlyTotalSalesById(req *requests.MonthTotalSalesCashier) ([]*response.CashierResponseMonthTotalSales, *response.ErrorResponse) {
	const method = "FindMonthlyTotalSalesById"

	month := req.Month
	year := req.Year

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.Int("month", month))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthlyTotalSalesByIdCache(req); found {
		logSuccess("Successfully fetched monthly total sales by ID from cache", zap.Int("year", year), zap.Int("month", month))

		return data, nil
	}

	res, err := s.cashierStats.GetMonthlyTotalSalesById(req)

	if err != nil {
		return s.errorhandler.HandleMonthlyTotalSalesByIdError(err, method, "FAILED_FIND_MONTHLY_TOTAL_SALES_BY_ID", span, &status, zap.Error(err))
	}

	so := s.mapping.ToCashierMonthlyTotalSales(res)

	s.mencache.SetMonthlyTotalSalesByIdCache(req, so)

	logSuccess("Successfully fetched monthly total sales by ID", zap.Int("year", year), zap.Int("month", month))

	return so, nil
}

func (s *cashierStatsByIdService) FindYearlyTotalSalesById(req *requests.YearTotalSalesCashier) ([]*response.CashierResponseYearTotalSales, *response.ErrorResponse) {
	const method = "FindMonthlyTotalSalesById"

	year := req.Year
	cashier_id := req.CashierID

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.Int("cashier_id", cashier_id))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearlyTotalSalesByIdCache(req); found {
		logSuccess("Successfully fetched yearly total sales by ID from cache", zap.Int("year", year), zap.Int("cashier_id", cashier_id))

		return data, nil
	}

	res, err := s.cashierStats.GetYearlyTotalSalesById(req)

	if err != nil {
		return s.errorhandler.HandleYearlyTotalSalesByIdError(err, method, "FAILED_FIND_YEARLY_TOTAL_SALES_BY_ID", span, &status, zap.Error(err))
	}

	so := s.mapping.ToCashierYearlyTotalSales(res)

	s.mencache.SetYearlyTotalSalesByIdCache(req, so)

	logSuccess("Successfully fetched yearly total sales by ID", zap.Int("year", year), zap.Int("cashier_id", cashier_id))

	return so, nil
}

func (s *cashierStatsByIdService) FindMonthlyCashierById(req *requests.MonthCashierId) ([]*response.CashierResponseMonthSales, *response.ErrorResponse) {
	const method = "FindMonthlyCashierById"

	year := req.Year
	cashier_id := req.CashierID

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.Int("cashier.id", cashier_id))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthlyCashierByIdCache(req); found {
		logSuccess("Successfully fetched monthly cashier sales by ID from cache", zap.Int("year", year), zap.Int("cashier_id", cashier_id))

		return data, nil
	}

	res, err := s.cashierStats.GetMonthlyCashierById(req)

	if err != nil {
		return s.errorhandler.HandleMonthlySalesByIdError(err, method, "FAILED_FIND_MONTHLY_CASHIER_BY_ID", span, &status, zap.Error(err))
	}

	so := s.mapping.ToCashierMonthlySales(res)

	s.mencache.SetMonthlyCashierByIdCache(req, so)

	logSuccess("Successfully fetched monthly cashier sales by ID", zap.Int("year", year), zap.Int("cashier_id", cashier_id))

	return so, nil
}

func (s *cashierStatsByIdService) FindYearlyCashierById(req *requests.YearCashierId) ([]*response.CashierResponseYearSales, *response.ErrorResponse) {
	const method = "FindMonthlyCashierById"

	year := req.Year
	cashier_id := req.CashierID

	span, end, status, logSuccess := s.startTracingAndLogging(method, attribute.Int("year", year), attribute.Int("cashier.id", cashier_id))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearlyCashierByIdCache(req); found {
		logSuccess("Successfully fetched yearly cashier sales by ID from cache", zap.Int("year", year), zap.Int("cashier_id", cashier_id))

		return data, nil
	}

	res, err := s.cashierStats.GetYearlyCashierById(req)

	if err != nil {
		return s.errorhandler.HandleYearlySalesByIdError(err, method, "FAILED_FIND_YEARLY_CASHIER_BY_ID", span, &status, zap.Error(err))
	}

	so := s.mapping.ToCashierYearlySales(res)

	s.mencache.SetYearlyCashierByIdCache(req, so)

	logSuccess("Successfully fetched yearly cashier sales by ID", zap.Int("year", year), zap.Int("cashier_id", cashier_id))

	return so, nil
}

func (s *cashierStatsByIdService) startTracingAndLogging(method string, attrs ...attribute.KeyValue) (
	trace.Span,
	func(string),
	string,
	func(string, ...zap.Field),
) {
	start := time.Now()
	status := "success"

	_, span := s.trace.Start(s.ctx, method)

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

	return span, end, status, logSuccess
}

func (s *cashierStatsByIdService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
