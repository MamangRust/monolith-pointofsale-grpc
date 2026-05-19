package handler

import (
	"context"
	"math"

	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/transaction_errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/pb"
	"github.com/MamangRust/monolith-point-of-sale-transacton/internal/service"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type transactionHandleGrpc struct {
	pb.UnimplementedTransactionServiceServer
	transactionQuery           service.TransactionQueryService
	transactionCommand         service.TransactionCommandService
	transactionStats           service.TransactionStatsService
	transactionStatsByMerchant service.TransactionStatsByMerchantService
	logger                     logger.LoggerInterface
}

func NewTransactionHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pb.TransactionServiceServer {
	return &transactionHandleGrpc{
		transactionQuery:           service.TransactionQuery,
		transactionCommand:         service.TransactionCommand,
		transactionStats:           service.TransactionStats,
		transactionStatsByMerchant: service.TransactionStatsByMerchant,
		logger:                     logger,
	}
}

func (s *transactionHandleGrpc) FindAll(ctx context.Context, request *pb.FindAllTransactionRequest) (*pb.ApiResponsePaginationTransaction, error) {
	s.logger.Info("FindAll transactions called", zap.Int32("page", request.GetPage()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllTransaction{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	transaction, totalRecords, err := s.transactionQuery.FindAllTransactions(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindAll transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindAll transactions success")

	return &pb.ApiResponsePaginationTransaction{
		Status:     "success",
		Message:    "Successfully fetched transaction",
		Data:       mapResponsesTransaction(transaction),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *transactionHandleGrpc) FindByMerchant(ctx context.Context, request *pb.FindAllTransactionMerchantRequest) (*pb.ApiResponsePaginationTransaction, error) {
	s.logger.Info("FindByMerchant transactions called", zap.Int32("merchantId", request.GetMerchantId()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()
	merchant_id := int(request.GetMerchantId())

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllTransactionByMerchant{
		MerchantID: merchant_id,
		Page:       page,
		PageSize:   pageSize,
		Search:     search,
	}

	transaction, totalRecords, err := s.transactionQuery.FindByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByMerchant transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByMerchant transactions success")

	return &pb.ApiResponsePaginationTransaction{
		Status:     "success",
		Message:    "Successfully fetched transaction",
		Data:       mapResponsesTransactionByMerchant(transaction),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *transactionHandleGrpc) FindById(ctx context.Context, request *pb.FindByIdTransactionRequest) (*pb.ApiResponseTransaction, error) {
	s.logger.Info("FindById transaction called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id == 0 {
		return nil, transaction_errors.ErrGrpcInvalidID
	}

	transaction, err := s.transactionQuery.FindById(ctx, id)
	if err != nil {
		s.logger.Error("FindById transaction failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindById transaction success")

	return &pb.ApiResponseTransaction{
		Status:  "success",
		Message: "Successfully fetched transaction",
		Data:    mapResponseTransaction(transaction),
	}, nil
}

func (s *transactionHandleGrpc) FindByActive(ctx context.Context, request *pb.FindAllTransactionRequest) (*pb.ApiResponsePaginationTransactionDeleteAt, error) {
	s.logger.Info("FindByActive transactions called", zap.Int32("page", request.GetPage()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllTransaction{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	transaction, totalRecords, err := s.transactionQuery.FindByActive(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByActive transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByActive transactions success")

	return &pb.ApiResponsePaginationTransactionDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active transaction",
		Data:       mapResponsesTransactionActive(transaction),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *transactionHandleGrpc) FindByTrashed(ctx context.Context, request *pb.FindAllTransactionRequest) (*pb.ApiResponsePaginationTransactionDeleteAt, error) {
	s.logger.Info("FindByTrashed transactions called")

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllTransaction{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	transaction, totalRecords, err := s.transactionQuery.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByTrashed transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pb.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByTrashed transactions success")

	return &pb.ApiResponsePaginationTransactionDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed transaction",
		Data:       mapResponsesTransactionTrashed(transaction),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *transactionHandleGrpc) FindMonthStatusSuccess(ctx context.Context, request *pb.FindMonthlyTransactionStatus) (*pb.ApiResponseTransactionMonthAmountSuccess, error) {
	s.logger.Info("FindMonthStatusSuccess transactions called", zap.Int32("year", request.GetYear()))

	year := int(request.GetYear())
	month := int(request.GetMonth())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}
	if month <= 0 || month >= 12 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	reqService := requests.MonthAmountTransaction{
		Year:  year,
		Month: month,
	}

	res, err := s.transactionStats.FindMonthlyAmountSuccess(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthStatusSuccess transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthStatusSuccess transactions success")

	return &pb.ApiResponseTransactionMonthAmountSuccess{
		Status:  "success",
		Message: "Monthly success data retrieved successfully",
		Data:    mapResponsesTransactionMonthlyAmountSuccess(res),
	}, nil
}

func (s *transactionHandleGrpc) FindYearStatusSuccess(ctx context.Context, request *pb.FindYearlyTransactionStatus) (*pb.ApiResponseTransactionYearAmountSuccess, error) {
	s.logger.Info("FindYearStatusSuccess transactions called", zap.Int32("year", request.GetYear()))

	year := int(request.GetYear())
	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	res, err := s.transactionStats.FindYearlyAmountSuccess(ctx, year)
	if err != nil {
		s.logger.Error("FindYearStatusSuccess transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearStatusSuccess transactions success")

	return &pb.ApiResponseTransactionYearAmountSuccess{
		Status:  "success",
		Message: "Yearly success data retrieved successfully",
		Data:    mapResponsesTransactionYearlyAmountSuccess(res),
	}, nil
}

func (s *transactionHandleGrpc) FindMonthStatusFailed(ctx context.Context, request *pb.FindMonthlyTransactionStatus) (*pb.ApiResponseTransactionMonthAmountFailed, error) {
	s.logger.Info("FindMonthStatusFailed transactions called", zap.Int32("year", request.GetYear()))

	year := int(request.GetYear())
	month := int(request.GetMonth())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}
	if month <= 0 || month >= 12 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	reqService := requests.MonthAmountTransaction{
		Year:  year,
		Month: month,
	}

	res, err := s.transactionStats.FindMonthlyAmountFailed(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthStatusFailed transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthStatusFailed transactions success")

	return &pb.ApiResponseTransactionMonthAmountFailed{
		Status:  "success",
		Message: "Monthly failed data retrieved successfully",
		Data:    mapResponsesTransactionMonthlyAmountFailed(res),
	}, nil
}

func (s *transactionHandleGrpc) FindYearStatusFailed(ctx context.Context, request *pb.FindYearlyTransactionStatus) (*pb.ApiResponseTransactionYearAmountFailed, error) {
	s.logger.Info("FindYearStatusFailed transactions called", zap.Int32("year", request.GetYear()))

	year := int(request.GetYear())
	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	res, err := s.transactionStats.FindYearlyAmountFailed(ctx, year)
	if err != nil {
		s.logger.Error("FindYearStatusFailed transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearStatusFailed transactions success")

	return &pb.ApiResponseTransactionYearAmountFailed{
		Status:  "success",
		Message: "Yearly failed data retrieved successfully",
		Data:    mapResponsesTransactionYearlyAmountFailed(res),
	}, nil
}

func (s *transactionHandleGrpc) FindMonthStatusSuccessByMerchant(ctx context.Context, request *pb.FindMonthlyTransactionStatusByMerchant) (*pb.ApiResponseTransactionMonthAmountSuccess, error) {
	s.logger.Info("FindMonthStatusSuccessByMerchant transactions called", zap.Int32("merchantId", request.GetMerchantId()))

	year := int(request.GetYear())
	month := int(request.GetMonth())
	id := int(request.GetMerchantId())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}
	if month <= 0 || month >= 12 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}
	if id <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidMerchantId
	}

	reqService := requests.MonthAmountTransactionMerchant{
		Year:       year,
		Month:      month,
		MerchantID: id,
	}

	res, err := s.transactionStatsByMerchant.FindMonthlyAmountSuccessByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthStatusSuccessByMerchant transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthStatusSuccessByMerchant transactions success")

	return &pb.ApiResponseTransactionMonthAmountSuccess{
		Status:  "success",
		Message: "Merchant monthly success data retrieved successfully",
		Data:    mapResponsesTransactionMonthlyAmountSuccessByMerchant(res),
	}, nil
}

func (s *transactionHandleGrpc) FindYearStatusSuccessByMerchant(ctx context.Context, request *pb.FindYearlyTransactionStatusByMerchant) (*pb.ApiResponseTransactionYearAmountSuccess, error) {
	s.logger.Info("FindYearStatusSuccessByMerchant transactions called", zap.Int32("merchantId", request.GetMerchantId()))

	year := int(request.GetYear())
	id := int(request.GetMerchantId())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}
	if id <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidMerchantId
	}

	reqService := requests.YearAmountTransactionMerchant{
		Year:       year,
		MerchantID: id,
	}

	res, err := s.transactionStatsByMerchant.FindYearlyAmountSuccessByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindYearStatusSuccessByMerchant transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearStatusSuccessByMerchant transactions success")

	return &pb.ApiResponseTransactionYearAmountSuccess{
		Status:  "success",
		Message: "Merchant yearly success data retrieved successfully",
		Data:    mapResponsesTransactionYearlyAmountSuccessByMerchant(res),
	}, nil
}

func (s *transactionHandleGrpc) FindMonthStatusFailedByMerchant(ctx context.Context, request *pb.FindMonthlyTransactionStatusByMerchant) (*pb.ApiResponseTransactionMonthAmountFailed, error) {
	s.logger.Info("FindMonthStatusFailedByMerchant transactions called", zap.Int32("merchantId", request.GetMerchantId()))

	year := int(request.GetYear())
	month := int(request.GetMonth())
	id := int(request.GetMerchantId())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}
	if month <= 0 || month >= 12 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}
	if id <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidMerchantId
	}

	reqService := requests.MonthAmountTransactionMerchant{
		Year:       year,
		Month:      month,
		MerchantID: id,
	}

	res, err := s.transactionStatsByMerchant.FindMonthlyAmountFailedByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthStatusFailedByMerchant transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthStatusFailedByMerchant transactions success")

	return &pb.ApiResponseTransactionMonthAmountFailed{
		Status:  "success",
		Message: "Merchant monthly failed data retrieved successfully",
		Data:    mapResponsesTransactionMonthlyAmountFailedByMerchant(res),
	}, nil
}

func (s *transactionHandleGrpc) FindYearStatusFailedByMerchant(ctx context.Context, request *pb.FindYearlyTransactionStatusByMerchant) (*pb.ApiResponseTransactionYearAmountFailed, error) {
	s.logger.Info("FindYearStatusFailedByMerchant transactions called", zap.Int32("merchantId", request.GetMerchantId()))

	year := int(request.GetYear())
	id := int(request.GetMerchantId())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}
	if id <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidMerchantId
	}

	reqService := requests.YearAmountTransactionMerchant{
		Year:       year,
		MerchantID: id,
	}

	res, err := s.transactionStatsByMerchant.FindYearlyAmountFailedByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindYearStatusFailedByMerchant transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearStatusFailedByMerchant transactions success")

	return &pb.ApiResponseTransactionYearAmountFailed{
		Status:  "success",
		Message: "Merchant yearly failed data retrieved successfully",
		Data:    mapResponsesTransactionYearlyAmountFailedByMerchant(res),
	}, nil
}

func (s *transactionHandleGrpc) FindMonthMethodSuccess(ctx context.Context, req *pb.MonthTransactionMethod) (*pb.ApiResponseTransactionMonthPaymentMethod, error) {
	s.logger.Info("FindMonthMethodSuccess transactions called", zap.Int32("year", req.GetYear()))

	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}
	if month <= 0 || month >= 12 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	methods, err := s.transactionStats.FindMonthlyMethodSuccess(ctx, &requests.MonthMethodTransaction{
		Year:  year,
		Month: month,
	})
	if err != nil {
		s.logger.Error("FindMonthMethodSuccess transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthMethodSuccess transactions success")

	return &pb.ApiResponseTransactionMonthPaymentMethod{
		Status:  "success",
		Message: "Monthly payment methods retrieved successfully",
		Data:    mapResponsesTransactionMonthlyMethodSuccess(methods),
	}, nil
}

func (s *transactionHandleGrpc) FindYearMethodSuccess(ctx context.Context, req *pb.YearTransactionMethod) (*pb.ApiResponseTransactionYearPaymentmethod, error) {
	s.logger.Info("FindYearMethodSuccess transactions called", zap.Int32("year", req.GetYear()))

	year := int(req.GetYear())
	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	methods, err := s.transactionStats.FindYearlyMethodSuccess(ctx, year)
	if err != nil {
		s.logger.Error("FindYearMethodSuccess transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearMethodSuccess transactions success")

	return &pb.ApiResponseTransactionYearPaymentmethod{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    mapResponsesTransactionYearlyMethodSuccess(methods),
	}, nil
}

func (s *transactionHandleGrpc) FindMonthMethodByMerchantSuccess(ctx context.Context, req *pb.MonthTransactionMethodByMerchant) (*pb.ApiResponseTransactionMonthPaymentMethod, error) {
	s.logger.Info("FindMonthMethodByMerchantSuccess transactions called", zap.Int32("merchantId", req.GetMerchantId()))

	year := int(req.GetYear())
	id := int(req.GetMerchantId())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}
	if id <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidMerchantId
	}
	if month <= 0 || month >= 12 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	reqService := requests.MonthMethodTransactionMerchant{
		Year:       year,
		MerchantID: id,
		Month:      month,
	}

	methods, err := s.transactionStatsByMerchant.FindMonthlyMethodByMerchantSuccess(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthMethodByMerchantSuccess transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthMethodByMerchantSuccess transactions success")

	return &pb.ApiResponseTransactionMonthPaymentMethod{
		Status:  "success",
		Message: "Merchant monthly payment methods retrieved successfully",
		Data:    mapResponsesTransactionMonthlyMethodByMerchantSuccess(methods),
	}, nil
}

func (s *transactionHandleGrpc) FindYearMethodByMerchantSuccess(ctx context.Context, req *pb.YearTransactionMethodByMerchant) (*pb.ApiResponseTransactionYearPaymentmethod, error) {
	s.logger.Info("FindYearMethodByMerchantSuccess transactions called", zap.Int32("merchantId", req.GetMerchantId()))

	year := int(req.GetYear())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}
	if id <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidMerchantId
	}

	reqService := requests.YearMethodTransactionMerchant{
		Year:       year,
		MerchantID: id,
	}

	methods, err := s.transactionStatsByMerchant.FindYearlyMethodByMerchantSuccess(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindYearMethodByMerchantSuccess transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearMethodByMerchantSuccess transactions success")

	return &pb.ApiResponseTransactionYearPaymentmethod{
		Status:  "success",
		Message: "Merchant yearly payment methods retrieved successfully",
		Data:    mapResponsesTransactionYearlyMethodByMerchantSuccess(methods),
	}, nil
}

func (s *transactionHandleGrpc) FindMonthMethodFailed(ctx context.Context, req *pb.MonthTransactionMethod) (*pb.ApiResponseTransactionMonthPaymentMethod, error) {
	s.logger.Info("FindMonthMethodFailed transactions called", zap.Int32("year", req.GetYear()))

	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}
	if month <= 0 || month >= 12 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	methods, err := s.transactionStats.FindMonthlyMethodFailed(ctx, &requests.MonthMethodTransaction{
		Year:  year,
		Month: month,
	})
	if err != nil {
		s.logger.Error("FindMonthMethodFailed transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthMethodFailed transactions success")

	return &pb.ApiResponseTransactionMonthPaymentMethod{
		Status:  "Failed",
		Message: "Monthly payment methods retrieved Failedfully",
		Data:    mapResponsesTransactionMonthlyMethodFailed(methods),
	}, nil
}

func (s *transactionHandleGrpc) FindYearMethodFailed(ctx context.Context, req *pb.YearTransactionMethod) (*pb.ApiResponseTransactionYearPaymentmethod, error) {
	s.logger.Info("FindYearMethodFailed transactions called", zap.Int32("year", req.GetYear()))

	year := int(req.GetYear())
	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	methods, err := s.transactionStats.FindYearlyMethodFailed(ctx, year)
	if err != nil {
		s.logger.Error("FindYearMethodFailed transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearMethodFailed transactions success")

	return &pb.ApiResponseTransactionYearPaymentmethod{
		Status:  "Failed",
		Message: "Yearly payment methods retrieved Failedfully",
		Data:    mapResponsesTransactionYearlyMethodFailed(methods),
	}, nil
}

func (s *transactionHandleGrpc) FindMonthMethodByMerchantFailed(ctx context.Context, req *pb.MonthTransactionMethodByMerchant) (*pb.ApiResponseTransactionMonthPaymentMethod, error) {
	s.logger.Info("FindMonthMethodByMerchantFailed transactions called", zap.Int32("merchantId", req.GetMerchantId()))

	year := int(req.GetYear())
	id := int(req.GetMerchantId())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}
	if id <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidMerchantId
	}
	if month <= 0 || month >= 12 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	reqService := requests.MonthMethodTransactionMerchant{
		Year:       year,
		MerchantID: id,
		Month:      month,
	}

	methods, err := s.transactionStatsByMerchant.FindMonthlyMethodByMerchantFailed(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthMethodByMerchantFailed transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthMethodByMerchantFailed transactions success")

	return &pb.ApiResponseTransactionMonthPaymentMethod{
		Status:  "Failed",
		Message: "Merchant monthly payment methods retrieved Failedfully",
		Data:    mapResponsesTransactionMonthlyMethodByMerchantFailed(methods),
	}, nil
}

func (s *transactionHandleGrpc) FindYearlyMethodByMerchantFailed(ctx context.Context, req *pb.YearTransactionMethodByMerchant) (*pb.ApiResponseTransactionYearPaymentmethod, error) {
	s.logger.Info("FindYearlyMethodByMerchantFailed transactions called", zap.Int32("merchantId", req.GetMerchantId()))

	year := int(req.GetYear())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}
	if id <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidMerchantId
	}

	reqService := requests.YearMethodTransactionMerchant{
		Year:       year,
		MerchantID: id,
	}

	methods, err := s.transactionStatsByMerchant.FindYearlyMethodByMerchantFailed(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindYearlyMethodByMerchantFailed transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearlyMethodByMerchantFailed transactions success")

	return &pb.ApiResponseTransactionYearPaymentmethod{
		Status:  "Failed",
		Message: "Merchant yearly payment methods retrieved Failedfully",
		Data:    mapResponsesTransactionYearlyMethodByMerchantFailed(methods),
	}, nil
}

func (s *transactionHandleGrpc) Create(ctx context.Context, request *pb.CreateTransactionRequest) (*pb.ApiResponseTransaction, error) {
	s.logger.Info("Create transaction called", zap.Int32("orderId", request.GetOrderId()))

	req := &requests.CreateTransactionRequest{
		CashierID:     int(request.GetCashierId()),
		OrderID:       int(request.GetOrderId()),
		PaymentMethod: request.GetPaymentMethod(),
		Amount:        int(request.GetAmount()),
	}

	if err := req.Validate(); err != nil {
		return nil, transaction_errors.ErrGrpcValidateCreateTransaction
	}

	transaction, err := s.transactionCommand.CreateTransaction(ctx, req)
	if err != nil {
		s.logger.Error("Create transaction failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Create transaction success")

	return &pb.ApiResponseTransaction{
		Status:  "success",
		Message: "Successfully created transaction",
		Data:    mapResponseTransaction(transaction),
	}, nil
}

func (s *transactionHandleGrpc) Update(ctx context.Context, request *pb.UpdateTransactionRequest) (*pb.ApiResponseTransaction, error) {
	s.logger.Info("Update transaction called", zap.Int32("id", request.GetTransactionId()))

	id := int(request.GetTransactionId())
	if id == 0 {
		return nil, transaction_errors.ErrGrpcInvalidID
	}

	req := &requests.UpdateTransactionRequest{
		TransactionID: &id,
		OrderID:       int(request.GetOrderId()),
		PaymentMethod: request.GetPaymentMethod(),
		Amount:        int(request.GetAmount()),
	}

	if err := req.Validate(); err != nil {
		return nil, transaction_errors.ErrGrpcValidateUpdateTransaction
	}

	transaction, err := s.transactionCommand.UpdateTransaction(ctx, req)
	if err != nil {
		s.logger.Error("Update transaction failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Update transaction success")

	return &pb.ApiResponseTransaction{
		Status:  "success",
		Message: "Successfully updated transaction",
		Data:    mapResponseTransaction(transaction),
	}, nil
}

func (s *transactionHandleGrpc) TrashedTransaction(ctx context.Context, request *pb.FindByIdTransactionRequest) (*pb.ApiResponseTransactionDeleteAt, error) {
	s.logger.Info("TrashedTransaction called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id == 0 {
		return nil, transaction_errors.ErrGrpcInvalidID
	}

	transaction, err := s.transactionCommand.TrashedTransaction(ctx, id)
	if err != nil {
		s.logger.Error("TrashedTransaction failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("TrashedTransaction success")

	return &pb.ApiResponseTransactionDeleteAt{
		Status:  "success",
		Message: "Successfully trashed transaction",
		Data:    mapResponseTransactionDeleteAt(transaction),
	}, nil
}

func (s *transactionHandleGrpc) RestoreTransaction(ctx context.Context, request *pb.FindByIdTransactionRequest) (*pb.ApiResponseTransactionDeleteAt, error) {
	s.logger.Info("RestoreTransaction called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id == 0 {
		return nil, transaction_errors.ErrGrpcInvalidID
	}

	transaction, err := s.transactionCommand.RestoreTransaction(ctx, id)
	if err != nil {
		s.logger.Error("RestoreTransaction failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreTransaction success")

	return &pb.ApiResponseTransactionDeleteAt{
		Status:  "success",
		Message: "Successfully restored transaction",
		Data:    mapResponseTransactionDeleteAt(transaction),
	}, nil
}

func (s *transactionHandleGrpc) DeleteTransactionPermanent(ctx context.Context, request *pb.FindByIdTransactionRequest) (*pb.ApiResponseTransactionDelete, error) {
	s.logger.Info("DeleteTransactionPermanent called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id == 0 {
		return nil, transaction_errors.ErrGrpcInvalidID
	}

	_, err := s.transactionCommand.DeleteTransactionPermanently(ctx, id)
	if err != nil {
		s.logger.Error("DeleteTransactionPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteTransactionPermanent success")

	return &pb.ApiResponseTransactionDelete{
		Status:  "success",
		Message: "Successfully deleted Transaction permanently",
	}, nil
}

func (s *transactionHandleGrpc) RestoreAllTransaction(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseTransactionAll, error) {
	s.logger.Info("RestoreAllTransaction called")

	_, err := s.transactionCommand.RestoreAllTransactions(ctx)
	if err != nil {
		s.logger.Error("RestoreAllTransaction failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreAllTransaction success")

	return &pb.ApiResponseTransactionAll{
		Status:  "success",
		Message: "Successfully restore all Transaction",
	}, nil
}

func (s *transactionHandleGrpc) DeleteAllTransactionPermanent(ctx context.Context, _ *emptypb.Empty) (*pb.ApiResponseTransactionAll, error) {
	s.logger.Info("DeleteAllTransactionPermanent called")

	_, err := s.transactionCommand.DeleteAllTransactionPermanent(ctx)
	if err != nil {
		s.logger.Error("DeleteAllTransactionPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteAllTransactionPermanent success")

	return &pb.ApiResponseTransactionAll{
		Status:  "success",
		Message: "Successfully delete Transaction permanen",
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

func mapResponseTransaction(transaction *db.Transaction) *pb.TransactionResponse {
	if transaction == nil {
		return nil
	}
	var createdAtStr string
	if transaction.CreatedAt.Valid {
		createdAtStr = transaction.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}
	var updatedAtStr string
	if transaction.UpdatedAt.Valid {
		updatedAtStr = transaction.UpdatedAt.Time.Format("2006-01-02 15:04:05")
	}
	var changeAmount int32
	if transaction.ChangeAmount != nil {
		changeAmount = *transaction.ChangeAmount
	}
	return &pb.TransactionResponse{
		Id:            transaction.TransactionID,
		OrderId:       transaction.OrderID,
		MerchantId:    transaction.MerchantID,
		PaymentMethod: transaction.PaymentMethod,
		Amount:        transaction.Amount,
		ChangeAmount:  changeAmount,
		PaymentStatus: transaction.PaymentStatus,
		CreatedAt:     createdAtStr,
		UpdatedAt:     updatedAtStr,
	}
}

func mapResponsesTransaction(transactions []*db.GetTransactionsRow) []*pb.TransactionResponse {
	var mappedTransactions []*pb.TransactionResponse
	for _, t := range transactions {
		if t == nil {
			continue
		}
		var createdAtStr string
		if t.CreatedAt.Valid {
			createdAtStr = t.CreatedAt.Time.Format("2006-01-02 15:04:05")
		}
		var updatedAtStr string
		if t.UpdatedAt.Valid {
			updatedAtStr = t.UpdatedAt.Time.Format("2006-01-02 15:04:05")
		}
		var changeAmount int32
		if t.ChangeAmount != nil {
			changeAmount = *t.ChangeAmount
		}
		mappedTransactions = append(mappedTransactions, &pb.TransactionResponse{
			Id:            t.TransactionID,
			OrderId:       t.OrderID,
			MerchantId:    t.MerchantID,
			PaymentMethod: t.PaymentMethod,
			Amount:        t.Amount,
			ChangeAmount:  changeAmount,
			PaymentStatus: t.PaymentStatus,
			CreatedAt:     createdAtStr,
			UpdatedAt:     updatedAtStr,
		})
	}
	return mappedTransactions
}

func mapResponsesTransactionByMerchant(transactions []*db.GetTransactionByMerchantRow) []*pb.TransactionResponse {
	var mappedTransactions []*pb.TransactionResponse
	for _, t := range transactions {
		if t == nil {
			continue
		}
		var createdAtStr string
		if t.CreatedAt.Valid {
			createdAtStr = t.CreatedAt.Time.Format("2006-01-02 15:04:05")
		}
		var updatedAtStr string
		if t.UpdatedAt.Valid {
			updatedAtStr = t.UpdatedAt.Time.Format("2006-01-02 15:04:05")
		}
		var changeAmount int32
		if t.ChangeAmount != nil {
			changeAmount = *t.ChangeAmount
		}
		mappedTransactions = append(mappedTransactions, &pb.TransactionResponse{
			Id:            t.TransactionID,
			OrderId:       t.OrderID,
			MerchantId:    t.MerchantID,
			PaymentMethod: t.PaymentMethod,
			Amount:        t.Amount,
			ChangeAmount:  changeAmount,
			PaymentStatus: t.PaymentStatus,
			CreatedAt:     createdAtStr,
			UpdatedAt:     updatedAtStr,
		})
	}
	return mappedTransactions
}

func mapResponseTransactionDeleteAt(transaction *db.Transaction) *pb.TransactionResponseDeleteAt {
	if transaction == nil {
		return nil
	}
	var createdAtStr string
	if transaction.CreatedAt.Valid {
		createdAtStr = transaction.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}
	var updatedAtStr string
	if transaction.UpdatedAt.Valid {
		updatedAtStr = transaction.UpdatedAt.Time.Format("2006-01-02 15:04:05")
	}
	var deletedAt *wrapperspb.StringValue
	if transaction.DeletedAt.Valid {
		deletedAt = wrapperspb.String(transaction.DeletedAt.Time.Format("2006-01-02 15:04:05"))
	}
	var changeAmount int32
	if transaction.ChangeAmount != nil {
		changeAmount = *transaction.ChangeAmount
	}

	return &pb.TransactionResponseDeleteAt{
		Id:            transaction.TransactionID,
		OrderId:       transaction.OrderID,
		MerchantId:    transaction.MerchantID,
		PaymentMethod: transaction.PaymentMethod,
		Amount:        transaction.Amount,
		ChangeAmount:  changeAmount,
		PaymentStatus: transaction.PaymentStatus,
		CreatedAt:     createdAtStr,
		UpdatedAt:     updatedAtStr,
		DeletedAt:     deletedAt,
	}
}

func mapResponsesTransactionActive(transactions []*db.GetTransactionsActiveRow) []*pb.TransactionResponseDeleteAt {
	var mappedTransactions []*pb.TransactionResponseDeleteAt
	for _, t := range transactions {
		if t == nil {
			continue
		}
		var createdAtStr string
		if t.CreatedAt.Valid {
			createdAtStr = t.CreatedAt.Time.Format("2006-01-02 15:04:05")
		}
		var updatedAtStr string
		if t.UpdatedAt.Valid {
			updatedAtStr = t.UpdatedAt.Time.Format("2006-01-02 15:04:05")
		}
		var deletedAt *wrapperspb.StringValue
		if t.DeletedAt.Valid {
			deletedAt = wrapperspb.String(t.DeletedAt.Time.Format("2006-01-02 15:04:05"))
		}
		var changeAmount int32
		if t.ChangeAmount != nil {
			changeAmount = *t.ChangeAmount
		}
		mappedTransactions = append(mappedTransactions, &pb.TransactionResponseDeleteAt{
			Id:            t.TransactionID,
			OrderId:       t.OrderID,
			MerchantId:    t.MerchantID,
			PaymentMethod: t.PaymentMethod,
			Amount:        t.Amount,
			ChangeAmount:  changeAmount,
			PaymentStatus: t.PaymentStatus,
			CreatedAt:     createdAtStr,
			UpdatedAt:     updatedAtStr,
			DeletedAt:     deletedAt,
		})
	}
	return mappedTransactions
}

func mapResponsesTransactionTrashed(transactions []*db.GetTransactionsTrashedRow) []*pb.TransactionResponseDeleteAt {
	var mappedTransactions []*pb.TransactionResponseDeleteAt
	for _, t := range transactions {
		if t == nil {
			continue
		}
		var createdAtStr string
		if t.CreatedAt.Valid {
			createdAtStr = t.CreatedAt.Time.Format("2006-01-02 15:04:05")
		}
		var updatedAtStr string
		if t.UpdatedAt.Valid {
			updatedAtStr = t.UpdatedAt.Time.Format("2006-01-02 15:04:05")
		}
		var deletedAt *wrapperspb.StringValue
		if t.DeletedAt.Valid {
			deletedAt = wrapperspb.String(t.DeletedAt.Time.Format("2006-01-02 15:04:05"))
		}
		var changeAmount int32
		if t.ChangeAmount != nil {
			changeAmount = *t.ChangeAmount
		}
		mappedTransactions = append(mappedTransactions, &pb.TransactionResponseDeleteAt{
			Id:            t.TransactionID,
			OrderId:       t.OrderID,
			MerchantId:    t.MerchantID,
			PaymentMethod: t.PaymentMethod,
			Amount:        t.Amount,
			ChangeAmount:  changeAmount,
			PaymentStatus: t.PaymentStatus,
			CreatedAt:     createdAtStr,
			UpdatedAt:     updatedAtStr,
			DeletedAt:     deletedAt,
		})
	}
	return mappedTransactions
}

func mapResponseTransactionMonthAmountSuccess(row *db.GetMonthlyAmountTransactionSuccessRow) *pb.TransactionMonthlyAmountSuccess {
	if row == nil {
		return nil
	}
	return &pb.TransactionMonthlyAmountSuccess{
		Year:         row.Year,
		Month:        row.Month,
		TotalSuccess: int32(row.TotalSuccess),
		TotalAmount:  int32(row.TotalAmount),
	}
}

func mapResponsesTransactionMonthlyAmountSuccess(rows []*db.GetMonthlyAmountTransactionSuccessRow) []*pb.TransactionMonthlyAmountSuccess {
	var transaction []*pb.TransactionMonthlyAmountSuccess
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionMonthAmountSuccess(row))
	}
	return transaction
}

func mapResponseTransactionYearAmountSuccess(row *db.GetYearlyAmountTransactionSuccessRow) *pb.TransactionYearlyAmountSuccess {
	if row == nil {
		return nil
	}
	return &pb.TransactionYearlyAmountSuccess{
		Year:         row.Year,
		TotalSuccess: int32(row.TotalSuccess),
		TotalAmount:  int32(row.TotalAmount),
	}
}

func mapResponsesTransactionYearlyAmountSuccess(rows []*db.GetYearlyAmountTransactionSuccessRow) []*pb.TransactionYearlyAmountSuccess {
	var transaction []*pb.TransactionYearlyAmountSuccess
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionYearAmountSuccess(row))
	}
	return transaction
}

func mapResponseTransactionMonthAmountFailed(row *db.GetMonthlyAmountTransactionFailedRow) *pb.TransactionMonthlyAmountFailed {
	if row == nil {
		return nil
	}
	return &pb.TransactionMonthlyAmountFailed{
		Year:        row.Year,
		Month:       row.Month,
		TotalFailed: int32(row.TotalFailed),
		TotalAmount: int32(row.TotalAmount),
	}
}

func mapResponsesTransactionMonthlyAmountFailed(rows []*db.GetMonthlyAmountTransactionFailedRow) []*pb.TransactionMonthlyAmountFailed {
	var transaction []*pb.TransactionMonthlyAmountFailed
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionMonthAmountFailed(row))
	}
	return transaction
}

func mapResponseTransactionYearAmountFailed(row *db.GetYearlyAmountTransactionFailedRow) *pb.TransactionYearlyAmountFailed {
	if row == nil {
		return nil
	}
	return &pb.TransactionYearlyAmountFailed{
		Year:        row.Year,
		TotalFailed: int32(row.TotalFailed),
		TotalAmount: int32(row.TotalAmount),
	}
}

func mapResponsesTransactionYearlyAmountFailed(rows []*db.GetYearlyAmountTransactionFailedRow) []*pb.TransactionYearlyAmountFailed {
	var transaction []*pb.TransactionYearlyAmountFailed
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionYearAmountFailed(row))
	}
	return transaction
}

func mapResponseTransactionMonthAmountSuccessByMerchant(row *db.GetMonthlyAmountTransactionSuccessByMerchantRow) *pb.TransactionMonthlyAmountSuccess {
	if row == nil {
		return nil
	}
	return &pb.TransactionMonthlyAmountSuccess{
		Year:         row.Year,
		Month:        row.Month,
		TotalSuccess: int32(row.TotalSuccess),
		TotalAmount:  int32(row.TotalAmount),
	}
}

func mapResponsesTransactionMonthlyAmountSuccessByMerchant(rows []*db.GetMonthlyAmountTransactionSuccessByMerchantRow) []*pb.TransactionMonthlyAmountSuccess {
	var transaction []*pb.TransactionMonthlyAmountSuccess
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionMonthAmountSuccessByMerchant(row))
	}
	return transaction
}

func mapResponseTransactionYearAmountSuccessByMerchant(row *db.GetYearlyAmountTransactionSuccessByMerchantRow) *pb.TransactionYearlyAmountSuccess {
	if row == nil {
		return nil
	}
	return &pb.TransactionYearlyAmountSuccess{
		Year:         row.Year,
		TotalSuccess: int32(row.TotalSuccess),
		TotalAmount:  int32(row.TotalAmount),
	}
}

func mapResponsesTransactionYearlyAmountSuccessByMerchant(rows []*db.GetYearlyAmountTransactionSuccessByMerchantRow) []*pb.TransactionYearlyAmountSuccess {
	var transaction []*pb.TransactionYearlyAmountSuccess
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionYearAmountSuccessByMerchant(row))
	}
	return transaction
}

func mapResponseTransactionMonthAmountFailedByMerchant(row *db.GetMonthlyAmountTransactionFailedByMerchantRow) *pb.TransactionMonthlyAmountFailed {
	if row == nil {
		return nil
	}
	return &pb.TransactionMonthlyAmountFailed{
		Year:        row.Year,
		Month:       row.Month,
		TotalFailed: int32(row.TotalFailed),
		TotalAmount: int32(row.TotalAmount),
	}
}

func mapResponsesTransactionMonthlyAmountFailedByMerchant(rows []*db.GetMonthlyAmountTransactionFailedByMerchantRow) []*pb.TransactionMonthlyAmountFailed {
	var transaction []*pb.TransactionMonthlyAmountFailed
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionMonthAmountFailedByMerchant(row))
	}
	return transaction
}

func mapResponseTransactionYearAmountFailedByMerchant(row *db.GetYearlyAmountTransactionFailedByMerchantRow) *pb.TransactionYearlyAmountFailed {
	if row == nil {
		return nil
	}
	return &pb.TransactionYearlyAmountFailed{
		Year:        row.Year,
		TotalFailed: int32(row.TotalFailed),
		TotalAmount: int32(row.TotalAmount),
	}
}

func mapResponsesTransactionYearlyAmountFailedByMerchant(rows []*db.GetYearlyAmountTransactionFailedByMerchantRow) []*pb.TransactionYearlyAmountFailed {
	var transaction []*pb.TransactionYearlyAmountFailed
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionYearAmountFailedByMerchant(row))
	}
	return transaction
}

