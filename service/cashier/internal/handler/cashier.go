package handler

import (
	"context"
	"math"

	"github.com/MamangRust/monolith-point-of-sale-cashier/internal/service"
	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/cashier_errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/pb"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type cashierHandleGrpc struct {
	pb.UnimplementedCashierServiceServer
	cashierQuery           service.CashierQueryService
	cashierCommand         service.CashierCommandService
	cashierStats           service.CashierStatsService
	cashierStatsById       service.CashierStatsByIdService
	cashierStatsByMerchant service.CashierStatsByMerchant
	logger                 logger.LoggerInterface
}

func NewCashierHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pb.CashierServiceServer {
	return &cashierHandleGrpc{
		cashierQuery:           service.CashierQuery,
		cashierCommand:         service.CashierCommand,
		cashierStats:           service.CashierStats,
		cashierStatsById:       service.CashierStatsById,
		cashierStatsByMerchant: service.CashierStatsByMerchant,
		logger:                 logger,
	}
}

func (s *cashierHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllCashierRequest) (*pb.ApiResponsePaginationCashier, error) {
	s.logger.Info("FindAll cashier called", zap.Int32("page", request.GetPage()), zap.Int32("pageSize", request.GetPageSize()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllCashiers{
		Search:   search,
		Page:     page,
		PageSize: pageSize,
	}

	cashier, totalRecords, err := s.cashierQuery.FindAll(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindAll cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindAll cashier success", zap.Int("count", len(cashier)))

	return &pb.ApiResponsePaginationCashier{
		Status:     "success",
		Message:    "Successfully fetched cashier",
		Data:       mapResponsesCashier(cashier),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *cashierHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdCashierRequest) (*pb.ApiResponseCashier, error) {
	s.logger.Info("FindById cashier called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id == 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidId
	}

	cashier, err := s.cashierQuery.FindById(ctx, id)
	if err != nil {
		s.logger.Error("FindById cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindById cashier success", zap.Int("id", id))

	return &pb.ApiResponseCashier{
		Status:  "success",
		Message: "Successfully fetched categories",
		Data:    mapResponseCashier(cashier),
	}, nil
}

func (s *cashierHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllCashierRequest) (*pb.ApiResponsePaginationCashierDeleteAt, error) {
	s.logger.Info("FindByActive cashier called", zap.Int32("page", request.GetPage()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllCashiers{
		Search:   search,
		Page:     page,
		PageSize: pageSize,
	}

	cashier, totalRecords, err := s.cashierQuery.FindByActive(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByActive cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByActive cashier success", zap.Int("count", len(cashier)))

	return &pb.ApiResponsePaginationCashierDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active cashier",
		Data:       mapResponsesCashierActive(cashier),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *cashierHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllCashierRequest) (*pb.ApiResponsePaginationCashierDeleteAt, error) {
	s.logger.Info("FindByTrashed cashier called")

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllCashiers{
		Search:   search,
		Page:     page,
		PageSize: pageSize,
	}

	users, totalRecords, err := s.cashierQuery.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByTrashed cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByTrashed cashier success", zap.Int("count", len(users)))

	return &pb.ApiResponsePaginationCashierDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed cashier",
		Data:       mapResponsesCashierTrashed(users),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *cashierHandleGrpc) FindByMerchant(ctx context.Context, request *pb.FindByMerchantCashierRequest) (*pb.ApiResponsePaginationCashier, error) {
	s.logger.Info("FindByMerchant cashier called", zap.Int32("merchantId", request.GetMerchantId()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()
	merchant_id := int(request.GetMerchantId())

	if merchant_id <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidMerchantId
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllCashierMerchant{
		Search:     search,
		Page:       page,
		PageSize:   pageSize,
		MerchantID: merchant_id,
	}

	cashier, totalRecords, err := s.cashierQuery.FindByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByMerchant cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByMerchant cashier success", zap.Int("count", len(cashier)))

	return &pb.ApiResponsePaginationCashier{
		Status:     "success",
		Message:    "Successfully fetched cashier",
		Data:       mapResponsesCashierByMerchant(cashier),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *cashierHandleGrpc) FindMonthlyTotalSales(ctx context.Context, req *pb.FindYearMonthTotalSales) (*pb.ApiResponseCashierMonthlyTotalSales, error) {
	s.logger.Info("FindMonthlyTotalSales cashier called", zap.Int32("year", req.GetYear()), zap.Int32("month", req.GetMonth()))

	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidYear
	}
	if month <= 0 || month >= 12 {
		return nil, cashier_errors.ErrGrpcFailedInvalidMonth
	}

	reqService := requests.MonthTotalSales{
		Year:  year,
		Month: month,
	}

	methods, err := s.cashierStats.FindMonthlyTotalSales(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthlyTotalSales cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthlyTotalSales cashier success")

	return &pb.ApiResponseCashierMonthlyTotalSales{
		Status:  "success",
		Message: "Monthly sales retrieved successfully",
		Data:    mapResponseCashierMonthlyTotalSales(methods),
	}, nil
}

func (s *cashierHandleGrpc) FindYearlyTotalSales(ctx context.Context, req *pb.FindYearTotalSales) (*pb.ApiResponseCashierYearlyTotalSales, error) {
	s.logger.Info("FindYearlyTotalSales cashier called", zap.Int32("year", req.GetYear()))

	year := int(req.GetYear())
	if year <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidYear
	}

	methods, err := s.cashierStats.FindYearlyTotalSales(ctx, year)
	if err != nil {
		s.logger.Error("FindYearlyTotalSales cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearlyTotalSales cashier success")

	return &pb.ApiResponseCashierYearlyTotalSales{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    mapResponseCashierYearlyTotalSales(methods),
	}, nil
}

func (s *cashierHandleGrpc) FindMonthlyTotalSalesById(ctx context.Context, req *pb.FindYearMonthTotalSalesById) (*pb.ApiResponseCashierMonthlyTotalSales, error) {
	s.logger.Info("FindMonthlyTotalSalesById cashier called", zap.Int32("id", req.GetCashierId()))

	year := int(req.GetYear())
	month := int(req.GetMonth())
	id := int(req.GetCashierId())

	if year <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidYear
	}
	if month <= 0 || month >= 12 {
		return nil, cashier_errors.ErrGrpcFailedInvalidMonth
	}
	if id <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidId
	}

	reqService := requests.MonthTotalSalesCashier{
		Year:      year,
		Month:     month,
		CashierID: id,
	}

	methods, err := s.cashierStatsById.FindMonthlyTotalSalesById(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthlyTotalSalesById cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthlyTotalSalesById cashier success")

	return &pb.ApiResponseCashierMonthlyTotalSales{
		Status:  "success",
		Message: "Monthly sales retrieved successfully",
		Data:    mapResponseCashierMonthlyTotalSalesById(methods),
	}, nil
}

func (s *cashierHandleGrpc) FindYearlyTotalSalesById(ctx context.Context, req *pb.FindYearTotalSalesById) (*pb.ApiResponseCashierYearlyTotalSales, error) {
	s.logger.Info("FindYearlyTotalSalesById cashier called", zap.Int32("id", req.GetCashierId()))

	year := int(req.GetYear())
	id := int(req.GetCashierId())

	if year <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidYear
	}
	if id <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidId
	}

	reqService := requests.YearTotalSalesCashier{
		Year:      year,
		CashierID: id,
	}

	methods, err := s.cashierStatsById.FindYearlyTotalSalesById(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindYearlyTotalSalesById cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearlyTotalSalesById cashier success")

	return &pb.ApiResponseCashierYearlyTotalSales{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    mapResponseCashierYearlyTotalSalesById(methods),
	}, nil
}

func (s *cashierHandleGrpc) FindMonthlyTotalSalesByMerchant(ctx context.Context, req *pb.FindYearMonthTotalSalesByMerchant) (*pb.ApiResponseCashierMonthlyTotalSales, error) {
	s.logger.Info("FindMonthlyTotalSalesByMerchant cashier called", zap.Int32("merchantId", req.GetMerchantId()))

	year := int(req.GetYear())
	month := int(req.GetMonth())
	merchantId := int(req.GetMerchantId())

	if year <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidYear
	}
	if month <= 0 || month >= 12 {
		return nil, cashier_errors.ErrGrpcFailedInvalidMonth
	}
	if merchantId <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidMerchantId
	}

	reqService := requests.MonthTotalSalesMerchant{
		Year:       year,
		Month:      month,
		MerchantID: merchantId,
	}

	methods, err := s.cashierStatsByMerchant.FindMonthlyTotalSalesByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthlyTotalSalesByMerchant cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthlyTotalSalesByMerchant cashier success")

	return &pb.ApiResponseCashierMonthlyTotalSales{
		Status:  "success",
		Message: "Monthly sales retrieved successfully",
		Data:    mapResponseCashierMonthlyTotalSalesByMerchant(methods),
	}, nil
}

func (s *cashierHandleGrpc) FindYearlyTotalSalesByMerchant(ctx context.Context, req *pb.FindYearTotalSalesByMerchant) (*pb.ApiResponseCashierYearlyTotalSales, error) {
	s.logger.Info("FindYearlyTotalSalesByMerchant cashier called", zap.Int32("merchantId", req.GetMerchantId()))

	year := int(req.GetYear())
	merchantId := int(req.GetMerchantId())

	if year <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidYear
	}
	if merchantId <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidMerchantId
	}

	reqService := requests.YearTotalSalesMerchant{
		Year:       year,
		MerchantID: merchantId,
	}

	methods, err := s.cashierStatsByMerchant.FindYearlyTotalSalesByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindYearlyTotalSalesByMerchant cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearlyTotalSalesByMerchant cashier success")

	return &pb.ApiResponseCashierYearlyTotalSales{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    mapResponseCashierYearlyTotalSalesByMerchant(methods),
	}, nil
}

func (s *cashierHandleGrpc) FindMonthSales(ctx context.Context, req *pb.FindYearCashier) (*pb.ApiResponseCashierMonthSales, error) {
	s.logger.Info("FindMonthSales cashier called", zap.Int32("year", req.GetYear()))

	year := int(req.GetYear())
	if year <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidYear
	}

	methods, err := s.cashierStats.FindMonthlySales(ctx, year)
	if err != nil {
		s.logger.Error("FindMonthSales cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthSales cashier success")

	return &pb.ApiResponseCashierMonthSales{
		Status:  "success",
		Message: "Monthly sales retrieved successfully",
		Data:    mapResponsesCashierMonthlySales(methods),
	}, nil
}

func (s *cashierHandleGrpc) FindYearSales(ctx context.Context, req *pb.FindYearCashier) (*pb.ApiResponseCashierYearSales, error) {
	s.logger.Info("FindYearSales cashier called", zap.Int32("year", req.GetYear()))

	year := int(req.GetYear())
	if year <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidYear
	}

	methods, err := s.cashierStats.FindYearlySales(ctx, year)
	if err != nil {
		s.logger.Error("FindYearSales cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearSales cashier success")

	return &pb.ApiResponseCashierYearSales{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    mapResponsesCashierYearlySales(methods),
	}, nil
}

func (s *cashierHandleGrpc) FindMonthSalesByMerchant(ctx context.Context, req *pb.FindYearCashierByMerchant) (*pb.ApiResponseCashierMonthSales, error) {
	s.logger.Info("FindMonthSalesByMerchant cashier called", zap.Int32("merchantId", req.GetMerchantId()))

	year := int(req.GetYear())
	merchantId := int(req.GetMerchantId())

	if year <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidYear
	}
	if merchantId <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidMerchantId
	}

	reqService := requests.MonthCashierMerchant{
		Year:       year,
		MerchantID: merchantId,
	}

	methods, err := s.cashierStatsByMerchant.FindMonthlyCashierByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthSalesByMerchant cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthSalesByMerchant cashier success")

	return &pb.ApiResponseCashierMonthSales{
		Status:  "success",
		Message: "Merchant monthly revenue retrieved successfully",
		Data:    mapResponsesCashierMonthlySalesByMerchant(methods),
	}, nil
}

func (s *cashierHandleGrpc) FindYearSalesByMerchant(ctx context.Context, req *pb.FindYearCashierByMerchant) (*pb.ApiResponseCashierYearSales, error) {
	s.logger.Info("FindYearSalesByMerchant cashier called", zap.Int32("merchantId", req.GetMerchantId()))

	year := int(req.GetYear())
	merchantId := int(req.GetMerchantId())

	if year <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidYear
	}
	if merchantId <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidMerchantId
	}

	reqService := requests.YearCashierMerchant{
		Year:       year,
		MerchantID: merchantId,
	}

	methods, err := s.cashierStatsByMerchant.FindYearlyCashierByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindYearSalesByMerchant cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearSalesByMerchant cashier success")

	return &pb.ApiResponseCashierYearSales{
		Status:  "success",
		Message: "Merchant yearly payment methods retrieved successfully",
		Data:    mapResponsesCashierYearlySalesByMerchant(methods),
	}, nil
}

func (s *cashierHandleGrpc) FindMonthSalesById(ctx context.Context, req *pb.FindYearCashierById) (*pb.ApiResponseCashierMonthSales, error) {
	s.logger.Info("FindMonthSalesById cashier called", zap.Int32("id", req.GetCashierId()))

	year := int(req.GetYear())
	cashierId := int(req.GetCashierId())

	if year <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidYear
	}
	if cashierId <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidId
	}

	reqService := requests.MonthCashierId{
		Year:      year,
		CashierID: cashierId,
	}

	methods, err := s.cashierStatsById.FindMonthlyCashierById(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthSalesById cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthSalesById cashier success")

	return &pb.ApiResponseCashierMonthSales{
		Status:  "success",
		Message: "Cashier monthly sales retrieved successfully",
		Data:    mapResponsesCashierMonthlySalesById(methods),
	}, nil
}

