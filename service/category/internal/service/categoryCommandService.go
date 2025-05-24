package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-category/internal/repository"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	traceunic "github.com/MamangRust/monolith-point-of-sale-pkg/trace_unic"
	"github.com/MamangRust/monolith-point-of-sale-pkg/utils"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/category_errors"
	response_service "github.com/MamangRust/monolith-point-of-sale-shared/mapper/response/service"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type categoryCommandService struct {
	ctx                       context.Context
	trace                     trace.Tracer
	categoryQueryRepository   repository.CategoryQueryRepository
	categoryCommandRepository repository.CategoryCommandRepository
	logger                    logger.LoggerInterface
	mapping                   response_service.CategoryResponseMapper
	requestCounter            *prometheus.CounterVec
	requestDuration           *prometheus.HistogramVec
}

func NewCategoryCommandService(ctx context.Context, categoryCommandRepository repository.CategoryCommandRepository,
	categoryQueryRepository repository.CategoryQueryRepository,
	logger logger.LoggerInterface, mapping response_service.CategoryResponseMapper) *categoryCommandService {
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "category_command_service_request_total",
			Help: "Total number of requests to the CategoryCommandService",
		},
		[]string{"method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "category_command_service_request_duration_seconds",
			Help:    "Duration of requests to the CategoryCommandService",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	prometheus.MustRegister(requestCounter, requestDuration)

	return &categoryCommandService{
		ctx:                       ctx,
		trace:                     otel.Tracer("category-command-service"),
		categoryCommandRepository: categoryCommandRepository,
		categoryQueryRepository:   categoryQueryRepository,
		logger:                    logger,
		mapping:                   mapping,
		requestCounter:            requestCounter,
		requestDuration:           requestDuration,
	}
}

