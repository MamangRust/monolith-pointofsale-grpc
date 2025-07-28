package service

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-order/internal/errorhandler"
	mencache "github.com/MamangRust/monolith-point-of-sale-order/internal/redis"
	"github.com/MamangRust/monolith-point-of-sale-order/internal/repository"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/cashier_errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/merchant_errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/order_errors"
	orderitem_errors "github.com/MamangRust/monolith-point-of-sale-shared/errors/order_item_errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/product_errors"
	response_service "github.com/MamangRust/monolith-point-of-sale-shared/mapper/response/service"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type orderCommandService struct {
	errorhandler               errorhandler.OrderCommandError
	mencache                   mencache.OrderCommandCache
	trace                      trace.Tracer
	cashierQueryRepository     repository.CashierQueryRepository
	orderQueryRepository       repository.OrderQueryRepository
	orderCommandRepository     repository.OrderCommandRepository
	orderItemQueryRepository   repository.OrderItemQueryRepository
	orderItemCommandRepository repository.OrderItemCommandRepository
	merchantQueryRepository    repository.MerchantQueryRepository
	productQueryRepository     repository.ProductQueryRepository
	productCommandRepository   repository.ProductCommandRepository
	logger                     logger.LoggerInterface
	mapping                    response_service.OrderResponseMapper
	requestCounter             *prometheus.CounterVec
	requestDuration            *prometheus.HistogramVec
}

func NewOrderCommandService(
	errorhandler errorhandler.OrderCommandError,
	mencache mencache.OrderCommandCache,
	cashierQueryRepository repository.CashierQueryRepository,
	orderItemQueryRepository repository.OrderItemQueryRepository,
	orderItemCommandRepository repository.OrderItemCommandRepository,
	orderQueryRepository repository.OrderQueryRepository,
	orderCommandRepository repository.OrderCommandRepository,
	productQueryRepository repository.ProductQueryRepository,
	productCommandRepository repository.ProductCommandRepository,
	merchantQueryRepository repository.MerchantQueryRepository,
	logger logger.LoggerInterface,
	mapping response_service.OrderResponseMapper,

) *orderCommandService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "order_command_service_request_count",
			Help: "Total number of requests to the OrderCommandService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "order_command_service_request_duration",
			Help:    "Histogram of request durations for the OrderCommandService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &orderCommandService{
		errorhandler:               errorhandler,
		mencache:                   mencache,
		trace:                      otel.Tracer("order-command-service"),
		cashierQueryRepository:     cashierQueryRepository,
		orderQueryRepository:       orderQueryRepository,
		orderCommandRepository:     orderCommandRepository,
		orderItemQueryRepository:   orderItemQueryRepository,
		orderItemCommandRepository: orderItemCommandRepository,
		merchantQueryRepository:    merchantQueryRepository,
		productQueryRepository:     productQueryRepository,
		productCommandRepository:   productCommandRepository,
		logger:                     logger,
		mapping:                    mapping,
		requestCounter:             requestCounter,
		requestDuration:            requestDuration,
	}
}