func (s *cashierHandleGrpc) FindYearSalesById(ctx context.Context, req *pb.FindYearCashierById) (*pb.ApiResponseCashierYearSales, error) {
	s.logger.Info("FindYearSalesById cashier called", zap.Int32("id", req.GetCashierId()))

	year := int(req.GetYear())
	cashierId := int(req.GetCashierId())

	if year <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidYear
	}
	if cashierId <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidId
	}

	reqService := requests.YearCashierId{
		Year:      year,
		CashierID: cashierId,
	}

	methods, err := s.cashierStatsById.FindYearlyCashierById(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindYearSalesById cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearSalesById cashier success")

	return &pb.ApiResponseCashierYearSales{
		Status:  "success",
		Message: "Cashier yearly sales retrieved successfully",
		Data:    mapResponsesCashierYearlySalesById(methods),
	}, nil
}

func (s *cashierHandleGrpc) CreateCashier(ctx context.Context, request *pb.CreateCashierRequest) (*pb.ApiResponseCashier, error) {
	s.logger.Info("CreateCashier called", zap.String("name", request.GetName()))

	req := &requests.CreateCashierRequest{
		Name:       request.GetName(),
		MerchantID: int(request.GetMerchantId()),
		UserID:     int(request.GetUserId()),
	}

	if err := req.Validate(); err != nil {
		return nil, cashier_errors.ErrGrpcValidateCreateCashier
	}

	cashier, err := s.cashierCommand.CreateCashier(ctx, req)
	if err != nil {
		s.logger.Error("CreateCashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("CreateCashier success")

	return &pb.ApiResponseCashier{
		Status:  "success",
		Message: "Successfully created cashier",
		Data:    mapResponseCashier(cashier),
	}, nil
}

func (s *cashierHandleGrpc) UpdateCashier(ctx context.Context, request *pb.UpdateCashierRequest) (*pb.ApiResponseCashier, error) {
	s.logger.Info("UpdateCashier called", zap.Int32("id", request.GetCashierId()))

	id := int(request.GetCashierId())
	if id == 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidId
	}

	req := &requests.UpdateCashierRequest{
		CashierID: &id,
		Name:      request.GetName(),
	}

	if err := req.Validate(); err != nil {
		return nil, cashier_errors.ErrGrpcValidateUpdateCashier
	}

	cashier, err := s.cashierCommand.UpdateCashier(ctx, req)
	if err != nil {
		s.logger.Error("UpdateCashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("UpdateCashier success")

	return &pb.ApiResponseCashier{
		Status:  "success",
		Message: "Successfully updated cashier",
		Data:    mapResponseCashier(cashier),
	}, nil
}

func (s *cashierHandleGrpc) TrashedCashier(ctx context.Context, request *pb.FindByIdCashierRequest) (*pb.ApiResponseCashierDeleteAt, error) {
	s.logger.Info("TrashedCashier called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id == 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidId
	}

	cashier, err := s.cashierCommand.TrashedCashier(ctx, id)
	if err != nil {
		s.logger.Error("TrashedCashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("TrashedCashier success")

	return &pb.ApiResponseCashierDeleteAt{
		Status:  "success",
		Message: "Successfully trashed cashier",
		Data:    mapResponseCashierDeleteAt(cashier),
	}, nil
}

func (s *cashierHandleGrpc) RestoreCashier(ctx context.Context, request *pb.FindByIdCashierRequest) (*pb.ApiResponseCashierDeleteAt, error) {
	s.logger.Info("RestoreCashier called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id == 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidId
	}

	cashier, err := s.cashierCommand.RestoreCashier(ctx, id)
	if err != nil {
		s.logger.Error("RestoreCashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreCashier success")

	return &pb.ApiResponseCashierDeleteAt{
		Status:  "success",
		Message: "Successfully restored cashier",
		Data:    mapResponseCashierDeleteAt(cashier),
	}, nil
}

func (s *cashierHandleGrpc) DeleteCashierPermanent(ctx context.Context, request *pb.FindByIdCashierRequest) (*pb.ApiResponseCashierDelete, error) {
	s.logger.Info("DeleteCashierPermanent called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id == 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidId
	}

	_, err := s.cashierCommand.DeleteCashierPermanent(ctx, id)
	if err != nil {
		s.logger.Error("DeleteCashierPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteCashierPermanent success")

	return &pb.ApiResponseCashierDelete{
		Status:  "success",
		Message: "Successfully deleted cashier permanently",
	}, nil
}

func (s *cashierHandleGrpc) RestoreAllCashier(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseCashierAll, error) {
	s.logger.Info("RestoreAllCashier called")

	_, err := s.cashierCommand.RestoreAllCashier(ctx)
	if err != nil {
		s.logger.Error("RestoreAllCashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreAllCashier success")

	return &pb.ApiResponseCashierAll{
		Status:  "success",
		Message: "Successfully restore all cashier",
	}, nil
}

func (s *cashierHandleGrpc) DeleteAllCashierPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseCashierAll, error) {
	s.logger.Info("DeleteAllCashierPermanent called")

	_, err := s.cashierCommand.DeleteAllCashierPermanent(ctx)
	if err != nil {
		s.logger.Error("DeleteAllCashierPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteAllCashierPermanent success")

	return &pb.ApiResponseCashierAll{
		Status:  "success",
		Message: "Successfully delete cashier permanen",
	}, nil
}

// Map helpers
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

func mapResponseCashier(cashier *db.Cashier) *pb.CashierResponse {
	if cashier == nil {
		return nil
	}
	var createdAtStr, updatedAtStr string
	if cashier.CreatedAt.Valid {
		createdAtStr = cashier.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}
	if cashier.UpdatedAt.Valid {
		updatedAtStr = cashier.UpdatedAt.Time.Format("2006-01-02 15:04:05")
	}
	return &pb.CashierResponse{
		Id:         int32(cashier.CashierID),
		MerchantId: int32(cashier.MerchantID),
		Name:       cashier.Name,
		CreatedAt:  createdAtStr,
		UpdatedAt:  updatedAtStr,
	}
}

func mapResponsesCashier(cashiers []*db.GetCashiersRow) []*pb.CashierResponse {
	var mappedCashiers []*pb.CashierResponse
	for _, cashier := range cashiers {
		var createdAtStr, updatedAtStr string
		if cashier.CreatedAt.Valid {
			createdAtStr = cashier.CreatedAt.Time.Format("2006-01-02 15:04:05")
		}
		if cashier.UpdatedAt.Valid {
			updatedAtStr = cashier.UpdatedAt.Time.Format("2006-01-02 15:04:05")
		}
		mappedCashiers = append(mappedCashiers, &pb.CashierResponse{
			Id:         int32(cashier.CashierID),
			MerchantId: int32(cashier.MerchantID),
			Name:       cashier.Name,
			CreatedAt:  createdAtStr,
			UpdatedAt:  updatedAtStr,
		})
	}
	return mappedCashiers
}

func mapResponsesCashierByMerchant(cashiers []*db.GetCashiersByMerchantRow) []*pb.CashierResponse {
	var mappedCashiers []*pb.CashierResponse
	for _, cashier := range cashiers {
		var createdAtStr, updatedAtStr string
		if cashier.CreatedAt.Valid {
			createdAtStr = cashier.CreatedAt.Time.Format("2006-01-02 15:04:05")
		}
		if cashier.UpdatedAt.Valid {
			updatedAtStr = cashier.UpdatedAt.Time.Format("2006-01-02 15:04:05")
		}
		mappedCashiers = append(mappedCashiers, &pb.CashierResponse{
			Id:         int32(cashier.CashierID),
			MerchantId: int32(cashier.MerchantID),
			Name:       cashier.Name,
			CreatedAt:  createdAtStr,
			UpdatedAt:  updatedAtStr,
		})
	}
	return mappedCashiers
}

func mapResponseCashierDeleteAt(cashier *db.Cashier) *pb.CashierResponseDeleteAt {
	if cashier == nil {
		return nil
	}
	var createdAtStr, updatedAtStr string
	if cashier.CreatedAt.Valid {
		createdAtStr = cashier.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}
	if cashier.UpdatedAt.Valid {
		updatedAtStr = cashier.UpdatedAt.Time.Format("2006-01-02 15:04:05")
	}
	var deletedAt *wrapperspb.StringValue
	if cashier.DeletedAt.Valid {
		deletedAt = wrapperspb.String(cashier.DeletedAt.Time.Format("2006-01-02 15:04:05"))
	}

	return &pb.CashierResponseDeleteAt{
		Id:         int32(cashier.CashierID),
		MerchantId: int32(cashier.MerchantID),
		Name:       cashier.Name,
		CreatedAt:  createdAtStr,
		UpdatedAt:  updatedAtStr,
		DeletedAt:  deletedAt,
	}
}

func mapResponsesCashierActive(cashiers []*db.GetCashiersActiveRow) []*pb.CashierResponseDeleteAt {
	var mappedCashiers []*pb.CashierResponseDeleteAt
	for _, cashier := range cashiers {
		var createdAtStr, updatedAtStr string
		if cashier.CreatedAt.Valid {
			createdAtStr = cashier.CreatedAt.Time.Format("2006-01-02 15:04:05")
		}
		if cashier.UpdatedAt.Valid {
			updatedAtStr = cashier.UpdatedAt.Time.Format("2006-01-02 15:04:05")
		}
		var deletedAt *wrapperspb.StringValue
		if cashier.DeletedAt.Valid {
			deletedAt = wrapperspb.String(cashier.DeletedAt.Time.Format("2006-01-02 15:04:05"))
		}
		mappedCashiers = append(mappedCashiers, &pb.CashierResponseDeleteAt{
			Id:         int32(cashier.CashierID),
			MerchantId: int32(cashier.MerchantID),
			Name:       cashier.Name,
			CreatedAt:  createdAtStr,
			UpdatedAt:  updatedAtStr,
			DeletedAt:  deletedAt,
		})
	}
	return mappedCashiers
}

func mapResponsesCashierTrashed(cashiers []*db.GetCashiersTrashedRow) []*pb.CashierResponseDeleteAt {
	var mappedCashiers []*pb.CashierResponseDeleteAt
	for _, cashier := range cashiers {
		var createdAtStr, updatedAtStr string
		if cashier.CreatedAt.Valid {
			createdAtStr = cashier.CreatedAt.Time.Format("2006-01-02 15:04:05")
		}
		if cashier.UpdatedAt.Valid {
			updatedAtStr = cashier.UpdatedAt.Time.Format("2006-01-02 15:04:05")
		}
		var deletedAt *wrapperspb.StringValue
		if cashier.DeletedAt.Valid {
			deletedAt = wrapperspb.String(cashier.DeletedAt.Time.Format("2006-01-02 15:04:05"))
		}
		mappedCashiers = append(mappedCashiers, &pb.CashierResponseDeleteAt{
			Id:         int32(cashier.CashierID),
			MerchantId: int32(cashier.MerchantID),
			Name:       cashier.Name,
			CreatedAt:  createdAtStr,
			UpdatedAt:  updatedAtStr,
			DeletedAt:  deletedAt,
		})
	}
	return mappedCashiers
}

func mapResponseCashierMonthlySale(cashier *db.GetMonthlyCashierRow) *pb.CashierResponseMonthSales {
	if cashier == nil {
		return nil
	}
	return &pb.CashierResponseMonthSales{
		Month:       cashier.Month,
		CashierId:   int32(cashier.CashierID),
		CashierName: cashier.CashierName,
		OrderCount:  int32(cashier.OrderCount),
		TotalSales:  int32(cashier.TotalSales),
	}
}

func mapResponsesCashierMonthlySales(c []*db.GetMonthlyCashierRow) []*pb.CashierResponseMonthSales {
	var cashierRecords []*pb.CashierResponseMonthSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, mapResponseCashierMonthlySale(cashier))
	}
	return cashierRecords
}

func mapResponseCashierMonthlySaleById(cashier *db.GetMonthlyCashierByCashierIdRow) *pb.CashierResponseMonthSales {
	if cashier == nil {
		return nil
	}
	return &pb.CashierResponseMonthSales{
		Month:       cashier.Month,
		CashierId:   int32(cashier.CashierID),
		CashierName: cashier.CashierName,
		OrderCount:  int32(cashier.OrderCount),
		TotalSales:  int32(cashier.TotalSales),
	}
}

func mapResponsesCashierMonthlySalesById(c []*db.GetMonthlyCashierByCashierIdRow) []*pb.CashierResponseMonthSales {
	var cashierRecords []*pb.CashierResponseMonthSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, mapResponseCashierMonthlySaleById(cashier))
	}
	return cashierRecords
}

func mapResponseCashierMonthlySaleByMerchant(cashier *db.GetMonthlyCashierByMerchantRow) *pb.CashierResponseMonthSales {
	if cashier == nil {
		return nil
	}
	return &pb.CashierResponseMonthSales{
		Month:       cashier.Month,
		CashierId:   int32(cashier.CashierID),
		CashierName: cashier.CashierName,
		OrderCount:  int32(cashier.OrderCount),
		TotalSales:  int32(cashier.TotalSales),
	}
}

func mapResponsesCashierMonthlySalesByMerchant(c []*db.GetMonthlyCashierByMerchantRow) []*pb.CashierResponseMonthSales {
	var cashierRecords []*pb.CashierResponseMonthSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, mapResponseCashierMonthlySaleByMerchant(cashier))
	}
	return cashierRecords
}

func mapResponseCashierYearlySale(cashier *db.GetYearlyCashierRow) *pb.CashierResponseYearSales {
	if cashier == nil {
		return nil
	}
	return &pb.CashierResponseYearSales{
		Year:        cashier.Year,
		CashierId:   int32(cashier.CashierID),
		CashierName: cashier.CashierName,
		OrderCount:  int32(cashier.OrderCount),
		TotalSales:  int32(cashier.TotalSales),
	}
}

func mapResponsesCashierYearlySales(c []*db.GetYearlyCashierRow) []*pb.CashierResponseYearSales {
	var cashierRecords []*pb.CashierResponseYearSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, mapResponseCashierYearlySale(cashier))
	}
	return cashierRecords
}

