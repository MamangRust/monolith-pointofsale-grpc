package service

import (
	"context"

	mencache "github.com/MamangRust/monolith-point-of-sale-order/internal/redis"
	"github.com/MamangRust/monolith-point-of-sale-order/internal/repository"
	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	sharederrorhandler "github.com/MamangRust/monolith-point-of-sale-shared/errorhandler"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/cashier_errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/merchant_errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/order_errors"
	orderitem_errors "github.com/MamangRust/monolith-point-of-sale-shared/errors/order_item_errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/product_errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type orderCommandDeps struct {
	Cache                      mencache.OrderCommandCache
	CashierQueryRepository     repository.CashierQueryRepository
	OrderQueryRepository       repository.OrderQueryRepository
	OrderCommandRepository     repository.OrderCommandRepository
	OrderItemQueryRepository   repository.OrderItemQueryRepository
	OrderItemCommandRepository repository.OrderItemCommandRepository
	MerchantQueryRepository    repository.MerchantQueryRepository
	ProductQueryRepository     repository.ProductQueryRepository
	ProductCommandRepository   repository.ProductCommandRepository
	Logger                     logger.LoggerInterface
	Observability              observability.TraceLoggerObservability
}

type orderCommandService struct {
	mencache                   mencache.OrderCommandCache
	cashierQueryRepository     repository.CashierQueryRepository
	orderQueryRepository       repository.OrderQueryRepository
	orderCommandRepository     repository.OrderCommandRepository
	orderItemQueryRepository   repository.OrderItemQueryRepository
	orderItemCommandRepository repository.OrderItemCommandRepository
	merchantQueryRepository    repository.MerchantQueryRepository
	productQueryRepository     repository.ProductQueryRepository
	productCommandRepository   repository.ProductCommandRepository
	logger                     logger.LoggerInterface
	observability              observability.TraceLoggerObservability
}

func NewOrderCommandService(params *orderCommandDeps) OrderCommandService {
	return &orderCommandService{
		mencache:                   params.Cache,
		cashierQueryRepository:     params.CashierQueryRepository,
		orderQueryRepository:       params.OrderQueryRepository,
		orderCommandRepository:     params.OrderCommandRepository,
		orderItemQueryRepository:   params.OrderItemQueryRepository,
		orderItemCommandRepository: params.OrderItemCommandRepository,
		merchantQueryRepository:    params.MerchantQueryRepository,
		productQueryRepository:     params.ProductQueryRepository,
		productCommandRepository:   params.ProductCommandRepository,
		logger:                     params.Logger,
		observability:              params.Observability,
	}
}

