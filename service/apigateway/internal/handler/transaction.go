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

type transactionHandleApi struct {
	client     pb.TransactionServiceClient
	logger     logger.LoggerInterface
	mapping    response_api.TransactionResponseMapper
	apiHandler errors.ApiHandler
	cache      *gateway_cache.GatewayCache
}

func NewHandlerTransaction(
	router *echo.Echo,
	client pb.TransactionServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.TransactionResponseMapper,
	apiHandler errors.ApiHandler,
	cache *gateway_cache.GatewayCache,
) *transactionHandleApi {
	transactionHandle := &transactionHandleApi{
		client:     client,
		logger:     logger,
		mapping:    mapping,
		apiHandler: apiHandler,
		cache:      cache,
	}

	routerTransaction := router.Group("/api/transaction")

	routerTransaction.GET("", apiHandler.Handle("get-transaction-findalltransaction", transactionHandle.FindAllTransaction))
	routerTransaction.GET("/:id", apiHandler.Handle("get-transaction-findbyid", transactionHandle.FindById))
	routerTransaction.GET("/merchant/:merchant_id", apiHandler.Handle("get-transaction-findbymerchant", transactionHandle.FindByMerchant))
	routerTransaction.GET("/active", apiHandler.Handle("get-transaction-findbyactive", transactionHandle.FindByActive))
	routerTransaction.GET("/trashed", apiHandler.Handle("get-transaction-findbytrashed", transactionHandle.FindByTrashed))

	routerTransaction.GET("/monthly-success", apiHandler.Handle("get-transaction-findmonthstatussuccess", transactionHandle.FindMonthStatusSuccess))
	routerTransaction.GET("/yearly-success", apiHandler.Handle("get-transaction-findyearstatussuccess", transactionHandle.FindYearStatusSuccess))
	routerTransaction.GET("/monthly-failed", apiHandler.Handle("get-transaction-findmonthstatusfailed", transactionHandle.FindMonthStatusFailed))
	routerTransaction.GET("/yearly-failed", apiHandler.Handle("get-transaction-findyearstatusfailed", transactionHandle.FindYearStatusFailed))

	routerTransaction.GET("/merchant/monthly-success", apiHandler.Handle("get-transaction-findmonthstatussuccessbymerchant", transactionHandle.FindMonthStatusSuccessByMerchant))
	routerTransaction.GET("/merchant/yearly-success", apiHandler.Handle("get-transaction-findyearstatussuccessbymerchant", transactionHandle.FindYearStatusSuccessByMerchant))
	routerTransaction.GET("/merchant/monthly-failed", apiHandler.Handle("get-transaction-findmonthstatusfailedbymerchant", transactionHandle.FindMonthStatusFailedByMerchant))
	routerTransaction.GET("/merchant/yearly-failed", apiHandler.Handle("get-transaction-findyearstatusfailedbymerchant", transactionHandle.FindYearStatusFailedByMerchant))

	routerTransaction.GET("/monthly-method-success", apiHandler.Handle("get-transaction-findmonthmethodsuccess", transactionHandle.FindMonthMethodSuccess))
	routerTransaction.GET("/yearly-method-success", apiHandler.Handle("get-transaction-findyearmethodsuccess", transactionHandle.FindYearMethodSuccess))

	routerTransaction.GET("/merchant/monthly-method-success/:merchant_id", apiHandler.Handle("get-transaction-findmonthmethodbymerchantsuccess", transactionHandle.FindMonthMethodByMerchantSuccess))
	routerTransaction.GET("/merchant/yearly-method-success/:merchant_id", apiHandler.Handle("get-transaction-findyearmethodbymerchantsuccess", transactionHandle.FindYearMethodByMerchantSuccess))

	routerTransaction.GET("/monthly-method-failed", apiHandler.Handle("get-transaction-findmonthmethodfailed", transactionHandle.FindMonthMethodFailed))
	routerTransaction.GET("/yearly-method-failed", apiHandler.Handle("get-transaction-findyearmethodfailed", transactionHandle.FindYearMethodFailed))

	routerTransaction.GET("/merchant/monthly-method-failed/:merchant_id", apiHandler.Handle("get-transaction-findmonthmethodbymerchantfailed", transactionHandle.FindMonthMethodByMerchantFailed))
	routerTransaction.GET("/merchant/yearly-method-failed/:merchant_id", apiHandler.Handle("get-transaction-findyearmethodbymerchantfailed", transactionHandle.FindYearMethodByMerchantFailed))

	routerTransaction.POST("/create", apiHandler.Handle("post-transaction-create", transactionHandle.Create))
	routerTransaction.POST("/update/:id", apiHandler.Handle("post-transaction-update", transactionHandle.Update))

	routerTransaction.POST("/trashed/:id", apiHandler.Handle("post-transaction-trashedtransaction", transactionHandle.TrashedTransaction))
	routerTransaction.POST("/restore/:id", apiHandler.Handle("post-transaction-restoretransaction", transactionHandle.RestoreTransaction))
	routerTransaction.DELETE("/permanent/:id", apiHandler.Handle("delete-transaction-deletetransactionpermanent", transactionHandle.DeleteTransactionPermanent))

	routerTransaction.POST("/restore/all", apiHandler.Handle("post-transaction-restorealltransaction", transactionHandle.RestoreAllTransaction))
	routerTransaction.POST("/permanent/all", apiHandler.Handle("delete-transaction-deletealltransactionpermanent", transactionHandle.DeleteAllTransactionPermanent))

	return transactionHandle
}