func mapResponseTransactionMonthMethodSuccess(row *db.GetMonthlyTransactionMethodsSuccessRow) *pb.TransactionMonthlyMethod {
	if row == nil {
		return nil
	}
	return &pb.TransactionMonthlyMethod{
		Month:             row.Month,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int32(row.TotalTransactions),
		TotalAmount:       int32(row.TotalAmount),
	}
}

func mapResponsesTransactionMonthlyMethodSuccess(rows []*db.GetMonthlyTransactionMethodsSuccessRow) []*pb.TransactionMonthlyMethod {
	var transaction []*pb.TransactionMonthlyMethod
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionMonthMethodSuccess(row))
	}
	return transaction
}

func mapResponseTransactionYearlyMethodSuccess(row *db.GetYearlyTransactionMethodsSuccessRow) *pb.TransactionYearlyMethod {
	if row == nil {
		return nil
	}
	return &pb.TransactionYearlyMethod{
		Year:              row.Year,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int32(row.TotalTransactions),
		TotalAmount:       int32(row.TotalAmount),
	}
}

func mapResponsesTransactionYearlyMethodSuccess(rows []*db.GetYearlyTransactionMethodsSuccessRow) []*pb.TransactionYearlyMethod {
	var transaction []*pb.TransactionYearlyMethod
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionYearlyMethodSuccess(row))
	}
	return transaction
}