func mapResponseCashierYearlySaleById(cashier *db.GetYearlyCashierByCashierIdRow) *pb.CashierResponseYearSales {
	if cashier == nil {
		return nil
	}
	return &pb.CashierResponseYearSales{
		Year:        cashier.Year,
		CashierId:   int32(cashier.CashierID),
		CashierName: cashier.CashierName,
		OrderCount:  int32(cashier.OrderCount),
		TotalSales:  int32(cashier.TotalSales),
	}
}

func mapResponsesCashierYearlySalesById(c []*db.GetYearlyCashierByCashierIdRow) []*pb.CashierResponseYearSales {
	var cashierRecords []*pb.CashierResponseYearSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, mapResponseCashierYearlySaleById(cashier))
	}
	return cashierRecords
}

func mapResponseCashierYearlySaleByMerchant(cashier *db.GetYearlyCashierByMerchantRow) *pb.CashierResponseYearSales {
	if cashier == nil {
		return nil
	}
	return &pb.CashierResponseYearSales{
		Year:        cashier.Year,
		CashierId:   int32(cashier.CashierID),
		CashierName: cashier.CashierName,
		OrderCount:  int32(cashier.OrderCount),
		TotalSales:  int32(cashier.TotalSales),
	}
}

func mapResponsesCashierYearlySalesByMerchant(c []*db.GetYearlyCashierByMerchantRow) []*pb.CashierResponseYearSales {
	var cashierRecords []*pb.CashierResponseYearSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, mapResponseCashierYearlySaleByMerchant(cashier))
	}
	return cashierRecords
}

