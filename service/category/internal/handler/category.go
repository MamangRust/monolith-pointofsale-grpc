package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-point-of-sale-category/internal/service"
	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/category_errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/pb"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type categoryHandleGrpc struct {
	pb.UnimplementedCategoryServiceServer
	categoryQuery           service.CategoryQueryService
	categoryCommand         service.CategoryCommandService
	categoryStats           service.CategoryStatsService
	categoryStatsById       service.CategoryStatsByIdService
	categoryStatsByMerchant service.CategoryStatsByMerchantService
	logger                  logger.LoggerInterface
}

func NewCategoryHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pb.CategoryServiceServer {
	return &categoryHandleGrpc{
		categoryQuery:           service.CategoryQuery,
		categoryCommand:         service.CategoryCommand,
		categoryStats:           service.CategoryStats,
		categoryStatsById:       service.CategoryStatsById,
		categoryStatsByMerchant: service.CategoryStatsByMerchant,
		logger:                  logger,
	}
}

func (s *categoryHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllCategoryRequest) (*pb.ApiResponsePaginationCategory, error) {
	s.logger.Info("FindAll categories called", zap.Int32("page", request.GetPage()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllCategory{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	category, totalRecords, err := s.categoryQuery.FindAll(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindAll categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindAll categories success", zap.Int("count", len(category)))

	return &pb.ApiResponsePaginationCategory{
		Status:     "success",
		Message:    "Successfully fetched categories",
		Data:       mapResponsesCategory(category),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *categoryHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdCategoryRequest) (*pb.ApiResponseCategory, error) {
	s.logger.Info("FindById category called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id == 0 {
		return nil, category_errors.ErrGrpcFailedInvalidId
	}

	category, err := s.categoryQuery.FindById(ctx, id)
	if err != nil {
		s.logger.Error("FindById category failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindById category success", zap.Int("id", id))

	return &pb.ApiResponseCategory{
		Status:  "success",
		Message: "Successfully fetched category",
		Data:    mapResponseCategory(category),
	}, nil
}

func (s *categoryHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllCategoryRequest) (*pb.ApiResponsePaginationCategoryDeleteAt, error) {
	s.logger.Info("FindByActive categories called", zap.Int32("page", request.GetPage()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllCategory{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	categories, totalRecords, err := s.categoryQuery.FindByActive(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByActive categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByActive categories success")

	return &pb.ApiResponsePaginationCategoryDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active categories",
		Data:       mapResponsesCategoryActive(categories),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *categoryHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllCategoryRequest) (*pb.ApiResponsePaginationCategoryDeleteAt, error) {
	s.logger.Info("FindByTrashed categories called")

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllCategory{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	categories, totalRecords, err := s.categoryQuery.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByTrashed categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByTrashed categories success")

	return &pb.ApiResponsePaginationCategoryDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed categories",
		Data:       mapResponsesCategoryTrashed(categories),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *categoryHandleGrpc) FindMonthlyTotalPrices(ctx context.Context, req *pb.FindYearMonthTotalPrices) (*pb.ApiResponseCategoryMonthlyTotalPrice, error) {
	s.logger.Info("FindMonthlyTotalPrices categories called", zap.Int32("year", req.GetYear()))

	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidYear
	}
	if month <= 0 || month > 12 {
		return nil, category_errors.ErrGrpcFailedInvalidMonth
	}

	reqService := requests.MonthTotalPrice{
		Year:  year,
		Month: month,
	}

	methods, err := s.categoryStats.FindMonthlyTotalPrice(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthlyTotalPrices categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthlyTotalPrices categories success")

	return &pb.ApiResponseCategoryMonthlyTotalPrice{
		Status:  "success",
		Message: "Monthly sales retrieved successfully",
		Data:    mapResponseCategoryMonthlyTotalPrices(methods),
	}, nil
}

func (s *categoryHandleGrpc) FindYearlyTotalPrices(ctx context.Context, req *pb.FindYearTotalPrices) (*pb.ApiResponseCategoryYearlyTotalPrice, error) {
	s.logger.Info("FindYearlyTotalPrices categories called", zap.Int32("year", req.GetYear()))

	year := int(req.GetYear())
	if year <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidYear
	}

	methods, err := s.categoryStats.FindYearlyTotalPrice(ctx, year)
	if err != nil {
		s.logger.Error("FindYearlyTotalPrices categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearlyTotalPrices categories success")

	return &pb.ApiResponseCategoryYearlyTotalPrice{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    mapResponseCategoryYearlyTotalPrices(methods),
	}, nil
}

func (s *categoryHandleGrpc) FindMonthlyTotalPricesById(ctx context.Context, req *pb.FindYearMonthTotalPriceById) (*pb.ApiResponseCategoryMonthlyTotalPrice, error) {
	s.logger.Info("FindMonthlyTotalPricesById categories called", zap.Int32("id", req.GetCategoryId()))

	year := int(req.GetYear())
	month := int(req.GetMonth())
	id := int(req.GetCategoryId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidYear
	}
	if month <= 0 || month > 12 {
		return nil, category_errors.ErrGrpcFailedInvalidMonth
	}
	if id <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidId
	}

	reqService := requests.MonthTotalPriceCategory{
		Year:       year,
		Month:      month,
		CategoryID: id,
	}

	methods, err := s.categoryStatsById.FindMonthlyTotalPriceById(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthlyTotalPricesById categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthlyTotalPricesById categories success")

	return &pb.ApiResponseCategoryMonthlyTotalPrice{
		Status:  "success",
		Message: "Monthly sales retrieved successfully",
		Data:    mapResponseCategoryMonthlyTotalPricesById(methods),
	}, nil
}

func (s *categoryHandleGrpc) FindYearlyTotalPricesById(ctx context.Context, req *pb.FindYearTotalPriceById) (*pb.ApiResponseCategoryYearlyTotalPrice, error) {
	s.logger.Info("FindYearlyTotalPricesById categories called", zap.Int32("id", req.GetCategoryId()))

	year := int(req.GetYear())
	id := int(req.GetCategoryId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidYear
	}
	if id <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidId
	}

	reqService := requests.YearTotalPriceCategory{
		Year:       year,
		CategoryID: id,
	}

	methods, err := s.categoryStatsById.FindYearlyTotalPriceById(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindYearlyTotalPricesById categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearlyTotalPricesById categories success")

	return &pb.ApiResponseCategoryYearlyTotalPrice{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    mapResponseCategoryYearlyTotalPricesById(methods),
	}, nil
}

func (s *categoryHandleGrpc) FindMonthlyTotalPricesByMerchant(ctx context.Context, req *pb.FindYearMonthTotalPriceByMerchant) (*pb.ApiResponseCategoryMonthlyTotalPrice, error) {
	s.logger.Info("FindMonthlyTotalPricesByMerchant categories called", zap.Int32("merchantId", req.GetMerchantId()))

	year := int(req.GetYear())
	month := int(req.GetMonth())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidYear
	}
	if month <= 0 || month > 12 {
		return nil, category_errors.ErrGrpcFailedInvalidMonth
	}
	if id <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidMerchantId
	}

	reqService := requests.MonthTotalPriceMerchant{
		Year:       year,
		Month:      month,
		MerchantID: id,
	}

	methods, err := s.categoryStatsByMerchant.FindMonthlyTotalPriceByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthlyTotalPricesByMerchant categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthlyTotalPricesByMerchant categories success")

	return &pb.ApiResponseCategoryMonthlyTotalPrice{
		Status:  "success",
		Message: "Monthly sales retrieved successfully",
		Data:    mapResponseCategoryMonthlyTotalPricesByMerchant(methods),
	}, nil
}

func (s *categoryHandleGrpc) FindYearlyTotalPricesByMerchant(ctx context.Context, req *pb.FindYearTotalPriceByMerchant) (*pb.ApiResponseCategoryYearlyTotalPrice, error) {
	s.logger.Info("FindYearlyTotalPricesByMerchant categories called", zap.Int32("merchantId", req.GetMerchantId()))

	year := int(req.GetYear())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidYear
	}
	if id <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidMerchantId
	}

	reqService := requests.YearTotalPriceMerchant{
		Year:       year,
		MerchantID: id,
	}

	methods, err := s.categoryStatsByMerchant.FindYearlyTotalPriceByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindYearlyTotalPricesByMerchant categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearlyTotalPricesByMerchant categories success")

	return &pb.ApiResponseCategoryYearlyTotalPrice{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    mapResponseCategoryYearlyTotalPricesByMerchant(methods),
	}, nil
}

func (s *categoryHandleGrpc) FindMonthPrice(ctx context.Context, req *pb.FindYearCategory) (*pb.ApiResponseCategoryMonthPrice, error) {
	s.logger.Info("FindMonthPrice categories called", zap.Int32("year", req.GetYear()))

	year := int(req.GetYear())
	if year <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidYear
	}

	methods, err := s.categoryStats.FindMonthPrice(ctx, year)
	if err != nil {
		s.logger.Error("FindMonthPrice categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthPrice categories success")

	return &pb.ApiResponseCategoryMonthPrice{
		Status:  "success",
		Message: "Monthly payment methods retrieved successfully",
		Data:    mapResponsesCategoryMonthlyPrices(methods),
	}, nil
}

func (s *categoryHandleGrpc) FindYearPrice(ctx context.Context, req *pb.FindYearCategory) (*pb.ApiResponseCategoryYearPrice, error) {
	s.logger.Info("FindYearPrice categories called", zap.Int32("year", req.GetYear()))

	year := int(req.GetYear())
	if year <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidYear
	}

	methods, err := s.categoryStats.FindYearPrice(ctx, year)
	if err != nil {
		s.logger.Error("FindYearPrice categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearPrice categories success")

	return &pb.ApiResponseCategoryYearPrice{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    mapResponsesCategoryYearlyPrices(methods),
	}, nil
}

func (s *categoryHandleGrpc) FindMonthPriceByMerchant(ctx context.Context, req *pb.FindYearCategoryByMerchant) (*pb.ApiResponseCategoryMonthPrice, error) {
	s.logger.Info("FindMonthPriceByMerchant categories called", zap.Int32("merchantId", req.GetMerchantId()))

	year := int(req.GetYear())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidYear
	}
	if id <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidMerchantId
	}

	reqService := requests.MonthPriceMerchant{
		Year:       year,
		MerchantID: id,
	}

	methods, err := s.categoryStatsByMerchant.FindMonthPriceByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthPriceByMerchant categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthPriceByMerchant categories success")

	return &pb.ApiResponseCategoryMonthPrice{
		Status:  "success",
		Message: "Merchant monthly payment methods retrieved successfully",
		Data:    mapResponsesCategoryMonthlyPricesByMerchant(methods),
	}, nil
}

func (s *categoryHandleGrpc) FindYearPriceByMerchant(ctx context.Context, req *pb.FindYearCategoryByMerchant) (*pb.ApiResponseCategoryYearPrice, error) {
	s.logger.Info("FindYearPriceByMerchant categories called", zap.Int32("merchantId", req.GetMerchantId()))

	year := int(req.GetYear())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidYear
	}
	if id <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidMerchantId
	}

	reqService := requests.YearPriceMerchant{
		Year:       year,
		MerchantID: id,
	}

	methods, err := s.categoryStatsByMerchant.FindYearPriceByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindYearPriceByMerchant categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearPriceByMerchant categories success")

	return &pb.ApiResponseCategoryYearPrice{
		Status:  "success",
		Message: "Merchant yearly payment methods retrieved successfully",
		Data:    mapResponsesCategoryYearlyPricesByMerchant(methods),
	}, nil
}

func (s *categoryHandleGrpc) FindMonthPriceById(ctx context.Context, req *pb.FindYearCategoryById) (*pb.ApiResponseCategoryMonthPrice, error) {
	s.logger.Info("FindMonthPriceById categories called", zap.Int32("id", req.GetCategoryId()))

	year := int(req.GetYear())
	id := int(req.GetCategoryId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidYear
	}
	if id <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidId
	}

	reqService := requests.MonthPriceId{
		Year:       year,
		CategoryID: id,
	}

	methods, err := s.categoryStatsById.FindMonthPriceById(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthPriceById categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthPriceById categories success")

	return &pb.ApiResponseCategoryMonthPrice{
		Status:  "success",
		Message: "Merchant monthly payment methods retrieved successfully",
		Data:    mapResponsesCategoryMonthlyPricesById(methods),
	}, nil
}

func (s *categoryHandleGrpc) FindYearPriceById(ctx context.Context, req *pb.FindYearCategoryById) (*pb.ApiResponseCategoryYearPrice, error) {
	s.logger.Info("FindYearPriceById categories called", zap.Int32("id", req.GetCategoryId()))

	year := int(req.GetYear())
	id := int(req.GetCategoryId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidYear
	}
	if id <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidId
	}

	reqService := requests.YearPriceId{
		Year:       year,
		CategoryID: id,
	}

	methods, err := s.categoryStatsById.FindYearPriceById(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindYearPriceById categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearPriceById categories success")

	return &pb.ApiResponseCategoryYearPrice{
		Status:  "success",
		Message: "Merchant yearly payment methods retrieved successfully",
		Data:    mapResponsesCategoryYearlyPricesById(methods),
	}, nil
}

func (s *categoryHandleGrpc) Create(ctx context.Context, request *pb.CreateCategoryRequest) (*pb.ApiResponseCategory, error) {
	s.logger.Info("Create category called", zap.String("name", request.GetName()))

	req := &requests.CreateCategoryRequest{
		Name:        request.GetName(),
		Description: request.GetDescription(),
	}

	if err := req.Validate(); err != nil {
		return nil, category_errors.ErrGrpcValidateCreateCategory
	}

	category, err := s.categoryCommand.CreateCategory(ctx, req)
	if err != nil {
		s.logger.Error("Create category failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Create category success")

	return &pb.ApiResponseCategory{
		Status:  "success",
		Message: "Successfully created category",
		Data:    mapResponseCategory(category),
	}, nil
}

func (s *categoryHandleGrpc) Update(ctx context.Context, request *pb.UpdateCategoryRequest) (*pb.ApiResponseCategory, error) {
	s.logger.Info("Update category called", zap.Int32("id", request.GetCategoryId()))

	id := int(request.GetCategoryId())
	if id == 0 {
		return nil, category_errors.ErrGrpcFailedInvalidId
	}

	req := &requests.UpdateCategoryRequest{
		CategoryID:  &id,
		Name:        request.GetName(),
		Description: request.GetDescription(),
	}

	if err := req.Validate(); err != nil {
		return nil, category_errors.ErrGrpcValidateUpdateCategory
	}

	category, err := s.categoryCommand.UpdateCategory(ctx, req)
	if err != nil {
		s.logger.Error("Update category failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Update category success")

	return &pb.ApiResponseCategory{
		Status:  "success",
		Message: "Successfully updated category",
		Data:    mapResponseCategory(category),
	}, nil
}

func (s *categoryHandleGrpc) TrashedCategory(ctx context.Context, request *pb.FindByIdCategoryRequest) (*pb.ApiResponseCategoryDeleteAt, error) {
	s.logger.Info("TrashedCategory called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id == 0 {
		return nil, category_errors.ErrGrpcFailedInvalidId
	}

	category, err := s.categoryCommand.TrashedCategory(ctx, id)
	if err != nil {
		s.logger.Error("TrashedCategory failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("TrashedCategory success")

	return &pb.ApiResponseCategoryDeleteAt{
		Status:  "success",
		Message: "Successfully trashed category",
		Data:    mapResponseCategoryDeleteAt(category),
	}, nil
}

func (s *categoryHandleGrpc) RestoreCategory(ctx context.Context, request *pb.FindByIdCategoryRequest) (*pb.ApiResponseCategoryDeleteAt, error) {
	s.logger.Info("RestoreCategory called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id == 0 {
		return nil, category_errors.ErrGrpcFailedInvalidId
	}

	category, err := s.categoryCommand.RestoreCategory(ctx, id)
	if err != nil {
		s.logger.Error("RestoreCategory failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreCategory success")

	return &pb.ApiResponseCategoryDeleteAt{
		Status:  "success",
		Message: "Successfully restored category",
		Data:    mapResponseCategoryDeleteAt(category),
	}, nil
}

func (s *categoryHandleGrpc) DeleteCategoryPermanent(ctx context.Context, request *pb.FindByIdCategoryRequest) (*pb.ApiResponseCategoryDelete, error) {
	s.logger.Info("DeleteCategoryPermanent called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id == 0 {
		return nil, category_errors.ErrGrpcFailedInvalidId
	}

	_, err := s.categoryCommand.DeleteCategoryPermanent(ctx, id)
	if err != nil {
		s.logger.Error("DeleteCategoryPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteCategoryPermanent success")

	return &pb.ApiResponseCategoryDelete{
		Status:  "success",
		Message: "Successfully deleted category permanently",
	}, nil
}

func (s *categoryHandleGrpc) RestoreAllCategory(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseCategoryAll, error) {
	s.logger.Info("RestoreAllCategory called")

	_, err := s.categoryCommand.RestoreAllCategories(ctx)
	if err != nil {
		s.logger.Error("RestoreAllCategory failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreAllCategory success")

	return &pb.ApiResponseCategoryAll{
		Status:  "success",
		Message: "Successfully restore all category",
	}, nil
}

func (s *categoryHandleGrpc) DeleteAllCategoryPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseCategoryAll, error) {
	s.logger.Info("DeleteAllCategoryPermanent called")

	_, err := s.categoryCommand.DeleteAllCategoriesPermanent(ctx)
	if err != nil {
		s.logger.Error("DeleteAllCategoryPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteAllCategoryPermanent success")

	return &pb.ApiResponseCategoryAll{
		Status:  "success",
		Message: "Successfully delete category permanen",
	}, nil
}

// Internal map helpers
func mapPaginationMeta(meta *pb.PaginationMeta) *pb.PaginationMeta {
	if meta == nil {
		return nil
	}
	return &pb.PaginationMeta{
		CurrentPage:  meta.CurrentPage,
		PageSize:     meta.PageSize,
		TotalPages:   meta.TotalPages,
		TotalRecords: meta.TotalRecords,
	}
}

func mapResponseCategory(category *db.Category) *pb.CategoryResponse {
	if category == nil {
		return nil
	}
	var description, slugCategory string
	if category.Description != nil {
		description = *category.Description
	}
	if category.SlugCategory != nil {
		slugCategory = *category.SlugCategory
	}
	var createdAtStr, updatedAtStr string
	if category.CreatedAt.Valid {
		createdAtStr = category.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}
	if category.UpdatedAt.Valid {
		updatedAtStr = category.UpdatedAt.Time.Format("2006-01-02 15:04:05")
	}
	return &pb.CategoryResponse{
		Id:            int32(category.CategoryID),
		Name:          category.Name,
		Description:   description,
		SlugCategory:  slugCategory,
		CreatedAt:     createdAtStr,
		UpdatedAt:     updatedAtStr,
	}
}

func mapResponseGetCategory(category *db.GetCategoriesRow) *pb.CategoryResponse {
	if category == nil {
		return nil
	}
	var description, slugCategory string
	if category.Description != nil {
		description = *category.Description
	}
	if category.SlugCategory != nil {
		slugCategory = *category.SlugCategory
	}
	var createdAtStr, updatedAtStr string
	if category.CreatedAt.Valid {
		createdAtStr = category.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}
	if category.UpdatedAt.Valid {
		updatedAtStr = category.UpdatedAt.Time.Format("2006-01-02 15:04:05")
	}
	return &pb.CategoryResponse{
		Id:            int32(category.CategoryID),
		Name:          category.Name,
		Description:   description,
		SlugCategory:  slugCategory,
		CreatedAt:     createdAtStr,
		UpdatedAt:     updatedAtStr,
	}
}

func mapResponsesCategory(categories []*db.GetCategoriesRow) []*pb.CategoryResponse {
	var mappedCategories []*pb.CategoryResponse
	for _, category := range categories {
		mappedCategories = append(mappedCategories, mapResponseGetCategory(category))
	}
	return mappedCategories
}

func mapResponseCategoryDeleteAt(category *db.Category) *pb.CategoryResponseDeleteAt {
	if category == nil {
		return nil
	}
	var description, slugCategory string
	if category.Description != nil {
		description = *category.Description
	}
	if category.SlugCategory != nil {
		slugCategory = *category.SlugCategory
	}
	var createdAtStr, updatedAtStr string
	if category.CreatedAt.Valid {
		createdAtStr = category.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}
	if category.UpdatedAt.Valid {
		updatedAtStr = category.UpdatedAt.Time.Format("2006-01-02 15:04:05")
	}
	var deletedAt *wrapperspb.StringValue
	if category.DeletedAt.Valid {
		deletedAt = wrapperspb.String(category.DeletedAt.Time.Format("2006-01-02 15:04:05"))
	}

	return &pb.CategoryResponseDeleteAt{
		Id:            int32(category.CategoryID),
		Name:          category.Name,
		Description:   description,
		SlugCategory:  slugCategory,
		CreatedAt:     createdAtStr,
		UpdatedAt:     updatedAtStr,
		DeletedAt:     deletedAt,
	}
}

func mapResponseGetCategoryActive(category *db.GetCategoriesActiveRow) *pb.CategoryResponseDeleteAt {
	if category == nil {
		return nil
	}
	var description, slugCategory string
	if category.Description != nil {
		description = *category.Description
	}
	if category.SlugCategory != nil {
		slugCategory = *category.SlugCategory
	}
	var createdAtStr, updatedAtStr string
	if category.CreatedAt.Valid {
		createdAtStr = category.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}
	if category.UpdatedAt.Valid {
		updatedAtStr = category.UpdatedAt.Time.Format("2006-01-02 15:04:05")
	}
	var deletedAt *wrapperspb.StringValue
	if category.DeletedAt.Valid {
		deletedAt = wrapperspb.String(category.DeletedAt.Time.Format("2006-01-02 15:04:05"))
	}

	return &pb.CategoryResponseDeleteAt{
		Id:            int32(category.CategoryID),
		Name:          category.Name,
		Description:   description,
		SlugCategory:  slugCategory,
		CreatedAt:     createdAtStr,
		UpdatedAt:     updatedAtStr,
		DeletedAt:     deletedAt,
	}
}

func mapResponsesCategoryActive(categories []*db.GetCategoriesActiveRow) []*pb.CategoryResponseDeleteAt {
	var mappedCategories []*pb.CategoryResponseDeleteAt
	for _, category := range categories {
		mappedCategories = append(mappedCategories, mapResponseGetCategoryActive(category))
	}
	return mappedCategories
}

func mapResponseGetCategoryTrashed(category *db.GetCategoriesTrashedRow) *pb.CategoryResponseDeleteAt {
	if category == nil {
		return nil
	}
	var description, slugCategory string
	if category.Description != nil {
		description = *category.Description
	}
	if category.SlugCategory != nil {
		slugCategory = *category.SlugCategory
	}
	var createdAtStr, updatedAtStr string
	if category.CreatedAt.Valid {
		createdAtStr = category.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}
	if category.UpdatedAt.Valid {
		updatedAtStr = category.UpdatedAt.Time.Format("2006-01-02 15:04:05")
	}
	var deletedAt *wrapperspb.StringValue
	if category.DeletedAt.Valid {
		deletedAt = wrapperspb.String(category.DeletedAt.Time.Format("2006-01-02 15:04:05"))
	}

	return &pb.CategoryResponseDeleteAt{
		Id:            int32(category.CategoryID),
		Name:          category.Name,
		Description:   description,
		SlugCategory:  slugCategory,
		CreatedAt:     createdAtStr,
		UpdatedAt:     updatedAtStr,
		DeletedAt:     deletedAt,
	}
}

func mapResponsesCategoryTrashed(categories []*db.GetCategoriesTrashedRow) []*pb.CategoryResponseDeleteAt {
	var mappedCategories []*pb.CategoryResponseDeleteAt
	for _, category := range categories {
		mappedCategories = append(mappedCategories, mapResponseGetCategoryTrashed(category))
	}
	return mappedCategories
}

func mapResponseCategoryMonthlyPrice(category *db.GetMonthlyCategoryRow) *pb.CategoryMonthPriceResponse {
	if category == nil {
		return nil
	}
	return &pb.CategoryMonthPriceResponse{
		Month:        category.Month,
		CategoryId:   category.CategoryID,
		CategoryName: category.CategoryName,
		OrderCount:   int32(category.OrderCount),
		ItemsSold:    int32(category.ItemsSold),
		TotalRevenue: category.TotalRevenue,
	}
}

func mapResponsesCategoryMonthlyPrices(c []*db.GetMonthlyCategoryRow) []*pb.CategoryMonthPriceResponse {
	var categoryRecords []*pb.CategoryMonthPriceResponse
	for _, category := range c {
		categoryRecords = append(categoryRecords, mapResponseCategoryMonthlyPrice(category))
	}
	return categoryRecords
}

func mapResponseCategoryMonthlyPriceById(category *db.GetMonthlyCategoryByIdRow) *pb.CategoryMonthPriceResponse {
	if category == nil {
		return nil
	}
	return &pb.CategoryMonthPriceResponse{
		Month:        category.Month,
		CategoryId:   category.CategoryID,
		CategoryName: category.CategoryName,
		OrderCount:   int32(category.OrderCount),
		ItemsSold:    int32(category.ItemsSold),
		TotalRevenue: category.TotalRevenue,
	}
}

func mapResponsesCategoryMonthlyPricesById(c []*db.GetMonthlyCategoryByIdRow) []*pb.CategoryMonthPriceResponse {
	var categoryRecords []*pb.CategoryMonthPriceResponse
	for _, category := range c {
		categoryRecords = append(categoryRecords, mapResponseCategoryMonthlyPriceById(category))
	}
	return categoryRecords
}

func mapResponseCategoryMonthlyPriceByMerchant(category *db.GetMonthlyCategoryByMerchantRow) *pb.CategoryMonthPriceResponse {
	if category == nil {
		return nil
	}
	return &pb.CategoryMonthPriceResponse{
		Month:        category.Month,
		CategoryId:   category.CategoryID,
		CategoryName: category.CategoryName,
		OrderCount:   int32(category.OrderCount),
		ItemsSold:    int32(category.ItemsSold),
		TotalRevenue: category.TotalRevenue,
	}
}

func mapResponsesCategoryMonthlyPricesByMerchant(c []*db.GetMonthlyCategoryByMerchantRow) []*pb.CategoryMonthPriceResponse {
	var categoryRecords []*pb.CategoryMonthPriceResponse
	for _, category := range c {
		categoryRecords = append(categoryRecords, mapResponseCategoryMonthlyPriceByMerchant(category))
	}
	return categoryRecords
}

func mapResponseCategoryYearlyPrice(category *db.GetYearlyCategoryRow) *pb.CategoryYearPriceResponse {
	if category == nil {
		return nil
	}
	return &pb.CategoryYearPriceResponse{
		Year:               category.Year,
		CategoryId:         category.CategoryID,
		CategoryName:       category.CategoryName,
		OrderCount:         int32(category.OrderCount),
		ItemsSold:          int32(category.ItemsSold),
		TotalRevenue:       category.TotalRevenue,
		UniqueProductsSold: int32(category.UniqueProductsSold),
	}
}

func mapResponsesCategoryYearlyPrices(c []*db.GetYearlyCategoryRow) []*pb.CategoryYearPriceResponse {
	var categoryRecords []*pb.CategoryYearPriceResponse
	for _, category := range c {
		categoryRecords = append(categoryRecords, mapResponseCategoryYearlyPrice(category))
	}
	return categoryRecords
}

func mapResponseCategoryYearlyPriceById(category *db.GetYearlyCategoryByIdRow) *pb.CategoryYearPriceResponse {
	if category == nil {
		return nil
	}
	return &pb.CategoryYearPriceResponse{
		Year:               category.Year,
		CategoryId:         category.CategoryID,
		CategoryName:       category.CategoryName,
		OrderCount:         int32(category.OrderCount),
		ItemsSold:          int32(category.ItemsSold),
		TotalRevenue:       category.TotalRevenue,
		UniqueProductsSold: int32(category.UniqueProductsSold),
	}
}

func mapResponsesCategoryYearlyPricesById(c []*db.GetYearlyCategoryByIdRow) []*pb.CategoryYearPriceResponse {
	var categoryRecords []*pb.CategoryYearPriceResponse
	for _, category := range c {
		categoryRecords = append(categoryRecords, mapResponseCategoryYearlyPriceById(category))
	}
	return categoryRecords
}

func mapResponseCategoryYearlyPriceByMerchant(category *db.GetYearlyCategoryByMerchantRow) *pb.CategoryYearPriceResponse {
	if category == nil {
		return nil
	}
	return &pb.CategoryYearPriceResponse{
		Year:               category.Year,
		CategoryId:         category.CategoryID,
		CategoryName:       category.CategoryName,
		OrderCount:         int32(category.OrderCount),
		ItemsSold:          int32(category.ItemsSold),
		TotalRevenue:       category.TotalRevenue,
		UniqueProductsSold: int32(category.UniqueProductsSold),
	}
}

func mapResponsesCategoryYearlyPricesByMerchant(c []*db.GetYearlyCategoryByMerchantRow) []*pb.CategoryYearPriceResponse {
	var categoryRecords []*pb.CategoryYearPriceResponse
	for _, category := range c {
		categoryRecords = append(categoryRecords, mapResponseCategoryYearlyPriceByMerchant(category))
	}
	return categoryRecords
}

func mapResponseCashierMonthlyTotalPrice(c *db.GetMonthlyTotalPriceRow) *pb.CategoriesMonthlyTotalPriceResponse {
	if c == nil {
		return nil
	}
	return &pb.CategoriesMonthlyTotalPriceResponse{
		Year:         c.Year,
		Month:        c.Month,
		TotalRevenue: c.TotalRevenue,
	}
}

func mapResponseCategoryMonthlyTotalPrices(c []*db.GetMonthlyTotalPriceRow) []*pb.CategoriesMonthlyTotalPriceResponse {
	var CategoryRecords []*pb.CategoriesMonthlyTotalPriceResponse
	for _, Category := range c {
		CategoryRecords = append(CategoryRecords, mapResponseCashierMonthlyTotalPrice(Category))
	}
	return CategoryRecords
}

func mapResponseCashierMonthlyTotalPriceById(c *db.GetMonthlyTotalPriceByIdRow) *pb.CategoriesMonthlyTotalPriceResponse {
	if c == nil {
		return nil
	}
	return &pb.CategoriesMonthlyTotalPriceResponse{
		Year:         c.Year,
		Month:        c.Month,
		TotalRevenue: c.TotalRevenue,
	}
}

func mapResponseCategoryMonthlyTotalPricesById(c []*db.GetMonthlyTotalPriceByIdRow) []*pb.CategoriesMonthlyTotalPriceResponse {
	var CategoryRecords []*pb.CategoriesMonthlyTotalPriceResponse
	for _, Category := range c {
		CategoryRecords = append(CategoryRecords, mapResponseCashierMonthlyTotalPriceById(Category))
	}
	return CategoryRecords
}

func mapResponseCashierMonthlyTotalPriceByMerchant(c *db.GetMonthlyTotalPriceByMerchantRow) *pb.CategoriesMonthlyTotalPriceResponse {
	if c == nil {
		return nil
	}
	return &pb.CategoriesMonthlyTotalPriceResponse{
		Year:         c.Year,
		Month:        c.Month,
		TotalRevenue: c.TotalRevenue,
	}
}

func mapResponseCategoryMonthlyTotalPricesByMerchant(c []*db.GetMonthlyTotalPriceByMerchantRow) []*pb.CategoriesMonthlyTotalPriceResponse {
	var CategoryRecords []*pb.CategoriesMonthlyTotalPriceResponse
	for _, Category := range c {
		CategoryRecords = append(CategoryRecords, mapResponseCashierMonthlyTotalPriceByMerchant(Category))
	}
	return CategoryRecords
}

func mapResponseCategoryYearlyTotalSale(c *db.GetYearlyTotalPriceRow) *pb.CategoriesYearlyTotalPriceResponse {
	if c == nil {
		return nil
	}
	return &pb.CategoriesYearlyTotalPriceResponse{
		Year:         c.Year,
		TotalRevenue: c.TotalRevenue,
	}
}

func mapResponseCategoryYearlyTotalPrices(c []*db.GetYearlyTotalPriceRow) []*pb.CategoriesYearlyTotalPriceResponse {
	var CategoryRecords []*pb.CategoriesYearlyTotalPriceResponse
	for _, Category := range c {
		CategoryRecords = append(CategoryRecords, mapResponseCategoryYearlyTotalSale(Category))
	}
	return CategoryRecords
}

func mapResponseCategoryYearlyTotalSaleById(c *db.GetYearlyTotalPriceByIdRow) *pb.CategoriesYearlyTotalPriceResponse {
	if c == nil {
		return nil
	}
	return &pb.CategoriesYearlyTotalPriceResponse{
		Year:         c.Year,
		TotalRevenue: c.TotalRevenue,
	}
}

func mapResponseCategoryYearlyTotalPricesById(c []*db.GetYearlyTotalPriceByIdRow) []*pb.CategoriesYearlyTotalPriceResponse {
	var CategoryRecords []*pb.CategoriesYearlyTotalPriceResponse
	for _, Category := range c {
		CategoryRecords = append(CategoryRecords, mapResponseCategoryYearlyTotalSaleById(Category))
	}
	return CategoryRecords
}

func mapResponseCategoryYearlyTotalSaleByMerchant(c *db.GetYearlyTotalPriceByMerchantRow) *pb.CategoriesYearlyTotalPriceResponse {
	if c == nil {
		return nil
	}
	return &pb.CategoriesYearlyTotalPriceResponse{
		Year:         c.Year,
		TotalRevenue: c.TotalRevenue,
	}
}

func mapResponseCategoryYearlyTotalPricesByMerchant(c []*db.GetYearlyTotalPriceByMerchantRow) []*pb.CategoriesYearlyTotalPriceResponse {
	var CategoryRecords []*pb.CategoriesYearlyTotalPriceResponse
	for _, Category := range c {
		CategoryRecords = append(CategoryRecords, mapResponseCategoryYearlyTotalSaleByMerchant(Category))
	}
	return CategoryRecords
}