func mapResponseTransactionMonthMethodFailed(row *db.GetMonthlyTransactionMethodsFailedRow) *pb.TransactionMonthlyMethod {
	if row == nil {
		return nil
	}
	return &pb.TransactionMonthlyMethod{
		Month:             row.Month,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int32(row.TotalTransactions),
		TotalAmount:       int32(row.TotalAmount),
	}
}

func mapResponsesTransactionMonthlyMethodFailed(rows []*db.GetMonthlyTransactionMethodsFailedRow) []*pb.TransactionMonthlyMethod {
	var transaction []*pb.TransactionMonthlyMethod
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionMonthMethodFailed(row))
	}
	return transaction
}

func mapResponseTransactionYearlyMethodFailed(row *db.GetYearlyTransactionMethodsFailedRow) *pb.TransactionYearlyMethod {
	if row == nil {
		return nil
	}
	return &pb.TransactionYearlyMethod{
		Year:              row.Year,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int32(row.TotalTransactions),
		TotalAmount:       int32(row.TotalAmount),
	}
}

func mapResponsesTransactionYearlyMethodFailed(rows []*db.GetYearlyTransactionMethodsFailedRow) []*pb.TransactionYearlyMethod {
	var transaction []*pb.TransactionYearlyMethod
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionYearlyMethodFailed(row))
	}
	return transaction
}

