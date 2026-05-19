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

type merchantDocumentHandleApi struct {
	merchantDocument pb.MerchantDocumentServiceClient
	logger           logger.LoggerInterface
	mapping          response_api.MerchantDocumentResponseMapper
	apiHandler       errors.ApiHandler
	cache            *gateway_cache.GatewayCache
}

func NewHandlerMerchantDocument(
	router *echo.Echo,
	merchantDocument pb.MerchantDocumentServiceClient,
	logger logger.LoggerInterface,
	ma response_api.MerchantDocumentResponseMapper,
	apiHandler errors.ApiHandler,
	cache *gateway_cache.GatewayCache,
) *merchantDocumentHandleApi {
	merchantDocumentHandler := &merchantDocumentHandleApi{
		merchantDocument: merchantDocument,
		logger:           logger,
		mapping:          ma,
		apiHandler:       apiHandler,
		cache:            cache,
	}

	routerMerchantDocument := router.Group("/api/merchant-documents")

	routerMerchantDocument.GET("", apiHandler.Handle("get-merchantdocument-findall", merchantDocumentHandler.FindAll))
	routerMerchantDocument.GET("/:id", apiHandler.Handle("get-merchantdocument-findbyid", merchantDocumentHandler.FindById))
	routerMerchantDocument.GET("/active", apiHandler.Handle("get-merchantdocument-findallactive", merchantDocumentHandler.FindAllActive))
	routerMerchantDocument.GET("/trashed", apiHandler.Handle("get-merchantdocument-findalltrashed", merchantDocumentHandler.FindAllTrashed))

	routerMerchantDocument.POST("/create", apiHandler.Handle("post-merchantdocument-create", merchantDocumentHandler.Create))
	routerMerchantDocument.POST("/updates/:id", apiHandler.Handle("post-merchantdocument-update", merchantDocumentHandler.Update))
	routerMerchantDocument.POST("/update-status/:id", apiHandler.Handle("post-merchantdocument-updatestatus", merchantDocumentHandler.UpdateStatus))

	routerMerchantDocument.POST("/trashed/:id", apiHandler.Handle("post-merchantdocument-trasheddocument", merchantDocumentHandler.TrashedDocument))
	routerMerchantDocument.POST("/restore/:id", apiHandler.Handle("post-merchantdocument-restoredocument", merchantDocumentHandler.RestoreDocument))
	routerMerchantDocument.DELETE("/permanent/:id", apiHandler.Handle("delete-merchantdocument-delete", merchantDocumentHandler.Delete))

	routerMerchantDocument.POST("/restore/all", apiHandler.Handle("post-merchantdocument-restorealldocuments", merchantDocumentHandler.RestoreAllDocuments))
	routerMerchantDocument.POST("/permanent/all", apiHandler.Handle("delete-merchantdocument-deletealldocumentspermanent", merchantDocumentHandler.DeleteAllDocumentsPermanent))

	return merchantDocumentHandler
}