func (s *orderCommandService) CreateOrder(ctx context.Context, req *requests.CreateOrderRequest) (*db.Order, error) {
	const method = "CreateOrder"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant.id", req.MerchantID),
		attribute.Int("cashier.id", req.CashierID),
	)
	defer func() {
		end(status)
	}()

	_, err := s.merchantQueryRepository.FindById(ctx, req.MerchantID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Order](
			s.logger,
			merchant_errors.ErrFailedFindMerchantById.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	_, err = s.cashierQueryRepository.FindById(ctx, req.CashierID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Order](
			s.logger,
			cashier_errors.ErrFailedFindCashierById.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	order, err := s.orderCommandRepository.CreateOrder(ctx, &requests.CreateOrderRecordRequest{
		MerchantID: req.MerchantID,
		CashierID:  req.CashierID,
	})
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Order](
			s.logger,
			order_errors.ErrFailedCreateOrder.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	span.SetAttributes(attribute.Int("order.id", int(order.OrderID)))

	for _, item := range req.Items {
		product, err := s.productQueryRepository.FindById(ctx, item.ProductID)
		if err != nil {
			status = "error"
			return sharederrorhandler.HandleError[*db.Order](
				s.logger,
				product_errors.ErrFailedFindProductById.WithInternal(err),
				method,
				span,
				zap.Error(err),
			)
		}

		if product.CountInStock < int32(item.Quantity) {
			status = "error"
			return sharederrorhandler.HandleError[*db.Order](
				s.logger,
				order_errors.ErrInsufficientProductStock,
				method,
				span,
			)
		}

		_, err = s.orderItemCommandRepository.CreateOrderItem(ctx, &requests.CreateOrderItemRecordRequest{
			OrderID:   int(order.OrderID),
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     int(product.Price),
		})
		if err != nil {
			status = "error"
			return sharederrorhandler.HandleError[*db.Order](
				s.logger,
				orderitem_errors.ErrFailedCreateOrderItem.WithInternal(err),
				method,
				span,
				zap.Error(err),
			)
		}

		newStock := int(product.CountInStock) - item.Quantity
		_, err = s.productCommandRepository.UpdateProductCountStock(ctx, int(product.ProductID), newStock)
		if err != nil {
			status = "error"
			return sharederrorhandler.HandleError[*db.Order](
				s.logger,
				product_errors.ErrFailedUpdateProduct.WithInternal(err),
				method,
				span,
				zap.Error(err),
			)
		}
	}

	totalPrice, err := s.orderItemQueryRepository.CalculateTotalPrice(ctx, int(order.OrderID))
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Order](
			s.logger,
			orderitem_errors.ErrFailedCalculateTotal.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	res, err := s.orderCommandRepository.UpdateOrder(ctx, &requests.UpdateOrderRecordRequest{
		OrderID:    int(order.OrderID),
		TotalPrice: int(*totalPrice),
	})
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Order](
			s.logger,
			order_errors.ErrFailedUpdateOrder.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	logSuccess("Successfully created order", zap.Int("order.id", int(res.OrderID)))
	return res, nil
}

func (s *orderCommandService) UpdateOrder(ctx context.Context, req *requests.UpdateOrderRequest) (*db.Order, error) {
	const method = "UpdateOrder"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("order.id", *req.OrderID))
	defer func() {
		end(status)
	}()

	_, err := s.orderQueryRepository.FindById(ctx, *req.OrderID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Order](
			s.logger,
			order_errors.ErrFailedFindOrderById.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	for _, item := range req.Items {
		product, err := s.productQueryRepository.FindById(ctx, item.ProductID)
		if err != nil {
			status = "error"
			return sharederrorhandler.HandleError[*db.Order](
				s.logger,
				product_errors.ErrFailedFindProductById.WithInternal(err),
				method,
				span,
				zap.Error(err),
			)
		}

		if item.OrderItemID > 0 {
			_, err := s.orderItemCommandRepository.UpdateOrderItem(ctx, &requests.UpdateOrderItemRecordRequest{
				OrderItemID: item.OrderItemID,
				ProductID:   item.ProductID,
				Quantity:    item.Quantity,
				Price:       int(product.Price),
			})
			if err != nil {
				status = "error"
				return sharederrorhandler.HandleError[*db.Order](
					s.logger,
					orderitem_errors.ErrFailedUpdateOrderItem.WithInternal(err),
					method,
					span,
					zap.Error(err),
				)
			}
		} else {
			if product.CountInStock < int32(item.Quantity) {
				status = "error"
				return sharederrorhandler.HandleError[*db.Order](
					s.logger,
					order_errors.ErrInsufficientProductStock,
					method,
					span,
				)
			}

			_, err := s.orderItemCommandRepository.CreateOrderItem(ctx, &requests.CreateOrderItemRecordRequest{
				OrderID:   *req.OrderID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     int(product.Price),
			})
			if err != nil {
				status = "error"
				return sharederrorhandler.HandleError[*db.Order](
					s.logger,
					orderitem_errors.ErrFailedCreateOrderItem.WithInternal(err),
					method,
					span,
					zap.Error(err),
				)
			}

			newStock := int(product.CountInStock) - item.Quantity
			_, err = s.productCommandRepository.UpdateProductCountStock(ctx, int(product.ProductID), newStock)
			if err != nil {
				status = "error"
				return sharederrorhandler.HandleError[*db.Order](
					s.logger,
					product_errors.ErrFailedUpdateProduct.WithInternal(err),
					method,
					span,
					zap.Error(err),
				)
			}
		}
	}

	totalPrice, err := s.orderItemQueryRepository.CalculateTotalPrice(ctx, *req.OrderID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Order](
			s.logger,
			orderitem_errors.ErrFailedCalculateTotal.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	res, err := s.orderCommandRepository.UpdateOrder(ctx, &requests.UpdateOrderRecordRequest{
		OrderID:    *req.OrderID,
		TotalPrice: int(*totalPrice),
	})
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Order](
			s.logger,
			order_errors.ErrFailedUpdateOrder.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.DeleteOrderCache(ctx, *req.OrderID)
	logSuccess("Successfully updated order", zap.Int("order.id", *req.OrderID))
	return res, nil
}

func (s *orderCommandService) TrashedOrder(ctx context.Context, orderID int) (*db.Order, error) {
	const method = "TrashedOrder"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("order.id", orderID))
	defer func() {
		end(status)
	}()

	order, err := s.orderQueryRepository.FindById(ctx, orderID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Order](
			s.logger,
			order_errors.ErrFailedFindOrderById.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	if order.DeletedAt.Valid {
		status = "error"
		return sharederrorhandler.HandleError[*db.Order](
			s.logger,
			order_errors.ErrFailedNotDeleteAtOrder,
			method,
			span,
		)
	}

	orderItems, err := s.orderItemQueryRepository.FindOrderItemByOrder(ctx, orderID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Order](
			s.logger,
			orderitem_errors.ErrFailedFindOrderItemByOrder.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	for _, item := range orderItems {
		if item.DeletedAt.Valid {
			status = "error"
			return sharederrorhandler.HandleError[*db.Order](
				s.logger,
				orderitem_errors.ErrFailedNotDeleteAtOrderItem,
				method,
				span,
			)
		}

		trashedItem, err := s.orderItemCommandRepository.TrashedOrderItem(ctx, int(item.OrderItemID))
		if err != nil {
			status = "error"
			return sharederrorhandler.HandleError[*db.Order](
				s.logger,
				orderitem_errors.ErrFailedTrashedOrderItem.WithInternal(err),
				method,
				span,
				zap.Error(err),
			)
		}

		s.logger.Debug("Order item trashed successfully",
			zap.Int("order_item_id", int(trashedItem.OrderItemID)),
			zap.Time("deleted_at", trashedItem.DeletedAt.Time))
	}

	trashedOrder, err := s.orderCommandRepository.TrashedOrder(ctx, orderID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Order](
			s.logger,
			order_errors.ErrFailedTrashOrder.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.DeleteOrderCache(ctx, orderID)
	logSuccess("Successfully trashed order", zap.Int("order.id", orderID))
	return trashedOrder, nil
}

func (s *orderCommandService) RestoreOrder(ctx context.Context, orderID int) (*db.Order, error) {
	const method = "RestoreOrder"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("order.id", orderID))
	defer func() {
		end(status)
	}()

	orderItems, err := s.orderItemQueryRepository.FindOrderItemByOrder(ctx, orderID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Order](
			s.logger,
			orderitem_errors.ErrFailedFindOrderItemByOrder.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	for _, item := range orderItems {
		_, err := s.orderItemCommandRepository.RestoreOrderItem(ctx, int(item.OrderItemID))
		if err != nil {
			status = "error"
			return sharederrorhandler.HandleError[*db.Order](
				s.logger,
				orderitem_errors.ErrFailedRestoreOrderItem.WithInternal(err),
				method,
				span,
				zap.Error(err),
			)
		}
	}

	order, err := s.orderCommandRepository.RestoreOrder(ctx, orderID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*db.Order](
			s.logger,
			order_errors.ErrFailedRestoreOrder.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	logSuccess("Successfully restored order", zap.Int("order.id", orderID))
	return order, nil
}

func (s *orderCommandService) DeleteOrderPermanent(ctx context.Context, orderID int) (bool, error) {
	const method = "DeleteOrderPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("order.id", orderID))
	defer func() {
		end(status)
	}()

	orderItems, err := s.orderItemQueryRepository.FindOrderItemByOrder(ctx, orderID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](
			s.logger,
			orderitem_errors.ErrFailedFindOrderItemByOrder.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	for _, item := range orderItems {
		_, err := s.orderItemCommandRepository.DeleteOrderItemPermanent(ctx, int(item.OrderItemID))
		if err != nil {
			status = "error"
			return sharederrorhandler.HandleError[bool](
				s.logger,
				orderitem_errors.ErrFailedDeleteOrderItem.WithInternal(err),
				method,
				span,
				zap.Error(err),
			)
		}
	}

	success, err := s.orderCommandRepository.DeleteOrderPermanent(ctx, orderID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](
			s.logger,
			order_errors.ErrFailedDeleteOrderPermanent.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	logSuccess("Successfully deleted order permanently", zap.Int("order.id", orderID))
	return success, nil
}

func (s *orderCommandService) RestoreAllOrder(ctx context.Context) (bool, error) {
	const method = "RestoreAllOrder"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() {
		end(status)
	}()

	successItems, err := s.orderItemCommandRepository.RestoreAllOrderItem(ctx)
	if err != nil || !successItems {
		status = "error"
		return sharederrorhandler.HandleError[bool](
			s.logger,
			orderitem_errors.ErrFailedRestoreAllOrderItem.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	success, err := s.orderCommandRepository.RestoreAllOrder(ctx)
	if err != nil || !success {
		status = "error"
		return sharederrorhandler.HandleError[bool](
			s.logger,
			order_errors.ErrFailedRestoreAllOrder.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	logSuccess("Successfully restored all orders", zap.Bool("success", success))
	return success, nil
}

func (s *orderCommandService) DeleteAllOrderPermanent(ctx context.Context) (bool, error) {
	const method = "DeleteAllOrderPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() {
		end(status)
	}()

	successItems, err := s.orderItemCommandRepository.DeleteAllOrderPermanent(ctx)
	if err != nil || !successItems {
		status = "error"
		return sharederrorhandler.HandleError[bool](
			s.logger,
			orderitem_errors.ErrFailedDeleteAllOrderItem.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	success, err := s.orderCommandRepository.DeleteAllOrderPermanent(ctx)
	if err != nil || !success {
		status = "error"
		return sharederrorhandler.HandleError[bool](
			s.logger,
			order_errors.ErrFailedDeleteAllOrderPermanent.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	logSuccess("Successfully deleted all orders permanently", zap.Bool("success", success))
	return success, nil
}