func mapResponseTransactionMonthMethodByMerchantSuccess(row *db.GetMonthlyTransactionMethodsByMerchantSuccessRow) *pb.TransactionMonthlyMethod {
	if row == nil {
		return nil
	}
	return &pb.TransactionMonthlyMethod{
		Month:             row.Month,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int32(row.TotalTransactions),
		TotalAmount:       int32(row.TotalAmount),
	}
}

func mapResponsesTransactionMonthlyMethodByMerchantSuccess(rows []*db.GetMonthlyTransactionMethodsByMerchantSuccessRow) []*pb.TransactionMonthlyMethod {
	var transaction []*pb.TransactionMonthlyMethod
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionMonthMethodByMerchantSuccess(row))
	}
	return transaction
}

func mapResponseTransactionYearlyMethodByMerchantSuccess(row *db.GetYearlyTransactionMethodsByMerchantSuccessRow) *pb.TransactionYearlyMethod {
	if row == nil {
		return nil
	}
	return &pb.TransactionYearlyMethod{
		Year:              row.Year,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int32(row.TotalTransactions),
		TotalAmount:       int32(row.TotalAmount),
	}
}

func mapResponsesTransactionYearlyMethodByMerchantSuccess(rows []*db.GetYearlyTransactionMethodsByMerchantSuccessRow) []*pb.TransactionYearlyMethod {
	var transaction []*pb.TransactionYearlyMethod
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionYearlyMethodByMerchantSuccess(row))
	}
	return transaction
}

