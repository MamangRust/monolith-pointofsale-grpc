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
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/cashier_errors"
	response_service "github.com/MamangRust/monolith-point-of-sale-shared/mapper/response/service"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type cashierQueryService struct {
	errorhandler    errorhandler.CashierQueryError
	mencache        mencache.CashierQueryCache
	trace           trace.Tracer
	cashierQuery    repository.CashierQueryRepository
	logger          logger.LoggerInterface
	mapping         response_service.CashierResponseMapper
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewCashierQueryService(
	errorhandler errorhandler.CashierQueryError,
	mencache mencache.CashierQueryCache,
	cashierQuery repository.CashierQueryRepository,
	logger logger.LoggerInterface, mapping response_service.CashierResponseMapper,
) *cashierQueryService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cashier_	uery_service_requests_total",
			Help: "Total number of requests to the CashierQueryService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cashier_query_service_request_duration_seconds",
			Help:    "Histogram of request durations for the CashierQueryService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &cashierQueryService{
		errorhandler:    errorhandler,
		mencache:        mencache,
		trace:           otel.Tracer("cashier-query-service"),
		cashierQuery:    cashierQuery,
		logger:          logger,
		mapping:         mapping,
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}
}

func (s *cashierQueryService) FindAll(ctx context.Context, req *requests.FindAllCashiers) ([]*response.CashierResponse, *int, *response.ErrorResponse) {
	const method = "FindAll"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedCashiersCache(ctx, req); found {
		logSuccess("Successfully fetched cashier from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	cashier, totalRecords, err := s.cashierQuery.FindAllCashiers(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_FIND_ALL_cashier", span, &status, zap.Error(err))
	}

	categoriesResponse := s.mapping.ToCashiersResponse(cashier)

	s.mencache.SetCachedCashiersCache(ctx, req, categoriesResponse, totalRecords)

	logSuccess("Successfully fetched cashier", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return categoriesResponse, totalRecords, nil
}

func (s *cashierQueryService) FindByMerchant(ctx context.Context, req *requests.FindAllCashierMerchant) ([]*response.CashierResponse, *int, *response.ErrorResponse) {
	const method = "FindByMerchant"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search), attribute.Int("merchant.id", req.MerchantID))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedCashiersByMerchant(ctx, req); found {
		logSuccess("Successfully fetched cashier from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	cashier, totalRecords, err := s.cashierQuery.FindByMerchant(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationError(err, method, "FAILED_FIND_ALL_cashier", span, &status, zap.Error(err))
	}

	categoriesResponse := s.mapping.ToCashiersResponse(cashier)

	s.mencache.SetCachedCashiersByMerchant(ctx, req, categoriesResponse, totalRecords)

	logSuccess("Successfully fetched cashier", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return categoriesResponse, totalRecords, nil
}

func (s *cashierQueryService) FindByActive(ctx context.Context, req *requests.FindAllCashiers) ([]*response.CashierResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByActive"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedCashiersActive(ctx, req); found {
		logSuccess("Successfully fetched cashier from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))

		return data, total, nil
	}

	cashier, totalRecords, err := s.cashierQuery.FindByActive(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_BY_ACTIVE_cashier", span, &status, cashier_errors.ErrFailedFindCashierByActive, zap.Error(err))
	}

	so := s.mapping.ToCashiersResponseDeleteAt(cashier)

	s.mencache.SetCachedCashiersActive(ctx, req, so, totalRecords)

	logSuccess("Successfully fetched cashier", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *cashierQueryService) FindByTrashed(ctx context.Context, req *requests.FindAllCashiers) ([]*response.CashierResponseDeleteAt, *int, *response.ErrorResponse) {
	const method = "FindByTrashed"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedCashiersTrashed(ctx, req); found {
		logSuccess("Successfully fetched categories from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))

		return data, total, nil
	}

	categories, totalRecords, err := s.cashierQuery.FindByTrashed(ctx, req)

	if err != nil {
		return s.errorhandler.HandleRepositoryPaginationDeleteAtError(err, method, "FAILED_FIND_BY_TRASHED_cashier", span, &status, cashier_errors.ErrFailedFindCashierByTrashed, zap.Error(err))
	}

	so := s.mapping.ToCashiersResponseDeleteAt(categories)

	s.mencache.SetCachedCashiersTrashed(ctx, req, so, totalRecords)

	logSuccess("Successfully fetched categories", zap.Int("totalRecords", *totalRecords), zap.Int("page", page), zap.Int("pageSize", pageSize))

	return so, totalRecords, nil
}

func (s *cashierQueryService) FindById(ctx context.Context, cashier_id int) (*response.CashierResponse, *response.ErrorResponse) {
	const method = "FindById"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("cashier.id", cashier_id))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedCashier(ctx, cashier_id); found {
		logSuccess("Successfully fetched cashier from cache", zap.Int("cashier.id", cashier_id))

		return data, nil
	}

	cashier, err := s.cashierQuery.FindById(ctx, cashier_id)

	if err != nil {
		return s.errorhandler.HandleRepositorySingleError(err, method, "FAILED_FIND_cashier_BY_ID", span, &status, cashier_errors.ErrFailedFindCashierById, zap.Error(err))
	}

	so := s.mapping.ToCashierResponse(cashier)

	s.mencache.SetCachedCashier(ctx, so)

	logSuccess("Successfully fetched cashier", zap.Int("cashier.id", cashier_id))

	return so, nil
}

func (s *cashierQueryService) startTracingAndLogging(ctx context.Context, method string, attrs ...attribute.KeyValue) (
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

	s.logger.Info("Start: " + method)

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

func (s *cashierQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}

func (s *cashierQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
