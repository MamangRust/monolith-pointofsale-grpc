package handler

import (
	"net/http"
	"strconv"

	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/cashier_errors"
	response_api "github.com/MamangRust/monolith-point-of-sale-shared/mapper/response/api"
	"github.com/MamangRust/monolith-point-of-sale-shared/pb"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type cashierHandleApi struct {
	client  pb.CashierServiceClient
	logger  logger.LoggerInterface
	mapping response_api.CashierResponseMapper
}

func NewHandlerCashier(
	router *echo.Echo,
	client pb.CashierServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.CashierResponseMapper,
) *cashierHandleApi {
	cashierHandler := &cashierHandleApi{
		client:  client,
		logger:  logger,
		mapping: mapping,
	}

	routerCashier := router.Group("/api/cashier")

	routerCashier.GET("", cashierHandler.FindAllCashier)
	routerCashier.GET("/:id", cashierHandler.FindById)
	routerCashier.GET("/active", cashierHandler.FindByActive)
	routerCashier.GET("/trashed", cashierHandler.FindByTrashed)

	routerCashier.GET("/monthly-total-sales", cashierHandler.FindMonthlyTotalSales)
	routerCashier.GET("/yearly-total-sales", cashierHandler.FindYearTotalSales)

	routerCashier.GET("/merchant/monthly-total-sales", cashierHandler.FindMonthlyTotalSalesByMerchant)
	routerCashier.GET("/merchant/yearly-total-sales", cashierHandler.FindYearTotalSalesByMerchant)

	routerCashier.GET("/mycashier/monthly-total-sales", cashierHandler.FindMonthlyTotalSalesById)
	routerCashier.GET("/mycashier/yearly-total-sales", cashierHandler.FindYearTotalSalesById)

	routerCashier.GET("/monthly-sales", cashierHandler.FindMonthSales)
	routerCashier.GET("/yearly-sales", cashierHandler.FindYearSales)
	routerCashier.GET("/merchant/monthly-sales", cashierHandler.FindMonthSalesByMerchant)
	routerCashier.GET("/merchant/yearly-sales", cashierHandler.FindYearSalesByMerchant)
	routerCashier.GET("/mycashier/monthly-sales", cashierHandler.FindMonthSalesById)
	routerCashier.GET("/mycashier/yearly-sales", cashierHandler.FindYearSalesById)

	routerCashier.POST("/create", cashierHandler.CreateCashier)
	routerCashier.POST("/update/:id", cashierHandler.UpdateCashier)

	routerCashier.POST("/trashed/:id", cashierHandler.TrashedCashier)
	routerCashier.POST("/restore/:id", cashierHandler.RestoreCashier)
	routerCashier.DELETE("/permanent/:id", cashierHandler.DeleteCashierPermanent)

	routerCashier.POST("/restore/all", cashierHandler.RestoreAllCashier)
	routerCashier.POST("/permanent/all", cashierHandler.DeleteAllCashierPermanent)

	return cashierHandler
}

// @Security Bearer
// @Summary Find all cashiers
// @Tags Cashier
// @Description Retrieve a list of all cashiers
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationCashier "List of cashiers"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve cashier data"
// @Router /api/cashier [get]
func (h *cashierHandleApi) FindAllCashier(c echo.Context) error {
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

	req := &pb.FindAllCashierRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, req)
	if err != nil {
		h.logger.Error("Failed to fetch cashiers", zap.Error(err))
		return cashier_errors.ErrApiCashierFailedFindAll(c)
	}

	return c.JSON(http.StatusOK, h.mapping.ToApiResponsePaginationCashier(res))
}

