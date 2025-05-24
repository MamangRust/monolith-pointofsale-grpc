package handler

import (
	"net/http"
	"strconv"

	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/category_errors"
	response_api "github.com/MamangRust/monolith-point-of-sale-shared/mapper/response/api"
	"github.com/MamangRust/monolith-point-of-sale-shared/pb"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type categoryHandleApi struct {
	client  pb.CategoryServiceClient
	logger  logger.LoggerInterface
	mapping response_api.CategoryResponseMapper
}

func NewHandlerCategory(
	router *echo.Echo,
	client pb.CategoryServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.CategoryResponseMapper,
) *categoryHandleApi {
	categoryHandler := &categoryHandleApi{
		client:  client,
		logger:  logger,
		mapping: mapping,
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

	req := &pb.FindAllCategoryRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch categories", zap.Error(err))
		return category_errors.ErrApiCategoryFailedFindAll(c)
	}

	so := h.mapping.ToApiResponsePaginationCategory(res)

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

	req := &pb.FindAllCategoryRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch active categories", zap.Error(err))
		return category_errors.ErrApiCategoryFailedFindByActive(c)
	}

	so := h.mapping.ToApiResponsePaginationCategoryDeleteAt(res)

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

	req := &pb.FindAllCategoryRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, req)

	if err != nil {
		h.logger.Error("Failed to fetch archived categories", zap.Error(err))
		return category_errors.ErrApiCategoryFailedFindByTrashed(c)
	}

	so := h.mapping.ToApiResponsePaginationCategoryDeleteAt(res)

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
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))

		return category_errors.ErrApiCategoryInvalidYear(c)
	}

	monthStr := c.QueryParam("month")

	month, err := strconv.Atoi(monthStr)

	if err != nil {
		h.logger.Debug("Invalid month parameter", zap.Error(err))

		return category_errors.ErrApiCategoryInvalidMonth(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthlyTotalPrices(ctx, &pb.FindYearMonthTotalPrices{
		Year:  int32(year),
		Month: int32(month),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve monthly category price", zap.Error(err))

		return category_errors.ErrApiCategoryFailedMonthTotalPrice(c)
	}

	so := h.mapping.ToApiResponseCategoryMonthlyTotalPrice(res)

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
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))

		return category_errors.ErrApiCategoryInvalidYear(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearlyTotalPrices(ctx, &pb.FindYearTotalPrices{
		Year: int32(year),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve yearly category price", zap.Error(err))

		return category_errors.ErrApiCategoryFailedYearTotalPrice(c)
	}

	so := h.mapping.ToApiResponseCategoryYearlyTotalPrice(res)

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
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)

	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidYear(c)
	}

	monthStr := c.QueryParam("month")

	month, err := strconv.Atoi(monthStr)

	if err != nil {
		h.logger.Debug("Invalid month parameter", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidMonth(c)
	}

	categoryStr := c.QueryParam("category_id")

	category, err := strconv.Atoi(categoryStr)

	if err != nil {
		h.logger.Debug("Invalid category id parameter", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidId(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthlyTotalPricesById(ctx, &pb.FindYearMonthTotalPriceById{
		Year:       int32(year),
		Month:      int32(month),
		CategoryId: int32(category),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve monthly category price", zap.Error(err))

		return category_errors.ErrApiCategoryFailedMonthTotalPriceById(c)
	}

	so := h.mapping.ToApiResponseCategoryMonthlyTotalPrice(res)

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
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidYear(c)
	}

	categoryStr := c.QueryParam("category_id")

	category, err := strconv.Atoi(categoryStr)

	if err != nil {
		h.logger.Debug("Invalid category id parameter", zap.Error(err))

		return category_errors.ErrApiCategoryInvalidId(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearlyTotalPricesById(ctx, &pb.FindYearTotalPriceById{
		Year:       int32(year),
		CategoryId: int32(category),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve yearly category price", zap.Error(err))

		return category_errors.ErrApiCategoryFailedYearTotalPriceById(c)
	}

	so := h.mapping.ToApiResponseCategoryYearlyTotalPrice(res)

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
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)

	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidYear(c)
	}

	monthStr := c.QueryParam("month")

	month, err := strconv.Atoi(monthStr)

	if err != nil {
		h.logger.Debug("Invalid month parameter", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidMonth(c)
	}

	merchantStr := c.QueryParam("merchant_id")

	merchant, err := strconv.Atoi(merchantStr)

	if err != nil {
		h.logger.Debug("Invalid merchant id parameter", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidMerchantId(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthlyTotalPricesByMerchant(ctx, &pb.FindYearMonthTotalPriceByMerchant{
		Year:       int32(year),
		Month:      int32(month),
		MerchantId: int32(merchant),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve monthly category price", zap.Error(err))

		return category_errors.ErrApiCategoryFailedMonthTotalPriceByMerchant(c)
	}

	so := h.mapping.ToApiResponseCategoryMonthlyTotalPrice(res)

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
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidYear(c)
	}

	merchantStr := c.QueryParam("merchant_id")

	merchant, err := strconv.Atoi(merchantStr)

	if err != nil {
		h.logger.Debug("Invalid merchant id parameter", zap.Error(err))

		return category_errors.ErrApiCategoryInvalidMerchantId(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearlyTotalPricesByMerchant(ctx, &pb.FindYearTotalPriceByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchant),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve yearly category price", zap.Error(err))

		return category_errors.ErrApiCategoryFailedYearTotalPriceByMerchant(c)
	}

	so := h.mapping.ToApiResponseCategoryYearlyTotalPrice(res)

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
	yearStr := c.QueryParam("year")

	year, err := strconv.Atoi(yearStr)

	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidYear(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthPrice(ctx, &pb.FindYearCategory{
		Year: int32(year),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve monthly category price", zap.Error(err))
		return category_errors.ErrApiCategoryFailedMonthPrice(c)
	}

	so := h.mapping.ToApiResponseCategoryMonthlyPrice(res)

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
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidYear(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearPrice(ctx, &pb.FindYearCategory{
		Year: int32(year),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly category price", zap.Error(err))
		return category_errors.ErrApiCategoryFailedYearPrice(c)
	}

	so := h.mapping.ToApiResponseCategoryYearlyPrice(res)

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
	yearStr := c.QueryParam("year")
	merchantIdStr := c.QueryParam("merchant_id")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidYear(c)
	}

	merchant_id, err := strconv.Atoi(merchantIdStr)

	if err != nil {
		h.logger.Debug("Invalid merchant id parameter", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidMerchantId(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthPriceByMerchant(ctx, &pb.FindYearCategoryByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchant_id),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve monthly category price", zap.Error(err))
		return category_errors.ErrApiCategoryFailedMonthPriceByMerchant(c)
	}

	so := h.mapping.ToApiResponseCategoryMonthlyPrice(res)

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
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)

	merchantIdStr := c.QueryParam("merchant_id")

	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidYear(c)
	}

	merchant_id, err := strconv.Atoi(merchantIdStr)

	if err != nil {
		h.logger.Debug("Invalid merchant id parameter", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidMerchantId(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearPriceByMerchant(ctx, &pb.FindYearCategoryByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchant_id),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly category price", zap.Error(err))
		return category_errors.ErrApiCategoryFailedYearPriceByMerchant(c)
	}

	so := h.mapping.ToApiResponseCategoryYearlyPrice(res)

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
	yearStr := c.QueryParam("year")
	categoryIdStr := c.QueryParam("category_id")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidYear(c)
	}

	category_id, err := strconv.Atoi(categoryIdStr)

	if err != nil {
		h.logger.Debug("Invalid category id parameter", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidId(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindMonthPriceById(ctx, &pb.FindYearCategoryById{
		Year:       int32(year),
		CategoryId: int32(category_id),
	})

	if err != nil {
		h.logger.Debug("Failed to retrieve monthly category price", zap.Error(err))
		return category_errors.ErrApiCategoryFailedMonthPriceById(c)
	}

	so := h.mapping.ToApiResponseCategoryMonthlyPrice(res)

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
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)

	categoryIdStr := c.QueryParam("category_id")

	if err != nil {
		h.logger.Debug("Invalid year parameter", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidYear(c)
	}

	category_id, err := strconv.Atoi(categoryIdStr)

	if err != nil {
		h.logger.Debug("Invalid category id parameter", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidId(c)
	}

	ctx := c.Request().Context()

	res, err := h.client.FindYearPriceById(ctx, &pb.FindYearCategoryById{
		Year:       int32(year),
		CategoryId: int32(category_id),
	})
	if err != nil {
		h.logger.Debug("Failed to retrieve yearly category price", zap.Error(err))
		return category_errors.ErrApiCategoryFailedYearPriceById(c)
	}

	so := h.mapping.ToApiResponseCategoryYearlyPrice(res)

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
	var body requests.CreateCategoryRequest

	if err := c.Bind(&body); err != nil {
		h.logger.Debug("Invalid request format", zap.Error(err))
		return category_errors.ErrApiBindCreateCategory(c)
	}

	if err := body.Validate(); err != nil {
		h.logger.Debug("Validation failed", zap.Error(err))
		return category_errors.ErrApiValidateCreateCategory(c)
	}

	ctx := c.Request().Context()

	req := &pb.CreateCategoryRequest{
		Name:        body.Name,
		Description: body.Description,
	}

	h.logger.Debug("info", zap.Any("req", req))

	res, err := h.client.Create(ctx, req)
	if err != nil {
		h.logger.Error("Category creation failed",
			zap.Error(err),
			zap.Any("request", req),
		)

		return category_errors.ErrApiCategoryFailedCreate(c)
	}

	so := h.mapping.ToApiResponseCategory(res)

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
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		h.logger.Debug("Invalid id parameter", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidId(c)
	}

	var body requests.UpdateCategoryRequest

	if err := c.Bind(&body); err != nil {
		return category_errors.ErrApiBindUpdateCategory(c)
	}

	if err := body.Validate(); err != nil {
		h.logger.Debug("Validation failed", zap.Error(err))
		return category_errors.ErrApiValidateUpdateCategory(c)
	}

	ctx := c.Request().Context()

	req := &pb.UpdateCategoryRequest{
		CategoryId:  int32(idInt),
		Name:        body.Name,
		Description: body.Description,
	}

	res, err := h.client.Update(ctx, req)
	if err != nil {
		h.logger.Error("Category update failed",
			zap.Error(err),
			zap.Any("request", req),
		)

		return category_errors.ErrApiCategoryFailedUpdate(c)
	}

	so := h.mapping.ToApiResponseCategory(res)

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
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid category ID format", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidId(c)
	}

	ctx := c.Request().Context()
	req := &pb.FindByIdCategoryRequest{Id: int32(id)}

	res, err := h.client.TrashedCategory(ctx, req)

	if err != nil {
		h.logger.Error("Failed to archive category", zap.Error(err))
		return category_errors.ErrApiCategoryFailedTrashed(c)
	}

	so := h.mapping.ToApiResponseCategoryDeleteAt(res)

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
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Debug("Invalid category ID format", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidId(c)
	}

	ctx := c.Request().Context()
	req := &pb.FindByIdCategoryRequest{Id: int32(id)}

	res, err := h.client.RestoreCategory(ctx, req)
	if err != nil {
		h.logger.Error("Failed to restore category", zap.Error(err))
		return category_errors.ErrApiCategoryFailedRestore(c)
	}

	so := h.mapping.ToApiResponseCategoryDeleteAt(res)

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
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.logger.Debug("Invalid category ID format", zap.Error(err))
		return category_errors.ErrApiCategoryInvalidId(c)
	}

	ctx := c.Request().Context()
	req := &pb.FindByIdCategoryRequest{Id: int32(id)}

	res, err := h.client.DeleteCategoryPermanent(ctx, req)
	if err != nil {
		h.logger.Error("Failed to permanently delete category", zap.Error(err))
		return category_errors.ErrApiCategoryFailedDeletePermanent(c)
	}

	so := h.mapping.ToApiResponseCategoryDelete(res)

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
	ctx := c.Request().Context()

	res, err := h.client.RestoreAllCategory(ctx, &emptypb.Empty{})
	if err != nil {
		h.logger.Error("Bulk category restoration failed", zap.Error(err))
		return category_errors.ErrApiCategoryFailedRestoreAll(c)
	}

	so := h.mapping.ToApiResponseCategoryAll(res)

	h.logger.Debug("Successfully restored all category")

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
	ctx := c.Request().Context()
	res, err := h.client.DeleteAllCategoryPermanent(ctx, &emptypb.Empty{})
	if err != nil {
		h.logger.Error("Bulk category deletion failed", zap.Error(err))
		return category_errors.ErrApiCategoryFailedDeleteAllPermanent(c)
	}

	h.logger.Debug("All categories permanently deleted")

	so := h.mapping.ToApiResponseCategoryAll(res)
	return c.JSON(http.StatusOK, so)
}
