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

type cashierHandleApi struct {
	client     pb.CashierServiceClient
	logger     logger.LoggerInterface
	mapping    response_api.CashierResponseMapper
	apiHandler errors.ApiHandler
	cache      *gateway_cache.GatewayCache
}

func NewHandlerCashier(
	router *echo.Echo,
	client pb.CashierServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.CashierResponseMapper,
	apiHandler errors.ApiHandler,
	cache *gateway_cache.GatewayCache,
) *cashierHandleApi {
	cashierHandler := &cashierHandleApi{
		client:     client,
		logger:     logger,
		mapping:    mapping,
		apiHandler: apiHandler,
		cache:      cache,
	}

	routerCashier := router.Group("/api/cashier")

	routerCashier.GET("", apiHandler.Handle("get-cashier-findallcashier", cashierHandler.FindAllCashier))
	routerCashier.GET("/:id", apiHandler.Handle("get-cashier-findbyid", cashierHandler.FindById))
	routerCashier.GET("/active", apiHandler.Handle("get-cashier-findbyactive", cashierHandler.FindByActive))
	routerCashier.GET("/trashed", apiHandler.Handle("get-cashier-findbytrashed", cashierHandler.FindByTrashed))

	routerCashier.GET("/monthly-total-sales", apiHandler.Handle("get-cashier-findmonthlytotalsales", cashierHandler.FindMonthlyTotalSales))
	routerCashier.GET("/yearly-total-sales", apiHandler.Handle("get-cashier-findyeartotalsales", cashierHandler.FindYearTotalSales))

	routerCashier.GET("/merchant/monthly-total-sales", apiHandler.Handle("get-cashier-findmonthlytotalsalesbymerchant", cashierHandler.FindMonthlyTotalSalesByMerchant))
	routerCashier.GET("/merchant/yearly-total-sales", apiHandler.Handle("get-cashier-findyeartotalsalesbymerchant", cashierHandler.FindYearTotalSalesByMerchant))

	routerCashier.GET("/mycashier/monthly-total-sales", apiHandler.Handle("get-cashier-findmonthlytotalsalesbyid", cashierHandler.FindMonthlyTotalSalesById))
	routerCashier.GET("/mycashier/yearly-total-sales", apiHandler.Handle("get-cashier-findyeartotalsalesbyid", cashierHandler.FindYearTotalSalesById))

	routerCashier.GET("/monthly-sales", apiHandler.Handle("get-cashier-findmonthsales", cashierHandler.FindMonthSales))
	routerCashier.GET("/yearly-sales", apiHandler.Handle("get-cashier-findyearsales", cashierHandler.FindYearSales))
	routerCashier.GET("/merchant/monthly-sales", apiHandler.Handle("get-cashier-findmonthsalesbymerchant", cashierHandler.FindMonthSalesByMerchant))
	routerCashier.GET("/merchant/yearly-sales", apiHandler.Handle("get-cashier-findyearsalesbymerchant", cashierHandler.FindYearSalesByMerchant))
	routerCashier.GET("/mycashier/monthly-sales", apiHandler.Handle("get-cashier-findmonthsalesbyid", cashierHandler.FindMonthSalesById))
	routerCashier.GET("/mycashier/yearly-sales", apiHandler.Handle("get-cashier-findyearsalesbyid", cashierHandler.FindYearSalesById))

	routerCashier.POST("/create", apiHandler.Handle("post-cashier-createcashier", cashierHandler.CreateCashier))
	routerCashier.POST("/update/:id", apiHandler.Handle("post-cashier-updatecashier", cashierHandler.UpdateCashier))

	routerCashier.POST("/trashed/:id", apiHandler.Handle("post-cashier-trashedcashier", cashierHandler.TrashedCashier))
	routerCashier.POST("/restore/:id", apiHandler.Handle("post-cashier-restorecashier", cashierHandler.RestoreCashier))
	routerCashier.DELETE("/permanent/:id", apiHandler.Handle("delete-cashier-deletecashierpermanent", cashierHandler.DeleteCashierPermanent))

	routerCashier.POST("/restore/all", apiHandler.Handle("post-cashier-restoreallcashier", cashierHandler.RestoreAllCashier))
	routerCashier.POST("/permanent/all", apiHandler.Handle("post-cashier-deleteallcashierpermanent", cashierHandler.DeleteAllCashierPermanent))

	return cashierHandler
}