func (s *categoryCommandService) CreateCategory(req *requests.CreateCategoryRequest) (*response.CategoryResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("CreateCategory", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "CreateCategory")
	defer span.End()

	s.logger.Debug("Creating new category")

	slug := utils.GenerateSlug(req.Name)

	req.Name = slug

	cashier, err := s.categoryCommandRepository.CreateCategory(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_CREATE_CATEGORY")

		s.logger.Error("Failed to create category",
			zap.Error(err),
			zap.Any("request", req),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to create category")

		status = "failed_create_category"

		return nil, category_errors.ErrFailedCreateCategory
	}

	return s.mapping.ToCategoryResponse(cashier), nil
}

func (s *categoryCommandService) UpdateCategory(req *requests.UpdateCategoryRequest) (*response.CategoryResponse, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("UpdateCategory", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "UpdateCategory")
	defer span.End()

	s.logger.Debug("Updating category",
		zap.Int("category_id", *req.CategoryID),
		zap.String("category_name", req.Name))

	span.SetAttributes(
		attribute.Int("category.id", *req.CategoryID),
		attribute.String("category.name", req.Name),
	)

	slug := utils.GenerateSlug(req.Name)
	req.SlugCategory = &slug

	_, err := s.categoryQueryRepository.FindByNameAndId(&requests.CategoryNameAndId{
		Name:       req.Name,
		CategoryID: *req.CategoryID,
	})

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_FIND_CATEGORY_BY_NAME")

		s.logger.Error("Error retrieving category",
			zap.Error(err),
			zap.Int("category_id", *req.CategoryID),
			zap.String("category_name", req.Name),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.Bool("error.exists_check", true),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to find category by name")

		status = "failed_find_category_by_name"

		return nil, category_errors.ErrFailedFindCategoryByName
	}

	category, err := s.categoryCommandRepository.UpdateCategory(req)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_UPDATE_CATEGORY")

		s.logger.Error("Failed to update category",
			zap.Error(err),
			zap.Any("request", req),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("error.trace_id", traceID),
			attribute.Bool("error.update_operation", true),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to update category")

		status = "failed_update_category"

		return nil, category_errors.ErrFailedUpdateCategory
	}

	span.SetAttributes(
		attribute.String("category.new_slug", slug),
		attribute.Bool("operation.success", true),
	)

	return s.mapping.ToCategoryResponse(category), nil
}

func (s *categoryCommandService) TrashedCategory(category_id int) (*response.CategoryResponseDeleteAt, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("TrashedCategory", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "TrashedCategory")
	defer span.End()

	s.logger.Debug("Trashing category", zap.Int("category", category_id))

	category, err := s.categoryCommandRepository.TrashedCategory(category_id)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_TRASHED_CATEGORY")

		s.logger.Error("Failed to move category to trash",
			zap.Error(err),
			zap.Int("category_id", category_id),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
			attribute.Int("category_id", category_id),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to move category to trash")

		status = "failed_trashed_category"

		return nil, category_errors.ErrFailedTrashedCategory
	}

	span.SetAttributes(attribute.Int("category_id", category_id))
	return s.mapping.ToCategoryResponseDeleteAt(category), nil
}

func (s *categoryCommandService) RestoreCategory(categoryID int) (*response.CategoryResponseDeleteAt, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreCategory", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreCategory")
	defer span.End()

	s.logger.Debug("Restoring category", zap.Int("categoryID", categoryID))

	category, err := s.categoryCommandRepository.RestoreCategory(categoryID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_CATEGORY")

		s.logger.Error("Failed to restore category from trash",
			zap.Error(err),
			zap.Int("category_id", categoryID),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
			attribute.Int("category_id", categoryID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore category from trash")

		status = "failed_restore_category"

		return nil, category_errors.ErrFailedRestoreCategory
	}

	span.SetAttributes(attribute.Int("category_id", categoryID))
	return s.mapping.ToCategoryResponseDeleteAt(category), nil
}

func (s *categoryCommandService) DeleteCategoryPermanent(categoryID int) (bool, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteCategoryPermanent", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteCategoryPermanent")
	defer span.End()

	s.logger.Debug("Permanently deleting category", zap.Int("categoryID", categoryID))

	success, err := s.categoryCommandRepository.DeleteCategoryPermanently(categoryID)

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_CATEGORY_PERMANENT")

		s.logger.Error("Failed to permanently delete category",
			zap.Error(err),
			zap.Int("category_id", categoryID),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
			attribute.Int("category_id", categoryID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to permanently delete category")

		status = "failed_delete_category_permanent"

		return false, category_errors.ErrFailedDeleteCategoryPermanent
	}

	span.SetAttributes(attribute.Int("category_id", categoryID))
	return success, nil
}

func (s *categoryCommandService) RestoreAllCategories() (bool, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("RestoreAllCategories", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "RestoreAllCategories")
	defer span.End()

	s.logger.Debug("Restoring all trashed categories")

	success, err := s.categoryCommandRepository.RestoreAllCategories()
	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_RESTORE_ALL_CATEGORIES")

		s.logger.Error("Failed to restore all trashed categories",
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to restore all trashed categories")

		status = "failed_restore_all_categories"

		return false, category_errors.ErrFailedRestoreAllCategories
	}

	return success, nil
}

func (s *categoryCommandService) DeleteAllCategoriesPermanent() (bool, *response.ErrorResponse) {
	startTime := time.Now()
	status := "success"

	defer func() {
		s.recordMetrics("DeleteAllCategoriesPermanent", status, startTime)
	}()

	_, span := s.trace.Start(s.ctx, "DeleteAllCategoriesPermanent")
	defer span.End()

	s.logger.Debug("Permanently deleting all categories")

	success, err := s.categoryCommandRepository.DeleteAllPermanentCategories()

	if err != nil {
		traceID := traceunic.GenerateTraceID("FAILED_DELETE_ALL_CATEGORIES_PERMANENT")

		s.logger.Error("Failed to permanently delete all trashed categories",
			zap.Error(err),
			zap.String("traceID", traceID))

		span.SetAttributes(
			attribute.String("traceID", traceID),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to permanently delete all trashed categories")

		status = "failed_delete_all_categories_permanent"

		return false, category_errors.ErrFailedDeleteAllCategoriesPermanent
	}

	return success, nil
}

func (s *categoryCommandService) recordMetrics(method string, status string, start time.Time) {
	s.requestCounter.WithLabelValues(method, status).Inc()
	s.requestDuration.WithLabelValues(method).Observe(time.Since(start).Seconds())
}
