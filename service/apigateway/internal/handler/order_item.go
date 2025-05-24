package handler

import (
	"net/http"
	"strconv"

	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
	orderitem_errors "github.com/MamangRust/monolith-point-of-sale-shared/errors/order_item_errors"
	response_api "github.com/MamangRust/monolith-point-of-sale-shared/mapper/response/api"
	"github.com/MamangRust/monolith-point-of-sale-shared/pb"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type orderItemHandleApi struct {
	client  pb.OrderItemServiceClient
	logger  logger.LoggerInterface
	mapping response_api.OrderItemResponseMapper
}

func NewHandlerOrderItem(
	router *echo.Echo,
	client pb.OrderItemServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.OrderItemResponseMapper,
) *orderItemHandleApi {
	categoryHandler := &orderItemHandleApi{
		client:  client,
		logger:  logger,
		mapping: mapping,
	}

	routercategory := router.Group("/api/order-item")

	routercategory.GET("", categoryHandler.FindAllOrderItems)
	routercategory.GET("/:order_id", categoryHandler.FindOrderItemByOrder)
	routercategory.GET("/active", categoryHandler.FindByActive)
	routercategory.GET("/trashed", categoryHandler.FindByTrashed)

	return categoryHandler
}

// @Security Bearer
// @Summary Find all order items
// @Tags Order-Item
// @Description Retrieve a list of all order items
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationOrderItem "List of order items"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve order item data"
// @Router /api/order-item [get]
func (h *orderItemHandleApi) FindAllOrderItems(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	search := c.QueryParam("search")

	ctx := c.Request().Context()

	req := &pb.FindAllOrderItemRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch order-items", zap.Error(err))
		return orderitem_errors.ErrApiOrderItemFailedFindAll(c)
	}

	so := h.mapping.ToApiResponsePaginationOrderItem(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve active order items
// @Tags Order-Item
// @Description Retrieve a list of active order items
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationOrderItemDeleteAt "List of active order items"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve order item data"
// @Router /api/order-item/active [get]
func (h *orderItemHandleApi) FindByActive(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	search := c.QueryParam("search")

	ctx := c.Request().Context()

	req := &pb.FindAllOrderItemRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch active order-items", zap.Error(err))
		return orderitem_errors.ErrApiOrderItemFailedFindByActive(c)
	}

	so := h.mapping.ToApiResponsePaginationOrderItemDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve trashed order items
// @Tags Order-Item
// @Description Retrieve a list of trashed order items
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationOrderItemDeleteAt "List of trashed order items"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve order item data"
// @Router /api/order-item/trashed [get]
func (h *orderItemHandleApi) FindByTrashed(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	search := c.QueryParam("search")

	ctx := c.Request().Context()

	req := &pb.FindAllOrderItemRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch archived order-items", zap.Error(err))
		return orderitem_errors.ErrApiOrderItemFailedFindByTrashed(c)
	}

	so := h.mapping.ToApiResponsePaginationOrderItemDeleteAt(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Find order items by order ID
// @Tags Order-Item
// @Description Retrieve order items by order ID
// @Accept json
// @Produce json
// @Param order_id path int true "Order ID"
// @Success 200 {object} response.ApiResponsesOrderItem "List of order items by order ID"
// @Failure 400 {object} response.ErrorResponse "Invalid order ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve order item data"
// @Router /api/order-item/order/{order_id} [get]
func (h *orderItemHandleApi) FindOrderItemByOrder(c echo.Context) error {
	orderID, err := strconv.Atoi(c.Param("order_id"))

	if err != nil {
		h.logger.Debug("Invalid order ID format", zap.Error(err))
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "invalid_input",
			Message: "Please provide a valid order ID.",
			Code:    http.StatusBadRequest,
		})
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdOrderItemRequest{
		Id: int32(orderID),
	}

	res, err := h.client.FindOrderItemByOrder(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch order item details", zap.Error(err))
		return orderitem_errors.ErrApiOrderItemFailedFindByOrderId(c)
	}

	so := h.mapping.ToApiResponsesOrderItem(res)

	return c.JSON(http.StatusOK, so)
}
