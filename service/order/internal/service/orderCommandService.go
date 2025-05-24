package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	traceunic "github.com/MamangRust/monolith-point-of-sale-pkg/trace_unic"

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
	ctx                        context.Context
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
	ctx context.Context,
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
		ctx:                        ctx,
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

func (s *orderCommandService) CreateOrder(req *requests.CreateOrderRequest) (*response.OrderResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("CreateOrder", status, startTime)
	}()

	ctx, span := s.trace.Start(s.ctx, "CreateOrder")
	defer span.End()

	span.SetAttributes(
		attribute.Int("merchant.id", req.MerchantID),
		attribute.Int("cashier.id", req.CashierID),
		attribute.Int("items.count", len(req.Items)),
	)

	s.logger.Debug("Creating new order",
		zap.Int("merchantID", req.MerchantID),
		zap.Int("cashierID", req.CashierID))

	_, err := s.merchantQueryRepository.FindById(req.MerchantID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_MERCHANT")
		status = "failed_find_merchant"

		s.logger.Error("Merchant not found for order creation",
			zap.Int("merchantID", req.MerchantID),
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "merchant_not_found"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Merchant not found")

		return nil, merchant_errors.ErrFailedFindMerchantById
	}

	_, err = s.cashierQueryRepository.FindById(req.CashierID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_CASHIER")
		status = "failed_find_cashier"

		s.logger.Error("Cashier not found for order creation",
			zap.Int("cashierID", req.CashierID),
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "cashier_not_found"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Cashier not found")

		return nil, cashier_errors.ErrFailedFindCashierById
	}

	order, err := s.orderCommandRepository.CreateOrder(&requests.CreateOrderRecordRequest{
		MerchantID: req.MerchantID,
		CashierID:  req.CashierID,
	})
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_CREATE_ORDER")
		status = "failed_create_order"

		s.logger.Error("Failed to create order record",
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "create_order_failed"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to create order")

		return nil, order_errors.ErrFailedCreateOrder
	}

	span.SetAttributes(attribute.Int("order.id", order.ID))

	for i, item := range req.Items {
		_, itemSpan := s.trace.Start(ctx, fmt.Sprintf("ProcessItem-%d", i))

		product, err := s.productQueryRepository.FindById(item.ProductID)
		if err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_FIND_PRODUCT")
			status = "failed_find_product"

			s.logger.Error("Product not found for order item",
				zap.Int("productID", item.ProductID),
				zap.Error(err),
				zap.String("traceID", traceID))

			itemSpan.SetAttributes(
				attribute.String("error.trace_id", traceID),
				attribute.Int("product.id", item.ProductID),
				attribute.String("error.type", "product_not_found"),
			)
			itemSpan.RecordError(err)
			itemSpan.SetStatus(codes.Error, "Product not found")
			itemSpan.End()

			return nil, product_errors.ErrFailedFindProductById
		}

		if product.CountInStock < item.Quantity {
			traceID := traceunic.GenerateTraceID("INSUFFICIENT_STOCK")
			status = "insufficient_stock"

			s.logger.Error("Insufficient product stock",
				zap.Int("productID", item.ProductID),
				zap.Int("requested", item.Quantity),
				zap.Int("available", product.CountInStock),
				zap.String("traceID", traceID))

			itemSpan.SetAttributes(
				attribute.String("error.trace_id", traceID),
				attribute.Int("product.id", item.ProductID),
				attribute.Int("requested.quantity", item.Quantity),
				attribute.Int("available.quantity", product.CountInStock),
				attribute.String("error.type", "insufficient_stock"),
			)
			itemSpan.RecordError(err)
			itemSpan.SetStatus(codes.Error, "Insufficient stock")
			itemSpan.End()

			return nil, order_errors.ErrFailedInvalidCountInStock
		}

		_, err = s.orderItemCommandRepository.CreateOrderItem(&requests.CreateOrderItemRecordRequest{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		})
		if err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_CREATE_ORDER_ITEM")
			status = "failed_create_order_item"

			s.logger.Error("Failed to create order item",
				zap.Error(err),
				zap.String("traceID", traceID))

			itemSpan.SetAttributes(
				attribute.String("error.trace_id", traceID),
				attribute.Int("product.id", item.ProductID),
				attribute.String("error.type", "create_order_item_failed"),
			)
			itemSpan.RecordError(err)
			itemSpan.SetStatus(codes.Error, "Failed to create order item")
			itemSpan.End()

			return nil, orderitem_errors.ErrFailedCreateOrderItem
		}

		product.CountInStock -= item.Quantity
		_, err = s.productCommandRepository.UpdateProductCountStock(product.ID, product.CountInStock)
		if err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_UPDATE_PRODUCT_STOCK")
			status = "failed_update_product_stock"

			s.logger.Error("Failed to update product stock",
				zap.Error(err),
				zap.String("traceID", traceID))

			itemSpan.SetAttributes(
				attribute.String("error.trace_id", traceID),
				attribute.Int("product.id", product.ID),
				attribute.Int("new.stock", product.CountInStock),
				attribute.String("error.type", "update_product_stock_failed"),
			)
			itemSpan.RecordError(err)
			itemSpan.SetStatus(codes.Error, "Failed to update product stock")
			itemSpan.End()

			return nil, product_errors.ErrFailedUpdateProduct
		}

		itemSpan.SetAttributes(
			attribute.Int("product.id", product.ID),
			attribute.Int("quantity", item.Quantity),
			attribute.Float64("price", float64(product.Price)),
		)
		itemSpan.End()
	}

	// Calculate total price
	totalPrice, err := s.orderItemQueryRepository.CalculateTotalPrice(order.ID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_CALCULATE_TOTAL")
		status = "failed_calculate_total"

		s.logger.Error("Failed to calculate order total price",
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "calculate_total_failed"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to calculate total price")

		return nil, orderitem_errors.ErrFailedCalculateTotal
	}

	res, err := s.orderCommandRepository.UpdateOrder(&requests.UpdateOrderRecordRequest{
		OrderID:    order.ID,
		TotalPrice: int(*totalPrice),
	})
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_ORDER")
		status = "failed_update_order"

		s.logger.Error("Failed to update order total price",
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "update_order_failed"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update order")

		return nil, order_errors.ErrFailedUpdateOrder
	}

	span.SetAttributes(
		attribute.Int("order.total_price", int(*totalPrice)),
		attribute.Bool("order.completed", true),
	)

	return s.mapping.ToOrderResponse(res), nil
}

