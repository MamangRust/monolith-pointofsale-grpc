package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/cashier_errors"
	response_api "github.com/MamangRust/monolith-point-of-sale-shared/mapper/response/api"
	"github.com/MamangRust/monolith-point-of-sale-shared/pb"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelcode "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type cashierHandleApi struct {
	client          pb.CashierServiceClient
	logger          logger.LoggerInterface
	mapping         response_api.CashierResponseMapper
	trace           trace.Tracer
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewHandlerCashier(
	router *echo.Echo,
	client pb.CashierServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.CashierResponseMapper,
) *cashierHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cashier_handler_requests_total",
			Help: "Total number of cashier requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cashier_handler_request_duration_seconds",
			Help:    "Duration of cashier requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter)

	cashierHandler := &cashierHandleApi{
		client:          client,
		logger:          logger,
		mapping:         mapping,
		trace:           otel.Tracer("cashier-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
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
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindAllCashier"
	)

	page := parseQueryInt(c, "page", defaultPage)
	pageSize := parseQueryInt(c, "page_size", defaultPageSize)

	search := c.QueryParam("search")

	ctx := c.Request().Context()
	end, logSuccess, logError := h.startTracingAndLogging(
		ctx,
		method,
		attribute.Int("page", page),
		attribute.Int("page_size", pageSize),
		attribute.String("search", search),
	)

	defer func() {
		end()
	}()

	req := &pb.FindAllCashierRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, req)
	if err != nil {
		logError("Failed to retrieve cashier data", err)
		return cashier_errors.ErrApiCashierFailedFindAll(c)
	}

	so := h.mapping.ToApiResponsePaginationCashier(res)

	logSuccess("Successfully retrieved cashier data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
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
	const method = "FindById"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logError("Invalid cashier ID", err, zap.Error(err))

		return cashier_errors.ErrApiCashierInvalidId(c)
	}

	req := &pb.FindByIdCashierRequest{Id: int32(id)}

	cashier, err := h.client.FindById(ctx, req)

	if err != nil {
		logError("Failed to retrieve cashier data", err, zap.Error(err))

		return cashier_errors.ErrApiCashierFailedFindById(c)
	}

	so := h.mapping.ToApiResponseCashier(cashier)

	logSuccess("Successfully retrieve cashier data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
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
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindAllCashier"
	)

	page := parseQueryInt(c, "page", defaultPage)
	pageSize := parseQueryInt(c, "page_size", defaultPageSize)

	search := c.QueryParam("search")

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(
		ctx,
		method,
		attribute.Int("page", page),
		attribute.Int("page_size", pageSize),
		attribute.String("search", search),
	)

	defer func() { end() }()

	req := &pb.FindAllCashierRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, req)

	if err != nil {
		logError("Failed to retrieve cashier data", err, zap.Error(err))

		return cashier_errors.ErrApiCashierFailedFindByActive(c)
	}

	so := h.mapping.ToApiResponsePaginationCashierDeleteAt(res)

	logSuccess("Successfully retrieve cashier data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
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
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindAllCashier"
	)

	page := parseQueryInt(c, "page", defaultPage)
	pageSize := parseQueryInt(c, "page_size", defaultPageSize)

	search := c.QueryParam("search")

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(
		ctx,
		method,
		attribute.Int("page", page),
		attribute.Int("page_size", pageSize),
		attribute.String("search", search),
	)

	defer func() { end() }()

	req := &pb.FindAllCashierRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, req)

	if err != nil {
		logError("Failed to retrieve cashier data", err, zap.Error(err))

		return cashier_errors.ErrApiCashierFailedFindByTrashed(c)
	}

	so := h.mapping.ToApiResponsePaginationCashierDeleteAt(res)

	logSuccess("Successfully retrieve cashier data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
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
	const method = "FindMonthlyTotalSales"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		logError("Invalid year parameter", err, zap.String("year", c.QueryParam("year")))
		return cashier_errors.ErrApiCashierInvalidYear(c)
	}

	month, err := parseQueryIntWithValidation(c, "month", 1, 12)
	if err != nil {
		logError("Invalid month parameter", err, zap.String("month", c.QueryParam("month")))
		return cashier_errors.ErrApiCashierInvalidMonth(c)
	}

	res, err := h.client.FindMonthlyTotalSales(ctx, &pb.FindYearMonthTotalSales{
		Year:  int32(year),
		Month: int32(month),
	})

	if err != nil {
		logError("Failed to retrieve monthly sales data", err, zap.Error(err))

		return cashier_errors.ErrApiCashierFailedMonthlyTotalSales(c)
	}

	so := h.mapping.ToApiResponseMonthlyTotalSales(res)

	logSuccess("Successfully retrieve monthly sales data", zap.Bool("success", true))

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
	const method = "FindYearTotalSales"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		logError("Invalid year parameter", err, zap.String("year", c.QueryParam("year")))
		return cashier_errors.ErrApiCashierInvalidYear(c)
	}

	res, err := h.client.FindYearlyTotalSales(ctx, &pb.FindYearTotalSales{
		Year: int32(year),
	})

	if err != nil {
		logError("Failed to retrieve yearly cashier sales", err, zap.Error(err))

		return cashier_errors.ErrApiCashierFailedYearlyTotalSales(c)
	}

	so := h.mapping.ToApiResponseYearlyTotalSales(res)

	logSuccess("Successfully retrieve yearly cashier sales", zap.Bool("success", true))

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
	const method = "FindMonthlyTotalSalesById"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		logError("Invalid year parameter", err, zap.String("year", c.QueryParam("year")))
		return cashier_errors.ErrApiCashierInvalidYear(c)
	}

	month, err := parseQueryIntWithValidation(c, "month", 1, 12)
	if err != nil {
		logError("Invalid month parameter", err, zap.String("month", c.QueryParam("month")))
		return cashier_errors.ErrApiCashierInvalidMonth(c)
	}

	cashier, err := parseQueryIntWithValidation(c, "cashier_id", 1, 9999)

	if err != nil {
		logError("Invalid cashier_id parameter", err, zap.String("cashier_id", c.QueryParam("cashier_id")))
		return cashier_errors.ErrApiCashierInvalidId(c)
	}

	res, err := h.client.FindMonthlyTotalSalesById(ctx, &pb.FindYearMonthTotalSalesById{
		Year:      int32(year),
		Month:     int32(month),
		CashierId: int32(cashier),
	})

	if err != nil {
		logError("Failed to retrieve monthly sales data", err, zap.Error(err))

		return cashier_errors.ErrApiCashierFailedMonthlyTotalSalesById(c)
	}

	so := h.mapping.ToApiResponseMonthlyTotalSales(res)

	logSuccess("Successfully retrieve monthly sales data", zap.Bool("success", true))

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
	const method = "FindYearTotalSalesById"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		logError("Invalid year parameter", err, zap.String("year", c.QueryParam("year")))
		return cashier_errors.ErrApiCashierInvalidYear(c)
	}

	cashier, err := parseQueryIntWithValidation(c, "cashier_id", 1, 9999)

	if err != nil {
		logError("Invalid cashier_id parameter", err, zap.String("cashier_id", c.QueryParam("cashier_id")))
		return cashier_errors.ErrApiCashierInvalidId(c)
	}

	res, err := h.client.FindYearlyTotalSalesById(ctx, &pb.FindYearTotalSalesById{
		Year:      int32(year),
		CashierId: int32(cashier),
	})

	if err != nil {
		logError("Failed to retrieve yearly sales data", err, zap.Error(err))

		return cashier_errors.ErrApiCashierFailedYearlyTotalSalesById(c)
	}

	so := h.mapping.ToApiResponseYearlyTotalSales(res)

	logSuccess("Successfully retrieve yearly sales data", zap.Bool("success", true))

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
	const method = "FindMonthlyTotalSalesByMerchant"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)

	if err != nil {
		logError("Invalid year parameter", err, zap.String("year", c.QueryParam("year")))
		return cashier_errors.ErrApiCashierInvalidYear(c)
	}

	month, err := parseQueryIntWithValidation(c, "month", 1, 12)

	if err != nil {
		logError("Invalid month parameter", err, zap.String("month", c.QueryParam("month")))
		return cashier_errors.ErrApiCashierInvalidMonth(c)
	}

	merchant, err := parseQueryIntWithValidation(c, "merchant_id", 1, 9999)

	if err != nil {
		logError("Invalid merchant_id parameter", err, zap.String("merchant_id", c.QueryParam("merchant_id")))
		return cashier_errors.ErrApiCashierInvalidMerchantId(c)
	}

	res, err := h.client.FindMonthlyTotalSalesByMerchant(ctx, &pb.FindYearMonthTotalSalesByMerchant{
		Year:       int32(year),
		Month:      int32(month),
		MerchantId: int32(merchant),
	})

	if err != nil {
		logError("Failed to retrieve monthly sales data", err, zap.Error(err))

		return cashier_errors.ErrApiCashierFailedMonthlyTotalSalesByMerchant(c)
	}

	so := h.mapping.ToApiResponseMonthlyTotalSales(res)

	logSuccess("Successfully retrieve monthly sales data", zap.Bool("success", true))

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
	const method = "FindYearTotalSalesByMerchant"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)

	if err != nil {
		logError("Invalid year parameter", err, zap.String("year", c.QueryParam("year")))
		return cashier_errors.ErrApiCashierInvalidYear(c)
	}

	merchant, err := parseQueryIntWithValidation(c, "merchant_id", 1, 9999)

	if err != nil {
		logError("Invalid merchant_id parameter", err, zap.String("merchant_id", c.QueryParam("merchant_id")))
		return cashier_errors.ErrApiCashierInvalidMerchantId(c)
	}

	res, err := h.client.FindYearlyTotalSalesByMerchant(ctx, &pb.FindYearTotalSalesByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchant),
	})

	if err != nil {
		logError("Failed to retrieve yearly sales data", err, zap.Error(err))

		return cashier_errors.ErrApiCashierFailedYearlyTotalSalesByMerchant(c)
	}

	so := h.mapping.ToApiResponseYearlyTotalSales(res)

	logSuccess("Successfully retrieve yearly sales data", zap.Bool("success", true))

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
	const method = "FindMonthSales"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)

	if err != nil {
		logError("Invalid year parameter", err, zap.String("year", c.QueryParam("year")))
		return cashier_errors.ErrApiCashierInvalidYear(c)
	}

	res, err := h.client.FindMonthSales(ctx, &pb.FindYearCashier{
		Year: int32(year),
	})

	if err != nil {
		logError("Failed to retrieve monthly sales data", err, zap.Error(err))

		return cashier_errors.ErrApiCashierFailedMonthSales(c)
	}

	so := h.mapping.ToApiResponseCashierMonthlySale(res)

	logSuccess("Successfully retrieve monthly sales data", zap.Bool("success", true))

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
	const method = "FindYearSales"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)

	if err != nil {
		logError("Invalid year parameter", err, zap.String("year", c.QueryParam("year")))
		return cashier_errors.ErrApiCashierInvalidYear(c)
	}

	res, err := h.client.FindYearSales(ctx, &pb.FindYearCashier{
		Year: int32(year),
	})
	if err != nil {
		logError("Failed to retrieve yearly sales data", err, zap.Error(err))

		return cashier_errors.ErrApiCashierFailedYearSales(c)
	}

	so := h.mapping.ToApiResponseCashierYearlySale(res)

	logSuccess("Successfully retrieve yearly sales data", zap.Bool("success", true))

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
	const method = "FindById"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		logError("Invalid year parameter", err, zap.String("year", c.QueryParam("year")))
		return cashier_errors.ErrApiCashierInvalidYear(c)
	}

	merchant, err := parseQueryIntWithValidation(c, "merchant_id", 1, 9999)

	if err != nil {
		logError("Invalid merchant_id parameter", err, zap.String("merchant_id", c.QueryParam("merchant_id")))
		return cashier_errors.ErrApiCashierInvalidMerchantId(c)
	}

	res, err := h.client.FindMonthSalesByMerchant(ctx, &pb.FindYearCashierByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchant),
	})
	if err != nil {
		logError("Failed to retrieve monthly sales data", err, zap.Error(err))

		return cashier_errors.ErrApiCashierFailedMonthSalesByMerchant(c)
	}

	so := h.mapping.ToApiResponseCashierMonthlySale(res)

	logSuccess("Successfully retrieve monthly sales data", zap.Bool("success", true))

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
	const method = "FindYearSalesByMerchant"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		logError("Invalid year parameter", err, zap.String("year", c.QueryParam("year")))
		return cashier_errors.ErrApiCashierInvalidYear(c)
	}

	merchant, err := parseQueryIntWithValidation(c, "merchant_id", 1, 9999)

	if err != nil {
		logError("Invalid merchant_id parameter", err, zap.String("merchant_id", c.QueryParam("merchant_id")))
		return cashier_errors.ErrApiCashierInvalidMerchantId(c)
	}

	res, err := h.client.FindYearSalesByMerchant(ctx, &pb.FindYearCashierByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchant),
	})
	if err != nil {

		logError("Failed to retrieve yearly sales data", err, zap.Error(err))

		return cashier_errors.ErrApiCashierFailedYearSalesByMerchant(c)
	}

	so := h.mapping.ToApiResponseCashierYearlySale(res)

	logSuccess("Successfully retrieve yearly sales data", zap.Bool("success", true))

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
	const method = "FindMonthSalesById"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		logError("Invalid year parameter", err, zap.String("year", c.QueryParam("year")))
		return cashier_errors.ErrApiCashierInvalidYear(c)
	}

	cashier, err := parseQueryIntWithValidation(c, "cashier_id", 1, 9999)

	if err != nil {
		logError("Invalid cashier_id parameter", err, zap.String("cashier_id", c.QueryParam("cashier_id")))
		return cashier_errors.ErrApiCashierInvalidMerchantId(c)
	}

	res, err := h.client.FindMonthSalesById(ctx, &pb.FindYearCashierById{
		Year:      int32(year),
		CashierId: int32(cashier),
	})

	if err != nil {
		logError("Failed to retrieve monthly sales data", err, zap.Error(err))

		return cashier_errors.ErrApiCashierFailedMonthSalesById(c)
	}

	so := h.mapping.ToApiResponseCashierMonthlySale(res)

	logSuccess("Successfully retrieve monthly sales data", zap.Bool("success", true))

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
	const method = "FindYearSalesById"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		logError("Invalid year parameter", err, zap.String("year", c.QueryParam("year")))
		return cashier_errors.ErrApiCashierInvalidYear(c)
	}

	cashier_id, err := parseQueryIntWithValidation(c, "cashier_id", 1, 9999)

	if err != nil {
		logError("Invalid cashier_id parameter", err, zap.String("cashier_id", c.QueryParam("cashier_id")))
		return cashier_errors.ErrApiCashierInvalidMerchantId(c)
	}

	res, err := h.client.FindYearSalesById(ctx, &pb.FindYearCashierById{
		Year:      int32(year),
		CashierId: int32(cashier_id),
	})
	if err != nil {
		logError("Failed to retrieve yearly sales data", err, zap.Error(err))

		return cashier_errors.ErrApiCashierFailedYearSalesById(c)
	}

	so := h.mapping.ToApiResponseCashierYearlySale(res)

	logSuccess("Successfully retrieve yearly sales data", zap.Bool("success", true))

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
	const method = "CreateCashier"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	var body requests.CreateCashierRequest

	if err := c.Bind(&body); err != nil {
		logError("Failed to bind request body", err, zap.Error(err))

		return cashier_errors.ErrApiBindCreateCashier(c)
	}

	if err := body.Validate(); err != nil {
		logError("Invalid request body", err, zap.Error(err))

		return cashier_errors.ErrApiValidateCreateCashier(c)
	}

	req := &pb.CreateCashierRequest{
		MerchantId: int32(body.MerchantID),
		UserId:     int32(body.UserID),
		Name:       body.Name,
	}

	res, err := h.client.CreateCashier(ctx, req)
	if err != nil {
		logError("Failed to create cashier", err, zap.Error(err))

		return cashier_errors.ErrApiCashierFailedCreate(c)
	}

	so := h.mapping.ToApiResponseCashier(res)

	logSuccess("Successfully created cashier", zap.Bool("success", true))

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
	const method = "UpdateCashier"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	id := c.Param("id")

	idStr, err := strconv.Atoi(id)

	if err != nil {
		logError("Invalid cashier ID format", err, zap.Error(err))

		return cashier_errors.ErrApiCashierInvalidId(c)
	}

	var body requests.UpdateCashierRequest

	if err := c.Bind(&body); err != nil {
		logError("Failed to bind request body", err, zap.Error(err))

		return cashier_errors.ErrApiBindUpdateCashier(c)
	}

	if err := body.Validate(); err != nil {
		logError("Invalid request body", err, zap.Error(err))

		return cashier_errors.ErrApiValidateUpdateCashier(c)
	}

	req := &pb.UpdateCashierRequest{
		CashierId: int32(idStr),
		Name:      body.Name,
	}

	res, err := h.client.UpdateCashier(ctx, req)
	if err != nil {
		logError("Failed to update cashier", err, zap.Error(err))

		return cashier_errors.ErrApiCashierFailedUpdate(c)
	}

	so := h.mapping.ToApiResponseCashier(res)

	logSuccess("Successfully updated cashier", zap.Bool("success", true))

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
	const method = "TrashedCashier"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		logError("Invalid cashier ID format", err, zap.Error(err))

		return cashier_errors.ErrApiCashierInvalidId(c)
	}

	req := &pb.FindByIdCashierRequest{Id: int32(id)}

	cashier, err := h.client.TrashedCashier(ctx, req)

	if err != nil {
		logError("Failed to retrieve trashed cashier", err, zap.Error(err))

		return cashier_errors.ErrApiCashierFailedTrashed(c)
	}

	so := h.mapping.ToApiResponseCashierDeleteAt(cashier)

	logSuccess("Successfully retrieved trashed cashier", zap.Bool("success", true))

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
	const method = "RestoreCashier"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logError("Invalid cashier ID format", err, zap.Error(err))

		return cashier_errors.ErrApiCashierInvalidId(c)
	}

	req := &pb.FindByIdCashierRequest{Id: int32(id)}

	cashier, err := h.client.RestoreCashier(ctx, req)
	if err != nil {
		logError("Failed to restore cashier", err, zap.Error(err))

		return cashier_errors.ErrApiCashierFailedRestore(c)
	}

	so := h.mapping.ToApiResponseCashierDeleteAt(cashier)

	logSuccess("Successfully restored cashier", zap.Bool("success", true))

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
	const method = "DeleteCashierPermanent"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logError("Invalid cashier ID format", err, zap.Error(err))

		return cashier_errors.ErrApiCashierInvalidId(c)
	}

	req := &pb.FindByIdCashierRequest{Id: int32(id)}

	cashier, err := h.client.DeleteCashierPermanent(ctx, req)
	if err != nil {
		logError("Failed to delete cashier", err, zap.Error(err))

		return cashier_errors.ErrApiCashierFailedDeletePermanent(c)
	}

	so := h.mapping.ToApiResponseCashierDelete(cashier)

	logSuccess("Successfully deleted cashier record permanently", zap.Bool("success", true))

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
	const method = "RestoreAllCashier"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	res, err := h.client.RestoreAllCashier(ctx, &emptypb.Empty{})
	if err != nil {
		logError("Failed to restore all cashier", err, zap.Error(err))

		return cashier_errors.ErrApiCashierFailedRestoreAll(c)
	}

	so := h.mapping.ToApiResponseCashierAll(res)

	logSuccess("Successfully restored all cashier", zap.Bool("success", true))

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
	const method = "DeleteAllCashierPermanent"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	res, err := h.client.DeleteAllCashierPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		logError("Failed to permanently delete all cashier", err, zap.Error(err))

		return cashier_errors.ErrApiCashierFailedDeleteAllPermanent(c)
	}

	so := h.mapping.ToApiResponseCashierAll(res)

	logSuccess("Successfully deleted all cashier permanently", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *cashierHandleApi) startTracingAndLogging(
	ctx context.Context,
	method string,
	attrs ...attribute.KeyValue,
) (
	end func(),
	logSuccess func(string, ...zap.Field),
	logError func(string, error, ...zap.Field),
) {
	start := time.Now()
	_, span := s.trace.Start(ctx, method)

	if len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}

	span.AddEvent("Start: " + method)
	s.logger.Debug("Start: " + method)

	status := "success"

	end = func() {
		s.recordMetrics(method, status, start)
		code := otelcode.Ok
		if status != "success" {
			code = otelcode.Error
		}
		span.SetStatus(code, status)
		span.End()
	}

	logSuccess = func(msg string, fields ...zap.Field) {
		status = "success"
		span.AddEvent(msg)
		s.logger.Debug(msg, fields...)
	}

	logError = func(msg string, err error, fields ...zap.Field) {
		status = "error"
		span.RecordError(err)
		span.SetStatus(otelcode.Error, msg)
		span.AddEvent(msg)
		allFields := append([]zap.Field{zap.Error(err)}, fields...)
		s.logger.Error(msg, allFields...)
	}

	return end, logSuccess, logError
}

func (s *cashierHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