func (s *orderCommandService) CreateOrder(ctx context.Context, req *requests.CreateOrderRequest) (*response.OrderResponse, *response.ErrorResponse) {
	const method = "CreateOrder"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("merchant.id", req.MerchantID), attribute.Int("cashier.id", req.CashierID))

	defer func() {
		end(status)
	}()

	_, err := s.merchantQueryRepository.FindById(ctx, req.MerchantID)
	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.OrderResponse](s.logger, err, method, "FAILED_FIND_MERCHANT_BY_ID", span, &status, merchant_errors.ErrFailedFindMerchantById, zap.Error(err))
	}

	_, err = s.cashierQueryRepository.FindById(ctx, req.CashierID)
	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.OrderResponse](s.logger, err, method, "FAILED_FIND_CASHIER_BY_ID", span, &status, cashier_errors.ErrFailedFindCashierById, zap.Error(err))
	}

	order, err := s.orderCommandRepository.CreateOrder(ctx, &requests.CreateOrderRecordRequest{
		MerchantID: req.MerchantID,
		CashierID:  req.CashierID,
	})
	if err != nil {
		return s.errorhandler.HandleCreateOrderError(err, method, "FAILED_CREATE_ORDER", span, &status, zap.Error(err))
	}

	span.SetAttributes(attribute.Int("order.id", order.ID))

	for _, item := range req.Items {

		product, err := s.productQueryRepository.FindById(ctx, item.ProductID)
		if err != nil {
			return errorhandler.HandleRepositorySingleError[*response.OrderResponse](s.logger, err, method, "FAILED_FIND_PRODUCT_BY_ID", span, &status, product_errors.ErrFailedFindProductById, zap.Error(err))
		}

		if product.CountInStock < item.Quantity {
			return errorhandler.HandleRepositorySingleError[*response.OrderResponse](s.logger, err, method, "FAILED_FIND_PRODUCT_BY_ID", span, &status, product_errors.ErrFailedFindProductById, zap.Error(err))
		}

		_, err = s.orderItemCommandRepository.CreateOrderItem(ctx, &requests.CreateOrderItemRecordRequest{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		})
		if err != nil {
			return errorhandler.HandleRepositorySingleError[*response.OrderResponse](s.logger, err, method, "FAILED_CREATE_ORDER_ITEM", span, &status, orderitem_errors.ErrFailedCreateOrderItem, zap.Error(err))
		}

		product.CountInStock -= item.Quantity
		_, err = s.productCommandRepository.UpdateProductCountStock(ctx, product.ID, product.CountInStock)
		if err != nil {
			return s.errorhandler.HandleErrorInvalidCountStockTemplate(err, method, "FAILED_UPDATE_PRODUCT_COUNT_STOCK", span, &status, product_errors.ErrFailedUpdateProduct, zap.Error(err))
		}

	}

	totalPrice, err := s.orderItemQueryRepository.CalculateTotalPrice(ctx, order.ID)
	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.OrderResponse](s.logger, err, method, "FAILED_CALCULATE_TOTAL_PRICE", span, &status, orderitem_errors.ErrFailedCalculateTotal, zap.Error(err))
	}

	_, err = s.orderCommandRepository.UpdateOrder(ctx, &requests.UpdateOrderRecordRequest{
		OrderID:    order.ID,
		TotalPrice: int(*totalPrice),
	})
	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.OrderResponse](s.logger, err, method, "FAILED_UPDATE_ORDER", span, &status, order_errors.ErrFailedUpdateOrder, zap.Error(err))
	}

	so := s.mapping.ToOrderResponse(order)

	logSuccess("Successfully create order", zap.Int("order.id", order.ID))

	return so, nil
}

func (s *orderCommandService) UpdateOrder(ctx context.Context, req *requests.UpdateOrderRequest) (*response.OrderResponse, *response.ErrorResponse) {
	const method = "UpdateOrder"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("order.id", *req.OrderID))

	defer func() {
		end(status)
	}()

	_, err := s.orderQueryRepository.FindById(ctx, *req.OrderID)
	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.OrderResponse](s.logger, err, method, "FAILED_FIND_ORDER_BY_ID", span, &status, order_errors.ErrFailedFindOrderById, zap.Error(err))
	}

	for i, item := range req.Items {
		_, itemSpan := s.trace.Start(ctx, fmt.Sprintf("ProcessItem-%d", i))
		itemSpan.SetAttributes(
			attribute.Int("item.product_id", item.ProductID),
			attribute.Int("item.quantity", item.Quantity),
		)

		product, err := s.productQueryRepository.FindById(ctx, item.ProductID)
		if err != nil {
			return errorhandler.HandleRepositorySingleError[*response.OrderResponse](s.logger, err, method, "FAILED_FIND_PRODUCT_BY_ID", span, &status, product_errors.ErrFailedFindProductById, zap.Error(err))
		}

		if item.OrderItemID > 0 {
			_, err := s.orderItemCommandRepository.UpdateOrderItem(ctx, &requests.UpdateOrderItemRecordRequest{
				OrderItemID: item.OrderItemID,
				ProductID:   item.ProductID,
				Quantity:    item.Quantity,
				Price:       product.Price,
			})
			if err != nil {
				return errorhandler.HandleRepositorySingleError[*response.OrderResponse](s.logger, err, method, "FAILED_UPDATE_ORDER_ITEM", span, &status, orderitem_errors.ErrFailedUpdateOrderItem, zap.Error(err))
			}
		} else {
			if product.CountInStock < item.Quantity {
				return s.errorhandler.HandleErrorInsufficientStockTemplate(err, method, "FAILED_INSUFFICIENT_STOCK", span, &status, order_errors.ErrFailedInvalidCountInStock, zap.Error(err))
			}

			_, err := s.orderItemCommandRepository.CreateOrderItem(ctx, &requests.CreateOrderItemRecordRequest{
				OrderID:   *req.OrderID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     product.Price,
			})
			if err != nil {
				return errorhandler.HandleRepositorySingleError[*response.OrderResponse](s.logger, err, method, "FAILED_CREATE_ORDER_ITEM", span, &status, orderitem_errors.ErrFailedCreateOrderItem, zap.Error(err))
			}

			product.CountInStock -= item.Quantity
			_, err = s.productCommandRepository.UpdateProductCountStock(ctx, product.ID, product.CountInStock)
			if err != nil {
				return errorhandler.HandleRepositorySingleError[*response.OrderResponse](s.logger, err, method, "FAILED_UPDATE_PRODUCT_COUNT_STOCK", span, &status, product_errors.ErrFailedCountStock, zap.Error(err))
			}
		}
		itemSpan.End()
	}

	totalPrice, err := s.orderItemQueryRepository.CalculateTotalPrice(ctx, *req.OrderID)
	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.OrderResponse](s.logger, err, method, "FAILED_CALCULATE_TOTAL_PRICE", span, &status, orderitem_errors.ErrFailedCalculateTotal, zap.Error(err))
	}

	res, err := s.orderCommandRepository.UpdateOrder(ctx, &requests.UpdateOrderRecordRequest{
		OrderID:    *req.OrderID,
		TotalPrice: int(*totalPrice),
	})
	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.OrderResponse](s.logger, err, method, "FAILED_UPDATE_ORDER", span, &status, order_errors.ErrFailedUpdateOrder, zap.Error(err))
	}

	so := s.mapping.ToOrderResponse(res)

	s.mencache.DeleteOrderCache(ctx, *req.OrderID)

	logSuccess("Successfully updated order", zap.Int("order.id", *req.OrderID))

	return so, nil
}