// @Security Bearer
// @Summary Find cashier by ID
// @Tags Cashier
// @Description Retrieve a cashier by ID
// @Accept json
// @Produce json
// @Param id path int true "cashier ID"
// @Success 200 {object} response.ApiResponseCashier "cashier data"
// @Failure 400 {object} response.ErrorResponse "Invalid cashier ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve cashier data"
// @Router /api/cashier/{id} [get]
func (h *cashierHandleApi) FindById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Debug("Invalid cashier ID format", zap.Error(err))
		return cashier_errors.ErrApiCashierInvalidId(c)
	}

	ctx := c.Request().Context()
	req := &pb.FindByIdCashierRequest{Id: int32(id)}

	cashier, err := h.client.FindById(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch cashier details", zap.Error(err))
		return cashier_errors.ErrApiCashierFailedFindById(c)
	}

	return c.JSON(http.StatusOK, h.mapping.ToApiResponseCashier(cashier))
}

// @Security Bearer
// @Summary Retrieve active cashier
// @Tags Cashier
// @Description Retrieve a list of active cashier
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationCashierDeleteAt "List of active cashier"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve cashier data"
// @Router /api/cashier/active [get]
func (h *cashierHandleApi) FindByActive(c echo.Context) error {
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
	req := &pb.FindAllCashierRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, req)
	if err != nil {
		h.logger.Error("Failed to fetch active cashiers", zap.Error(err))
		return cashier_errors.ErrApiCashierFailedFindByActive(c)
	}

	return c.JSON(http.StatusOK, h.mapping.ToApiResponsePaginationCashierDeleteAt(res))
}

// @Security Bearer
// FindByTrashed retrieves a list of trashed cashier records.
// @Summary Retrieve trashed cashier
// @Tags Cashier
// @Description Retrieve a list of trashed cashier records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationCashierDeleteAt "List of trashed cashier data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve cashier data"
// @Router /api/cashier/trashed [get]
func (h *cashierHandleApi) FindByTrashed(c echo.Context) error {
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
	req := &pb.FindAllCashierRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, req)
	if err != nil {
		h.logger.Error("Failed to fetch archived cashiers", zap.Error(err))
		return cashier_errors.ErrApiCashierFailedFindByTrashed(c)
	}

	return c.JSON(http.StatusOK, h.mapping.ToApiResponsePaginationCashierDeleteAt(res))
}