func mapResponseCashierMonthlyTotalSales(c []*db.GetMonthlyTotalSalesCashierRow) []*pb.CashierResponseMonthTotalSales {
	var cashierRecords []*pb.CashierResponseMonthTotalSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, &pb.CashierResponseMonthTotalSales{
			Year:       cashier.Year,
			Month:      cashier.Month,
			TotalSales: cashier.TotalSales,
		})
	}
	return cashierRecords
}

func mapResponseCashierMonthlyTotalSalesById(c []*db.GetMonthlyTotalSalesByIdRow) []*pb.CashierResponseMonthTotalSales {
	var cashierRecords []*pb.CashierResponseMonthTotalSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, &pb.CashierResponseMonthTotalSales{
			Year:       cashier.Year,
			Month:      cashier.Month,
			TotalSales: cashier.TotalSales,
		})
	}
	return cashierRecords
}

func mapResponseCashierMonthlyTotalSalesByMerchant(c []*db.GetMonthlyTotalSalesByMerchantRow) []*pb.CashierResponseMonthTotalSales {
	var cashierRecords []*pb.CashierResponseMonthTotalSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, &pb.CashierResponseMonthTotalSales{
			Year:       cashier.Year,
			Month:      cashier.Month,
			TotalSales: cashier.TotalSales,
		})
	}
	return cashierRecords
}

func mapResponseCashierYearlyTotalSales(c []*db.GetYearlyTotalSalesCashierRow) []*pb.CashierResponseYearTotalSales {
	var cashierRecords []*pb.CashierResponseYearTotalSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, &pb.CashierResponseYearTotalSales{
			Year:       cashier.Year,
			TotalSales: cashier.TotalSales,
		})
	}
	return cashierRecords
}

func mapResponseCashierYearlyTotalSalesById(c []*db.GetYearlyTotalSalesByIdRow) []*pb.CashierResponseYearTotalSales {
	var cashierRecords []*pb.CashierResponseYearTotalSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, &pb.CashierResponseYearTotalSales{
			Year:       cashier.Year,
			TotalSales: cashier.TotalSales,
		})
	}
	return cashierRecords
}

func mapResponseCashierYearlyTotalSalesByMerchant(c []*db.GetYearlyTotalSalesByMerchantRow) []*pb.CashierResponseYearTotalSales {
	var cashierRecords []*pb.CashierResponseYearTotalSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, &pb.CashierResponseYearTotalSales{
			Year:       cashier.Year,
			TotalSales: cashier.TotalSales,
		})
	}
	return cashierRecords
}