func (h *merchantDocumentHandleApi) FindAll(c echo.Context) error {
	ctx := c.Request().Context()
	page := parseQueryInt(c, "page", 1)
	pageSize := parseQueryInt(c, "page_size", 10)
	search := c.QueryParam("search")

	cacheKey := fmt.Sprintf("merchant_document:findall:page_%d:size_%d:search_%s", page, pageSize, search)
	if cached, found := gateway_cache.Get[response.ApiResponsePaginationMerchantDocument](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	req := &pb.FindAllMerchantDocumentsRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.merchantDocument.FindAll(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponsePaginationMerchantDocument(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *merchantDocumentHandleApi) FindById(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("invalid document ID")
	}

	cacheKey := fmt.Sprintf("merchant_document:findbyid:id_%d", id)
	if cached, found := gateway_cache.Get[response.ApiResponseMerchantDocument](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.merchantDocument.FindById(ctx, &pb.FindMerchantDocumentByIdRequest{
		DocumentId: int32(id),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseMerchantDocument(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *merchantDocumentHandleApi) FindAllActive(c echo.Context) error {
	ctx := c.Request().Context()
	page := parseQueryInt(c, "page", 1)
	pageSize := parseQueryInt(c, "page_size", 10)
	search := c.QueryParam("search")

	cacheKey := fmt.Sprintf("merchant_document:findallactive:page_%d:size_%d:search_%s", page, pageSize, search)
	if cached, found := gateway_cache.Get[response.ApiResponsePaginationMerchantDocument](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	req := &pb.FindAllMerchantDocumentsRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.merchantDocument.FindAllActive(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponsePaginationMerchantDocument(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *merchantDocumentHandleApi) FindAllTrashed(c echo.Context) error {
	ctx := c.Request().Context()
	page := parseQueryInt(c, "page", 1)
	pageSize := parseQueryInt(c, "page_size", 10)
	search := c.QueryParam("search")

	cacheKey := fmt.Sprintf("merchant_document:findalltrashed:page_%d:size_%d:search_%s", page, pageSize, search)
	if cached, found := gateway_cache.Get[response.ApiResponsePaginationMerchantDocumentDeleteAt](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	req := &pb.FindAllMerchantDocumentsRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.merchantDocument.FindAllTrashed(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponsePaginationMerchantDocumentDeleteAt(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *merchantDocumentHandleApi) Create(c echo.Context) error {
	ctx := c.Request().Context()
	var body requests.CreateMerchantDocumentRequest
	if err := c.Bind(&body); err != nil {
		return errors.NewBadRequestError("bind failed")
	}
	if err := body.Validate(); err != nil {
		return errors.NewBadRequestError("validation failed")
	}

	req := &pb.CreateMerchantDocumentRequest{
		MerchantId:   int32(body.MerchantID),
		DocumentType: body.DocumentType,
		DocumentUrl:  body.DocumentUrl,
	}

	res, err := h.merchantDocument.Create(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseMerchantDocument(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "merchant_document:*")
	return c.JSON(http.StatusOK, so)
}

func (h *merchantDocumentHandleApi) Update(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("invalid document ID")
	}

	var body requests.UpdateMerchantDocumentRequest
	if err := c.Bind(&body); err != nil {
		return errors.NewBadRequestError("bind failed")
	}
	if err := body.Validate(); err != nil {
		return errors.NewBadRequestError("validation failed")
	}

	req := &pb.UpdateMerchantDocumentRequest{
		DocumentId:   int32(id),
		MerchantId:   int32(body.MerchantID),
		DocumentType: body.DocumentType,
		DocumentUrl:  body.DocumentUrl,
	}

	res, err := h.merchantDocument.Update(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseMerchantDocument(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "merchant_document:*")
	return c.JSON(http.StatusOK, so)
}

func (h *merchantDocumentHandleApi) UpdateStatus(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("invalid document ID")
	}

	var body requests.UpdateMerchantDocumentStatusRequest
	if err := c.Bind(&body); err != nil {
		return errors.NewBadRequestError("bind failed")
	}
	if err := body.Validate(); err != nil {
		return errors.NewBadRequestError("validation failed")
	}

	req := &pb.UpdateMerchantDocumentStatusRequest{
		DocumentId: int32(id),
		Status:     body.Status,
	}

	res, err := h.merchantDocument.UpdateStatus(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseMerchantDocument(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "merchant_document:*")
	return c.JSON(http.StatusOK, so)
}

func (h *merchantDocumentHandleApi) TrashedDocument(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("invalid document ID")
	}

	req := &pb.TrashedMerchantDocumentRequest{DocumentId: int32(id)}
	res, err := h.merchantDocument.Trashed(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseMerchantDocument(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "merchant_document:*")
	return c.JSON(http.StatusOK, so)
}

func (h *merchantDocumentHandleApi) RestoreDocument(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("invalid document ID")
	}

	req := &pb.RestoreMerchantDocumentRequest{DocumentId: int32(id)}
	res, err := h.merchantDocument.Restore(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseMerchantDocument(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "merchant_document:*")
	return c.JSON(http.StatusOK, so)
}

func (h *merchantDocumentHandleApi) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("invalid document ID")
	}

	req := &pb.DeleteMerchantDocumentPermanentRequest{DocumentId: int32(id)}
	res, err := h.merchantDocument.DeletePermanent(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseMerchantDocumentDeleteAt(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "merchant_document:*")
	return c.JSON(http.StatusOK, so)
}

func (h *merchantDocumentHandleApi) RestoreAllDocuments(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := h.merchantDocument.RestoreAll(ctx, &emptypb.Empty{})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseMerchantDocumentAll(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "merchant_document:*")
	return c.JSON(http.StatusOK, so)
}

func (h *merchantDocumentHandleApi) DeleteAllDocumentsPermanent(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := h.merchantDocument.DeleteAllPermanent(ctx, &emptypb.Empty{})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseMerchantDocumentAll(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "merchant_document:*")
	return c.JSON(http.StatusOK, so)
}
