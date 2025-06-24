package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/category_errors"
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

type categoryHandleApi struct {
	client          pb.CategoryServiceClient
	logger          logger.LoggerInterface
	mapping         response_api.CategoryResponseMapper
	trace           trace.Tracer
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewHandlerCategory(
	router *echo.Echo,
	client pb.CategoryServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.CategoryResponseMapper,
) *categoryHandleApi {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "category_handler_requests_total",
			Help: "Total number of category requests",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "category_handler_request_duration_seconds",
			Help:    "Duration of category requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter)

	categoryHandler := &categoryHandleApi{
		client:          client,
		logger:          logger,
		mapping:         mapping,
		trace:           otel.Tracer("category-handler"),
		requestCounter:  requestCounter,
		requestDuration: requestDuration,
	}

	routercategory := router.Group("/api/category")

	routercategory.GET("", categoryHandler.FindAllCategory)
	routercategory.GET("/:id", categoryHandler.FindById)
	routercategory.GET("/active", categoryHandler.FindByActive)
	routercategory.GET("/trashed", categoryHandler.FindByTrashed)

	routercategory.GET("/monthly-total-pricing", categoryHandler.FindMonthTotalPrice)
	routercategory.GET("/yearly-total-pricing", categoryHandler.FindYearTotalPrice)
	routercategory.GET("/merchant/monthly-total-pricing", categoryHandler.FindMonthTotalPriceByMerchant)
	routercategory.GET("/merchant/yearly-total-pricing", categoryHandler.FindYearTotalPriceByMerchant)
	routercategory.GET("/mycategory/monthly-total-pricing", categoryHandler.FindMonthTotalPriceById)
	routercategory.GET("/mycategory/yearly-total-pricing", categoryHandler.FindYearTotalPriceById)

	routercategory.GET("/monthly-pricing", categoryHandler.FindMonthPrice)
	routercategory.GET("/yearly-pricing", categoryHandler.FindYearPrice)
	routercategory.GET("/merchant/monthly-pricing", categoryHandler.FindMonthPriceByMerchant)
	routercategory.GET("/merchant/yearly-pricing", categoryHandler.FindYearPriceByMerchant)
	routercategory.GET("/mycategory/monthly-pricing", categoryHandler.FindMonthPriceById)
	routercategory.GET("/mycategory/yearly-pricing", categoryHandler.FindYearPriceById)

	routercategory.POST("/create", categoryHandler.Create)
	routercategory.POST("/update/:id", categoryHandler.Update)

	routercategory.POST("/trashed/:id", categoryHandler.TrashedCategory)
	routercategory.POST("/restore/:id", categoryHandler.RestoreCategory)
	routercategory.DELETE("/permanent/:id", categoryHandler.DeleteCategoryPermanent)

	routercategory.POST("/restore/all", categoryHandler.RestoreAllCategory)
	routercategory.POST("/permanent/all", categoryHandler.DeleteAllCategoryPermanent)

	return categoryHandler
}

// @Security Bearer
// @Summary Find all category
// @Tags Category
// @Description Retrieve a list of all category
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationCategory "List of category"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve category data"
// @Router /api/category [get]
func (h *categoryHandleApi) FindAllCategory(c echo.Context) error {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		method          = "FindAllCategory"
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

	req := &pb.FindAllCategoryRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, req)

	if err != nil {

		logError("Failed to retrieve category data", err, zap.Error(err))

		return category_errors.ErrApiCategoryFailedFindAll(c)
	}

	so := h.mapping.ToApiResponsePaginationCategory(res)

	logSuccess("Successfully retrieve category data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Find category by ID
// @Tags Category
// @Description Retrieve a category by ID
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} response.ApiResponseCategory "Category data"
// @Failure 400 {object} response.ErrorResponse "Invalid category ID"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve category data"
// @Router /api/category/{id} [get]
func (h *categoryHandleApi) FindById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid category ID", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidId(c)
	}

	ctx := c.Request().Context()

	req := &pb.FindByIdCategoryRequest{
		Id: int32(id),
	}

	res, err := h.client.FindById(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch category details", zap.Error(err))
		return category_errors.ErrApiCategoryFailedFindById(c)
	}

	so := h.mapping.ToApiResponseCategory(res)

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// @Summary Retrieve active category
// @Tags Category
// @Description Retrieve a list of active category
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationCategoryDeleteAt "List of active category"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve category data"
// @Router /api/category/active [get]
func (h *categoryHandleApi) FindByActive(c echo.Context) error {
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

	req := &pb.FindAllCategoryRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, req)

	if err != nil {
		logError("Failed to retrieve category data", err, zap.Error(err))

		return category_errors.ErrApiCategoryFailedFindByActive(c)
	}

	so := h.mapping.ToApiResponsePaginationCategoryDeleteAt(res)

	logSuccess("Successfully retrieve category data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// FindByTrashed retrieves a list of trashed category records.
// @Summary Retrieve trashed category
// @Tags Category
// @Description Retrieve a list of trashed category records
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} response.ApiResponsePaginationCategoryDeleteAt "List of trashed category data"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve category data"
// @Router /api/category/trashed [get]
func (h *categoryHandleApi) FindByTrashed(c echo.Context) error {
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

	req := &pb.FindAllCategoryRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, req)

	if err != nil {
		logError("Failed to retrieve category data", err, zap.Error(err))

		return category_errors.ErrApiCategoryFailedFindByTrashed(c)
	}

	so := h.mapping.ToApiResponsePaginationCategoryDeleteAt(res)

	logSuccess("Successfully retrieve category data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthTotalPrice retrieves monthly category pricing statistics
// @Summary Get monthly category pricing
// @Tags Category
// @Security Bearer
// @Description Retrieve monthly pricing statistics for all categories
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponseCategoryMonthPrice "Monthly category pricing data"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/category/monthly-total-pricing [get]
func (h *categoryHandleApi) FindMonthTotalPrice(c echo.Context) error {
	const method = "FindMonthTotalPrice"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		logError("Invalid year parameter", err, zap.String("year", c.QueryParam("year")))
		return category_errors.ErrApiCategoryInvalidYear(c)
	}

	month, err := parseQueryIntWithValidation(c, "month", 1, 12)
	if err != nil {
		logError("Invalid month parameter", err, zap.String("month", c.QueryParam("month")))
		return category_errors.ErrApiCategoryInvalidMonth(c)
	}

	res, err := h.client.FindMonthlyTotalPrices(ctx, &pb.FindYearMonthTotalPrices{
		Year:  int32(year),
		Month: int32(month),
	})

	if err != nil {
		logError("Failed to retrieve category data", err, zap.Error(err))

		return category_errors.ErrApiCategoryFailedMonthTotalPrice(c)
	}

	so := h.mapping.ToApiResponseCategoryMonthlyTotalPrice(res)

	logSuccess("Successfully retrieve category data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearTotalPrice retrieves yearly category pricing statistics
// @Summary Get yearly category pricing
// @Tags Category
// @Security Bearer
// @Description Retrieve yearly pricing statistics for all categories
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponseCategoryYearPrice "Yearly category pricing data"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/category/yearly-total-pricing [get]
func (h *categoryHandleApi) FindYearTotalPrice(c echo.Context) error {
	const method = "FindYearTotalPrice"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		logError("Invalid year parameter", err, zap.String("year", c.QueryParam("year")))
		return category_errors.ErrApiCategoryInvalidYear(c)
	}

	res, err := h.client.FindYearlyTotalPrices(ctx, &pb.FindYearTotalPrices{
		Year: int32(year),
	})

	if err != nil {
		logError("Failed to retrieve category data", err, zap.Error(err))

		return category_errors.ErrApiCategoryFailedYearTotalPrice(c)
	}

	so := h.mapping.ToApiResponseCategoryYearlyTotalPrice(res)

	logSuccess("Successfully retrieve category data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthTotalPriceById retrieves monthly category pricing statistics
// @Summary Get monthly category pricing
// @Tags Category
// @Security Bearer
// @Description Retrieve monthly pricing statistics for all categories
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param month query int true "Month"
// @Param category_id query int true "Category ID"
// @Success 200 {object} response.ApiResponseCategoryMonthPrice "Monthly category pricing data"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/category/mycategory/monthly-total-pricing [get]
func (h *categoryHandleApi) FindMonthTotalPriceById(c echo.Context) error {
	const method = "FindMonthTotalPriceById"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		logError("Invalid year parameter", err, zap.String("year", c.QueryParam("year")))
		return category_errors.ErrApiCategoryInvalidYear(c)
	}

	month, err := parseQueryIntWithValidation(c, "month", 1, 12)
	if err != nil {
		logError("Invalid month parameter", err, zap.String("month", c.QueryParam("month")))
		return category_errors.ErrApiCategoryInvalidMonth(c)
	}

	category, err := parseQueryIntWithValidation(c, "category_id", 1, 9999)

	if err != nil {
		logError("Invalid category_id parameter", err, zap.String("category_id", c.QueryParam("category_id")))
		return category_errors.ErrApiCategoryInvalidId(c)
	}

	res, err := h.client.FindMonthlyTotalPricesById(ctx, &pb.FindYearMonthTotalPriceById{
		Year:       int32(year),
		Month:      int32(month),
		CategoryId: int32(category),
	})

	if err != nil {
		logError("Failed to retrieve category data", err, zap.Error(err))

		return category_errors.ErrApiCategoryFailedMonthTotalPriceById(c)
	}

	so := h.mapping.ToApiResponseCategoryMonthlyTotalPrice(res)

	logSuccess("Successfully retrieve category data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearTotalPriceById retrieves yearly category pricing statistics
// @Summary Get yearly category pricing
// @Tags Category
// @Security Bearer
// @Description Retrieve yearly pricing statistics for all categories
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param category_id query int true "Category ID"
// @Success 200 {object} response.ApiResponseCategoryYearPrice "Yearly category pricing data"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/category/yearly-total-pricing [get]
func (h *categoryHandleApi) FindYearTotalPriceById(c echo.Context) error {
	const method = "FindYearTotalPriceById"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		logError("Invalid year parameter", err, zap.String("year", c.QueryParam("year")))
		return category_errors.ErrApiCategoryInvalidYear(c)
	}

	category, err := parseQueryIntWithValidation(c, "category_id", 1, 9999)

	if err != nil {
		logError("Invalid category_id parameter", err, zap.String("category_id", c.QueryParam("category_id")))
		return category_errors.ErrApiCategoryInvalidId(c)
	}

	res, err := h.client.FindYearlyTotalPricesById(ctx, &pb.FindYearTotalPriceById{
		Year:       int32(year),
		CategoryId: int32(category),
	})

	if err != nil {
		logError("Failed to retrieve category data", err, zap.Error(err))

		return category_errors.ErrApiCategoryFailedYearTotalPriceById(c)
	}

	so := h.mapping.ToApiResponseCategoryYearlyTotalPrice(res)

	logSuccess("Successfully retrieve category data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthTotalPriceByMerchant retrieves monthly category pricing statistics
// @Summary Get monthly category pricing
// @Tags Category
// @Security Bearer
// @Description Retrieve monthly pricing statistics for all categories
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param month query int true "Month"
// @Param category_id query int true "Category ID"
// @Success 200 {object} response.ApiResponseCategoryMonthPrice "Monthly category pricing data"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/category/merchant/monthly-total-pricing [get]
func (h *categoryHandleApi) FindMonthTotalPriceByMerchant(c echo.Context) error {
	const method = "FindMonthTotalPriceByMerchant"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		logError("Invalid year parameter", err, zap.String("year", c.QueryParam("year")))
		return category_errors.ErrApiCategoryInvalidYear(c)
	}

	month, err := parseQueryIntWithValidation(c, "month", 1, 12)
	if err != nil {
		logError("Invalid month parameter", err, zap.String("month", c.QueryParam("month")))
		return category_errors.ErrApiCategoryInvalidMonth(c)
	}

	merchant, err := parseQueryIntWithValidation(c, "merchant_id", 1, 9999)

	if err != nil {
		logError("Invalid merchant_id parameter", err, zap.String("merchant_id", c.QueryParam("merchant_id")))
		return category_errors.ErrApiCategoryInvalidId(c)
	}

	res, err := h.client.FindMonthlyTotalPricesByMerchant(ctx, &pb.FindYearMonthTotalPriceByMerchant{
		Year:       int32(year),
		Month:      int32(month),
		MerchantId: int32(merchant),
	})

	if err != nil {
		logError("Failed to retrieve category data", err, zap.Error(err))

		return category_errors.ErrApiCategoryFailedMonthTotalPriceByMerchant(c)
	}

	so := h.mapping.ToApiResponseCategoryMonthlyTotalPrice(res)

	logSuccess("Successfully retrieve category data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearTotalPriceByMerchant retrieves yearly category total pricing statistics
// @Summary Get yearly category pricing
// @Tags Category
// @Security Bearer
// @Description Retrieve yearly pricing statistics for all categories
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param category_id query int true "Category ID"
// @Success 200 {object} response.ApiResponseCategoryYearPrice "Yearly category pricing data"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/merchant/yearly-total-pricing [get]
func (h *categoryHandleApi) FindYearTotalPriceByMerchant(c echo.Context) error {
	const method = "FindYearTotalPriceByMerchant"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		logError("Invalid year parameter", err, zap.String("year", c.QueryParam("year")))
		return category_errors.ErrApiCategoryInvalidYear(c)
	}

	merchant, err := parseQueryIntWithValidation(c, "merchant_id", 1, 9999)

	if err != nil {
		logError("Invalid merchant_id parameter", err, zap.String("merchant_id", c.QueryParam("merchant_id")))
		return category_errors.ErrApiCategoryInvalidId(c)
	}

	res, err := h.client.FindYearlyTotalPricesByMerchant(ctx, &pb.FindYearTotalPriceByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchant),
	})

	if err != nil {
		logError("Failed to retrieve category data", err, zap.Error(err))

		return category_errors.ErrApiCategoryFailedYearTotalPriceByMerchant(c)
	}

	so := h.mapping.ToApiResponseCategoryYearlyTotalPrice(res)

	logSuccess("Successfully retrieve category data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthPrice retrieves monthly category pricing statistics
// @Summary Get monthly category pricing
// @Tags Category
// @Security Bearer
// @Description Retrieve monthly pricing statistics for all categories
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponseCategoryMonthPrice "Monthly category pricing data"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/category/monthly-pricing [get]
func (h *categoryHandleApi) FindMonthPrice(c echo.Context) error {
	const method = "FindMonthPrice"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		logError("Invalid year parameter", err, zap.String("year", c.QueryParam("year")))
		return category_errors.ErrApiCategoryInvalidYear(c)
	}

	res, err := h.client.FindMonthPrice(ctx, &pb.FindYearCategory{
		Year: int32(year),
	})

	if err != nil {
		logError("Failed to retrieve category data", err, zap.Error(err))

		return category_errors.ErrApiCategoryFailedMonthPrice(c)
	}

	so := h.mapping.ToApiResponseCategoryMonthlyPrice(res)

	logSuccess("Successfully retrieve category data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearPrice retrieves yearly category pricing statistics
// @Summary Get yearly category pricing
// @Tags Category
// @Security Bearer
// @Description Retrieve yearly pricing statistics for all categories
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Success 200 {object} response.ApiResponseCategoryYearPrice "Yearly category pricing data"
// @Failure 400 {object} response.ErrorResponse "Invalid year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/category/yearly-pricing [get]
func (h *categoryHandleApi) FindYearPrice(c echo.Context) error {
	const method = "FindYearPrice"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)

	if err != nil {
		logError("Invalid year parameter", err, zap.String("year", c.QueryParam("year")))
		return category_errors.ErrApiCategoryInvalidYear(c)
	}

	res, err := h.client.FindYearPrice(ctx, &pb.FindYearCategory{
		Year: int32(year),
	})

	if err != nil {
		logError("Failed to retrieve category data", err, zap.Error(err))
		return category_errors.ErrApiCategoryFailedYearPrice(c)
	}

	so := h.mapping.ToApiResponseCategoryYearlyPrice(res)

	logSuccess("Successfully retrieve category data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthPriceByMerchant retrieves monthly category pricing by merchant
// @Summary Get monthly category pricing by merchant
// @Tags Category
// @Security Bearer
// @Description Retrieve monthly pricing statistics for categories by specific merchant
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponseCategoryMonthPrice "Monthly category pricing by merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Merchant not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/category/merchant/monthly-pricing [get]
func (h *categoryHandleApi) FindMonthPriceByMerchant(c echo.Context) error {
	const method = "FindMonthPriceByMerchant"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		logError("Invalid year parameter", err, zap.String("year", c.QueryParam("year")))
		return category_errors.ErrApiCategoryInvalidYear(c)
	}

	merchant_id, err := parseQueryIntWithValidation(c, "merchant_id", 1, 9999)
	if err != nil {
		logError("Invalid merchant ID parameter", err, zap.String("merchant_id", c.QueryParam("merchant_id")))
		return category_errors.ErrApiCategoryInvalidMerchantId(c)
	}

	res, err := h.client.FindMonthPriceByMerchant(ctx, &pb.FindYearCategoryByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchant_id),
	})

	if err != nil {
		logError("Failed to retrieve category data", err, zap.Error(err))

		return category_errors.ErrApiCategoryFailedMonthPriceByMerchant(c)
	}

	so := h.mapping.ToApiResponseCategoryMonthlyPrice(res)

	logSuccess("Successfully retrieve category data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearPriceByMerchant retrieves yearly category pricing by merchant
// @Summary Get yearly category pricing by merchant
// @Tags Category
// @Security Bearer
// @Description Retrieve yearly pricing statistics for categories by specific merchant
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param merchant_id query int true "Merchant ID"
// @Success 200 {object} response.ApiResponseCategoryYearPrice "Yearly category pricing by merchant"
// @Failure 400 {object} response.ErrorResponse "Invalid merchant ID or year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Merchant not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/category/merchant/yearly-pricing [get]
func (h *categoryHandleApi) FindYearPriceByMerchant(c echo.Context) error {
	const method = "FindYearPriceByMerchant"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		logError("Invalid year parameter", err, zap.String("year", c.QueryParam("year")))
		return category_errors.ErrApiCategoryInvalidYear(c)
	}

	merchant_id, err := parseQueryIntWithValidation(c, "merchant_id", 1, 9999)
	if err != nil {
		logError("Invalid merchant ID parameter", err, zap.String("merchant_id", c.QueryParam("merchant_id")))
		return category_errors.ErrApiCategoryInvalidMerchantId(c)
	}

	res, err := h.client.FindYearPriceByMerchant(ctx, &pb.FindYearCategoryByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchant_id),
	})

	if err != nil {
		logError("Failed to retrieve category data", err, zap.Error(err))

		return category_errors.ErrApiCategoryFailedYearPriceByMerchant(c)
	}

	so := h.mapping.ToApiResponseCategoryYearlyPrice(res)

	logSuccess("Successfully retrieve category data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindMonthPriceById retrieves monthly pricing for specific category
// @Summary Get monthly pricing by category ID
// @Tags Category
// @Security Bearer
// @Description Retrieve monthly pricing statistics for specific category
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param category_id path int true "Category ID"
// @Success 200 {object} response.ApiResponseCategoryMonthPrice "Monthly pricing by category"
// @Failure 400 {object} response.ErrorResponse "Invalid category ID or year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Category not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/category/mycategory/monthly-pricing [get]
func (h *categoryHandleApi) FindMonthPriceById(c echo.Context) error {
	const method = "FindMonthPriceById"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		logError("Invalid year parameter", err, zap.String("year", c.QueryParam("year")))
		return category_errors.ErrApiCategoryInvalidYear(c)
	}

	category_id, err := parseQueryIntWithValidation(c, "category_id", 1, 9999)
	if err != nil {
		logError("Invalid category ID parameter", err, zap.String("category_id", c.Param("category_id")))
		return category_errors.ErrApiCategoryInvalidId(c)
	}

	res, err := h.client.FindMonthPriceById(ctx, &pb.FindYearCategoryById{
		Year:       int32(year),
		CategoryId: int32(category_id),
	})

	if err != nil {
		logError("Failed to retrieve category data", err, zap.Error(err))

		return category_errors.ErrApiCategoryFailedMonthPriceById(c)
	}

	so := h.mapping.ToApiResponseCategoryMonthlyPrice(res)

	logSuccess("Successfully retrieve category data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// FindYearPriceById retrieves yearly pricing for specific category
// @Summary Get yearly pricing by category ID
// @Tags Category
// @Security Bearer
// @Description Retrieve yearly pricing statistics for specific category
// @Accept json
// @Produce json
// @Param year query int true "Year in YYYY format (e.g., 2023)"
// @Param category_id path int true "Category ID"
// @Success 200 {object} response.ApiResponseCategoryYearPrice "Yearly pricing by category"
// @Failure 400 {object} response.ErrorResponse "Invalid category ID or year parameter"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Category not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/category/mycategory/yearly-pricing [get]
func (h *categoryHandleApi) FindYearPriceById(c echo.Context) error {
	const method = "FindYearPriceById"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	year, err := parseQueryIntWithValidation(c, "year", 1, 9999)
	if err != nil {
		logError("Invalid year parameter", err, zap.String("year", c.QueryParam("year")))
		return category_errors.ErrApiCategoryInvalidYear(c)
	}

	category_id, err := parseQueryIntWithValidation(c, "category_id", 1, 9999)
	if err != nil {
		logError("Invalid category ID parameter", err, zap.String("category_id", c.Param("category_id")))
		return category_errors.ErrApiCategoryInvalidId(c)
	}

	res, err := h.client.FindYearPriceById(ctx, &pb.FindYearCategoryById{
		Year:       int32(year),
		CategoryId: int32(category_id),
	})

	if err != nil {
		logError("Failed to retrieve category data", err, zap.Error(err))

		return category_errors.ErrApiCategoryFailedYearPriceById(c)
	}

	so := h.mapping.ToApiResponseCategoryYearlyPrice(res)

	logSuccess("Successfully retrieve category data", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// Create handles the creation of a new category without image upload.
// @Summary Create a new category
// @Tags Category
// @Description Create a new category with the provided details
// @Accept json
// @Produce json
// @Param request body requests.CreateCategoryRequest true "Category details"
// @Success 201 {object} response.ApiResponseCategory "Successfully created category"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to create category"
// @Router /api/category/create [post]
func (h *categoryHandleApi) Create(c echo.Context) error {
	const method = "Create"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	var body requests.CreateCategoryRequest

	if err := c.Bind(&body); err != nil {
		logError("Failed to bind request body", err, zap.Error(err))

		return category_errors.ErrApiBindCreateCategory(c)
	}

	if err := body.Validate(); err != nil {
		logError("Invalid request body", err, zap.Error(err))

		return category_errors.ErrApiValidateCreateCategory(c)
	}

	req := &pb.CreateCategoryRequest{
		Name:        body.Name,
		Description: body.Description,
	}

	res, err := h.client.Create(ctx, req)

	if err != nil {
		logError("Failed to create category", err, zap.Error(err))

		return category_errors.ErrApiCategoryFailedCreate(c)
	}

	so := h.mapping.ToApiResponseCategory(res)

	logSuccess("Successfully created category", zap.Bool("success", true))

	return c.JSON(http.StatusCreated, so)
}

// @Security Bearer
// Update handles the update of an existing category.
// @Summary Update an existing category
// @Tags Category
// @Description Update an existing category record with the provided details
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param request body requests.UpdateCategoryRequest true "Category update details"
// @Success 200 {object} response.ApiResponseCategory "Successfully updated category"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to update category"
// @Router /api/category/update/{id} [post]
func (h *categoryHandleApi) Update(c echo.Context) error {
	const method = "Update"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		logError("Invalid category ID parameter", err, zap.Error(err))

		return category_errors.ErrApiCategoryInvalidId(c)
	}

	var body requests.UpdateCategoryRequest

	if err := c.Bind(&body); err != nil {
		logError("Failed to bind request body", err, zap.Error(err))

		return category_errors.ErrApiBindUpdateCategory(c)
	}

	if err := body.Validate(); err != nil {
		logError("Invalid request body", err, zap.Error(err))

		return category_errors.ErrApiValidateUpdateCategory(c)
	}

	req := &pb.UpdateCategoryRequest{
		CategoryId:  int32(idInt),
		Name:        body.Name,
		Description: body.Description,
	}

	res, err := h.client.Update(ctx, req)

	if err != nil {
		logError("Failed to update category", err, zap.Error(err))

		return category_errors.ErrApiCategoryFailedUpdate(c)
	}

	so := h.mapping.ToApiResponseCategory(res)

	logSuccess("Successfully updated category", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// TrashedCategory retrieves a trashed category record by its ID.
// @Summary Retrieve a trashed category
// @Tags Category
// @Description Retrieve a trashed category record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} response.ApiResponseCategoryDeleteAt "Successfully retrieved trashed category"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve trashed category"
// @Router /api/category/trashed/{id} [get]
func (h *categoryHandleApi) TrashedCategory(c echo.Context) error {
	const method = "TrashedCategory"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		logError("Failed to parse category id", err, zap.Error(err))

		return category_errors.ErrApiCategoryInvalidId(c)
	}

	req := &pb.FindByIdCategoryRequest{Id: int32(id)}

	res, err := h.client.TrashedCategory(ctx, req)

	if err != nil {
		logError("Failed to retrieve trashed category", err, zap.Error(err))

		return category_errors.ErrApiCategoryFailedTrashed(c)
	}

	so := h.mapping.ToApiResponseCategoryDeleteAt(res)

	logSuccess("Successfully retrieved trashed category", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreCategory restores a category record from the trash by its ID.
// @Summary Restore a trashed category
// @Tags Category
// @Description Restore a trashed category record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} response.ApiResponseCategoryDeleteAt "Successfully restored category"
// @Failure 400 {object} response.ErrorResponse "Invalid category ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore category"
// @Router /api/category/restore/{id} [post]
func (h *categoryHandleApi) RestoreCategory(c echo.Context) error {
	const method = "RestoreCategory"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		logError("Failed to parse category id", err, zap.Error(err))

		return category_errors.ErrApiCategoryInvalidId(c)
	}

	req := &pb.FindByIdCategoryRequest{Id: int32(id)}

	res, err := h.client.RestoreCategory(ctx, req)

	if err != nil {
		logError("Failed to restore category", err, zap.Error(err))

		return category_errors.ErrApiCategoryFailedRestore(c)
	}

	so := h.mapping.ToApiResponseCategoryDeleteAt(res)

	logSuccess("Successfully restored category", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteCategoryPermanent permanently deletes a category record by its ID.
// @Summary Permanently delete a category
// @Tags Category
// @Description Permanently delete a category record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "category ID"
// @Success 200 {object} response.ApiResponseCategoryDelete "Successfully deleted category record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete category:"
// @Router /api/category/delete/{id} [delete]
func (h *categoryHandleApi) DeleteCategoryPermanent(c echo.Context) error {
	const method = "DeleteCategoryPermanent"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		logError("Failed to parse category id", err, zap.Error(err))

		return category_errors.ErrApiCategoryInvalidId(c)
	}

	req := &pb.FindByIdCategoryRequest{Id: int32(id)}

	res, err := h.client.DeleteCategoryPermanent(ctx, req)

	if err != nil {
		logError("Failed to delete category", err, zap.Error(err))

		return category_errors.ErrApiCategoryFailedDeletePermanent(c)
	}

	so := h.mapping.ToApiResponseCategoryDelete(res)

	logSuccess("Successfully deleted category", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// RestoreAllCategory restores a category record from the trash by its ID.
// @Summary Restore a trashed category
// @Tags Category
// @Description Restore a trashed category record by its ID.
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseCategoryAll "Successfully restored category all"
// @Failure 400 {object} response.ErrorResponse "Invalid category ID"
// @Failure 500 {object} response.ErrorResponse "Failed to restore category"
// @Router /api/category/restore/all [post]
func (h *categoryHandleApi) RestoreAllCategory(c echo.Context) error {
	const method = "RestoreAllCategory"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	res, err := h.client.RestoreAllCategory(ctx, &emptypb.Empty{})

	if err != nil {
		logError("Failed to restore all category", err, zap.Error(err))

		return category_errors.ErrApiCategoryFailedRestoreAll(c)
	}

	so := h.mapping.ToApiResponseCategoryAll(res)

	logSuccess("Successfully restored all category", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

// @Security Bearer
// DeleteAllCategoryPermanent permanently deletes a category record by its ID.
// @Summary Permanently delete a category
// @Tags Category
// @Description Permanently delete a category record by its ID.
// @Accept json
// @Produce json
// @Param id path int true "category ID"
// @Success 200 {object} response.ApiResponseCategoryAll "Successfully deleted category record permanently"
// @Failure 400 {object} response.ErrorResponse "Bad Request: Invalid ID"
// @Failure 500 {object} response.ErrorResponse "Failed to delete category:"
// @Router /api/category/delete/all [post]
func (h *categoryHandleApi) DeleteAllCategoryPermanent(c echo.Context) error {
	const method = "DeleteAllCategoryPermanent"

	ctx := c.Request().Context()

	end, logSuccess, logError := h.startTracingAndLogging(ctx, method)

	defer func() { end() }()

	res, err := h.client.DeleteAllCategoryPermanent(ctx, &emptypb.Empty{})

	if err != nil {
		logError("Failed to delete all category", err, zap.Error(err))

		return category_errors.ErrApiCategoryFailedDeleteAllPermanent(c)
	}

	so := h.mapping.ToApiResponseCategoryAll(res)

	logSuccess("Successfully deleted all category", zap.Bool("success", true))

	return c.JSON(http.StatusOK, so)
}

func (s *categoryHandleApi) startTracingAndLogging(
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

func (s *categoryHandleApi) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method, status).Observe(time.Since(start).Seconds())
}
