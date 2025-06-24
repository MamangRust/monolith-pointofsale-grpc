package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/transaction_errors"
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

type transactionHandleApi struct {
	client          pb.TransactionServiceClient
	logger          logger.LoggerInterface
	mapping         response_api.TransactionResponseMapper
	trace           trace.Tracer
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewHandlerTransaction(
	router *echo.Echo,
	client pb.TransactionServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.TransactionResponseMapper,
) *transactionHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transaction_handler_requests_total",
			Help: "Total number of transaction requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transaction_handler_request_duration_seconds",
			Help:    "Duration of transaction requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter)

	transactionHandle := &transactionHandleApi{
		client:          client,
		logger:          logger,
		mapping:         mapping,
		trace:           otel.Tracer("transaction-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routerTransaction := router.Group("/api/transaction")

	routerTransaction.GET("", transactionHandle.FindAllTransaction)
	routerTransaction.GET("/:id", transactionHandle.FindById)
	routerTransaction.GET("/merchant/:merchant_id", transactionHandle.FindByMerchant)
	routerTransaction.GET("/active", transactionHandle.FindByActive)
	routerTransaction.GET("/trashed", transactionHandle.FindByTrashed)

	routerTransaction.GET("/monthly-success", transactionHandle.FindMonthStatusSuccess)
	routerTransaction.GET("/yearly-success", transactionHandle.FindYearStatusSuccess)
	routerTransaction.GET("/monthly-failed", transactionHandle.FindMonthStatusFailed)
	routerTransaction.GET("/yearly-failed", transactionHandle.FindYearStatusFailed)

	routerTransaction.GET("/merchant/monthly-success", transactionHandle.FindMonthStatusSuccessByMerchant)
	routerTransaction.GET("/merchant/yearly-success", transactionHandle.FindYearStatusSuccessByMerchant)
	routerTransaction.GET("/merchant/monthly-failed", transactionHandle.FindMonthStatusFailedByMerchant)
	routerTransaction.GET("/merchant/yearly-failed", transactionHandle.FindYearStatusFailedByMerchant)

	routerTransaction.GET("/monthly-method-success", transactionHandle.FindMonthMethodSuccess)
	routerTransaction.GET("/yearly-method-success", transactionHandle.FindYearMethodSuccess)

	routerTransaction.GET("/merchant/monthly-method-success/:merchant_id", transactionHandle.FindMonthMethodByMerchantSuccess)
	routerTransaction.GET("/merchant/yearly-method-success/:merchant_id", transactionHandle.FindYearMethodByMerchantSuccess)

	routerTransaction.GET("/monthly-method-failed", transactionHandle.FindMonthMethodFailed)
	routerTransaction.GET("/yearly-method-failed", transactionHandle.FindYearMethodFailed)

	routerTransaction.GET("/merchant/monthly-method-failed/:merchant_id", transactionHandle.FindMonthMethodByMerchantFailed)
	routerTransaction.GET("/merchant/yearly-method-failed/:merchant_id", transactionHandle.FindYearMethodByMerchantFailed)

	routerTransaction.POST("/create", transactionHandle.Create)
	routerTransaction.POST("/update/:id", transactionHandle.Update)

	routerTransaction.POST("/trashed/:id", transactionHandle.TrashedTransaction)
	routerTransaction.POST("/restore/:id", transactionHandle.RestoreTransaction)
	routerTransaction.DELETE("/permanent/:id", transactionHandle.DeleteTransactionPermanent)

	routerTransaction.POST("/restore/all", transactionHandle.RestoreAllTransaction)
	routerTransaction.POST("/permanent/all", transactionHandle.DeleteAllTransactionPermanent)

	return transactionHandle
}

// @Security Bearer
// @Summary Find all transactions
// @Tags Transaction
// @Description Retrieve a list of all transactions
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationTransaction "List of transactions"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/transaction [get]
func (h *transactionHandleApi) FindAllTransaction(c echo.Context) error {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindAllTransaction"
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

	req := &pb.FindAllTransactionRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, req)

	if err != nil {
		logError("failed to retrieve transaction data", err, zap.Error(err))

		return transaction_errors.ErrApiTransactionFailedFindAll(c)
	}

	so := h.mapping.ToApiResponsePaginationTransaction(res)

	logSuccess("success to retrieve transaction data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Find all transactions by merchant
// @Tags Transaction
// @Description Retrieve a list of all transactions filtered by merchant
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationTransaction "List of transactions"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/transaction [get]
func (h *transactionHandleApi) FindByMerchant(c echo.Context) error {
	merchantID, err := strconv.Atoi(c.Param("merchant_id"))

	if err != nil || merchantID <= 0 {
		h.logger.Debug("Invalid merchant ID format", zap.Error(err))
		return transaction_errors.ErrApiTransactionInvalidMerchantId(c)
	}

	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindByMerchant"
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

	req := &pb.FindAllTransactionMerchantRequest{
		MerchantId: int32(merchantID),
		Page:       int32(page),
		PageSize:   int32(pageSize),
		Search:     search,
	}

	res, err := h.client.FindByMerchant(ctx, req)

	if err != nil {
		logError("failed to retrieve transaction data", err, zap.Error(err))

		return transaction_errors.ErrApiTransactionFailedFindByMerchant(c)
	}

	so := h.mapping.ToApiResponsePaginationTransaction(res)

	logSuccess("success to retrieve transaction data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Find transaction by ID
// @Tags Transaction
// @Description Retrieve a transaction by ID
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} response.ApiResponseTransaction "Transaction data"
// @Failure 400 {object} response.ErrorResponse "Invalid transaction ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/transaction/{id} [get]
func (h *transactionHandleApi) FindById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid transaction ID format", zap.Error(err))
		return transaction_errors.ErrApiTransactionInvalidId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdTransactionRequest{
		Id: int32(id),
	}

	res, err := h.client.FindById(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch transaction details", zap.Error(err))
		return transaction_errors.ErrApiTransactionFailedFindById(c)
	}

	so := h.mapping.ToApiResponseTransaction(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve active transactions
// @Tags Transaction
// @Description Retrieve a list of active transactions
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationTransactionDeleteAt "List of active transactions"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/transaction/active [get]
func (h *transactionHandleApi) FindByActive(c echo.Context) error {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindByActive"
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

	req := &pb.FindAllTransactionRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, req)

	if err != nil {
		logError("failed to retrieve transaction data", err, zap.Error(err))

		return transaction_errors.ErrApiTransactionFailedFindByActive(c)
	}

	so := h.mapping.ToApiResponsePaginationTransactionDeleteAt(res)

	logSuccess("success to retrieve transaction data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// FindByTrashed retrieves a list of trashed transaction records.
// @Summary Retrieve trashed transactions
// @Tags Transaction
// @Description Retrieve a list of trashed transaction records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationTransactionDeleteAt "List of trashed transaction data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve transaction data"
// @Router /api/transaction/trashed [get]
func (h *transactionHandleApi) FindByTrashed(c echo.Context) error {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindByTrashed"
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

	req := &pb.FindAllTransactionRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, req)

	if err != nil {
		logError("Failed to fetch trashed transactions", err, zap.Error(err))

		return transaction_errors.ErrApiTransactionFailedFindByTrashed(c)
	}

	so := h.mapping.ToApiResponsePaginationTransactionDeleteAt(res)

	logSuccess("Successfully fetched trashed transactions", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthStatusSuccess retrieves monthly successful transactions
// @Summary Get monthly successful transactions
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of successful transactions by month
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param month query int true "Month in MM format (1-12)"
// @Success 200 {object} response.ApiResponsesTransactionMonthSuccess
// @Failure 400 {object} response.ErrorResponse "Invalid year or month parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/monthly-success [get]
func (h *transactionHandleApi) FindMonthStatusSuccess(c echo.Context) error {
	const method = "FindMonthStatusSuccess"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)

	if err != nil {
		logError("Invalid year parameter", err, zap.String("year", c.QueryParam("year")))
		return transaction_errors.ErrApiTransactionInvalidYear(c)
	}

	month, err := parseQueryIntWithValidation(c, "month", 1, 12)

	if err != nil {
		logError("Invalid month parameter", err, zap.String("month", c.QueryParam("month")))
		return transaction_errors.ErrApiTransactionInvalidMonth(c)
	}

	res, err := h.client.FindMonthStatusSuccess(ctx, &pb.FindMonthlyTransactionStatus{
		Year:  int32(year),
		Month: int32(month),
	})

	if err != nil {
		logError("Failed to fetch monthly successful transactions", err, zap.Error(err))

		return transaction_errors.ErrApiTransactionFailedFindMonthSuccess(c)
	}

	so := h.mapping.ToApiResponseTransactionMonthAmountSuccess(res)

	logSuccess("Successfully fetched monthly successful transactions", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearStatusSuccess retrieves yearly successful transactions
// @Summary Get yearly successful transactions
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of successful transactions by year
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponsesTransactionYearSuccess
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/yearly-success [get]
func (h *transactionHandleApi) FindYearStatusSuccess(c echo.Context) error {
	const method = "FindMonthlyTotalSalesByMerchant"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)

	if err != nil {
		logError("Invalid year parameter", err, zap.String("year", c.QueryParam("year")))
		return transaction_errors.ErrApiTransactionInvalidYear(c)
	}

	res, err := h.client.FindYearStatusSuccess(ctx, &pb.FindYearlyTransactionStatus{
		Year: int32(year),
	})

	if err != nil {
		logError("Failed to fetch yearly successful transactions", err, zap.Error(err))

		return transaction_errors.ErrApiTransactionFailedFindYearSuccess(c)
	}

	so := h.mapping.ToApiResponseTransactionYearAmountSuccess(res)

	logSuccess("Successfully fetched yearly successful transactions", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthStatusFailed retrieves monthly failed transactions
// @Summary Get monthly failed transactions
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of failed transactions by month
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param month query int true "Month in MM format (1-12)"
// @Success 200 {object} response.ApiResponsesTransactionMonthFailed
// @Failure 400 {object} response.ErrorResponse "Invalid year or month parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/monthly-failed [get]
func (h *transactionHandleApi) FindMonthStatusFailed(c echo.Context) error {
	const method = "FindMonthlyTotalSalesByMerchant"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)

	if err != nil {
		logError("Invalid year parameter", err, zap.String("year", c.QueryParam("year")))
		return transaction_errors.ErrApiTransactionInvalidYear(c)
	}

	month, err := parseQueryIntWithValidation(c, "month", 1, 12)

	if err != nil {
		logError("Invalid month parameter", err, zap.String("month", c.QueryParam("month")))
		return transaction_errors.ErrApiTransactionInvalidMonth(c)
	}

	res, err := h.client.FindMonthStatusFailed(ctx, &pb.FindMonthlyTransactionStatus{
		Year:  int32(year),
		Month: int32(month),
	})

	if err != nil {
		logError("Failed to fetch monthly failed transactions", err, zap.Error(err))

		return transaction_errors.ErrApiTransactionFailedFindMonthFailed(c)
	}

	so := h.mapping.ToApiResponseTransactionMonthAmountFailed(res)

	logSuccess("Successfully fetched monthly failed transactions", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearStatusFailed retrieves yearly failed transactions
// @Summary Get yearly failed transactions
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of failed transactions by year
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponsesTransactionYearFailed
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/yearly-failed [get]
func (h *transactionHandleApi) FindYearStatusFailed(c echo.Context) error {
	const method = "FindMonthlyTotalSalesByMerchant"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)

	if err != nil {
		logError("Invalid year parameter", err, zap.String("year", c.QueryParam("year")))
		return transaction_errors.ErrApiTransactionInvalidYear(c)
	}

	res, err := h.client.FindYearStatusFailed(ctx, &pb.FindYearlyTransactionStatus{
		Year: int32(year),
	})

	if err != nil {
		logError("Failed to fetch yearly failed transactions", err, zap.Error(err))

		return transaction_errors.ErrApiTransactionFailedFindYearFailed(c)
	}

	so := h.mapping.ToApiResponseTransactionYearAmountFailed(res)

	logSuccess("Successfully fetched yearly failed transactions", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthStatusSuccessByMerchant retrieves monthly successful transactions by merchant
// @Summary Get monthly successful transactions by merchant
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of successful transactions by month for specific merchant
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param month query int true "Month in MM format (1-12)"
// @Success 200 {object} response.ApiResponsesTransactionMonthSuccess
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID, year or month parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Merchant not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/merchant/monthly-success [get]
func (h *transactionHandleApi) FindMonthStatusSuccessByMerchant(c echo.Context) error {
	const method = "FindMonthlyTotalSalesByMerchant"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)

	if err != nil {
		logError("Invalid year parameter", err, zap.String("year", c.QueryParam("year")))
		return transaction_errors.ErrApiTransactionInvalidYear(c)
	}

	month, err := parseQueryIntWithValidation(c, "month", 1, 12)

	if err != nil {
		logError("Invalid month parameter", err, zap.String("month", c.QueryParam("month")))
		return transaction_errors.ErrApiTransactionInvalidMonth(c)
	}

	merchant_id, err := parseQueryIntWithValidation(c, "merchant_id", 1, 9999)

	if err != nil {
		logError("Invalid merchant ID parameter", err, zap.String("merchant_id", c.QueryParam("merchant_id")))
		return transaction_errors.ErrApiTransactionInvalidMerchantId(c)
	}

	res, err := h.client.FindMonthStatusSuccessByMerchant(ctx, &pb.FindMonthlyTransactionStatusByMerchant{
		Year:       int32(year),
		Month:      int32(month),
		MerchantId: int32(merchant_id),
	})

	if err != nil {
		logError("Failed to fetch monthly successful transactions", err, zap.Error(err))

		return transaction_errors.ErrApiTransactionFailedFindMonthSuccessByMerchant(c)
	}

	so := h.mapping.ToApiResponseTransactionMonthAmountSuccess(res)

	logSuccess("Successfully fetched monthly successful transactions", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearStatusSuccessByMerchant retrieves yearly successful transactions by merchant
// @Summary Get yearly successful transactions by merchant
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of successful transactions by year for specific merchant
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponsesTransactionYearSuccess
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Merchant not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/merchant/yearly-success [get]
func (h *transactionHandleApi) FindYearStatusSuccessByMerchant(c echo.Context) error {
	const method = "FindMonthlyTotalSalesByMerchant"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)

	if err != nil {
		logError("Invalid year parameter", err, zap.String("year", c.QueryParam("year")))
		return transaction_errors.ErrApiTransactionInvalidYear(c)
	}

	merchant_id, err := parseQueryIntWithValidation(c, "merchant_id", 1, 9999)

	if err != nil {
		logError("Invalid merchant ID parameter", err, zap.String("merchant_id", c.QueryParam("merchant_id")))
		return transaction_errors.ErrApiTransactionInvalidMerchantId(c)
	}

	res, err := h.client.FindYearStatusSuccessByMerchant(ctx, &pb.FindYearlyTransactionStatusByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchant_id),
	})

	if err != nil {
		logError("Failed to fetch yearly successful transactions", err, zap.Error(err))

		return transaction_errors.ErrApiTransactionFailedFindYearSuccessByMerchant(c)
	}

	so := h.mapping.ToApiResponseTransactionYearAmountSuccess(res)

	logSuccess("Successfully fetched yearly successful transactions", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthStatusFailedByMerchant retrieves monthly failed transactions by merchant
// @Summary Get monthly failed transactions by merchant
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of failed transactions by month for specific merchant
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param month query int true "Month in MM format (1-12)"
// @Success 200 {object} response.ApiResponsesTransactionMonthFailed
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID, year or month parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Merchant not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/merchant/monthly-failed [get]
func (h *transactionHandleApi) FindMonthStatusFailedByMerchant(c echo.Context) error {
	const method = "FindMonthlyTotalSalesByMerchant"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)

	if err != nil {
		logError("Invalid year parameter", err, zap.String("year", c.QueryParam("year")))
		return transaction_errors.ErrApiTransactionInvalidYear(c)
	}

	month, err := parseQueryIntWithValidation(c, "month", 1, 12)

	if err != nil {
		logError("Invalid month parameter", err, zap.String("month", c.QueryParam("month")))
		return transaction_errors.ErrApiTransactionInvalidMonth(c)
	}

	merchant_id, err := parseQueryIntWithValidation(c, "merchant_id", 1, 9999)

	if err != nil {
		logError("Invalid merchant ID parameter", err, zap.String("merchant_id", c.QueryParam("merchant_id")))
		return transaction_errors.ErrApiTransactionInvalidMerchantId(c)
	}

	res, err := h.client.FindMonthStatusFailedByMerchant(ctx, &pb.FindMonthlyTransactionStatusByMerchant{
		Year:       int32(year),
		Month:      int32(month),
		MerchantId: int32(merchant_id),
	})

	if err != nil {
		logError("Failed to fetch monthly failed transactions", err, zap.Error(err))

		return transaction_errors.ErrApiTransactionFailedFindMonthFailedByMerchant(c)
	}

	so := h.mapping.ToApiResponseTransactionMonthAmountFailed(res)

	logSuccess("Successfully fetched monthly failed transactions", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearStatusFailedByMerchant retrieves yearly failed transactions by merchant
// @Summary Get yearly failed transactions by merchant
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of failed transactions by year for specific merchant
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponsesTransactionYearFailed
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Merchant not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/merchant/yearly-failed [get]
func (h *transactionHandleApi) FindYearStatusFailedByMerchant(c echo.Context) error {
	const method = "FindMonthlyTotalSalesByMerchant"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)

	if err != nil {
		logError("Invalid year parameter", err, zap.String("year", c.QueryParam("year")))
		return transaction_errors.ErrApiTransactionInvalidYear(c)
	}

	merchant_id, err := parseQueryIntWithValidation(c, "merchant_id", 1, 9999)

	if err != nil {
		logError("Invalid merchant ID parameter", err, zap.String("merchant_id", c.QueryParam("merchant_id")))
		return transaction_errors.ErrApiTransactionInvalidMerchantId(c)
	}

	res, err := h.client.FindYearStatusFailedByMerchant(ctx, &pb.FindYearlyTransactionStatusByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchant_id),
	})

	if err != nil {
		logError("Failed to fetch yearly failed transactions", err, zap.Error(err))

		return transaction_errors.ErrApiTransactionFailedFindYearFailedByMerchant(c)
	}

	so := h.mapping.ToApiResponseTransactionYearAmountFailed(res)

	logSuccess("Successfully fetched yearly failed transactions", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthMethod retrieves monthly payment method statistics
// @Summary Get monthly payment method distribution
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of payment methods used by month
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponsesTransactionMonthMethod
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/monthly-method-success [get]
func (h *transactionHandleApi) FindMonthMethodSuccess(c echo.Context) error {
	const method = "FindMonthMethodSuccess"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)

	if err != nil {
		logError("Invalid year parameter", err, zap.String("year", c.QueryParam("year")))
		return transaction_errors.ErrApiTransactionInvalidYear(c)
	}

	month, err := parseQueryIntWithValidation(c, "month", 1, 12)

	if err != nil {
		logError("Invalid month parameter", err, zap.String("month", c.QueryParam("month")))
		return transaction_errors.ErrApiTransactionInvalidMonth(c)
	}

	res, err := h.client.FindMonthMethodSuccess(ctx, &pb.MonthTransactionMethod{
		Year:  int32(year),
		Month: int32(month),
	})
	if err != nil {
		logError("Failed to fetch monthly success transactions", err, zap.Error(err))

		return transaction_errors.ErrApiTransactionFailedFindMonthMethod(c)
	}

	so := h.mapping.ToApiResponseTransactionMonthMethod(res)

	logSuccess("Successfully fetched monthly success transactions", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearMethod retrieves yearly payment method statistics
// @Summary Get yearly payment method distribution
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of payment methods used by year
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponsesTransactionYearMethod
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/yearly-method-success [get]
func (h *transactionHandleApi) FindYearMethodSuccess(c echo.Context) error {

	const method = "FindMonthlyTotalSalesByMerchant"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)

	if err != nil {
		logError("Invalid year parameter", err, zap.String("year", c.QueryParam("year")))
		return transaction_errors.ErrApiTransactionInvalidYear(c)
	}

	res, err := h.client.FindYearMethodSuccess(ctx, &pb.YearTransactionMethod{
		Year: int32(year),
	})
	if err != nil {
		logError("Failed to fetch yearly success transactions", err, zap.Error(err))

		return transaction_errors.ErrApiTransactionFailedFindYearMethod(c)
	}

	so := h.mapping.ToApiResponseTransactionYearMethod(res)

	logSuccess("Successfully fetched yearly success transactions", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthMethodByMerchant retrieves monthly payment method statistics by merchant
// @Summary Get monthly payment method distribution by merchant
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of payment methods used by month for specific merchant
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponsesTransactionMonthMethod
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Merchant not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/merchant/monthly-method-success/{merchant_id} [get]
func (h *transactionHandleApi) FindMonthMethodByMerchantSuccess(c echo.Context) error {
	const method = "FindMonthMethodByMerchantSuccess"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)

	if err != nil {
		logError("Invalid year parameter", err, zap.String("year", c.QueryParam("year")))
		return transaction_errors.ErrApiTransactionInvalidYear(c)
	}

	month, err := parseQueryIntWithValidation(c, "month", 1, 12)

	if err != nil {
		logError("Invalid month parameter", err, zap.String("month", c.QueryParam("month")))
		return transaction_errors.ErrApiTransactionInvalidMonth(c)
	}

	merchant_id, err := parseQueryIntWithValidation(c, "merchant_id", 1, 9999)

	if err != nil {
		logError("Invalid merchant ID parameter", err, zap.String("merchant_id", c.QueryParam("merchant_id")))
		return transaction_errors.ErrApiTransactionInvalidMerchantId(c)
	}

	res, err := h.client.FindMonthMethodByMerchantSuccess(ctx, &pb.MonthTransactionMethodByMerchant{
		Year:       int32(year),
		Month:      int32(month),
		MerchantId: int32(merchant_id),
	})

	if err != nil {
		logError("Failed to fetch monthly success transactions", err, zap.Error(err))

		return transaction_errors.ErrApiTransactionFailedFindMonthMethodByMerchant(c)
	}

	so := h.mapping.ToApiResponseTransactionMonthMethod(res)

	logSuccess("Successfully fetched monthly success transactions", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearMethodByMerchant retrieves yearly payment method statistics by merchant
// @Summary Get yearly payment method distribution by merchant
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of payment methods used by year for specific merchant
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponsesTransactionYearMethod
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Merchant not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/merchant/yearly-method-success/{merchant_id} [get]
func (h *transactionHandleApi) FindYearMethodByMerchantSuccess(c echo.Context) error {
	const method = "FindMonthlyTotalSalesByMerchant"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)

	if err != nil {
		logError("Invalid year parameter", err, zap.String("year", c.QueryParam("year")))
		return transaction_errors.ErrApiTransactionInvalidYear(c)
	}

	merchant_id, err := parseQueryIntWithValidation(c, "merchant_id", 1, 9999)

	if err != nil {
		logError("Invalid merchant ID parameter", err, zap.String("merchant_id", c.QueryParam("merchant_id")))
		return transaction_errors.ErrApiTransactionInvalidMerchantId(c)
	}

	res, err := h.client.FindYearMethodByMerchantSuccess(ctx, &pb.YearTransactionMethodByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchant_id),
	})
	if err != nil {
		logError("Failed to fetch yearly success transactions", err, zap.Error(err))

		return transaction_errors.ErrApiTransactionFailedFindYearMethodByMerchant(c)
	}

	so := h.mapping.ToApiResponseTransactionYearMethod(res)

	logSuccess("Successfully fetched yearly success transactions", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthMethod retrieves monthly payment method statistics
// @Summary Get monthly payment method distribution
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of payment methods used by month
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponsesTransactionMonthMethod
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/monthly-method-failed [get]
func (h *transactionHandleApi) FindMonthMethodFailed(c echo.Context) error {
	const method = "FindMonthMethodFailed"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)

	if err != nil {
		logError("Invalid year parameter", err, zap.String("year", c.QueryParam("year")))
		return transaction_errors.ErrApiTransactionInvalidYear(c)
	}

	month, err := parseQueryIntWithValidation(c, "month", 1, 12)

	if err != nil {
		logError("Invalid month parameter", err, zap.String("month", c.QueryParam("month")))
		return transaction_errors.ErrApiTransactionInvalidMonth(c)
	}

	res, err := h.client.FindMonthMethodFailed(ctx, &pb.MonthTransactionMethod{
		Year:  int32(year),
		Month: int32(month),
	})
	if err != nil {
		logError("Failed to fetch monthly failed transactions", err, zap.Error(err))

		return transaction_errors.ErrApiTransactionFailedFindMonthMethod(c)
	}

	so := h.mapping.ToApiResponseTransactionMonthMethod(res)

	logSuccess("Successfully fetched monthly failed transactions", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearMethod retrieves yearly payment method statistics
// @Summary Get yearly payment method distribution
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of payment methods used by year
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponsesTransactionYearMethod
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/yearly-method-failed [get]
func (h *transactionHandleApi) FindYearMethodFailed(c echo.Context) error {
	const method = "FindYearMethodFailed"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)

	if err != nil {
		logError("Invalid year parameter", err, zap.String("year", c.QueryParam("year")))
		return transaction_errors.ErrApiTransactionInvalidYear(c)
	}

	res, err := h.client.FindYearMethodFailed(ctx, &pb.YearTransactionMethod{
		Year: int32(year),
	})
	if err != nil {
		logError("Failed to fetch yearly failed transactions", err, zap.Error(err))

		return transaction_errors.ErrApiTransactionFailedFindYearMethod(c)
	}

	so := h.mapping.ToApiResponseTransactionYearMethod(res)

	logSuccess("Successfully fetched yearly failed transactions", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthMethodByMerchant retrieves monthly payment method statistics by merchant
// @Summary Get monthly payment method distribution by merchant
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of payment methods used by month for specific merchant
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponsesTransactionMonthMethod
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Merchant not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/merchant/monthly-method-failed/{merchant_id} [get]
func (h *transactionHandleApi) FindMonthMethodByMerchantFailed(c echo.Context) error {
	const method = "FindMonthlyTotalSalesByMerchant"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)

	if err != nil {
		logError("Invalid year parameter", err, zap.String("year", c.QueryParam("year")))
		return transaction_errors.ErrApiTransactionInvalidYear(c)
	}

	month, err := parseQueryIntWithValidation(c, "month", 1, 12)

	if err != nil {
		logError("Invalid month parameter", err, zap.String("month", c.QueryParam("month")))
		return transaction_errors.ErrApiTransactionInvalidMonth(c)
	}

	merchant_id, err := parseQueryIntWithValidation(c, "merchant_id", 1, 9999)

	if err != nil {
		logError("Invalid merchant ID parameter", err, zap.String("merchant_id", c.QueryParam("merchant_id")))
		return transaction_errors.ErrApiTransactionInvalidMerchantId(c)
	}

	res, err := h.client.FindMonthMethodByMerchantSuccess(ctx, &pb.MonthTransactionMethodByMerchant{
		Year:       int32(year),
		Month:      int32(month),
		MerchantId: int32(merchant_id),
	})

	if err != nil {
		logError("Failed to fetch monthly failed transactions by merchant", err, zap.Error(err))

		return transaction_errors.ErrApiTransactionFailedFindMonthMethodByMerchant(c)
	}

	so := h.mapping.ToApiResponseTransactionMonthMethod(res)

	logSuccess("Successfully fetched monthly failed transactions by merchant", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearMethodByMerchant retrieves yearly payment method statistics by merchant
// @Summary Get yearly payment method distribution by merchant
// @Tags Transaction
// @Security Bearer
// @Description Retrieve statistics of payment methods used by year for specific merchant
// @Accept json
// @Produce json
// @Param merchant_id query int true "Merchant ID"
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponsesTransactionYearMethod
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Merchant not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/transaction/merchant/yearly-method-failed/{merchant_id} [get]
func (h *transactionHandleApi) FindYearMethodByMerchantFailed(c echo.Context) error {
	const method = "FindMonthlyTotalSalesByMerchant"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)

	if err != nil {
		logError("Invalid year parameter", err, zap.String("year", c.QueryParam("year")))
		return transaction_errors.ErrApiTransactionInvalidYear(c)
	}

	merchant_id, err := parseQueryIntWithValidation(c, "merchant_id", 1, 9999)

	if err != nil {
		logError("Invalid merchant ID parameter", err, zap.String("merchant_id", c.QueryParam("merchant_id")))
		return transaction_errors.ErrApiTransactionInvalidMerchantId(c)
	}

	res, err := h.client.FindYearMethodByMerchantSuccess(ctx, &pb.YearTransactionMethodByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchant_id),
	})
	if err != nil {
		logError("Failed to fetch yearly failed transactions by merchant", err, zap.Error(err))

		return transaction_errors.ErrApiTransactionFailedFindYearMethodByMerchant(c)
	}

	so := h.mapping.ToApiResponseTransactionYearMethod(res)

	logSuccess("Successfully fetched yearly failed transactions by merchant", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Create a new transaction
// @Tags Transaction
// @Description Create a new transaction record
// @Accept json
// @Produce json
// @Param request body requests.CreateTransactionRequest true "Transaction details"
// @Success 200 {object} response.ApiResponseTransaction "Successfully created transaction"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create transaction"
// @Router /api/transaction/create [post]
func (h *transactionHandleApi) Create(c echo.Context) error {
	const method = "CreateTransaction"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	var body requests.CreateTransactionRequest

	if err := c.Bind(&body); err != nil {
		logError("Failed to bind request body", err, zap.Error(err))

		return transaction_errors.ErrApiBindCreateTransaction(c)
	}

	if err := body.Validate(); err != nil {
		logError("Failed to validate request body", err, zap.Error(err))

		return transaction_errors.ErrApiValidateCreateTransaction(c)
	}

	grpcReq := &pb.CreateTransactionRequest{
		OrderId:       int32(body.OrderID),
		CashierId:     int32(body.CashierID),
		PaymentMethod: body.PaymentMethod,
		Amount:        int32(body.Amount),
	}

	res, err := h.client.Create(ctx, grpcReq)

	if err != nil {
		logError("Failed to create transaction", err, zap.Error(err))

		return transaction_errors.ErrApiTransactionFailedCreate(c)
	}

	so := h.mapping.ToApiResponseTransaction(res)

	logSuccess("Success to create transaction", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Update an existing transaction
// @Tags Transaction
// @Description Update an existing transaction record
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Param request body requests.UpdateTransactionRequest true "Updated transaction details"
// @Success 200 {object} response.ApiResponseTransaction "Successfully updated transaction"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update transaction"
// @Router /api/transaction/update/{id} [post]
func (h *transactionHandleApi) Update(c echo.Context) error {
	const method = "Update"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		logError("Invalid transaction id", err, zap.Error(err))

		return transaction_errors.ErrApiTransactionInvalidId(c)
	}

	var body requests.UpdateTransactionRequest

	if err := c.Bind(&body); err != nil {
		logError("Invalid request format", err, zap.Error(err))

		return transaction_errors.ErrApiBindUpdateTransaction(c)
	}

	if err := body.Validate(); err != nil {
		logError("Validation failed", err, zap.Error(err))

		return transaction_errors.ErrApiValidateUpdateTransaction(c)
	}

	grpcReq := &pb.UpdateTransactionRequest{
		TransactionId: int32(idInt),
		OrderId:       int32(body.OrderID),
		CashierId:     int32(body.CashierID),
		PaymentMethod: body.PaymentMethod,
		Amount:        int32(body.Amount),
	}

	res, err := h.client.Update(ctx, grpcReq)

	if err != nil {
		logError("Transaction update failed", err, zap.Error(err))

		return transaction_errors.ErrApiTransactionFailedUpdate(c)
	}

	so := h.mapping.ToApiResponseTransaction(res)

	logSuccess("Successfully updated transaction", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// TrashedTransaction retrieves a trashed transaction record by its ID.
// @Summary Retrieve a trashed transaction
// @Tags Transaction
// @Description Retrieve a trashed transaction record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} response.ApiResponseTransactionDeleteAt "Successfully retrieved trashed transaction"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve trashed transaction"
// @Router /api/transaction/trashed/{id} [get]
func (h *transactionHandleApi) TrashedTransaction(c echo.Context) error {
	const method = "TrashedTransaction"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		logError("Invalid id parameter", err, zap.Error(err))

		return transaction_errors.ErrApiTransactionInvalidId(c)
	}

	req := &pb.FindByIdTransactionRequest{
		Id: int32(id),
	}

	res, err := h.client.TrashedTransaction(ctx, req)

	if err != nil {
		logError("Transaction trashed failed", err, zap.Error(err))

		return transaction_errors.ErrApiTransactionFailedTrashed(c)
	}

	so := h.mapping.ToApiResponseTransactionDeleteAt(res)

	logSuccess("Successfully retrieved trashed transaction", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreTransaction restores a transaction record from the trash by its ID.
// @Summary Restore a trashed transaction
// @Tags Transaction
// @Description Restore a trashed transaction record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} response.ApiResponseTransactionDeleteAt "Successfully restored transaction"
// @Failure 400 {object} response.ErrorResponse "Invalid transaction ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore transaction"
// @Router /api/transaction/restore/{id} [post]
func (h *transactionHandleApi) RestoreTransaction(c echo.Context) error {
	const method = "RestoreTransaction"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		logError("Invalid transaction ID format", err, zap.Error(err))

		return transaction_errors.ErrApiTransactionInvalidId(c)
	}

	req := &pb.FindByIdTransactionRequest{
		Id: int32(id),
	}

	res, err := h.client.RestoreTransaction(ctx, req)

	if err != nil {
		logError("Failed to restore transaction", err, zap.Error(err))

		return transaction_errors.ErrApiTransactionFailedRestore(c)
	}

	so := h.mapping.ToApiResponseTransactionDeleteAt(res)

	logSuccess("Successfully restored transaction", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteTransactionPermanent permanently deletes a transaction record by its ID.
// @Summary Permanently delete a transaction
// @Tags Transaction
// @Description Permanently delete a transaction record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} response.ApiResponseTransactionDelete "Successfully deleted transaction record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete transaction"
// @Router /api/transaction/delete/{id} [delete]
func (h *transactionHandleApi) DeleteTransactionPermanent(c echo.Context) error {
	const method = "DeleteTransactionPermanent"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		logError("Invalid transaction ID format", err, zap.Error(err))

		return transaction_errors.ErrApiTransactionInvalidId(c)
	}

	req := &pb.FindByIdTransactionRequest{
		Id: int32(id),
	}

	res, err := h.client.DeleteTransactionPermanent(ctx, req)

	if err != nil {
		logError("Failed to delete transaction", err, zap.Error(err))

		return transaction_errors.ErrApiTransactionFailedDeletePermanent(c)
	}

	so := h.mapping.ToApiResponseTransactionDelete(res)

	logSuccess("Successfully deleted transaction record permanently", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreAllTransaction restores all trashed transactions.
// @Summary Restore all trashed transactions
// @Tags Transaction
// @Description Restore all trashed transactions.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseTransactionAll "Successfully restored all transactions"
// @Failure 500 {object} response.ErrorResponse "Failed to restore transactions"
// @Router /api/transaction/restore/all [post]
func (h *transactionHandleApi) RestoreAllTransaction(c echo.Context) error {
	const method = "RestoreAllTransaction"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	res, err := h.client.RestoreAllTransaction(ctx, &emptypb.Empty{})

	if err != nil {
		logError("Failed to restore all transactions", err, zap.Error(err))

		return transaction_errors.ErrApiTransactionFailedRestoreAll(c)
	}

	so := h.mapping.ToApiResponseTransactionAll(res)

	logSuccess("Successfully restored all transactions", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteAllTransactionPermanent permanently deletes all transactions.
// @Summary Permanently delete all transactions
// @Tags Transaction
// @Description Permanently delete all transactions.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseTransactionAll "Successfully deleted all transactions permanently"
// @Failure 500 {object} response.ErrorResponse "Failed to delete transactions"
// @Router /api/transaction/delete/all [post]
func (h *transactionHandleApi) DeleteAllTransactionPermanent(c echo.Context) error {
	const method = "DeleteAllTransactionPermanent"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	res, err := h.client.DeleteAllTransactionPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		logError("Bulk transactions deletion failed", err, zap.Error(err))

		return transaction_errors.ErrApiTransactionFailedDeleteAllPermanent(c)
	}

	so := h.mapping.ToApiResponseTransactionAll(res)

	logSuccess("Successfully deleted all transactions permanently", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *transactionHandleApi) startTracingAndLogging(
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

func (s *transactionHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