func (h *cashierHandleApi) FindAllCashier(c echo.Context) error {
	ctx := c.Request().Context()
	page := parseQueryInt(c, "page", 1)
	pageSize := parseQueryInt(c, "page_size", 10)
	search := c.QueryParam("search")

	cacheKey := fmt.Sprintf("cashier:findallcashier:page_%d:size_%d:search_%s", page, pageSize, search)
	if cached, found := gateway_cache.Get[response.ApiResponsePaginationCashier](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	req := &pb.FindAllCashierRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponsePaginationCashier(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *cashierHandleApi) FindById(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("invalid cashier ID")
	}

	cacheKey := fmt.Sprintf("cashier:findbyid:id_%d", id)
	if cached, found := gateway_cache.Get[response.ApiResponseCashier](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	req := &pb.FindByIdCashierRequest{Id: int32(id)}
	res, err := h.client.FindById(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseCashier(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *cashierHandleApi) FindByActive(c echo.Context) error {
	ctx := c.Request().Context()
	page := parseQueryInt(c, "page", 1)
	pageSize := parseQueryInt(c, "page_size", 10)
	search := c.QueryParam("search")

	cacheKey := fmt.Sprintf("cashier:findbyactive:page_%d:size_%d:search_%s", page, pageSize, search)
	if cached, found := gateway_cache.Get[response.ApiResponsePaginationCashierDeleteAt](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	req := &pb.FindAllCashierRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponsePaginationCashierDeleteAt(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *cashierHandleApi) FindByTrashed(c echo.Context) error {
	ctx := c.Request().Context()
	page := parseQueryInt(c, "page", 1)
	pageSize := parseQueryInt(c, "page_size", 10)
	search := c.QueryParam("search")

	cacheKey := fmt.Sprintf("cashier:findbytrashed:page_%d:size_%d:search_%s", page, pageSize, search)
	if cached, found := gateway_cache.Get[response.ApiResponsePaginationCashierDeleteAt](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	req := &pb.FindAllCashierRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponsePaginationCashierDeleteAt(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *cashierHandleApi) FindMonthlyTotalSales(c echo.Context) error {
	ctx := c.Request().Context()
	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}
	month, err := parseQueryIntWithValidation(c, "month", 1, 12)
	if err != nil {
		return errors.NewBadRequestError("invalid month")
	}

	cacheKey := fmt.Sprintf("cashier:findmonthlytotalsales:year_%d:month_%d", year, month)
	if cached, found := gateway_cache.Get[response.ApiResponseCashierMonthlyTotalSales](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.client.FindMonthlyTotalSales(ctx, &pb.FindYearMonthTotalSales{
		Year:  int32(year),
		Month: int32(month),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseMonthlyTotalSales(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *cashierHandleApi) FindYearTotalSales(c echo.Context) error {
	ctx := c.Request().Context()
	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}

	cacheKey := fmt.Sprintf("cashier:findyeartotalsales:year_%d", year)
	if cached, found := gateway_cache.Get[response.ApiResponseCashierYearlyTotalSales](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.client.FindYearlyTotalSales(ctx, &pb.FindYearTotalSales{
		Year: int32(year),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseYearlyTotalSales(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *cashierHandleApi) FindMonthlyTotalSalesById(c echo.Context) error {
	ctx := c.Request().Context()
	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}
	month, err := parseQueryIntWithValidation(c, "month", 1, 12)
	if err != nil {
		return errors.NewBadRequestError("invalid month")
	}
	cashierID, err := parseQueryIntWithValidation(c, "cashier_id", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid cashier_id")
	}

	cacheKey := fmt.Sprintf("cashier:findmonthlytotalsalesbyid:year_%d:month_%d:cashierID_%d", year, month, cashierID)
	if cached, found := gateway_cache.Get[response.ApiResponseCashierMonthlyTotalSales](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.client.FindMonthlyTotalSalesById(ctx, &pb.FindYearMonthTotalSalesById{
		Year:      int32(year),
		Month:     int32(month),
		CashierId: int32(cashierID),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseMonthlyTotalSales(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *cashierHandleApi) FindYearTotalSalesById(c echo.Context) error {
	ctx := c.Request().Context()
	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}
	cashierID, err := parseQueryIntWithValidation(c, "cashier_id", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid cashierID")
	}

	cacheKey := fmt.Sprintf("cashier:findyeartotalsalesbyid:year_%d:cashierID_%d", year, cashierID)
	if cached, found := gateway_cache.Get[response.ApiResponseCashierYearlyTotalSales](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.client.FindYearlyTotalSalesById(ctx, &pb.FindYearTotalSalesById{
		Year:      int32(year),
		CashierId: int32(cashierID),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseYearlyTotalSales(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *cashierHandleApi) FindMonthlyTotalSalesByMerchant(c echo.Context) error {
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
		return errors.NewBadRequestError("invalid merchantID")
	}

	cacheKey := fmt.Sprintf("cashier:findmonthlytotalsalesbymerchant:year_%d:month_%d:merchantID_%d", year, month, merchantID)
	if cached, found := gateway_cache.Get[response.ApiResponseCashierMonthlyTotalSales](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.client.FindMonthlyTotalSalesByMerchant(ctx, &pb.FindYearMonthTotalSalesByMerchant{
		Year:       int32(year),
		Month:      int32(month),
		MerchantId: int32(merchantID),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseMonthlyTotalSales(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *cashierHandleApi) FindYearTotalSalesByMerchant(c echo.Context) error {
	ctx := c.Request().Context()
	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}
	merchantID, err := parseQueryIntWithValidation(c, "merchant_id", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid merchantID")
	}

	cacheKey := fmt.Sprintf("cashier:findyeartotalsalesbymerchant:year_%d:merchantID_%d", year, merchantID)
	if cached, found := gateway_cache.Get[response.ApiResponseCashierYearlyTotalSales](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.client.FindYearlyTotalSalesByMerchant(ctx, &pb.FindYearTotalSalesByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchantID),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseYearlyTotalSales(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *cashierHandleApi) FindMonthSales(c echo.Context) error {
	ctx := c.Request().Context()
	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}

	cacheKey := fmt.Sprintf("cashier:findmonthsales:year_%d", year)
	if cached, found := gateway_cache.Get[response.ApiResponseCashierMonthSales](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.client.FindMonthSales(ctx, &pb.FindYearCashier{
		Year: int32(year),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseCashierMonthlySale(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *cashierHandleApi) FindYearSales(c echo.Context) error {
	ctx := c.Request().Context()
	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}

	cacheKey := fmt.Sprintf("cashier:findyearsales:year_%d", year)
	if cached, found := gateway_cache.Get[response.ApiResponseCashierYearSales](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.client.FindYearSales(ctx, &pb.FindYearCashier{
		Year: int32(year),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseCashierYearlySale(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *cashierHandleApi) FindMonthSalesByMerchant(c echo.Context) error {
	ctx := c.Request().Context()
	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}
	merchant, err := parseQueryIntWithValidation(c, "merchant_id", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid merchant_id")
	}

	cacheKey := fmt.Sprintf("cashier:findmonthsalesbymerchant:year_%d:merchantID_%d", year, merchant)
	if cached, found := gateway_cache.Get[response.ApiResponseCashierMonthSales](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.client.FindMonthSalesByMerchant(ctx, &pb.FindYearCashierByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchant),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseCashierMonthlySale(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *cashierHandleApi) FindYearSalesByMerchant(c echo.Context) error {
	ctx := c.Request().Context()
	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}
	merchant, err := parseQueryIntWithValidation(c, "merchant_id", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid merchant_id")
	}

	cacheKey := fmt.Sprintf("cashier:findyearsalesbymerchant:year_%d:merchantID_%d", year, merchant)
	if cached, found := gateway_cache.Get[response.ApiResponseCashierYearSales](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.client.FindYearSalesByMerchant(ctx, &pb.FindYearCashierByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchant),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseCashierYearlySale(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *cashierHandleApi) FindMonthSalesById(c echo.Context) error {
	ctx := c.Request().Context()
	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}
	cashier, err := parseQueryIntWithValidation(c, "cashier_id", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid cashier_id")
	}

	cacheKey := fmt.Sprintf("cashier:findmonthsalesbyid:year_%d:cashierID_%d", year, cashier)
	if cached, found := gateway_cache.Get[response.ApiResponseCashierMonthSales](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.client.FindMonthSalesById(ctx, &pb.FindYearCashierById{
		Year:      int32(year),
		CashierId: int32(cashier),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseCashierMonthlySale(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *cashierHandleApi) FindYearSalesById(c echo.Context) error {
	ctx := c.Request().Context()
	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}
	cashier_id, err := parseQueryIntWithValidation(c, "cashier_id", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid cashier_id")
	}

	cacheKey := fmt.Sprintf("cashier:findyearsalesbyid:year_%d:cashierID_%d", year, cashier_id)
	if cached, found := gateway_cache.Get[response.ApiResponseCashierYearSales](ctx, h.cache, cacheKey); found && cached != nil {
		return c.JSON(http.StatusOK, cached)
	}

	res, err := h.client.FindYearSalesById(ctx, &pb.FindYearCashierById{
		Year:      int32(year),
		CashierId: int32(cashier_id),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseCashierYearlySale(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return c.JSON(http.StatusOK, so)
}

func (h *cashierHandleApi) CreateCashier(c echo.Context) error {
	ctx := c.Request().Context()
	var body requests.CreateCashierRequest
	if err := c.Bind(&body); err != nil {
		return errors.NewBadRequestError("bind failed")
	}
	if err := body.Validate(); err != nil {
		return errors.NewBadRequestError("validation failed")
	}

	req := &pb.CreateCashierRequest{
		MerchantId: int32(body.MerchantID),
		UserId:     int32(body.UserID),
		Name:       body.Name,
	}

	res, err := h.client.CreateCashier(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseCashier(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "cashier:*")
	return c.JSON(http.StatusCreated, so)
}

func (h *cashierHandleApi) UpdateCashier(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	idStr, err := strconv.Atoi(id)
	if err != nil {
		return errors.NewBadRequestError("invalid cashier ID")
	}

	var body requests.UpdateCashierRequest
	if err := c.Bind(&body); err != nil {
		return errors.NewBadRequestError("bind failed")
	}
	if err := body.Validate(); err != nil {
		return errors.NewBadRequestError("validation failed")
	}

	req := &pb.UpdateCashierRequest{
		CashierId: int32(idStr),
		Name:      body.Name,
	}

	res, err := h.client.UpdateCashier(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseCashier(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "cashier:*")
	return c.JSON(http.StatusOK, so)
}

func (h *cashierHandleApi) TrashedCashier(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("invalid cashier ID")
	}

	req := &pb.FindByIdCashierRequest{Id: int32(id)}
	res, err := h.client.TrashedCashier(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseCashierDeleteAt(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "cashier:*")
	return c.JSON(http.StatusOK, so)
}

func (h *cashierHandleApi) RestoreCashier(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("invalid cashier ID")
	}

	req := &pb.FindByIdCashierRequest{Id: int32(id)}
	res, err := h.client.RestoreCashier(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseCashierDeleteAt(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "cashier:*")
	return c.JSON(http.StatusOK, so)
}

func (h *cashierHandleApi) DeleteCashierPermanent(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("invalid cashier ID")
	}

	req := &pb.FindByIdCashierRequest{Id: int32(id)}
	res, err := h.client.DeleteCashierPermanent(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseCashierDelete(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "cashier:*")
	return c.JSON(http.StatusOK, so)
}

func (h *cashierHandleApi) RestoreAllCashier(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := h.client.RestoreAllCashier(ctx, &emptypb.Empty{})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseCashierAll(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "cashier:*")
	return c.JSON(http.StatusOK, so)
}

func (h *cashierHandleApi) DeleteAllCashierPermanent(c echo.Context) error {
	ctx := c.Request().Context()
	res, err := h.client.DeleteAllCashierPermanent(ctx, &emptypb.Empty{})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseCashierAll(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "cashier:*")
	return c.JSON(http.StatusOK, so)
}