// FindMonthlyTotalSales retrieves the monthly cashiers for a specific year.
// @Summary Get monthly cashiers statistics
// @Tags Cashier
// @Security Bearer
// @Description Retrieve monthly cashiers statistics for a given year
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param month query int true "Month"
// @Success 200 {object} response.ApiResponseCashierMonthSales "Successfully retrieved monthly sales data"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/cashier/monthly-total-sales [get]
func (h *cashierHandleApi) FindMonthlyTotalSales(c echo.Context) error {
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return cashier_errors.ErrApiCashierInvalidYear(c)
	}

	monthStr := c.QueryParam("month")

	month, err := strconv.Atoi(monthStr)

	if err != nil {
		h.logger.Debug("Invalid month parameter", zap.Error(err))
		return cashier_errors.ErrApiCashierInvalidMonth(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthlyTotalSales(ctx, &pb.FindYearMonthTotalSales{
		Year:  int32(year),
		Month: int32(month),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve monthly cashier sales", zap.Error(err))
		return cashier_errors.ErrApiCashierFailedMonthlyTotalSales(c)
	}

	so := h.mapping.ToApiResponseMonthlyTotalSales(res)

	return c.JSON(http.StatusOK, so)
}

// FindYearTotalSales retrieves the yearly cashiers for a specific year.
// @Summary Get yearly cashiers
// @Tags Cashier
// @Security Bearer
// @Description Retrieve the yearly cashiers for a specific year.
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponseCashierYearSales "Yearly cashiers"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly cashiers"
// @Router /api/cashier/yearly-total-sales [get]
func (h *cashierHandleApi) FindYearTotalSales(c echo.Context) error {
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return cashier_errors.ErrApiCashierInvalidYear(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearlyTotalSales(ctx, &pb.FindYearTotalSales{
		Year: int32(year),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve yearly cashier sales", zap.Error(err))
		return cashier_errors.ErrApiCashierFailedYearlyTotalSales(c)
	}

	so := h.mapping.ToApiResponseYearlyTotalSales(res)

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyTotalSalesById retrieves the monthly cashiers for a specific year.
// @Summary Get monthly cashiers statistics
// @Tags Cashier
// @Security Bearer
// @Description Retrieve monthly cashiers statistics for a given year
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param month query int true "Month"
// @Param cashier_id query int true "Cashier id"
// @Success 200 {object} response.ApiResponseCashierMonthSales "Successfully retrieved monthly sales data"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/cashier/mycashier/monthly-total-sales [get]
func (h *cashierHandleApi) FindMonthlyTotalSalesById(c echo.Context) error {
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return cashier_errors.ErrApiCashierInvalidYear(c)
	}

	monthStr := c.QueryParam("month")

	month, err := strconv.Atoi(monthStr)

	if err != nil {
		h.logger.Debug("Invalid month parameter", zap.Error(err))
		return cashier_errors.ErrApiCashierInvalidMonth(c)
	}

	cashierStr := c.QueryParam("cashier_id")

	cashier, err := strconv.Atoi(cashierStr)
	if err != nil {
		h.logger.Debug("Invalid cashier id parameter", zap.Error(err))
		return cashier_errors.ErrApiCashierInvalidId(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthlyTotalSalesById(ctx, &pb.FindYearMonthTotalSalesById{
		Year:      int32(year),
		Month:     int32(month),
		CashierId: int32(cashier),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve monthly cashier sales", zap.Error(err))
		return cashier_errors.ErrApiCashierFailedMonthlyTotalSalesById(c)
	}

	so := h.mapping.ToApiResponseMonthlyTotalSales(res)

	return c.JSON(http.StatusOK, so)
}

// FindYearTotalSalesById retrieves the yearly cashiers for a specific year.
// @Summary Get yearly cashiers
// @Tags Cashier
// @Security Bearer
// @Description Retrieve the yearly cashiers for a specific year.
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param cashier_id query int true "Cashier ID"
// @Success 200 {object} response.ApiResponseCashierYearSales "Yearly cashiers"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly cashiers"
// @Router /api/cashier/mycashier/yearly-total-sales [get]
func (h *cashierHandleApi) FindYearTotalSalesById(c echo.Context) error {
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return cashier_errors.ErrApiCashierInvalidYear(c)
	}

	cashierStr := c.QueryParam("cashier_id")

	cashier, err := strconv.Atoi(cashierStr)

	if err != nil {
		h.logger.Debug("Invalid cashier id parameter", zap.Error(err))
		return cashier_errors.ErrApiCashierInvalidId(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearlyTotalSalesById(ctx, &pb.FindYearTotalSalesById{
		Year:      int32(year),
		CashierId: int32(cashier),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve yearly cashier sales", zap.Error(err))
		return cashier_errors.ErrApiCashierFailedYearlyTotalSalesById(c)
	}

	so := h.mapping.ToApiResponseYearlyTotalSales(res)

	return c.JSON(http.StatusOK, so)
}

// FindMonthlyTotalSalesByMerchant retrieves the monthly cashiers for a specific year.
// @Summary Get monthly cashiers statistics
// @Tags Cashier
// @Security Bearer
// @Description Retrieve monthly cashiers statistics for a given year
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param month query int true "Month"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponseCashierMonthSales "Successfully retrieved monthly sales data"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/cashier/merchant/monthly-total-sales [get]
func (h *cashierHandleApi) FindMonthlyTotalSalesByMerchant(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return cashier_errors.ErrApiCashierInvalidYear(c)
	}

	monthStr := c.QueryParam("month")

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		h.logger.Debug("Invalid month parameter", zap.Error(err))
		return cashier_errors.ErrApiCashierInvalidMonth(c)
	}

	merchantStr := c.QueryParam("merchant_id")

	merchant, err := strconv.Atoi(merchantStr)
	if err != nil {
		h.logger.Debug("Invalid merchant id parameter", zap.Error(err))
		return cashier_errors.ErrApiCashierInvalidMerchantId(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthlyTotalSalesByMerchant(ctx, &pb.FindYearMonthTotalSalesByMerchant{
		Year:       int32(year),
		Month:      int32(month),
		MerchantId: int32(merchant),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve monthly cashier sales", zap.Error(err))
		return cashier_errors.ErrApiCashierFailedMonthlyTotalSalesByMerchant(c)
	}

	so := h.mapping.ToApiResponseMonthlyTotalSales(res)

	return c.JSON(http.StatusOK, so)
}

// FindYearTotalSalesByMerchant retrieves the yearly cashiers for a specific year.
// @Summary Get yearly cashiers
// @Tags Cashier
// @Security Bearer
// @Description Retrieve the yearly cashiers for a specific year.
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponseCashierYearSales "Yearly cashiers"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly cashiers"
// @Router /api/cashier/merchant/yearly-total-sales [get]
func (h *cashierHandleApi) FindYearTotalSalesByMerchant(c echo.Context) error {
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return cashier_errors.ErrApiCashierInvalidYear(c)
	}

	merchantStr := c.QueryParam("merchant_id")

	merchant, err := strconv.Atoi(merchantStr)

	if err != nil {
		h.logger.Debug("Invalid merchant id parameter", zap.Error(err))
		return cashier_errors.ErrApiCashierInvalidMerchantId(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearlyTotalSalesByMerchant(ctx, &pb.FindYearTotalSalesByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchant),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve yearly cashier sales", zap.Error(err))
		return cashier_errors.ErrApiCashierFailedYearlyTotalSalesByMerchant(c)
	}

	so := h.mapping.ToApiResponseYearlyTotalSales(res)

	return c.JSON(http.StatusOK, so)
}

// FindMonthSales retrieves the monthly cashiers for a specific year.
// @Summary Get monthly cashiers statistics
// @Tags Cashier
// @Security Bearer
// @Description Retrieve monthly cashiers statistics for a given year
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponseCashierMonthSales "Successfully retrieved monthly sales data"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/cashier/monthly-sales [get]
func (h *cashierHandleApi) FindMonthSales(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return cashier_errors.ErrApiCashierInvalidYear(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthSales(ctx, &pb.FindYearCashier{
		Year: int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly cashier sales", zap.Error(err))
		return cashier_errors.ErrApiCashierFailedMonthSales(c)
	}

	so := h.mapping.ToApiResponseCashierMonthlySale(res)

	return c.JSON(http.StatusOK, so)
}

// FindYearSales retrieves the yearly cashiers for a specific year.
// @Summary Get yearly cashiers
// @Tags Cashier
// @Security Bearer
// @Description Retrieve the yearly cashiers for a specific year.
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponseCashierYearSales "Yearly cashiers"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve yearly cashiers"
// @Router /api/cashier/yearly-sales [get]
func (h *cashierHandleApi) FindYearSales(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return cashier_errors.ErrApiCashierInvalidYear(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearSales(ctx, &pb.FindYearCashier{
		Year: int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly cashier sales", zap.Error(err))
		return cashier_errors.ErrApiCashierFailedYearSales(c)
	}

	so := h.mapping.ToApiResponseCashierYearlySale(res)

	return c.JSON(http.StatusOK, so)
}

// FindMonthSalesByMerchant retrieves monthly cashiers for a specific merchant.
// @Summary Get monthly sales by merchant
// @Tags Cashier
// @Security Bearer
// @Description Retrieve monthly cashiers statistics for a specific merchant
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponseCashierMonthSales "Successfully retrieved monthly sales by merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Merchant not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/cashier/merchant/monthly-sales [get]
func (h *cashierHandleApi) FindMonthSalesByMerchant(c echo.Context) error {
	yearStr := c.QueryParam("year")
	merchantIdStr := c.QueryParam("merchant_id")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return cashier_errors.ErrApiCashierInvalidYear(c)
	}

	merchant_id, err := strconv.Atoi(merchantIdStr)

	if err != nil {
		h.logger.Debug("Invalid merchant id parameter", zap.Error(err))
		return cashier_errors.ErrApiCashierInvalidMerchantId(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthSalesByMerchant(ctx, &pb.FindYearCashierByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchant_id),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly cashier sales", zap.Error(err))
		return cashier_errors.ErrApiCashierFailedMonthSalesByMerchant(c)
	}

	so := h.mapping.ToApiResponseCashierMonthlySale(res)

	return c.JSON(http.StatusOK, so)
}

// FindYearSalesByMerchant retrieves yearly cashier for a specific merchant.
// @Summary Get yearly sales by merchant
// @Tags Cashier
// @Security Bearer
// @Description Retrieve yearly cashier statistics for a specific merchant
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponseCashierYearSales "Successfully retrieved yearly sales by merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Merchant not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/cashier/merchant/yearly-sales [get]
func (h *cashierHandleApi) FindYearSalesByMerchant(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)

	merchantIdStr := c.QueryParam("merchant_id")

	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return cashier_errors.ErrApiCashierInvalidYear(c)
	}

	merchant_id, err := strconv.Atoi(merchantIdStr)

	if err != nil {
		h.logger.Debug("Invalid merchant id parameter", zap.Error(err))
		return cashier_errors.ErrApiCashierInvalidMerchantId(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearSalesByMerchant(ctx, &pb.FindYearCashierByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchant_id),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly cashier sales", zap.Error(err))
		return cashier_errors.ErrApiCashierFailedYearSalesByMerchant(c)
	}

	so := h.mapping.ToApiResponseCashierYearlySale(res)

	return c.JSON(http.StatusOK, so)
}

// FindMonthSalesById retrieves monthly cashier for a specific cashier.
// @Summary Get monthly sales by cashier
// @Tags Cashier
// @Security Bearer
// @Description Retrieve monthly cashier statistics for a specific cashier
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param cashier_id query int true "Cashier ID"
// @Success 200 {object} response.ApiResponseCashierMonthSales "Successfully retrieved monthly sales by cashier"
// @Failure 400 {object} response.ErrorResponse "Invalid cashier ID or year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Cashier not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/cashier/mycashier/monthly-sales [get]
func (h *cashierHandleApi) FindMonthSalesById(c echo.Context) error {
	yearStr := c.QueryParam("year")
	cashierIdStr := c.QueryParam("cashier_id")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return cashier_errors.ErrApiCashierInvalidYear(c)
	}

	cashier_id, err := strconv.Atoi(cashierIdStr)

	if err != nil {
		h.logger.Debug("Invalid cashier id parameter", zap.Error(err))
		return cashier_errors.ErrApiCashierInvalidId(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthSalesById(ctx, &pb.FindYearCashierById{
		Year:      int32(year),
		CashierId: int32(cashier_id),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve monthly cashier sales", zap.Error(err))
		return cashier_errors.ErrApiCashierFailedMonthSalesById(c)
	}

	so := h.mapping.ToApiResponseCashierMonthlySale(res)

	return c.JSON(http.StatusOK, so)
}

// FindYearSalesById retrieves yearly cashier for a specific cashier.
// @Summary Get yearly sales by cashier
// @Tags Cashier
// @Security Bearer
// @Description Retrieve yearly cashier statistics for a specific cashier
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param cashier_id query int true "Cashier ID"
// @Success 200 {object} response.ApiResponseCashierYearSales "Successfully retrieved yearly sales by cashier"
// @Failure 400 {object} response.ErrorResponse "Invalid cashier ID or year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Cashier not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/cashier/mycashier/yearly-sales [get]
func (h *cashierHandleApi) FindYearSalesById(c echo.Context) error {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)

	cashierIdStr := c.QueryParam("cashier_id")

	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return cashier_errors.ErrApiCashierInvalidYear(c)
	}

	cashier_id, err := strconv.Atoi(cashierIdStr)

	if err != nil {
		h.logger.Debug("Invalid cashier id parameter", zap.Error(err))
		return cashier_errors.ErrApiCashierInvalidId(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearSalesById(ctx, &pb.FindYearCashierById{
		Year:      int32(year),
		CashierId: int32(cashier_id),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly cashier sales", zap.Error(err))
		return cashier_errors.ErrApiCashierFailedYearSalesById(c)
	}

	so := h.mapping.ToApiResponseCashierYearlySale(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// Create handles the creation of a new cashier.
// @Summary Create a new cashier
// @Tags Cashier
// @Description Create a new cashier with the provided details
// @Accept json
// @Produce json
// @Param request body requests.CreateCashierRequest true "Create cashier request"
// @Success 200 {object} response.ApiResponseCashier "Successfully created cashier"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create cashier"
// @Router /api/cashier/create [post]
func (h *cashierHandleApi) CreateCashier(c echo.Context) error {
	var body requests.CreateCashierRequest

	if err := c.Bind(&body); err != nil {
		h.logger.Debug("Invalid request format", zap.Error(err))
		return cashier_errors.ErrApiBindCreateCashier(c)
	}

	if err := body.Validate(); err != nil {
		h.logger.Debug("Validation failed", zap.Error(err))
		return cashier_errors.ErrApiValidateCreateCashier(c)
	}

	ctx := c.Request().Context()
	req := &pb.CreateCashierRequest{
		MerchantId: int32(body.MerchantID),
		UserId:     int32(body.UserID),
		Name:       body.Name,
	}

	res, err := h.client.CreateCashier(ctx, req)
	if err != nil {
		h.logger.Error("Cashier creation failed", zap.Error(err))
		return cashier_errors.ErrApiCashierFailedCreate(c)
	}

	so := h.mapping.ToApiResponseCashier(res)

	return c.JSON(http.StatusCreated, so)
}

// @Security Bearer
// Update handles the update of an existing cashier record.
// @Summary Update an existing cashier
// @Tags Cashier
// @Description Update an existing cashier record with the provided details
// @Accept json
// @Produce json
// @Param id path int true "Cashier ID"
// @Param UpdateCashierRequest body requests.UpdateCashierRequest true "Update cashier request"
// @Success 200 {object} response.ApiResponseCashier "Successfully updated cashier"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update cashier"
// @Router /api/cashier/update/{id} [post]
func (h *cashierHandleApi) UpdateCashier(c echo.Context) error {
	id := c.Param("id")

	idStr, err := strconv.Atoi(id)

	if err != nil {
		h.logger.Debug("Invalid id parameter", zap.Error(err))

		return cashier_errors.ErrApiCashierInvalidId(c)
	}

	var body requests.UpdateCashierRequest

	if err := c.Bind(&body); err != nil {
		h.logger.Debug("Invalid request format", zap.Error(err))
		return cashier_errors.ErrApiBindUpdateCashier(c)
	}

	if err := body.Validate(); err != nil {
		h.logger.Debug("Validation failed", zap.Error(err))
		return cashier_errors.ErrApiValidateUpdateCashier(c)
	}

	ctx := c.Request().Context()
	req := &pb.UpdateCashierRequest{
		CashierId: int32(idStr),
		Name:      body.Name,
	}

	res, err := h.client.UpdateCashier(ctx, req)
	if err != nil {
		h.logger.Error("Cashier update failed", zap.Error(err))
		return cashier_errors.ErrApiCashierFailedUpdate(c)
	}

	so := h.mapping.ToApiResponseCashier(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// TrashedCasher retrieves a trashed casher record by its ID.
// @Summary Retrieve a trashed casher
// @Tags Cashier
// @Description Retrieve a trashed casher record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Cashier ID"
// @Success 200 {object} response.ApiResponseCashierDeleteAt "Successfully retrieved trashed cashier"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve trashed cashier"
// @Router /api/cashier/trashed/{id} [get]
func (h *cashierHandleApi) TrashedCashier(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Debug("Invalid cashier ID format", zap.Error(err))
		return cashier_errors.ErrApiCashierInvalidId(c)
	}

	ctx := c.Request().Context()
	req := &pb.FindByIdCashierRequest{Id: int32(id)}

	cashier, err := h.client.TrashedCashier(ctx, req)
	if err != nil {
		h.logger.Error("Failed to archive cashier", zap.Error(err))
		return cashier_errors.ErrApiCashierFailedTrashed(c)
	}

	so := h.mapping.ToApiResponseCashierDeleteAt(cashier)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreCashier restores a cashier record from the trash by its ID.
// @Summary Restore a trashed cashier
// @Tags Cashier
// @Description Restore a trashed cashier record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Cashier ID"
// @Success 200 {object} response.ApiResponseCashierDeleteAt "Successfully restored cashier"
// @Failure 400 {object} response.ErrorResponse "Invalid cashier ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore cashier"
// @Router /api/cashier/restore/{id} [post]
func (h *cashierHandleApi) RestoreCashier(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Debug("Invalid cashier ID format", zap.Error(err))
		return cashier_errors.ErrApiCashierInvalidId(c)
	}

	ctx := c.Request().Context()
	req := &pb.FindByIdCashierRequest{Id: int32(id)}

	cashier, err := h.client.RestoreCashier(ctx, req)
	if err != nil {
		h.logger.Error("Failed to restore cashier", zap.Error(err))
		return cashier_errors.ErrApiCashierFailedRestore(c)
	}

	so := h.mapping.ToApiResponseCashierDeleteAt(cashier)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteCashierPermanent permanently deletes a cashier record by its ID.
// @Summary Permanently delete a cashier
// @Tags Cashier
// @Description Permanently delete a cashier record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "cashier ID"
// @Success 200 {object} response.ApiResponseCashierDelete "Successfully deleted cashier record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete cashier:"
// @Router /api/cashier/delete/{id} [delete]
func (h *cashierHandleApi) DeleteCashierPermanent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Debug("Invalid cashier ID format", zap.Error(err))
		return cashier_errors.ErrApiCashierInvalidId(c)
	}

	ctx := c.Request().Context()
	req := &pb.FindByIdCashierRequest{Id: int32(id)}

	cashier, err := h.client.DeleteCashierPermanent(ctx, req)
	if err != nil {
		h.logger.Error("Failed to delete cashier", zap.Error(err))
		return cashier_errors.ErrApiCashierFailedDeletePermanent(c)
	}

	so := h.mapping.ToApiResponseCashierDelete(cashier)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreCashier restores a cashier record from the trash by its ID.
// @Summary Restore a trashed cashier
// @Tags Cashier
// @Description Restore a trashed cashier record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Cashier ID"
// @Success 200 {object} response.ApiResponseCashierAll "Successfully restored cashier all"
// @Failure 400 {object} response.ErrorResponse "Invalid cashier ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore cashier"
// @Router /api/cashier/restore/all [post]
func (h *cashierHandleApi) RestoreAllCashier(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.RestoreAllCashier(ctx, &emptypb.Empty{})
	if err != nil {
		h.logger.Error("Bulk cashier restoration failed", zap.Error(err))
		return cashier_errors.ErrApiCashierFailedRestoreAll(c)
	}

	h.logger.Info("All cashier accounts restored successfully")

	so := h.mapping.ToApiResponseCashierAll(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteCashierPermanent permanently deletes a cashier record by its ID.
// @Summary Permanently delete a cashier
// @Tags Cashier
// @Description Permanently delete a cashier record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "cashier ID"
// @Success 200 {object} response.ApiResponseCashierAll "Successfully deleted cashier record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete cashier:"
// @Router /api/cashier/delete/all [post]
func (h *cashierHandleApi) DeleteAllCashierPermanent(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := h.client.DeleteAllCashierPermanent(ctx, &emptypb.Empty{})
	if err != nil {
		h.logger.Error("Bulk cashier deletion failed", zap.Error(err))
		return cashier_errors.ErrApiCashierFailedDeleteAllPermanent(c)
	}

	h.logger.Info("All cashier accounts permanently deleted")

	so := h.mapping.ToApiResponseCashierAll(res)

	return c.JSON(http.StatusOK, so)
}