func (s *orderCommandService) TrashedOrder(ctx context.Context, orderID int) (*response.OrderResponseDeleteAt, *response.ErrorResponse) {
	const method = "TrashedOrder"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("order.id", orderID))

	defer func() {
		end(status)
	}()

	order, err := s.orderQueryRepository.FindById(ctx, orderID)
	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.OrderResponseDeleteAt](s.logger, err, method, "FAILED_FIND_ORDER_BY_ID", span, &status, order_errors.ErrFailedFindOrderById, zap.Error(err))
	}

	if order.DeletedAt != nil {
		return errorhandler.HandleRepositorySingleError[*response.OrderResponseDeleteAt](s.logger, err, method, "FAILED_NOT_DELETE_AT_ORDER", span, &status, order_errors.ErrFailedNotDeleteAtOrder, zap.Error(err))
	}

	orderItems, err := s.orderItemQueryRepository.FindOrderItemByOrder(ctx, orderID)
	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.OrderResponseDeleteAt](s.logger, err, method, "FAILED_FIND_ORDER_ITEM_BY_ORDER", span, &status, orderitem_errors.ErrFailedFindOrderItemByOrder, zap.Error(err))
	}

	for _, item := range orderItems {
		if item.DeletedAt != nil {
			return errorhandler.HandleRepositorySingleError[*response.OrderResponseDeleteAt](s.logger, err, method, "FAILED_NOT_DELETE_AT_ORDER_ITEM", span, &status, orderitem_errors.ErrFailedNotDeleteAtOrderItem, zap.Error(err))
		}

		trashedItem, err := s.orderItemCommandRepository.TrashedOrderItem(ctx, item.ID)

		if err != nil {
			return errorhandler.HandleRepositorySingleError[*response.OrderResponseDeleteAt](s.logger, err, method, "FAILED_TRASH_ORDER_ITEM", span, &status, orderitem_errors.ErrFailedTrashedOrderItem, zap.Error(err))
		}

		s.logger.Debug("Order item trashed successfully",
			zap.Int("order_item_id", trashedItem.ID),
			zap.String("deleted_at", *trashedItem.DeletedAt))
	}

	trashedOrder, err := s.orderCommandRepository.TrashedOrder(ctx, orderID)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.OrderResponseDeleteAt](s.logger, err, method, "FAILED_TRASH_ORDER", span, &status, order_errors.ErrFailedTrashOrder, zap.Error(err))
	}

	so := s.mapping.ToOrderResponseDeleteAt(trashedOrder)

	s.mencache.DeleteOrderCache(ctx, orderID)

	logSuccess("Successfully trashed order", zap.Int("order.id", orderID))

	return so, nil
}

