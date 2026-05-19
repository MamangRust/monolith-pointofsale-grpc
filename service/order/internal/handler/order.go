package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-point-of-sale-order/internal/service"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/order_errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/pb"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type orderHandleGrpc struct {
	pb.UnimplementedOrderServiceServer
	orderQuery           service.OrderQueryService
	orderCommand         service.OrderCommandService
	orderStats           service.OrderStatsService
	orderStatsByMerchant service.OrderStatByMerchantService
	logger               logger.LoggerInterface
}

func NewOrderHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pb.OrderServiceServer {
	return &orderHandleGrpc{
		orderQuery:           service.OrderQuery,
		orderCommand:         service.OrderCommand,
		orderStats:           service.OrderStats,
		orderStatsByMerchant: service.OrderStatsByMerchant,
		logger:               logger,
	}
}

func (s *orderHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllOrderRequest) (*pb.ApiResponsePaginationOrder, error) {
	s.logger.Info("FindAll orders called", zap.Int32("page", request.GetPage()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllOrders{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	merchant, totalRecords, err := s.orderQuery.FindAll(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindAll orders failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindAll orders success")

	return &pb.ApiResponsePaginationOrder{
		Status:     "success",
		Message:    "Successfully fetched order",
		Data:       mapResponsesOrder(merchant),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *orderHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdOrderRequest) (*pb.ApiResponseOrder, error) {
	s.logger.Info("FindById order called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id == 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	merchant, err := s.orderQuery.FindById(ctx, id)
	if err != nil {
		s.logger.Error("FindById order failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindById order success")

	return &pb.ApiResponseOrder{
		Status:  "success",
		Message: "Successfully fetched order",
		Data:    mapResponseOrder(merchant),
	}, nil
}

func (s *orderHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllOrderRequest) (*pb.ApiResponsePaginationOrderDeleteAt, error) {
	s.logger.Info("FindByActive orders called", zap.Int32("page", request.GetPage()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllOrders{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	merchant, totalRecords, err := s.orderQuery.FindByActive(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByActive orders failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByActive orders success")

	return &pb.ApiResponsePaginationOrderDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active order",
		Data:       mapResponsesOrderActive(merchant),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *orderHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllOrderRequest) (*pb.ApiResponsePaginationOrderDeleteAt, error) {
	s.logger.Info("FindByTrashed orders called")

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllOrders{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	users, totalRecords, err := s.orderQuery.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByTrashed orders failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByTrashed orders success")

	return &pb.ApiResponsePaginationOrderDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed order",
		Data:       mapResponsesOrderTrashed(users),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *orderHandleGrpc) FindMonthlyTotalRevenue(ctx context.Context, req *pb.FindYearMonthTotalRevenue) (*pb.ApiResponseOrderMonthlyTotalRevenue, error) {
	s.logger.Info("FindMonthlyTotalRevenue orders called", zap.Int32("year", req.GetYear()))

	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, order_errors.ErrGrpcInvalidYear
	}
	if month <= 0 || month >= 12 {
		return nil, order_errors.ErrGrpcInvalidMonth
	}

	reqService := requests.MonthTotalRevenue{
		Year:  year,
		Month: month,
	}

	methods, err := s.orderStats.FindMonthlyTotalRevenue(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthlyTotalRevenue orders failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthlyTotalRevenue orders success")

	return &pb.ApiResponseOrderMonthlyTotalRevenue{
		Status:  "success",
		Message: "Monthly sales retrieved successfully",
		Data:    mapResponseOrderMonthlyTotalRevenues(methods),
	}, nil
}

func (s *orderHandleGrpc) FindYearlyTotalRevenue(ctx context.Context, req *pb.FindYearTotalRevenue) (*pb.ApiResponseOrderYearlyTotalRevenue, error) {
	s.logger.Info("FindYearlyTotalRevenue orders called", zap.Int32("year", req.GetYear()))

	year := int(req.GetYear())
	if year <= 0 {
		return nil, order_errors.ErrGrpcInvalidYear
	}

	methods, err := s.orderStats.FindYearlyTotalRevenue(ctx, year)
	if err != nil {
		s.logger.Error("FindYearlyTotalRevenue orders failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearlyTotalRevenue orders success")

	return &pb.ApiResponseOrderYearlyTotalRevenue{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    mapResponseOrderYearlyTotalRevenues(methods),
	}, nil
}

func (s *orderHandleGrpc) FindMonthlyTotalRevenueByMerchant(ctx context.Context, req *pb.FindYearMonthTotalRevenueByMerchant) (*pb.ApiResponseOrderMonthlyTotalRevenue, error) {
	s.logger.Info("FindMonthlyTotalRevenueByMerchant orders called", zap.Int32("merchantId", req.GetMerchantId()))

	year := int(req.GetYear())
	month := int(req.GetMonth())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, order_errors.ErrGrpcInvalidYear
	}
	if month <= 0 || month >= 12 {
		return nil, order_errors.ErrGrpcInvalidMonth
	}
	if id <= 0 {
		return nil, order_errors.ErrGrpcFailedInvalidMerchantId
	}

	reqService := requests.MonthTotalRevenueMerchant{
		Year:       year,
		Month:      month,
		MerchantID: id,
	}

	methods, err := s.orderStatsByMerchant.FindMonthlyTotalRevenueByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthlyTotalRevenueByMerchant orders failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthlyTotalRevenueByMerchant orders success")

	return &pb.ApiResponseOrderMonthlyTotalRevenue{
		Status:  "success",
		Message: "Monthly sales retrieved successfully",
		Data:    mapResponseOrderMonthlyTotalRevenuesByMerchant(methods),
	}, nil
}

func (s *orderHandleGrpc) FindYearlyTotalRevenueByMerchant(ctx context.Context, req *pb.FindYearTotalRevenueByMerchant) (*pb.ApiResponseOrderYearlyTotalRevenue, error) {
	s.logger.Info("FindYearlyTotalRevenueByMerchant orders called", zap.Int32("merchantId", req.GetMerchantId()))

	year := int(req.GetYear())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, order_errors.ErrGrpcInvalidYear
	}
	if id <= 0 {
		return nil, order_errors.ErrGrpcFailedInvalidMerchantId
	}

	reqService := requests.YearTotalRevenueMerchant{
		Year:       year,
		MerchantID: id,
	}

	methods, err := s.orderStatsByMerchant.FindYearlyTotalRevenueByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindYearlyTotalRevenueByMerchant orders failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearlyTotalRevenueByMerchant orders success")

	return &pb.ApiResponseOrderYearlyTotalRevenue{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    mapResponseOrderYearlyTotalRevenuesByMerchant(methods),
	}, nil
}

func (s *orderHandleGrpc) FindMonthlyRevenue(ctx context.Context, request *pb.FindYearOrder) (*pb.ApiResponseOrderMonthly, error) {
	s.logger.Info("FindMonthlyRevenue orders called", zap.Int32("year", request.GetYear()))

	year := int(request.GetYear())
	if year <= 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	res, err := s.orderStats.FindMonthlyOrder(ctx, year)
	if err != nil {
		s.logger.Error("FindMonthlyRevenue orders failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthlyRevenue orders success")

	return &pb.ApiResponseOrderMonthly{
		Status:  "success",
		Message: "Monthly revenue data retrieved",
		Data:    mapResponsesOrderMonthlyPrices(res),
	}, nil
}

func (s *orderHandleGrpc) FindYearlyRevenue(ctx context.Context, request *pb.FindYearOrder) (*pb.ApiResponseOrderYearly, error) {
	s.logger.Info("FindYearlyRevenue orders called", zap.Int32("year", request.GetYear()))

	year := int(request.GetYear())
	if year <= 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	res, err := s.orderStats.FindYearlyOrder(ctx, year)
	if err != nil {
		s.logger.Error("FindYearlyRevenue orders failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearlyRevenue orders success")

	return &pb.ApiResponseOrderYearly{
		Status:  "success",
		Message: "Yearly revenue data retrieved",
		Data:    mapResponsesOrderYearlyPrices(res),
	}, nil
}

func (s *orderHandleGrpc) FindMonthlyRevenueByMerchant(ctx context.Context, request *pb.FindYearOrderByMerchant) (*pb.ApiResponseOrderMonthly, error) {
	s.logger.Info("FindMonthlyRevenueByMerchant orders called", zap.Int32("merchantId", request.GetMerchantId()))

	year := int(request.GetYear())
	id := int(request.GetMerchantId())

	if year <= 0 {
		return nil, order_errors.ErrGrpcInvalidYear
	}
	if id <= 0 {
		return nil, order_errors.ErrGrpcFailedInvalidMerchantId
	}

	reqService := requests.MonthOrderMerchant{
		Year:       year,
		MerchantID: id,
	}

	res, err := s.orderStatsByMerchant.FindMonthlyOrderByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthlyRevenueByMerchant orders failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthlyRevenueByMerchant orders success")

	return &pb.ApiResponseOrderMonthly{
		Status:  "success",
		Message: "Monthly revenue by merchant data retrieved",
		Data:    mapResponsesOrderMonthlyPricesByMerchant(res),
	}, nil
}

func (s *orderHandleGrpc) FindYearlyRevenueByMerchant(ctx context.Context, request *pb.FindYearOrderByMerchant) (*pb.ApiResponseOrderYearly, error) {
	s.logger.Info("FindYearlyRevenueByMerchant orders called", zap.Int32("merchantId", request.GetMerchantId()))

	year := int(request.GetYear())
	id := int(request.GetMerchantId())

	if year <= 0 {
		return nil, order_errors.ErrGrpcInvalidYear
	}
	if id <= 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	reqService := requests.YearOrderMerchant{
		Year:       year,
		MerchantID: id,
	}

	res, err := s.orderStatsByMerchant.FindYearlyOrderByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindYearlyRevenueByMerchant orders failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearlyRevenueByMerchant orders success")

	return &pb.ApiResponseOrderYearly{
		Status:  "success",
		Message: "Yearly revenue by merchant data retrieved",
		Data:    mapResponsesOrderYearlyPricesByMerchant(res),
	}, nil
}

func (s *orderHandleGrpc) Create(ctx context.Context, request *pb.CreateOrderRequest) (*pb.ApiResponseOrder, error) {
	s.logger.Info("Create order called", zap.Int32("merchantId", request.GetMerchantId()))

	req := &requests.CreateOrderRequest{
		MerchantID: int(request.GetMerchantId()),
		CashierID:  int(request.GetCashierId()),
	}

	for _, item := range request.GetItems() {
		req.Items = append(req.Items, requests.CreateOrderItemRequest{
			ProductID: int(item.GetProductId()),
			Quantity:  int(item.GetQuantity()),
		})
	}

	if err := req.Validate(); err != nil {
		return nil, order_errors.ErrGrpcValidateCreateOrder
	}

	order, err := s.orderCommand.CreateOrder(ctx, req)
	if err != nil {
		s.logger.Error("Create order failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Create order success")

	return &pb.ApiResponseOrder{
		Status:  "success",
		Message: "Successfully created order",
		Data:    mapResponseOrder(order),
	}, nil
}

func (s *orderHandleGrpc) Update(ctx context.Context, request *pb.UpdateOrderRequest) (*pb.ApiResponseOrder, error) {
	s.logger.Info("Update order called", zap.Int32("id", request.GetOrderId()))

	id := int(request.GetOrderId())
	if id == 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	req := &requests.UpdateOrderRequest{
		OrderID: &id,
	}

	for _, item := range request.GetItems() {
		req.Items = append(req.Items, requests.UpdateOrderItemRequest{
			OrderItemID: int(item.GetOrderItemId()),
			ProductID:   int(item.GetProductId()),
			Quantity:    int(item.GetQuantity()),
		})
	}

	if err := req.Validate(); err != nil {
		return nil, order_errors.ErrGrpcValidateUpdateOrder
	}

	order, err := s.orderCommand.UpdateOrder(ctx, req)
	if err != nil {
		s.logger.Error("Update order failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Update order success")

	return &pb.ApiResponseOrder{
		Status:  "success",
		Message: "Successfully updated order",
		Data:    mapResponseOrder(order),
	}, nil
}

func (s *orderHandleGrpc) TrashedOrder(ctx context.Context, request *pb.FindByIdOrderRequest) (*pb.ApiResponseOrderDeleteAt, error) {
	s.logger.Info("TrashedOrder called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id == 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	merchant, err := s.orderCommand.TrashedOrder(ctx, id)
	if err != nil {
		s.logger.Error("TrashedOrder failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("TrashedOrder success")

	return &pb.ApiResponseOrderDeleteAt{
		Status:  "success",
		Message: "Successfully trashed order",
		Data:    mapResponseOrderDeleteAt(merchant),
	}, nil
}

func (s *orderHandleGrpc) RestoreOrder(ctx context.Context, request *pb.FindByIdOrderRequest) (*pb.ApiResponseOrderDeleteAt, error) {
	s.logger.Info("RestoreOrder called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id == 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	merchant, err := s.orderCommand.RestoreOrder(ctx, id)
	if err != nil {
		s.logger.Error("RestoreOrder failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreOrder success")

	return &pb.ApiResponseOrderDeleteAt{
		Status:  "success",
		Message: "Successfully restored order",
		Data:    mapResponseOrderDeleteAt(merchant),
	}, nil
}

func (s *orderHandleGrpc) DeleteOrderPermanent(ctx context.Context, request *pb.FindByIdOrderRequest) (*pb.ApiResponseOrderDelete, error) {
	s.logger.Info("DeleteOrderPermanent called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id == 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	_, err := s.orderCommand.DeleteOrderPermanent(ctx, id)
	if err != nil {
		s.logger.Error("DeleteOrderPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteOrderPermanent success")

	return &pb.ApiResponseOrderDelete{
		Status:  "success",
		Message: "Successfully deleted order permanently",
	}, nil
}

func (s *orderHandleGrpc) RestoreAllOrder(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseOrderAll, error) {
	s.logger.Info("RestoreAllOrder called")

	_, err := s.orderCommand.RestoreAllOrder(ctx)
	if err != nil {
		s.logger.Error("RestoreAllOrder failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreAllOrder success")

	return &pb.ApiResponseOrderAll{
		Status:  "success",
		Message: "Successfully restore all order",
	}, nil
}

func (s *orderHandleGrpc) DeleteAllOrderPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseOrderAll, error) {
	s.logger.Info("DeleteAllOrderPermanent called")

	_, err := s.orderCommand.DeleteAllOrderPermanent(ctx)
	if err != nil {
		s.logger.Error("DeleteAllOrderPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteAllOrderPermanent success")

	return &pb.ApiResponseOrderAll{
		Status:  "success",
		Message: "Successfully delete order permanen",
	}, nil
}

// Map helpers
func mapPaginationMeta(meta *pb.PaginationMeta) *pb.PaginationMeta {
	if meta == nil {
		return nil
	}
	return &pb.PaginationMeta{
		CurrentPage:  meta.CurrentPage,
		PageSize:     meta.PageSize,
		TotalPages:   meta.TotalPages,
		TotalRecords: meta.TotalRecords,
	}
}

func mapResponseOrder(order *db.Order) *pb.OrderResponse {
	if order == nil {
		return nil
	}
	return &pb.OrderResponse{
		Id:         int32(order.OrderID),
		MerchantId: int32(order.MerchantID),
		CashierId:  int32(order.CashierID),
		TotalPrice: int32(order.TotalPrice),
		CreatedAt:  order.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:  order.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
	}
}

func mapResponsesOrder(orders []*db.GetOrdersRow) []*pb.OrderResponse {
	var mappedOrders []*pb.OrderResponse
	for _, order := range orders {
		if order == nil {
			continue
		}
		mappedOrders = append(mappedOrders, &pb.OrderResponse{
			Id:         int32(order.OrderID),
			MerchantId: int32(order.MerchantID),
			CashierId:  int32(order.CashierID),
			TotalPrice: int32(order.TotalPrice),
			CreatedAt:  order.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			UpdatedAt:  order.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		})
	}
	return mappedOrders
}

func mapResponseOrderDeleteAt(order *db.Order) *pb.OrderResponseDeleteAt {
	if order == nil {
		return nil
	}
	var deletedAt *wrapperspb.StringValue
	if order.DeletedAt.Valid {
		deletedAt = wrapperspb.String(order.DeletedAt.Time.Format("2006-01-02 15:04:05"))
	}

	return &pb.OrderResponseDeleteAt{
		Id:         int32(order.OrderID),
		MerchantId: int32(order.MerchantID),
		CashierId:  int32(order.CashierID),
		TotalPrice: int32(order.TotalPrice),
		CreatedAt:  order.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:  order.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		DeletedAt:  deletedAt,
	}
}

func mapResponsesOrderActive(orders []*db.GetOrdersActiveRow) []*pb.OrderResponseDeleteAt {
	var mappedOrders []*pb.OrderResponseDeleteAt
	for _, order := range orders {
		if order == nil {
			continue
		}
		var deletedAt *wrapperspb.StringValue
		if order.DeletedAt.Valid {
			deletedAt = wrapperspb.String(order.DeletedAt.Time.Format("2006-01-02 15:04:05"))
		}
		mappedOrders = append(mappedOrders, &pb.OrderResponseDeleteAt{
			Id:         int32(order.OrderID),
			MerchantId: int32(order.MerchantID),
			CashierId:  int32(order.CashierID),
			TotalPrice: int32(order.TotalPrice),
			CreatedAt:  order.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			UpdatedAt:  order.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
			DeletedAt:  deletedAt,
		})
	}
	return mappedOrders
}

func mapResponsesOrderTrashed(orders []*db.GetOrdersTrashedRow) []*pb.OrderResponseDeleteAt {
	var mappedOrders []*pb.OrderResponseDeleteAt
	for _, order := range orders {
		if order == nil {
			continue
		}
		var deletedAt *wrapperspb.StringValue
		if order.DeletedAt.Valid {
			deletedAt = wrapperspb.String(order.DeletedAt.Time.Format("2006-01-02 15:04:05"))
		}
		mappedOrders = append(mappedOrders, &pb.OrderResponseDeleteAt{
			Id:         int32(order.OrderID),
			MerchantId: int32(order.MerchantID),
			CashierId:  int32(order.CashierID),
			TotalPrice: int32(order.TotalPrice),
			CreatedAt:  order.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			UpdatedAt:  order.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
			DeletedAt:  deletedAt,
		})
	}
	return mappedOrders
}

func mapResponseOrderMonthlyTotalRevenue(row *db.GetMonthlyTotalRevenueRow) *pb.OrderMonthlyTotalRevenueResponse {
	if row == nil {
		return nil
	}
	return &pb.OrderMonthlyTotalRevenueResponse{
		Year:           row.Year,
		Month:          row.Month,
		TotalRevenue:   int32(row.TotalRevenue),
		TotalItemsSold: 0,
	}
}

func mapResponseOrderMonthlyTotalRevenues(c []*db.GetMonthlyTotalRevenueRow) []*pb.OrderMonthlyTotalRevenueResponse {
	var orderRecords []*pb.OrderMonthlyTotalRevenueResponse
	for _, row := range c {
		orderRecords = append(orderRecords, mapResponseOrderMonthlyTotalRevenue(row))
	}
	return orderRecords
}

func mapResponseOrderMonthlyTotalRevenueByMerchant(row *db.GetMonthlyTotalRevenueByMerchantRow) *pb.OrderMonthlyTotalRevenueResponse {
	if row == nil {
		return nil
	}
	return &pb.OrderMonthlyTotalRevenueResponse{
		Year:           row.Year,
		Month:          row.Month,
		TotalRevenue:   int32(row.TotalRevenue),
		TotalItemsSold: 0,
	}
}

func mapResponseOrderMonthlyTotalRevenuesByMerchant(c []*db.GetMonthlyTotalRevenueByMerchantRow) []*pb.OrderMonthlyTotalRevenueResponse {
	var orderRecords []*pb.OrderMonthlyTotalRevenueResponse
	for _, row := range c {
		orderRecords = append(orderRecords, mapResponseOrderMonthlyTotalRevenueByMerchant(row))
	}
	return orderRecords
}

func mapResponseOrderYearlyTotalRevenue(row *db.GetYearlyTotalRevenueRow) *pb.OrderYearlyTotalRevenueResponse {
	if row == nil {
		return nil
	}
	return &pb.OrderYearlyTotalRevenueResponse{
		Year:         row.Year,
		TotalRevenue: int32(row.TotalRevenue),
	}
}

func mapResponseOrderYearlyTotalRevenues(c []*db.GetYearlyTotalRevenueRow) []*pb.OrderYearlyTotalRevenueResponse {
	var orderRecords []*pb.OrderYearlyTotalRevenueResponse
	for _, row := range c {
		orderRecords = append(orderRecords, mapResponseOrderYearlyTotalRevenue(row))
	}
	return orderRecords
}

func mapResponseOrderYearlyTotalRevenueByMerchant(row *db.GetYearlyTotalRevenueByMerchantRow) *pb.OrderYearlyTotalRevenueResponse {
	if row == nil {
		return nil
	}
	return &pb.OrderYearlyTotalRevenueResponse{
		Year:         row.Year,
		TotalRevenue: int32(row.TotalRevenue),
	}
}

func mapResponseOrderYearlyTotalRevenuesByMerchant(c []*db.GetYearlyTotalRevenueByMerchantRow) []*pb.OrderYearlyTotalRevenueResponse {
	var orderRecords []*pb.OrderYearlyTotalRevenueResponse
	for _, row := range c {
		orderRecords = append(orderRecords, mapResponseOrderYearlyTotalRevenueByMerchant(row))
	}
	return orderRecords
}

func mapResponsesOrderMonthlyPrices(c []*db.GetMonthlyOrderRow) []*pb.OrderMonthlyResponse {
	var categoryRecords []*pb.OrderMonthlyResponse
	for _, category := range c {
		if category == nil {
			continue
		}
		categoryRecords = append(categoryRecords, &pb.OrderMonthlyResponse{
			Month:          category.Month,
			OrderCount:     int32(category.OrderCount),
			TotalRevenue:   int32(category.TotalRevenue),
			TotalItemsSold: int32(category.TotalItemsSold),
		})
	}
	return categoryRecords
}

func mapResponsesOrderMonthlyPricesByMerchant(c []*db.GetMonthlyOrderByMerchantRow) []*pb.OrderMonthlyResponse {
	var categoryRecords []*pb.OrderMonthlyResponse
	for _, category := range c {
		if category == nil {
			continue
		}
		categoryRecords = append(categoryRecords, &pb.OrderMonthlyResponse{
			Month:          category.Month,
			OrderCount:     int32(category.OrderCount),
			TotalRevenue:   int32(category.TotalRevenue),
			TotalItemsSold: int32(category.TotalItemsSold),
		})
	}
	return categoryRecords
}

func mapResponsesOrderYearlyPrices(c []*db.GetYearlyOrderRow) []*pb.OrderYearlyResponse {
	var categoryRecords []*pb.OrderYearlyResponse
	for _, category := range c {
		if category == nil {
			continue
		}
		categoryRecords = append(categoryRecords, &pb.OrderYearlyResponse{
			Year:               category.Year,
			OrderCount:         int32(category.OrderCount),
			TotalRevenue:       int32(category.TotalRevenue),
			TotalItemsSold:     int32(category.TotalItemsSold),
			ActiveCashiers:     int32(category.ActiveCashiers),
			UniqueProductsSold: int32(category.UniqueProductsSold),
		})
	}
	return categoryRecords
}

func mapResponsesOrderYearlyPricesByMerchant(c []*db.GetYearlyOrderByMerchantRow) []*pb.OrderYearlyResponse {
	var categoryRecords []*pb.OrderYearlyResponse
	for _, category := range c {
		if category == nil {
			continue
		}
		categoryRecords = append(categoryRecords, &pb.OrderYearlyResponse{
			Year:               category.Year,
			OrderCount:         int32(category.OrderCount),
			TotalRevenue:       int32(category.TotalRevenue),
			TotalItemsSold:     int32(category.TotalItemsSold),
			ActiveCashiers:     int32(category.ActiveCashiers),
			UniqueProductsSold: int32(category.UniqueProductsSold),
		})
	}
	return categoryRecords
}
