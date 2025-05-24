package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-order/internal/repository"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	traceunic "github.com/MamangRust/monolith-point-of-sale-pkg/trace_unic"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/order_errors"
	response_service "github.com/MamangRust/monolith-point-of-sale-shared/mapper/response/service"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type orderQueryService struct {
	ctx                  context.Context
	trace                trace.Tracer
	orderQueryRepository repository.OrderQueryRepository
	logger               logger.LoggerInterface
	mapping              response_service.OrderResponseMapper
	requestCounter       *prometheus.CounterVec
	requestDuration      *prometheus.HistogramVec
}

func NewOrderQueryService(
	ctx context.Context,
	orderQueryRepository repository.OrderQueryRepository,
	logger logger.LoggerInterface,
	mapping response_service.OrderResponseMapper,
) *orderQueryService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "order_query_service_request_count",
			Help: "Total number of requests to the OrderQueryService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "order_query_service_request_duration",
			Help:    "Histogram of request durations for the OrderQueryService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &orderQueryService{
		ctx:                  ctx,
		trace:                otel.Tracer("order-query-service"),
		orderQueryRepository: orderQueryRepository,
		logger:               logger,
		mapping:              mapping,
		requestCounter:       requestCounter,
		requestDuration:      requestDuration,
	}
}

func (s *orderQueryService) FindAll(req *requests.FindAllOrders) ([]*response.OrderResponse, *int, *response.ErrorResponse) {
	start := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("FindAll", status, start)
	}()

	_, span := s.trace.Start(s.ctx, "FindAllOrder")
	defer span.End()

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search),
	)

	s.logger.Debug("Fetching all orders",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	orders, totalRecords, err := s.orderQueryRepository.FindAllOrders(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ALL_ORDERS")

		s.logger.Error("Failed to retrieve orders",
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(attribute.String("traceID", traceID))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to retrieve orders")

		status = "failed_find_all_orders"

		return nil, nil, order_errors.ErrFailedFindAllOrders
	}

	orderResponse := s.mapping.ToOrdersResponse(orders)

	s.logger.Debug("Successfully fetched order",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return orderResponse, totalRecords, nil
}

func (s *orderQueryService) FindById(order_id int) (*response.OrderResponse, *response.ErrorResponse) {
	s.logger.Debug("Fetching order by ID", zap.Int("order_id", order_id))

	order, err := s.orderQueryRepository.FindById(order_id)

	if err != nil {
		s.logger.Error("Failed to retrieve order details",
			zap.Error(err),
			zap.Int("order_id", order_id))

		return nil, order_errors.ErrFailedFindOrderById
	}

	return s.mapping.ToOrderResponse(order), nil
}

func (s *orderQueryService) FindByActive(req *requests.FindAllOrders) ([]*response.OrderResponseDeleteAt, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching all order active",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	orders, totalRecords, err := s.orderQueryRepository.FindByActive(req)

	if err != nil {
		s.logger.Error("Failed to retrieve active orders",
			zap.Error(err),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))

		return nil, nil, order_errors.ErrFailedFindOrdersByActive
	}

	orderResponse := s.mapping.ToOrdersResponseDeleteAt(orders)

	s.logger.Debug("Successfully fetched order",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return orderResponse, totalRecords, nil
}

func (s *orderQueryService) FindByTrashed(req *requests.FindAllOrders) ([]*response.OrderResponseDeleteAt, *int, *response.ErrorResponse) {
	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	s.logger.Debug("Fetching all order trashed",
		zap.Int("page", page),
		zap.Int("pageSize", pageSize),
		zap.String("search", search))

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	orders, totalRecords, err := s.orderQueryRepository.FindByTrashed(req)

	if err != nil {
		s.logger.Error("Failed to retrieve trashed orders from database",
			zap.Error(err),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return nil, nil, order_errors.ErrFailedFindOrdersByTrashed
	}

	orderResponse := s.mapping.ToOrdersResponseDeleteAt(orders)

	s.logger.Debug("Successfully fetched order",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return orderResponse, totalRecords, nil
}

func (s *orderQueryService) FindByMerchant(req *requests.FindAllOrderMerchant) ([]*response.OrderResponse, *int, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("FindByMerchant", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "FindByMerchant")
	defer span.End()

	span.SetAttributes(
		attribute.Int("page.number", req.Page),
		attribute.Int("page.size", req.PageSize),
		attribute.String("search.term", req.Search),
	)

	s.logger.Debug("Fetching orders",
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.PageSize),
		zap.String("search", req.Search))

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	orders, totalRecords, err := s.orderQueryRepository.FindByMerchant(req)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ORDERS_BY_MERCHANT")
		status = "failed_find_orders"

		s.logger.Error("Failed to retrieve active orders",
			zap.Error(err),
			zap.String("search", req.Search),
			zap.Int("page", req.Page),
			zap.Int("pageSize", req.PageSize),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "repository_error"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find orders by merchant")

		return nil, nil, order_errors.ErrFailedFindOrdersByMerchant
	}

	orderResponse := s.mapping.ToOrdersResponse(orders)

	span.SetAttributes(
		attribute.Int("result.count", len(orderResponse)),
		attribute.Int("total.records", *totalRecords),
	)

	s.logger.Debug("Successfully fetched order",
		zap.Int("totalRecords", *totalRecords),
		zap.Int("page", req.Page),
		zap.Int("pageSize", req.PageSize))

	return orderResponse, totalRecords, nil
}
func (s *orderQueryService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
