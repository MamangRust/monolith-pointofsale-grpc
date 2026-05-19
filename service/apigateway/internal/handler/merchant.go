package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	gateway_cache "github.com/MamangRust/monolith-point-of-sale-apigateway/internal/redis/api/gateway_cache"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors"
	response_api "github.com/MamangRust/monolith-point-of-sale-shared/mapper/response/api"
	"github.com/MamangRust/monolith-point-of-sale-shared/pb"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/types/known/emptypb"
)

type merchantHandleApi struct {
	client     pb.MerchantServiceClient
	logger     logger.LoggerInterface
	mapping    response_api.MerchantResponseMapper
	apiHandler errors.ApiHandler
	cache      *gateway_cache.GatewayCache
}

func NewHandlerMerchant(
	router *echo.Echo,
	client pb.MerchantServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.MerchantResponseMapper,
	apiHandler errors.ApiHandler,
	cache *gateway_cache.GatewayCache,
) *merchantHandleApi {
	merchantHandler := &merchantHandleApi{
		client:     client,
		logger:     logger,
		mapping:    mapping,
		apiHandler: apiHandler,
		cache:      cache,
	}

	routerMerchant := router.Group("/api/merchant")

	routerMerchant.GET("", apiHandler.Handle("get-merchant-findallmerchant", merchantHandler.FindAllMerchant))
	routerMerchant.GET("/:id", apiHandler.Handle("get-merchant-findbyid", merchantHandler.FindById))
	routerMerchant.GET("/active", apiHandler.Handle("get-merchant-findbyactive", merchantHandler.FindByActive))
	routerMerchant.GET("/trashed", apiHandler.Handle("get-merchant-findbytrashed", merchantHandler.FindByTrashed))

	routerMerchant.POST("/create", apiHandler.Handle("post-merchant-create", merchantHandler.Create))
	routerMerchant.POST("/update/:id", apiHandler.Handle("post-merchant-update", merchantHandler.Update))
	routerMerchant.POST("/update-status/:id", apiHandler.Handle("post-merchant-updatestatus", merchantHandler.UpdateStatus))

	routerMerchant.POST("/trashed/:id", apiHandler.Handle("post-merchant-trashedmerchant", merchantHandler.TrashedMerchant))
	routerMerchant.POST("/restore/:id", apiHandler.Handle("post-merchant-restoremerchant", merchantHandler.RestoreMerchant))
	routerMerchant.DELETE("/permanent/:id", apiHandler.Handle("delete-merchant-deletemerchantpermanent", merchantHandler.DeleteMerchantPermanent))

	routerMerchant.POST("/restore/all", apiHandler.Handle("post-merchant-restoreallmerchant", merchantHandler.RestoreAllMerchant))
	routerMerchant.POST("/permanent/all", apiHandler.Handle("post-merchant-deleteallmerchantpermanent", merchantHandler.DeleteAllMerchantPermanent))

	return merchantHandler
}