func (h *transactionHandleApi) FindAllTransaction(c echo.Context) error {
	ctx := c.Request().Context()
	page := parseQueryInt(c, "page", 1)
	pageSize := parseQueryInt(c, "page_size", 10)
	search := c.QueryParam("search")

	cacheKey := fmt.Sprintf("transaction:findalltransaction:page_%d:size_%d:search_%s", page, pageSize, search)
	if cached, found := gateway_cache.Get[response.ApiResponsePaginationTransaction](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	req := &pb.FindAllTransactionRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponsePaginationTransaction(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *transactionHandleApi) FindByMerchant(c echo.Context) error {
	merchantID, err := strconv.Atoi(c.Param("merchant_id"))
	if err != nil || merchantID <= 0 {
		return errors.NewBadRequestError("invalid merchant ID")
	}

	ctx := c.Request().Context()
	page := parseQueryInt(c, "page", 1)
	pageSize := parseQueryInt(c, "page_size", 10)
	search := c.QueryParam("search")

	cacheKey := fmt.Sprintf("transaction:findbymerchant:merchant_%d:page_%d:size_%d:search_%s", merchantID, page, pageSize, search)
	if cached, found := gateway_cache.Get[response.ApiResponsePaginationTransaction](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	req := &pb.FindAllTransactionMerchantRequest{
		MerchantId: int32(merchantID),
		Page:       int32(page),
		PageSize:   int32(pageSize),
		Search:     search,
	}

	res, err := h.client.FindByMerchant(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponsePaginationTransaction(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *transactionHandleApi) FindById(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("invalid transaction ID")
	}

	cacheKey := fmt.Sprintf("transaction:findbyid:id_%d", id)
	if cached, found := gateway_cache.Get[response.ApiResponseTransaction](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	req := &pb.FindByIdTransactionRequest{
		Id: int32(id),
	}

	res, err := h.client.FindById(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseTransaction(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *transactionHandleApi) FindByActive(c echo.Context) error {
	ctx := c.Request().Context()
	page := parseQueryInt(c, "page", 1)
	pageSize := parseQueryInt(c, "page_size", 10)
	search := c.QueryParam("search")

	cacheKey := fmt.Sprintf("transaction:findbyactive:page_%d:size_%d:search_%s", page, pageSize, search)
	if cached, found := gateway_cache.Get[response.ApiResponsePaginationTransactionDeleteAt](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	req := &pb.FindAllTransactionRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponsePaginationTransactionDeleteAt(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *transactionHandleApi) FindByTrashed(c echo.Context) error {
	ctx := c.Request().Context()
	page := parseQueryInt(c, "page", 1)
	pageSize := parseQueryInt(c, "page_size", 10)
	search := c.QueryParam("search")

	cacheKey := fmt.Sprintf("transaction:findbytrashed:page_%d:size_%d:search_%s", page, pageSize, search)
	if cached, found := gateway_cache.Get[response.ApiResponsePaginationTransactionDeleteAt](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	req := &pb.FindAllTransactionRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponsePaginationTransactionDeleteAt(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *transactionHandleApi) FindMonthStatusSuccess(c echo.Context) error {
	ctx := c.Request().Context()
	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}
	month, err := parseQueryIntWithValidation(c, "month", 1, 12)
	if err != nil {
		return errors.NewBadRequestError("invalid month")
	}

	cacheKey := fmt.Sprintf("transaction:findmonthstatussuccess:year_%d:month_%d", year, month)
	if cached, found := gateway_cache.Get[response.ApiResponsesTransactionMonthSuccess](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.client.FindMonthStatusSuccess(ctx, &pb.FindMonthlyTransactionStatus{
		Year:  int32(year),
		Month: int32(month),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseTransactionMonthAmountSuccess(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *transactionHandleApi) FindYearStatusSuccess(c echo.Context) error {
	ctx := c.Request().Context()
	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}

	cacheKey := fmt.Sprintf("transaction:findyearstatussuccess:year_%d", year)
	if cached, found := gateway_cache.Get[response.ApiResponsesTransactionYearSuccess](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.client.FindYearStatusSuccess(ctx, &pb.FindYearlyTransactionStatus{
		Year: int32(year),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseTransactionYearAmountSuccess(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *transactionHandleApi) FindMonthStatusFailed(c echo.Context) error {
	ctx := c.Request().Context()
	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}
	month, err := parseQueryIntWithValidation(c, "month", 1, 12)
	if err != nil {
		return errors.NewBadRequestError("invalid month")
	}

	cacheKey := fmt.Sprintf("transaction:findmonthstatusfailed:year_%d:month_%d", year, month)
	if cached, found := gateway_cache.Get[response.ApiResponsesTransactionMonthFailed](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.client.FindMonthStatusFailed(ctx, &pb.FindMonthlyTransactionStatus{
		Year:  int32(year),
		Month: int32(month),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseTransactionMonthAmountFailed(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *transactionHandleApi) FindYearStatusFailed(c echo.Context) error {
	ctx := c.Request().Context()
	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}

	cacheKey := fmt.Sprintf("transaction:findyearstatusfailed:year_%d", year)
	if cached, found := gateway_cache.Get[response.ApiResponsesTransactionYearFailed](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.client.FindYearStatusFailed(ctx, &pb.FindYearlyTransactionStatus{
		Year: int32(year),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseTransactionYearAmountFailed(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *transactionHandleApi) FindMonthStatusSuccessByMerchant(c echo.Context) error {
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
		return errors.NewBadRequestError("invalid merchant ID")
	}

	cacheKey := fmt.Sprintf("transaction:findmonthstatussuccessbymerchant:year_%d:month_%d:merchant_%d", year, month, merchantID)
	if cached, found := gateway_cache.Get[response.ApiResponsesTransactionMonthSuccess](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.client.FindMonthStatusSuccessByMerchant(ctx, &pb.FindMonthlyTransactionStatusByMerchant{
		Year:       int32(year),
		Month:      int32(month),
		MerchantId: int32(merchantID),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseTransactionMonthAmountSuccess(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *transactionHandleApi) FindYearStatusSuccessByMerchant(c echo.Context) error {
	ctx := c.Request().Context()
	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}
	merchantID, err := parseQueryIntWithValidation(c, "merchant_id", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid merchant ID")
	}

	cacheKey := fmt.Sprintf("transaction:findyearstatussuccessbymerchant:year_%d:merchant_%d", year, merchantID)
	if cached, found := gateway_cache.Get[response.ApiResponsesTransactionYearSuccess](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.client.FindYearStatusSuccessByMerchant(ctx, &pb.FindYearlyTransactionStatusByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchantID),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseTransactionYearAmountSuccess(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *transactionHandleApi) FindMonthStatusFailedByMerchant(c echo.Context) error {
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
		return errors.NewBadRequestError("invalid merchant ID")
	}

	cacheKey := fmt.Sprintf("transaction:findmonthstatusfailedbymerchant:year_%d:month_%d:merchant_%d", year, month, merchantID)
	if cached, found := gateway_cache.Get[response.ApiResponsesTransactionMonthFailed](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.client.FindMonthStatusFailedByMerchant(ctx, &pb.FindMonthlyTransactionStatusByMerchant{
		Year:       int32(year),
		Month:      int32(month),
		MerchantId: int32(merchantID),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseTransactionMonthAmountFailed(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *transactionHandleApi) FindYearStatusFailedByMerchant(c echo.Context) error {
	ctx := c.Request().Context()
	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}
	merchantID, err := parseQueryIntWithValidation(c, "merchant_id", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid merchant ID")
	}

	cacheKey := fmt.Sprintf("transaction:findyearstatusfailedbymerchant:year_%d:merchant_%d", year, merchantID)
	if cached, found := gateway_cache.Get[response.ApiResponsesTransactionYearFailed](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.client.FindYearStatusFailedByMerchant(ctx, &pb.FindYearlyTransactionStatusByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchantID),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseTransactionYearAmountFailed(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *transactionHandleApi) FindMonthMethodSuccess(c echo.Context) error {
	ctx := c.Request().Context()
	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}
	month, err := parseQueryIntWithValidation(c, "month", 1, 12)
	if err != nil {
		return errors.NewBadRequestError("invalid month")
	}

	cacheKey := fmt.Sprintf("transaction:findmonthmethodsuccess:year_%d:month_%d", year, month)
	if cached, found := gateway_cache.Get[response.ApiResponsesTransactionMonthMethod](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.client.FindMonthMethodSuccess(ctx, &pb.MonthTransactionMethod{
		Year:  int32(year),
		Month: int32(month),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseTransactionMonthMethod(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *transactionHandleApi) FindYearMethodSuccess(c echo.Context) error {
	ctx := c.Request().Context()
	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}

	cacheKey := fmt.Sprintf("transaction:findyearmethodsuccess:year_%d", year)
	if cached, found := gateway_cache.Get[response.ApiResponsesTransactionYearMethod](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.client.FindYearMethodSuccess(ctx, &pb.YearTransactionMethod{
		Year: int32(year),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseTransactionYearMethod(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *transactionHandleApi) FindMonthMethodByMerchantSuccess(c echo.Context) error {
	ctx := c.Request().Context()
	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}
	month, err := parseQueryIntWithValidation(c, "month", 1, 12)
	if err != nil {
		return errors.NewBadRequestError("invalid month")
	}
	merchantID, err := strconv.Atoi(c.Param("merchant_id"))
	if err != nil || merchantID <= 0 {
		return errors.NewBadRequestError("invalid merchant ID")
	}

	cacheKey := fmt.Sprintf("transaction:findmonthmethodbymerchantsuccess:year_%d:month_%d:merchant_%d", year, month, merchantID)
	if cached, found := gateway_cache.Get[response.ApiResponsesTransactionMonthMethod](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.client.FindMonthMethodByMerchantSuccess(ctx, &pb.MonthTransactionMethodByMerchant{
		Year:       int32(year),
		Month:      int32(month),
		MerchantId: int32(merchantID),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseTransactionMonthMethod(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *transactionHandleApi) FindYearMethodByMerchantSuccess(c echo.Context) error {
	ctx := c.Request().Context()
	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}
	merchantID, err := strconv.Atoi(c.Param("merchant_id"))
	if err != nil || merchantID <= 0 {
		return errors.NewBadRequestError("invalid merchant ID")
	}

	cacheKey := fmt.Sprintf("transaction:findyearmethodbymerchantsuccess:year_%d:merchant_%d", year, merchantID)
	if cached, found := gateway_cache.Get[response.ApiResponsesTransactionYearMethod](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.client.FindYearMethodByMerchantSuccess(ctx, &pb.YearTransactionMethodByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchantID),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseTransactionYearMethod(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *transactionHandleApi) FindMonthMethodFailed(c echo.Context) error {
	ctx := c.Request().Context()
	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}
	month, err := parseQueryIntWithValidation(c, "month", 1, 12)
	if err != nil {
		return errors.NewBadRequestError("invalid month")
	}

	cacheKey := fmt.Sprintf("transaction:findmonthmethodfailed:year_%d:month_%d", year, month)
	if cached, found := gateway_cache.Get[response.ApiResponsesTransactionMonthMethod](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.client.FindMonthMethodFailed(ctx, &pb.MonthTransactionMethod{
		Year:  int32(year),
		Month: int32(month),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseTransactionMonthMethod(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *transactionHandleApi) FindYearMethodFailed(c echo.Context) error {
	ctx := c.Request().Context()
	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}

	cacheKey := fmt.Sprintf("transaction:findyearmethodfailed:year_%d", year)
	if cached, found := gateway_cache.Get[response.ApiResponsesTransactionYearMethod](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.client.FindYearMethodFailed(ctx, &pb.YearTransactionMethod{
		Year: int32(year),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseTransactionYearMethod(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *transactionHandleApi) FindMonthMethodByMerchantFailed(c echo.Context) error {
	ctx := c.Request().Context()
	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}
	month, err := parseQueryIntWithValidation(c, "month", 1, 12)
	if err != nil {
		return errors.NewBadRequestError("invalid month")
	}
	merchantID, err := strconv.Atoi(c.Param("merchant_id"))
	if err != nil || merchantID <= 0 {
		return errors.NewBadRequestError("invalid merchant ID")
	}

	cacheKey := fmt.Sprintf("transaction:findmonthmethodbymerchantfailed:year_%d:month_%d:merchant_%d", year, month, merchantID)
	if cached, found := gateway_cache.Get[response.ApiResponsesTransactionMonthMethod](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.client.FindMonthMethodByMerchantFailed(ctx, &pb.MonthTransactionMethodByMerchant{
		Year:       int32(year),
		Month:      int32(month),
		MerchantId: int32(merchantID),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseTransactionMonthMethod(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *transactionHandleApi) FindYearMethodByMerchantFailed(c echo.Context) error {
	ctx := c.Request().Context()
	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}
	merchantID, err := strconv.Atoi(c.Param("merchant_id"))
	if err != nil || merchantID <= 0 {
		return errors.NewBadRequestError("invalid merchant ID")
	}

	cacheKey := fmt.Sprintf("transaction:findyearmethodbymerchantfailed:year_%d:merchant_%d", year, merchantID)
	if cached, found := gateway_cache.Get[response.ApiResponsesTransactionYearMethod](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.client.FindYearMethodByMerchantFailed(ctx, &pb.YearTransactionMethodByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchantID),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseTransactionYearMethod(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *transactionHandleApi) Create(c echo.Context) error {
	ctx := c.Request().Context()
	var body requests.CreateTransactionRequest
	if err := c.Bind(&body); err != nil {
		return errors.NewBadRequestError("bind failed")
	}
	if err := body.Validate(); err != nil {
		return errors.NewBadRequestError("validation failed")
	}

	grpcReq := &pb.CreateTransactionRequest{
		OrderId:       int32(body.OrderID),
		CashierId:     int32(body.CashierID),
		PaymentMethod: body.PaymentMethod,
		Amount:        int32(body.Amount),
	}

	res, err := h.client.Create(ctx, grpcReq)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseTransaction(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "transaction:*")
	return c.JSON(http.StatusOK, so)
}

func (h *transactionHandleApi) Update(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("invalid transaction ID")
	}

	var body requests.UpdateTransactionRequest
	if err := c.Bind(&body); err != nil {
		return errors.NewBadRequestError("bind failed")
	}
	if err := body.Validate(); err != nil {
		return errors.NewBadRequestError("validation failed")
	}

	grpcReq := &pb.UpdateTransactionRequest{
		TransactionId: int32(id),
		OrderId:       int32(body.OrderID),
		CashierId:     int32(body.CashierID),
		PaymentMethod: body.PaymentMethod,
		Amount:        int32(body.Amount),
	}

	res, err := h.client.Update(ctx, grpcReq)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseTransaction(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "transaction:*")
	return c.JSON(http.StatusOK, so)
}

func (h *transactionHandleApi) TrashedTransaction(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("invalid transaction ID")
	}

	req := &pb.FindByIdTransactionRequest{
		Id: int32(id),
	}

	res, err := h.client.TrashedTransaction(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseTransactionDeleteAt(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "transaction:*")
	return c.JSON(http.StatusOK, so)
}

func (h *transactionHandleApi) RestoreTransaction(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("invalid transaction ID")
	}

	req := &pb.FindByIdTransactionRequest{
		Id: int32(id),
	}

	res, err := h.client.RestoreTransaction(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseTransactionDeleteAt(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "transaction:*")
	return c.JSON(http.StatusOK, so)
}

func (h *transactionHandleApi) DeleteTransactionPermanent(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("invalid transaction ID")
	}

	req := &pb.FindByIdTransactionRequest{
		Id: int32(id),
	}

	res, err := h.client.DeleteTransactionPermanent(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseTransactionDelete(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "transaction:*")
	return c.JSON(http.StatusOK, so)
}

func (h *transactionHandleApi) RestoreAllTransaction(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := h.client.RestoreAllTransaction(ctx, &emptypb.Empty{})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseTransactionAll(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "transaction:*")
	return c.JSON(http.StatusOK, so)
}

func (h *transactionHandleApi) DeleteAllTransactionPermanent(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := h.client.DeleteAllTransactionPermanent(ctx, &emptypb.Empty{})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseTransactionAll(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "transaction:*")
	return c.JSON(http.StatusOK, so)
}