func (s *orderCommandService) UpdateOrder(req *requests.UpdateOrderRequest) (*response.OrderResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("UpdateOrder", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "UpdateOrder")
	defer span.End()

	span.SetAttributes(
		attribute.Int("order.id", *req.OrderID),
		attribute.Int("items.count", len(req.Items)),
	)

	s.logger.Debug("Updating order",
		zap.Int("orderID", *req.OrderID))

	_, err := s.orderQueryRepository.FindById(*req.OrderID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ORDER")
		status = "failed_find_order"

		s.logger.Error("Order not found for update",
			zap.Int("orderID", *req.OrderID),
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "order_not_found"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Order not found")

		return nil, order_errors.ErrFailedFindOrderById
	}

	for i, item := range req.Items {
		_, itemSpan := s.trace.Start(s.ctx, fmt.Sprintf("ProcessItem-%d", i))
		itemSpan.SetAttributes(
			attribute.Int("item.product_id", item.ProductID),
			attribute.Int("item.quantity", item.Quantity),
		)

		product, err := s.productQueryRepository.FindById(item.ProductID)
		if err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_FIND_PRODUCT")
			status = "failed_find_product"

			s.logger.Error("Product not found for order update",
				zap.Int("productID", item.ProductID),
				zap.Error(err),
				zap.String("traceID", traceID))

			itemSpan.SetAttributes(
				attribute.String("error.trace_id", traceID),
				attribute.String("error.type", "product_not_found"),
			)
			itemSpan.RecordError(err)
			itemSpan.SetStatus(codes.Error, "Product not found")
			itemSpan.End()

			return nil, product_errors.ErrFailedFindProductById
		}

		if item.OrderItemID > 0 {
			// Update existing item
			_, err := s.orderItemCommandRepository.UpdateOrderItem(&requests.UpdateOrderItemRecordRequest{
				OrderItemID: item.OrderItemID,
				ProductID:   item.ProductID,
				Quantity:    item.Quantity,
				Price:       product.Price,
			})
			if err != nil {
				traceID := traceunic.GenerateTraceID("FAILED_UPDATE_ORDER_ITEM")
				status = "failed_update_order_item"

				s.logger.Error("Failed to update order item",
					zap.Error(err),
					zap.String("traceID", traceID))

				itemSpan.SetAttributes(
					attribute.String("error.trace_id", traceID),
					attribute.String("error.type", "update_order_item_failed"),
				)
				itemSpan.RecordError(err)
				itemSpan.SetStatus(codes.Error, "Failed to update order item")
				itemSpan.End()

				return nil, orderitem_errors.ErrFailedUpdateOrderItem
			}
		} else {
			if product.CountInStock < item.Quantity {
				traceID := traceunic.GenerateTraceID("INSUFFICIENT_STOCK")
				status = "insufficient_stock"

				s.logger.Error("Insufficient product stock for new order item",
					zap.Int("productID", item.ProductID),
					zap.Int("requested", item.Quantity),
					zap.Int("available", product.CountInStock),
					zap.String("traceID", traceID))

				itemSpan.SetAttributes(
					attribute.String("error.trace_id", traceID),
					attribute.Int("available.stock", product.CountInStock),
					attribute.String("error.type", "insufficient_stock"),
				)
				itemSpan.RecordError(err)
				itemSpan.SetStatus(codes.Error, "Insufficient stock")
				itemSpan.End()

				return nil, order_errors.ErrFailedInvalidCountInStock
			}

			_, err := s.orderItemCommandRepository.CreateOrderItem(&requests.CreateOrderItemRecordRequest{
				OrderID:   *req.OrderID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     product.Price,
			})
			if err != nil {
				traceID := traceunic.GenerateTraceID("FAILED_CREATE_ORDER_ITEM")
				status = "failed_create_order_item"

				s.logger.Error("Failed to add new order item",
					zap.Error(err),
					zap.String("traceID", traceID))

				itemSpan.SetAttributes(
					attribute.String("error.trace_id", traceID),
					attribute.String("error.type", "create_order_item_failed"),
				)
				itemSpan.RecordError(err)
				itemSpan.SetStatus(codes.Error, "Failed to create order item")
				itemSpan.End()

				return nil, orderitem_errors.ErrFailedCreateOrderItem
			}

			product.CountInStock -= item.Quantity
			_, err = s.productCommandRepository.UpdateProductCountStock(product.ID, product.CountInStock)
			if err != nil {
				traceID := traceunic.GenerateTraceID("FAILED_UPDATE_PRODUCT_STOCK")
				status = "failed_update_product_stock"

				s.logger.Error("Failed to update product stock",
					zap.Error(err),
					zap.String("traceID", traceID))

				itemSpan.SetAttributes(
					attribute.String("error.trace_id", traceID),
					attribute.Int("new.stock", product.CountInStock),
					attribute.String("error.type", "update_product_stock_failed"),
				)
				itemSpan.RecordError(err)
				itemSpan.SetStatus(codes.Error, "Failed to update product stock")
				itemSpan.End()

				return nil, product_errors.ErrFailedUpdateProduct
			}
		}
		itemSpan.End()
	}

	totalPrice, err := s.orderItemQueryRepository.CalculateTotalPrice(*req.OrderID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_CALCULATE_TOTAL")
		status = "failed_calculate_total"

		s.logger.Error("Failed to calculate updated order total",
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "calculate_total_failed"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to calculate total price")

		return nil, orderitem_errors.ErrFailedCalculateTotal
	}

	res, err := s.orderCommandRepository.UpdateOrder(&requests.UpdateOrderRecordRequest{
		OrderID:    *req.OrderID,
		TotalPrice: int(*totalPrice),
	})
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_ORDER_TOTAL")
		status = "failed_update_order_total"

		s.logger.Error("Failed to update order total price",
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "update_order_total_failed"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update order total")

		return nil, order_errors.ErrFailedUpdateOrder
	}

	span.SetAttributes(
		attribute.Int("order.total_price", int(*totalPrice)),
		attribute.Bool("order.updated", true),
	)

	return s.mapping.ToOrderResponse(res), nil
}

func (s *orderCommandService) TrashedOrder(orderID int) (*response.OrderResponseDeleteAt, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("TrashedOrder", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "TrashedOrder")
	defer span.End()

	span.SetAttributes(
		attribute.Int("order.id", orderID),
	)

	s.logger.Debug("Moving order to trash",
		zap.Int("order_id", orderID))

	order, err := s.orderQueryRepository.FindById(orderID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ORDER")
		status = "failed_find_order"

		s.logger.Error("Failed to fetch order",
			zap.Int("order_id", orderID),
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "order_not_found"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Order not found")

		return nil, order_errors.ErrFailedFindOrderById
	}

	if order.DeletedAt != nil {
		status = "already_trashed"
		span.SetAttributes(
			attribute.String("order.status", "already_trashed"),
		)
		return nil, &response.ErrorResponse{
			Status:  "already_trashed",
			Message: fmt.Sprintf("Order with ID %d is already trashed", orderID),
			Code:    http.StatusBadRequest,
		}
	}

	orderItems, err := s.orderItemQueryRepository.FindOrderItemByOrder(orderID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ORDER_ITEMS")
		status = "failed_find_order_items"

		s.logger.Error("Failed to retrieve order items for trashing",
			zap.Int("order_id", orderID),
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "order_items_not_found"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find order items")

		return nil, orderitem_errors.ErrFailedOrderItemNotFound
	}

	span.SetAttributes(
		attribute.Int("order.items.count", len(orderItems)),
	)

	for i, item := range orderItems {
		_, itemSpan := s.trace.Start(s.ctx, fmt.Sprintf("TrashItem-%d", i))
		itemSpan.SetAttributes(
			attribute.Int("order_item.id", item.ID),
		)

		if item.DeletedAt != nil {
			itemSpan.SetAttributes(
				attribute.String("order_item.status", "already_trashed"),
			)
			s.logger.Debug("Order item already trashed, skipping",
				zap.Int("order_item_id", item.ID))
			itemSpan.End()
			continue
		}

		trashedItem, err := s.orderItemCommandRepository.TrashedOrderItem(item.ID)
		if err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_TRASH_ORDER_ITEM")
			status = "failed_trash_order_item"

			s.logger.Error("Failed to move order item to trash",
				zap.Int("order_item_id", item.ID),
				zap.Error(err),
				zap.String("traceID", traceID))

			itemSpan.SetAttributes(
				attribute.String("error.trace_id", traceID),
				attribute.String("error.type", "trash_order_item_failed"),
			)
			itemSpan.RecordError(err)
			itemSpan.SetStatus(codes.Error, "Failed to trash order item")
			itemSpan.End()

			return nil, orderitem_errors.ErrFailedTrashedOrderItem
		}

		itemSpan.SetAttributes(
			attribute.String("order_item.deleted_at", *trashedItem.DeletedAt),
		)
		s.logger.Debug("Order item trashed successfully",
			zap.Int("order_item_id", trashedItem.ID),
			zap.String("deleted_at", *trashedItem.DeletedAt))
		itemSpan.End()
	}

	trashedOrder, err := s.orderCommandRepository.TrashedOrder(orderID)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_TRASH_ORDER")
		status = "failed_trash_order"

		s.logger.Error("Failed to move order to trash",
			zap.Int("order_id", orderID),
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "trash_order_failed"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to trash order")

		return nil, order_errors.ErrFailedTrashOrder
	}

	span.SetAttributes(
		attribute.String("order.deleted_at", *trashedOrder.DeletedAt),
	)
	s.logger.Debug("Order moved to trash successfully",
		zap.Int("order_id", orderID),
		zap.String("deleted_at", *trashedOrder.DeletedAt))

	return s.mapping.ToOrderResponseDeleteAt(trashedOrder), nil
}

func (s *orderCommandService) RestoreOrder(order_id int) (*response.OrderResponseDeleteAt, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("RestoreOrder", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreOrder")
	defer span.End()

	span.SetAttributes(
		attribute.Int("order.id", order_id),
	)

	s.logger.Debug("Restoring order from trash",
		zap.Int("order_id", order_id))

	orderItems, err := s.orderItemQueryRepository.FindOrderItemByOrder(order_id)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ORDER_ITEMS")
		status = "failed_find_order_items"

		s.logger.Error("Failed to retrieve order items for restoration",
			zap.Int("order_id", order_id),
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "order_items_not_found"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find order items")

		return nil, orderitem_errors.ErrFailedFindOrderItemByOrder
	}

	span.SetAttributes(
		attribute.Int("order.items.count", len(orderItems)),
	)

	for i, item := range orderItems {
		_, itemSpan := s.trace.Start(s.ctx, fmt.Sprintf("RestoreItem-%d", i))
		itemSpan.SetAttributes(
			attribute.Int("order_item.id", item.ID),
		)

		_, err := s.orderItemCommandRepository.RestoreOrderItem(item.ID)
		if err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_RESTORE_ORDER_ITEM")
			status = "failed_restore_order_item"

			s.logger.Error("Failed to restore order item from trash",
				zap.Int("order_item_id", item.ID),
				zap.Error(err),
				zap.String("traceID", traceID))

			itemSpan.SetAttributes(
				attribute.String("error.trace_id", traceID),
				attribute.String("error.type", "restore_order_item_failed"),
			)
			itemSpan.RecordError(err)
			itemSpan.SetStatus(codes.Error, "Failed to restore order item")
			itemSpan.End()

			return nil, orderitem_errors.ErrFailedRestoreOrderItem
		}
		itemSpan.End()
	}

	order, err := s.orderCommandRepository.RestoreOrder(order_id)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_ORDER")
		status = "failed_restore_order"

		s.logger.Error("Failed to restore order from trash",
			zap.Int("order_id", order_id),
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "restore_order_failed"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore order")

		return nil, order_errors.ErrFailedRestoreOrder
	}

	span.SetAttributes(
		attribute.Bool("order.restored", true),
	)

	return s.mapping.ToOrderResponseDeleteAt(order), nil
}

func (s *orderCommandService) DeleteOrderPermanent(order_id int) (bool, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("DeleteOrderPermanent", status, startTime)
	}()

	ctx, span := s.trace.Start(s.ctx, "DeleteOrderPermanent")
	defer span.End()

	span.SetAttributes(
		attribute.Int("order.id", order_id),
	)

	s.logger.Debug("Permanently deleting order",
		zap.Int("order_id", order_id))

	orderItems, err := s.orderItemQueryRepository.FindOrderItemByOrder(order_id)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_ORDER_ITEMS")
		status = "failed_find_order_items"

		s.logger.Error("Failed to retrieve order items for permanent deletion",
			zap.Int("order_id", order_id),
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "order_items_not_found"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find order items")

		return false, orderitem_errors.ErrFailedFindOrderItemByOrder
	}

	span.SetAttributes(
		attribute.Int("order.items.count", len(orderItems)),
	)

	for i, item := range orderItems {
		_, itemSpan := s.trace.Start(ctx, fmt.Sprintf("DeleteItem-%d", i))
		itemSpan.SetAttributes(
			attribute.Int("order_item.id", item.ID),
		)

		_, err := s.orderItemCommandRepository.DeleteOrderItemPermanent(item.ID)
		if err != nil {
			traceID := traceunic.GenerateTraceID("FAILED_DELETE_ORDER_ITEM")
			status = "failed_delete_order_item"

			s.logger.Error("Failed to permanently delete order item",
				zap.Int("order_item_id", item.ID),
				zap.Error(err),
				zap.String("traceID", traceID))

			itemSpan.SetAttributes(
				attribute.String("error.trace_id", traceID),
				attribute.String("error.type", "delete_order_item_failed"),
			)
			itemSpan.RecordError(err)
			itemSpan.SetStatus(codes.Error, "Failed to delete order item")
			itemSpan.End()

			return false, orderitem_errors.ErrFailedDeleteOrderItem
		}
		itemSpan.End()
	}

	success, err := s.orderCommandRepository.DeleteOrderPermanent(order_id)
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_ORDER")
		status = "failed_delete_order"

		s.logger.Error("Failed to permanently delete order",
			zap.Int("order_id", order_id),
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "delete_order_failed"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to delete order")

		return false, order_errors.ErrFailedDeleteOrderPermanent
	}

	span.SetAttributes(
		attribute.Bool("operation.success", success),
	)
	return success, nil
}

func (s *orderCommandService) RestoreAllOrder() (bool, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("RestoreAllOrder", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreAllOrder")
	defer span.End()

	s.logger.Debug("Restoring all trashed orders")

	successItems, err := s.orderItemCommandRepository.RestoreAllOrderItem()
	if err != nil || !successItems {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_ALL_ITEMS")
		status = "failed_restore_all_items"

		s.logger.Error("Failed to restore all order items",
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "restore_all_items_failed"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore all items")

		return false, orderitem_errors.ErrFailedRestoreAllOrderItem
	}

	success, err := s.orderCommandRepository.RestoreAllOrder()
	if err != nil || !success {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_ALL_ORDERS")
		status = "failed_restore_all_orders"

		s.logger.Error("Failed to restore all orders",
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "restore_all_orders_failed"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore all orders")

		return false, order_errors.ErrFailedRestoreAllOrder
	}

	span.SetAttributes(
		attribute.Bool("operation.success", success),
	)
	return success, nil
}

func (s *orderCommandService) DeleteAllOrderPermanent() (bool, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"
	defer func() {
		s.recordMetrics("DeleteAllOrderPermanent", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteAllOrderPermanent")
	defer span.End()

	s.logger.Debug("Permanently deleting all trashed orders")

	successItems, err := s.orderItemCommandRepository.DeleteAllOrderPermanent()
	if err != nil || !successItems {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_ALL_ITEMS")
		status = "failed_delete_all_items"

		s.logger.Error("Failed to permanently delete all order items",
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "delete_all_items_failed"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to delete all items")

		return false, orderitem_errors.ErrFailedDeleteAllOrderItem
	}

	success, err := s.orderCommandRepository.DeleteAllOrderPermanent()
	if err != nil || !success {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_ALL_ORDERS")
		status = "failed_delete_all_orders"

		s.logger.Error("Failed to permanently delete all orders",
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.String("error.type", "delete_all_orders_failed"),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to delete all orders")

		return false, order_errors.ErrFailedDeleteAllOrderPermanent
	}

	span.SetAttributes(
		attribute.Bool("operation.success", success),
	)
	return success, nil
}

func (s *orderCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