func mapResponseTransactionMonthMethodByMerchantFailed(row *db.GetMonthlyTransactionMethodsByMerchantFailedRow) *pb.TransactionMonthlyMethod {
	if row == nil {
		return nil
	}
	return &pb.TransactionMonthlyMethod{
		Month:             row.Month,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int32(row.TotalTransactions),
		TotalAmount:       int32(row.TotalAmount),
	}
}

func mapResponsesTransactionMonthlyMethodByMerchantFailed(rows []*db.GetMonthlyTransactionMethodsByMerchantFailedRow) []*pb.TransactionMonthlyMethod {
	var transaction []*pb.TransactionMonthlyMethod
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionMonthMethodByMerchantFailed(row))
	}
	return transaction
}

func mapResponseTransactionYearlyMethodByMerchantFailed(row *db.GetYearlyTransactionMethodsByMerchantFailedRow) *pb.TransactionYearlyMethod {
	if row == nil {
		return nil
	}
	return &pb.TransactionYearlyMethod{
		Year:              row.Year,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int32(row.TotalTransactions),
		TotalAmount:       int32(row.TotalAmount),
	}
}

func mapResponsesTransactionYearlyMethodByMerchantFailed(rows []*db.GetYearlyTransactionMethodsByMerchantFailedRow) []*pb.TransactionYearlyMethod {
	var transaction []*pb.TransactionYearlyMethod
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionYearlyMethodByMerchantFailed(row))
	}
	return transaction
}
