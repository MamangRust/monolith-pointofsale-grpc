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

type cashierQueryService struct {
	ctx             context.Context
	trace           trace.Tracer
	cashierQuery    repository.CashierQueryRepository
	logger          logger.LoggerInterface
	mapping         response_service.CashierResponseMapper
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewCashierQueryService(ctx context.Context, cashierQuery repository.CashierQueryRepository,
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
		ctx:             ctx,
		trace:           otel.Tracer("cashier-query-service"),
		cashierQuery:    cashierQuery,
		logger:          logger,
		mapping:         mapping,
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}
}

func (s *cashierQueryService) FindAll(req *requests.FindAllCashiers) ([]*response.CashierResponse, *int, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("FindAllCashiers", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindAllCashiers")
	defer span.End()

	span.SetAttributes(
		attribute.Int("page.number", req.Page),
		attribute.Int("page.size", req.PageSize),
		attribute.String("search.term", req.Search),
	)

	s.logger.Debug("Fetching all cashiers",
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.PageSize),
		zap.String("search", req.Search))

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	cashiers, totalRecords, err := s.cashierQuery.FindAllCashiers(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ALL_CASHIERS")
		status = "failed_find_all_cashiers"

		s.logger.Error("Failed to fetch cashiers",
			zap.Error(err),
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "find_all_cashiers_failed"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find all cashiers")

		return nil, nil, cashier_errors.ErrFailedFindAllCashiers
	}

	span.SetAttributes(
		attribute.Int("result.count", len(cashiers)),
		attribute.Int("total.records", *totalRecords),
	)

	cashierResponse := s.mapping.ToCashiersResponse(cashiers)

	s.logger.Debug("Successfully fetched cashiers",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.PageSize))

	return cashierResponse, totalRecords, nil
}

func (s *cashierQueryService) FindById(cashierID int) (*response.CashierResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("FindCashierById", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindCashierById")
	defer span.End()

	span.SetAttributes(
		attribute.Int("cashier.id", cashierID),
	)

	s.logger.Debug("Fetching cashier by ID",
		zap.Int("cashierID", cashierID))

	cashier, err := s.cashierQuery.FindById(cashierID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_CASHIER_BY_ID")
		status = "failed_find_cashier_by_id"

		s.logger.Error("Failed to retrieve cashier details",
			zap.Error(err),
			zap.Int("cashier_id", cashierID),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "find_cashier_by_id_failed"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find cashier by ID")

		return nil, cashier_errors.ErrFailedFindCashierById
	}

	return s.mapping.ToCashierResponse(cashier), nil
}

func (s *cashierQueryService) FindByActive(req *requests.FindAllCashiers) ([]*response.CashierResponseDeleteAt, *int, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("FindActiveCashiers", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindActiveCashiers")
	defer span.End()

	span.SetAttributes(
		attribute.Int("page.number", req.Page),
		attribute.Int("page.size", req.PageSize),
		attribute.String("search.term", req.Search),
	)

	s.logger.Debug("Fetching active cashiers",
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.PageSize),
		zap.String("search", req.Search))

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	cashiers, totalRecords, err := s.cashierQuery.FindByActive(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ACTIVE_CASHIERS")
		status = "failed_find_active_cashiers"

		s.logger.Error("Failed to retrieve active cashiers",
			zap.Error(err),
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "find_active_cashiers_failed"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find active cashiers")

		return nil, nil, cashier_errors.ErrFailedFindCashierByActive
	}

	span.SetAttributes(
		attribute.Int("result.count", len(cashiers)),
		attribute.Int("total.records", *totalRecords),
	)

	s.logger.Debug("Successfully fetched active cashiers",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.PageSize))

	return s.mapping.ToCashiersResponseDeleteAt(cashiers), totalRecords, nil
}

func (s *cashierQueryService) FindByTrashed(req *requests.FindAllCashiers) ([]*response.CashierResponseDeleteAt, *int, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("FindTrashedCashiers", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindTrashedCashiers")
	defer span.End()

	span.SetAttributes(
		attribute.Int("page.number", req.Page),
		attribute.Int("page.size", req.PageSize),
		attribute.String("search.term", req.Search),
	)

	s.logger.Debug("Fetching trashed cashiers",
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.PageSize),
		zap.String("search", req.Search))

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	cashiers, totalRecords, err := s.cashierQuery.FindByTrashed(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_TRASHED_CASHIERS")
		status = "failed_find_trashed_cashiers"

		s.logger.Error("Failed to retrieve trashed cashiers",
			zap.Error(err),
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "find_trashed_cashiers_failed"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find trashed cashiers")

		return nil, nil, cashier_errors.ErrFailedFindCashierByTrashed
	}

	span.SetAttributes(
		attribute.Int("result.count", len(cashiers)),
		attribute.Int("total.records", *totalRecords),
	)

	s.logger.Debug("Successfully fetched trashed cashiers",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.PageSize))

	return s.mapping.ToCashiersResponseDeleteAt(cashiers), totalRecords, nil
}

func (s *cashierQueryService) FindByMerchant(req *requests.FindAllCashierMerchant) ([]*response.CashierResponse, *int, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("FindCashiersByMerchant", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindCashiersByMerchant")
	defer span.End()

	span.SetAttributes(
		attribute.Int("merchant.id", req.MerchantID),
		attribute.Int("page.number", req.Page),
		attribute.Int("page.size", req.PageSize),
		attribute.String("search.term", req.Search),
	)

	s.logger.Debug("Fetching merchant's cashiers",
		zap.Int("merchant_id", req.MerchantID),
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.PageSize),
		zap.String("search", req.Search))

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	cashiers, totalRecords, err := s.cashierQuery.FindByMerchant(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_CASHIERS_BY_MERCHANT")
		status = "failed_find_cashiers_by_merchant"

		s.logger.Error("Failed to retrieve merchant's cashiers",
			zap.Error(err),
			zap.Int("merchant_id", req.MerchantID),
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("search", req.Search),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "find_cashiers_by_merchant_failed"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find cashiers by merchant")

		return nil, nil, cashier_errors.ErrFailedFindCashierByMerchant
	}

	span.SetAttributes(
		attribute.Int("result.count", len(cashiers)),
		attribute.Int("total.records", *totalRecords),
	)

	s.logger.Debug("Successfully fetched merchant's cashiers",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.PageSize))

	return s.mapping.ToCashiersResponse(cashiers), totalRecords, nil
}

func (s *cashierQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
