package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	gateway_cache "github.com/MamangRust/monolith-point-of-sale-apigateway/cache/gateway_cache"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors"
	response_api "github.com/MamangRust/monolith-point-of-sale-shared/mapper/response/api"
	"github.com/MamangRust/monolith-point-of-sale-shared/pb"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/types/known/emptypb"
)

type orderHandleApi struct {
	client     pb.OrderServiceClient
	logger     logger.LoggerInterface
	mapping    response_api.OrderResponseMapper
	apiHandler errors.ApiHandler
	cache      *gateway_cache.GatewayCache
}

func NewHandlerOrder(
	router *echo.Echo,
	client pb.OrderServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.OrderResponseMapper,
	apiHandler errors.ApiHandler,
	cache *gateway_cache.GatewayCache,
) *orderHandleApi {
	orderHandler := &orderHandleApi{
		client:     client,
		logger:     logger,
		mapping:    mapping,
		apiHandler: apiHandler,
		cache:      cache,
	}

	routerOrder := router.Group("/api/order")

	routerOrder.GET("", apiHandler.Handle("get-order-findallorders", orderHandler.FindAllOrders))
	routerOrder.GET("/:id", apiHandler.Handle("get-order-findbyid", orderHandler.FindById))
	routerOrder.GET("/active", apiHandler.Handle("get-order-findbyactive", orderHandler.FindByActive))
	routerOrder.GET("/trashed", apiHandler.Handle("get-order-findbytrashed", orderHandler.FindByTrashed))

	routerOrder.GET("/monthly-total-revenue", apiHandler.Handle("get-order-findmonthlytotalrevenue", orderHandler.FindMonthlyTotalRevenue))
	routerOrder.GET("/yearly-total-revenue", apiHandler.Handle("get-order-findyearlytotalrevenue", orderHandler.FindYearlyTotalRevenue))
	routerOrder.GET("/merchant/monthly-total-revenue", apiHandler.Handle("get-order-findmonthlytotalrevenuebymerchant", orderHandler.FindMonthlyTotalRevenueByMerchant))
	routerOrder.GET("/merchant/yearly-total-revenue", apiHandler.Handle("get-order-findyearlytotalrevenuebymerchant", orderHandler.FindYearlyTotalRevenueByMerchant))

	routerOrder.GET("/monthly-revenue", apiHandler.Handle("get-order-findmonthlyrevenue", orderHandler.FindMonthlyRevenue))
	routerOrder.GET("/yearly-revenue", apiHandler.Handle("get-order-findyearlyrevenue", orderHandler.FindYearlyRevenue))
	routerOrder.GET("/merchant/monthly-revenue", apiHandler.Handle("get-order-findmonthlyrevenuebymerchant", orderHandler.FindMonthlyRevenueByMerchant))
	routerOrder.GET("/merchant/yearly-revenue", apiHandler.Handle("get-order-findyearlyrevenuebymerchant", orderHandler.FindYearlyRevenueByMerchant))

	routerOrder.POST("/create", apiHandler.Handle("post-order-create", orderHandler.Create))
	routerOrder.POST("/update/:id", apiHandler.Handle("post-order-update", orderHandler.Update))

	routerOrder.POST("/trashed/:id", apiHandler.Handle("post-order-trashedorder", orderHandler.TrashedOrder))
	routerOrder.POST("/restore/:id", apiHandler.Handle("post-order-restoreorder", orderHandler.RestoreOrder))
	routerOrder.DELETE("/permanent/:id", apiHandler.Handle("delete-order-deleteorderpermanent", orderHandler.DeleteOrderPermanent))

	routerOrder.POST("/restore/all", apiHandler.Handle("post-order-restoreallorder", orderHandler.RestoreAllOrder))
	routerOrder.POST("/permanent/all", apiHandler.Handle("delete-order-deleteallorderpermanent", orderHandler.DeleteAllOrderPermanent))

	return orderHandler
}
// List all order (paginated)
// List all order (paginated)
// List all order (paginated)
// @Summary List all order (paginated)
// @Tags Order
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Param search query string false "Search keyword"
// @Success 200 {object} response.ApiResponsePaginationOrder
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/order [get]
// @Security BearerAuth
func (h *orderHandleApi) FindAllOrders(c echo.Context) error {
	ctx := c.Request().Context()
	page := parseQueryInt(c, "page", 1)
	pageSize := parseQueryInt(c, "page_size", 10)
	search := c.QueryParam("search")

	cacheKey := fmt.Sprintf("order:findallorders:page_%d:size_%d:search_%s", page, pageSize, search)
	if cached, found := gateway_cache.Get[response.ApiResponsePaginationOrder](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	req := &pb.FindAllOrderRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponsePaginationOrder(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}
// Get order by ID
// Get order by ID
// Get order by ID
// @Summary Get order by ID
// @Tags Order
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} response.ApiResponseOrder
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/order/:id [get]
// @Security BearerAuth
func (h *orderHandleApi) FindById(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("invalid order ID")
	}

	cacheKey := fmt.Sprintf("order:findbyid:id_%d", id)
	if cached, found := gateway_cache.Get[response.ApiResponseOrder](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	req := &pb.FindByIdOrderRequest{
		Id: int32(id),
	}

	res, err := h.client.FindById(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseOrder(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}
// List active order
// List active order
// List active order
// @Summary List active order
// @Tags Order
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Param search query string false "Search keyword"
// @Success 200 {object} response.ApiResponsePaginationOrderDeleteAt
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/order/active [get]
// @Security BearerAuth
func (h *orderHandleApi) FindByActive(c echo.Context) error {
	ctx := c.Request().Context()
	page := parseQueryInt(c, "page", 1)
	pageSize := parseQueryInt(c, "page_size", 10)
	search := c.QueryParam("search")

	cacheKey := fmt.Sprintf("order:findbyactive:page_%d:size_%d:search_%s", page, pageSize, search)
	if cached, found := gateway_cache.Get[response.ApiResponsePaginationOrderDeleteAt](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	req := &pb.FindAllOrderRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponsePaginationOrderDeleteAt(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}
// List trashed order
// List trashed order
// List trashed order
// @Summary List trashed order
// @Tags Order
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Param search query string false "Search keyword"
// @Success 200 {object} response.ApiResponsePaginationOrderDeleteAt
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/order/trashed [get]
// @Security BearerAuth
func (h *orderHandleApi) FindByTrashed(c echo.Context) error {
	ctx := c.Request().Context()
	page := parseQueryInt(c, "page", 1)
	pageSize := parseQueryInt(c, "page_size", 10)
	search := c.QueryParam("search")

	cacheKey := fmt.Sprintf("order:findbytrashed:page_%d:size_%d:search_%s", page, pageSize, search)
	if cached, found := gateway_cache.Get[response.ApiResponsePaginationOrderDeleteAt](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	req := &pb.FindAllOrderRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponsePaginationOrderDeleteAt(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}
// Get order statistics
// Get order statistics
// Get order statistics
// @Summary Get order statistics
// @Tags Order
// @Accept json
// @Produce json
// @Param year query int false "Year (e.g. 2026)"
// @Param month query int false "Month (1-12)"
// @Success 200 {object} response.ApiResponseOrderMonthlyTotalRevenue
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/order/monthly-total-revenue [get]
// @Security BearerAuth
func (h *orderHandleApi) FindMonthlyTotalRevenue(c echo.Context) error {
	ctx := c.Request().Context()
	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}
	month, err := parseQueryIntWithValidation(c, "month", 1, 12)
	if err != nil {
		return errors.NewBadRequestError("invalid month")
	}

	cacheKey := fmt.Sprintf("order:findmonthlytotalrevenue:year_%d:month_%d", year, month)
	if cached, found := gateway_cache.Get[response.ApiResponseOrderMonthly](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.client.FindMonthlyTotalRevenue(ctx, &pb.FindYearMonthTotalRevenue{
		Year:  int32(year),
		Month: int32(month),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseMonthlyTotalRevenue(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}
// Get order statistics
// Get order statistics
// Get order statistics
// @Summary Get order statistics
// @Tags Order
// @Accept json
// @Produce json
// @Param year query int false "Year (e.g. 2026)"
// @Param month query int false "Month (1-12)"
// @Success 200 {object} response.ApiResponseOrderYearlyTotalRevenue
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/order/yearly-total-revenue [get]
// @Security BearerAuth
func (h *orderHandleApi) FindYearlyTotalRevenue(c echo.Context) error {
	ctx := c.Request().Context()
	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}

	cacheKey := fmt.Sprintf("order:findyearlytotalrevenue:year_%d", year)
	if cached, found := gateway_cache.Get[response.ApiResponseOrderYearly](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.client.FindYearlyTotalRevenue(ctx, &pb.FindYearTotalRevenue{
		Year: int32(year),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseYearlyTotalRevenue(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}
// Get order statistics
// Get order statistics
// Get order statistics
// @Summary Get order statistics
// @Tags Order
// @Accept json
// @Produce json
// @Param year query int false "Year (e.g. 2026)"
// @Param month query int false "Month (1-12)"
// @Success 200 {object} response.ApiResponseOrderMonthlyTotalRevenue
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/order/merchant/monthly-total-revenue [get]
// @Security BearerAuth
func (h *orderHandleApi) FindMonthlyTotalRevenueByMerchant(c echo.Context) error {
	ctx := c.Request().Context()
	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}
	month, err := parseQueryIntWithValidation(c, "month", 1, 12)
	if err != nil {
		return errors.NewBadRequestError("invalid month")
	}
	merchantID, err := parseQueryIntWithValidation(c, "merchant_id", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid merchant_id")
	}

	cacheKey := fmt.Sprintf("order:findmonthlytotalrevenuebymerchant:year_%d:month_%d:merchantID_%d", year, month, merchantID)
	if cached, found := gateway_cache.Get[response.ApiResponseOrderMonthly](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.client.FindMonthlyTotalRevenueByMerchant(ctx, &pb.FindYearMonthTotalRevenueByMerchant{
		Year:       int32(year),
		Month:      int32(month),
		MerchantId: int32(merchantID),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseMonthlyTotalRevenue(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}
// Get order statistics
// Get order statistics
// Get order statistics
// @Summary Get order statistics
// @Tags Order
// @Accept json
// @Produce json
// @Param year query int false "Year (e.g. 2026)"
// @Param month query int false "Month (1-12)"
// @Success 200 {object} response.ApiResponseOrderYearlyTotalRevenue
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/order/merchant/yearly-total-revenue [get]
// @Security BearerAuth
func (h *orderHandleApi) FindYearlyTotalRevenueByMerchant(c echo.Context) error {
	ctx := c.Request().Context()
	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}
	merchantID, err := parseQueryIntWithValidation(c, "merchant_id", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid merchant_id")
	}

	cacheKey := fmt.Sprintf("order:findyearlytotalrevenuebymerchant:year_%d:merchantID_%d", year, merchantID)
	if cached, found := gateway_cache.Get[response.ApiResponseOrderYearly](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.client.FindYearlyTotalRevenueByMerchant(ctx, &pb.FindYearTotalRevenueByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchantID),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseYearlyTotalRevenue(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}
// Get order statistics
// Get order statistics
// Get order statistics
// @Summary Get order statistics
// @Tags Order
// @Accept json
// @Produce json
// @Param year query int false "Year (e.g. 2026)"
// @Param month query int false "Month (1-12)"
// @Success 200 {object} response.ApiResponseOrderMonthly
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/order/monthly-revenue [get]
// @Security BearerAuth
func (h *orderHandleApi) FindMonthlyRevenue(c echo.Context) error {
	ctx := c.Request().Context()
	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}

	cacheKey := fmt.Sprintf("order:findmonthlyrevenue:year_%d", year)
	if cached, found := gateway_cache.Get[response.ApiResponseOrderMonthly](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.client.FindMonthlyRevenue(ctx, &pb.FindYearOrder{
		Year: int32(year),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseMonthlyOrder(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}
// Get order statistics
// Get order statistics
// Get order statistics
// @Summary Get order statistics
// @Tags Order
// @Accept json
// @Produce json
// @Param year query int false "Year (e.g. 2026)"
// @Param month query int false "Month (1-12)"
// @Success 200 {object} response.ApiResponseOrderYearly
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/order/yearly-revenue [get]
// @Security BearerAuth
func (h *orderHandleApi) FindYearlyRevenue(c echo.Context) error {
	ctx := c.Request().Context()
	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}

	cacheKey := fmt.Sprintf("order:findyearlyrevenue:year_%d", year)
	if cached, found := gateway_cache.Get[response.ApiResponseOrderYearly](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.client.FindYearlyRevenue(ctx, &pb.FindYearOrder{
		Year: int32(year),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseYearlyOrder(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}
// Get order statistics
// Get order statistics
// Get order statistics
// @Summary Get order statistics
// @Tags Order
// @Accept json
// @Produce json
// @Param year query int false "Year (e.g. 2026)"
// @Param month query int false "Month (1-12)"
// @Success 200 {object} response.ApiResponseOrderMonthly
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/order/merchant/monthly-revenue [get]
// @Security BearerAuth
func (h *orderHandleApi) FindMonthlyRevenueByMerchant(c echo.Context) error {
	ctx := c.Request().Context()
	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}
	merchant, err := parseQueryIntWithValidation(c, "merchant_id", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid merchant_id")
	}

	cacheKey := fmt.Sprintf("order:findmonthlyrevenuebymerchant:year_%d:merchantID_%d", year, merchant)
	if cached, found := gateway_cache.Get[response.ApiResponseOrderMonthly](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.client.FindMonthlyRevenueByMerchant(ctx, &pb.FindYearOrderByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchant),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseMonthlyOrder(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}
// Get order statistics
// Get order statistics
// Get order statistics
// @Summary Get order statistics
// @Tags Order
// @Accept json
// @Produce json
// @Param year query int false "Year (e.g. 2026)"
// @Param month query int false "Month (1-12)"
// @Success 200 {object} response.ApiResponseOrderYearly
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/order/merchant/yearly-revenue [get]
// @Security BearerAuth
func (h *orderHandleApi) FindYearlyRevenueByMerchant(c echo.Context) error {
	ctx := c.Request().Context()
	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}
	merchant, err := parseQueryIntWithValidation(c, "merchant_id", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid merchant_id")
	}

	cacheKey := fmt.Sprintf("order:findyearlyrevenuebymerchant:year_%d:merchantID_%d", year, merchant)
	if cached, found := gateway_cache.Get[response.ApiResponseOrderYearly](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.client.FindYearlyRevenueByMerchant(ctx, &pb.FindYearOrderByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchant),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseYearlyOrder(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}
// Create order
// Create order
// Create order
// @Summary Create order
// @Tags Order
// @Accept json
// @Produce json
// @Param request body requests.CreateOrderRequest true "Request body"
// @Success 200 {object} response.ApiResponseOrder
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/order/create [post]
// @Security BearerAuth
func (h *orderHandleApi) Create(c echo.Context) error {
	ctx := c.Request().Context()
	var body requests.CreateOrderRequest
	if err := c.Bind(&body); err != nil {
		return errors.NewBadRequestError("bind failed")
	}
	if err := body.Validate(); err != nil {
		return errors.NewBadRequestError("validation failed")
	}

	grpcReq := &pb.CreateOrderRequest{
		MerchantId: int32(body.MerchantID),
		CashierId:  int32(body.CashierID),
	}

	for _, item := range body.Items {
		grpcReq.Items = append(grpcReq.Items, &pb.CreateOrderItemRequest{
			ProductId: int32(item.ProductID),
			Quantity:  int32(item.Quantity),
		})
	}

	res, err := h.client.Create(ctx, grpcReq)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseOrder(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "order:*")
	return c.JSON(http.StatusOK, so)
}
// Update order
// Update order
// Update order
// @Summary Update order
// @Tags Order
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Param request body requests.UpdateOrderRequest true "Request body"
// @Success 200 {object} response.ApiResponseOrder
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/order/update/:id [post]
// @Security BearerAuth
func (h *orderHandleApi) Update(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("invalid order ID")
	}

	var body requests.UpdateOrderRequest
	if err := c.Bind(&body); err != nil {
		return errors.NewBadRequestError("bind failed")
	}
	if err := body.Validate(); err != nil {
		return errors.NewBadRequestError("validation failed")
	}

	grpcReq := &pb.UpdateOrderRequest{
		OrderId: int32(id),
		Items:   []*pb.UpdateOrderItemRequest{},
	}

	for _, item := range body.Items {
		grpcReq.Items = append(grpcReq.Items, &pb.UpdateOrderItemRequest{
			OrderItemId: int32(item.OrderItemID),
			ProductId:   int32(item.ProductID),
			Quantity:    int32(item.Quantity),
		})
	}

	res, err := h.client.Update(ctx, grpcReq)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseOrder(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "order:*")
	return c.JSON(http.StatusOK, so)
}
// Trash order
// Trash order
// Trash order
// @Summary Trash order
// @Tags Order
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} response.ApiResponseOrderDeleteAt
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/order/trashed/:id [post]
// @Security BearerAuth
func (h *orderHandleApi) TrashedOrder(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("invalid order ID")
	}

	req := &pb.FindByIdOrderRequest{
		Id: int32(id),
	}

	res, err := h.client.TrashedOrder(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseOrderDeleteAt(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "order:*")
	return c.JSON(http.StatusOK, so)
}
// Restore order
// Restore order
// Restore order
// @Summary Restore order
// @Tags Order
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} response.ApiResponseOrderDeleteAt
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/order/restore/:id [post]
// @Security BearerAuth
func (h *orderHandleApi) RestoreOrder(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("invalid order ID")
	}

	req := &pb.FindByIdOrderRequest{
		Id: int32(id),
	}

	res, err := h.client.RestoreOrder(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseOrderDeleteAt(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "order:*")
	return c.JSON(http.StatusOK, so)
}
// Delete order permanently
// Delete order permanently
// Delete order permanently
// @Summary Delete order permanently
// @Tags Order
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} response.ApiResponseOrderDelete
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/order/permanent/:id [delete]
// @Security BearerAuth
func (h *orderHandleApi) DeleteOrderPermanent(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("invalid order ID")
	}

	req := &pb.FindByIdOrderRequest{
		Id: int32(id),
	}

	res, err := h.client.DeleteOrderPermanent(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseOrderDelete(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "order:*")
	return c.JSON(http.StatusOK, so)
}
// Restore all trashed order
// Restore all trashed order
// Restore all trashed order
// @Summary Restore all trashed order
// @Tags Order
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseOrderAll
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/order/restore/all [post]
// @Security BearerAuth
func (h *orderHandleApi) RestoreAllOrder(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := h.client.RestoreAllOrder(ctx, &emptypb.Empty{})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseOrderAll(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "order:*")
	return c.JSON(http.StatusOK, so)
}
// Delete all trashed order permanently
// Delete all trashed order permanently
// Delete all trashed order permanently
// @Summary Delete all trashed order permanently
// @Tags Order
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseOrderAll
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/order/permanent/all [post]
// @Security BearerAuth
func (h *orderHandleApi) DeleteAllOrderPermanent(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := h.client.DeleteAllOrderPermanent(ctx, &emptypb.Empty{})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseOrderAll(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "order:*")
	return c.JSON(http.StatusOK, so)
}