func (h *merchantHandleApi) FindAllMerchant(c echo.Context) error {
	ctx := c.Request().Context()
	page := parseQueryInt(c, "page", 1)
	pageSize := parseQueryInt(c, "page_size", 10)
	search := c.QueryParam("search")

	cacheKey := fmt.Sprintf("merchant:findallmerchant:page_%d:size_%d:search_%s", page, pageSize, search)
	if cached, found := gateway_cache.Get[response.ApiResponsePaginationMerchant](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	req := &pb.FindAllMerchantRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponsePaginationMerchant(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *merchantHandleApi) FindById(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("invalid merchant ID")
	}

	cacheKey := fmt.Sprintf("merchant:findbyid:id_%d", id)
	if cached, found := gateway_cache.Get[response.ApiResponseMerchant](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	req := &pb.FindByIdMerchantRequest{
		Id: int32(id),
	}

	res, err := h.client.FindById(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseMerchant(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *merchantHandleApi) FindByActive(c echo.Context) error {
	ctx := c.Request().Context()
	page := parseQueryInt(c, "page", 1)
	pageSize := parseQueryInt(c, "page_size", 10)
	search := c.QueryParam("search")

	cacheKey := fmt.Sprintf("merchant:findbyactive:page_%d:size_%d:search_%s", page, pageSize, search)
	if cached, found := gateway_cache.Get[response.ApiResponsePaginationMerchantDeleteAt](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	req := &pb.FindAllMerchantRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponsePaginationMerchantDeleteAt(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *merchantHandleApi) FindByTrashed(c echo.Context) error {
	ctx := c.Request().Context()
	page := parseQueryInt(c, "page", 1)
	pageSize := parseQueryInt(c, "page_size", 10)
	search := c.QueryParam("search")

	cacheKey := fmt.Sprintf("merchant:findbytrashed:page_%d:size_%d:search_%s", page, pageSize, search)
	if cached, found := gateway_cache.Get[response.ApiResponsePaginationMerchantDeleteAt](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	req := &pb.FindAllMerchantRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponsePaginationMerchantDeleteAt(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *merchantHandleApi) Create(c echo.Context) error {
	ctx := c.Request().Context()
	var body requests.CreateMerchantRequest
	if err := c.Bind(&body); err != nil {
		return errors.NewBadRequestError("bind failed")
	}
	if err := body.Validate(); err != nil {
		return errors.NewBadRequestError("validation failed")
	}

	req := &pb.CreateMerchantRequest{
		UserId:       int32(body.UserID),
		Name:         body.Name,
		Description:  body.Description,
		Address:      body.Address,
		ContactEmail: body.ContactEmail,
		ContactPhone: body.ContactPhone,
		Status:       body.Status,
	}

	res, err := h.client.Create(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseMerchant(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "merchant:*")
	return c.JSON(http.StatusOK, so)
}

func (h *merchantHandleApi) Update(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("invalid merchant ID")
	}

	var body requests.UpdateMerchantRequest
	if err := c.Bind(&body); err != nil {
		return errors.NewBadRequestError("bind failed")
	}
	if err := body.Validate(); err != nil {
		return errors.NewBadRequestError("validation failed")
	}

	req := &pb.UpdateMerchantRequest{
		MerchantId:   int32(id),
		UserId:       int32(body.UserID),
		Name:         body.Name,
		Description:  body.Description,
		Address:      body.Address,
		ContactEmail: body.ContactEmail,
		ContactPhone: body.ContactPhone,
		Status:       body.Status,
	}

	res, err := h.client.Update(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseMerchant(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "merchant:*")
	return c.JSON(http.StatusOK, so)
}

func (h *merchantHandleApi) UpdateStatus(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("invalid merchant ID")
	}

	var body requests.UpdateMerchantStatusRequest
	if err := c.Bind(&body); err != nil {
		return errors.NewBadRequestError("bind failed")
	}
	if err := body.Validate(); err != nil {
		return errors.NewBadRequestError("validation failed")
	}

	req := &pb.UpdateMerchantStatusRequest{
		MerchantId: int32(id),
		Status:     body.Status,
	}

	res, err := h.client.UpdateMerchantStatus(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseMerchant(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "merchant:*")
	return c.JSON(http.StatusOK, so)
}

func (h *merchantHandleApi) TrashedMerchant(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("invalid merchant ID")
	}

	req := &pb.FindByIdMerchantRequest{
		Id: int32(id),
	}

	res, err := h.client.TrashedMerchant(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseMerchantDeleteAt(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "merchant:*")
	return c.JSON(http.StatusOK, so)
}

func (h *merchantHandleApi) RestoreMerchant(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("invalid merchant ID")
	}

	req := &pb.FindByIdMerchantRequest{
		Id: int32(id),
	}

	res, err := h.client.RestoreMerchant(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseMerchant(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "merchant:*")
	return c.JSON(http.StatusOK, so)
}

func (h *merchantHandleApi) DeleteMerchantPermanent(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("invalid merchant ID")
	}

	req := &pb.FindByIdMerchantRequest{
		Id: int32(id),
	}

	res, err := h.client.DeleteMerchantPermanent(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseMerchantDelete(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "merchant:*")
	return c.JSON(http.StatusOK, so)
}

func (h *merchantHandleApi) RestoreAllMerchant(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := h.client.RestoreAllMerchant(ctx, &emptypb.Empty{})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseMerchantAll(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "merchant:*")
	return c.JSON(http.StatusOK, so)
}

func (h *merchantHandleApi) DeleteAllMerchantPermanent(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := h.client.DeleteAllMerchantPermanent(ctx, &emptypb.Empty{})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseMerchantAll(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "merchant:*")
	return c.JSON(http.StatusOK, so)
}