func (s *orderCommandService) RestoreOrder(ctx context.Context, order_id int) (*response.OrderResponseDeleteAt, *response.ErrorResponse) {
	const method = "RestoreOrder"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("order.id", order_id))

	defer func() {
		end(status)
	}()

	orderItems, err := s.orderItemQueryRepository.FindOrderItemByOrder(ctx, order_id)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[*response.OrderResponseDeleteAt](s.logger, err, method, "FAILED_FIND_ORDER_ITEM_BY_ORDER", span, &status, orderitem_errors.ErrFailedFindOrderItemByOrder, zap.Error(err))
	}

	for _, item := range orderItems {
		_, err := s.orderItemCommandRepository.RestoreOrderItem(ctx, item.ID)

		if err != nil {
			return errorhandler.HandleRepositorySingleError[*response.OrderResponseDeleteAt](s.logger, err, method, "FAILED_RESTORE_ORDER_ITEM", span, &status, orderitem_errors.ErrFailedRestoreOrderItem, zap.Error(err))
		}
	}

	order, err := s.orderCommandRepository.RestoreOrder(ctx, order_id)

	if err != nil {
		return s.errorhandler.HandleRestoreOrderError(err, method, "FAILED_RESTORE_ORDER", span, &status, zap.Error(err))
	}

	so := s.mapping.ToOrderResponseDeleteAt(order)

	logSuccess("Successfully restored order", zap.Int("order.id", order_id))

	return so, nil
}

func (s *orderCommandService) DeleteOrderPermanent(ctx context.Context, order_id int) (bool, *response.ErrorResponse) {
	const method = "DeleteOrderPermanent"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method, attribute.Int("order.id", order_id))

	defer func() {
		end(status)
	}()

	orderItems, err := s.orderItemQueryRepository.FindOrderItemByOrder(ctx, order_id)

	if err != nil {
		return errorhandler.HandleRepositorySingleError[bool](s.logger, err, method, "FAILED_FIND_ORDER_ITEM_BY_ORDER", span, &status, orderitem_errors.ErrFailedFindOrderItemByOrder, zap.Error(err))
	}

	for _, item := range orderItems {
		_, err := s.orderItemCommandRepository.
			DeleteOrderItemPermanent(ctx, item.ID)

		if err != nil {
			return errorhandler.HandleRepositorySingleError[bool](s.logger, err, method, "FAILED_DELETE_ORDER_ITEM_PERMANENT", span, &status, orderitem_errors.ErrFailedDeleteOrderItem, zap.Error(err))
		}
	}

	success, err := s.orderCommandRepository.DeleteOrderPermanent(ctx, order_id)

	if err != nil {
		return s.errorhandler.HandleDeleteOrderError(err, method, "FAILED_DELETE_ORDER_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted order permanently", zap.Int("order.id", order_id))

	return success, nil
}

func (s *orderCommandService) RestoreAllOrder(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "RestoreAllOrder"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	successItems, err := s.orderItemCommandRepository.RestoreAllOrderItem(ctx)

	if err != nil || !successItems {
		return errorhandler.HandleRepositorySingleError[bool](s.logger, err, method, "FAILED_RESTORE_ALL_ORDER_ITEM", span, &status, orderitem_errors.ErrFailedRestoreAllOrderItem, zap.Error(err))
	}

	success, err := s.orderCommandRepository.RestoreAllOrder(ctx)
	if err != nil || !success {
		return errorhandler.HandleRepositorySingleError[bool](s.logger, err, method, "FAILED_RESTORE_ALL_ORDER", span, &status, order_errors.ErrFailedRestoreAllOrder, zap.Error(err))
	}

	logSuccess("Successfully restored all orders", zap.Bool("success", success))

	return success, nil
}

func (s *orderCommandService) DeleteAllOrderPermanent(ctx context.Context) (bool, *response.ErrorResponse) {
	const method = "DeleteAllOrderPermanent"

	ctx, span, end, status, logSuccess := s.startTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	successItems, err := s.orderCommandRepository.DeleteAllOrderPermanent(ctx)

	if err != nil || !successItems {
		return errorhandler.HandleRepositorySingleError[bool](s.logger, err, method, "FAILED_DELETE_ALL_ORDER_ITEM_PERMANENT", span, &status, orderitem_errors.ErrFailedDeleteAllOrderItem, zap.Error(err))
	}

	success, err := s.orderCommandRepository.DeleteAllOrderPermanent(ctx)

	if err != nil || !success {
		return s.errorhandler.HandleDeleteAllOrderError(err, method, "FAILED_DELETE_ALL_ORDER_PERMANENT", span, &status, zap.Error(err))
	}

	logSuccess("Successfully deleted all orders permanently", zap.Bool("success", success))

	return success, nil
}

func (s *orderCommandService) startTracingAndLogging(ctx context.Context, method string, attrs ...attribute.KeyValue) (
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

func (s *orderCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
