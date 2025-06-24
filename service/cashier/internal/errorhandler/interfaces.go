package errorhandler

import (
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type CashierCommadError interface {
	HandleRepositorySingleError(err error, method, tracePrefix string, span trace.Span, status *string, errResp *response.ErrorResponse, fields ...zap.Field) (*response.CashierResponse, *response.ErrorResponse)
	HandleCreateCashierError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.CashierResponse, *response.ErrorResponse)
	HandleUpdateCashierError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.CashierResponse, *response.ErrorResponse)
	HandleTrashedCashierError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.CashierResponseDeleteAt, *response.ErrorResponse)
	HandleRestoreCashierError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.CashierResponseDeleteAt, *response.ErrorResponse)
	HandleDeleteCashierPermanentError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleRestoreAllCashierError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleDeleteAllCashierPermanentError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
}

type CashierQueryError interface {
	HandleRepositoryPaginationError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.CashierResponse, *int, *response.ErrorResponse)
	HandleRepositoryPaginationDeleteAtError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) ([]*response.CashierResponseDeleteAt, *int, *response.ErrorResponse)
	HandleRepositorySingleError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) (*response.CashierResponse, *response.ErrorResponse)
}

type CashierStatsError interface {
	HandleMonthlyTotalSalesError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.CashierResponseMonthTotalSales, *response.ErrorResponse)
	HandleYearlyTotalSalesError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.CashierResponseYearTotalSales, *response.ErrorResponse)
	HandleMonthlySalesError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.CashierResponseMonthSales, *response.ErrorResponse)
	HandleYearlySalesError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.CashierResponseYearSales, *response.ErrorResponse)
}

type CashierStatsByIdError interface {
	HandleMonthlyTotalSalesByIdError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.CashierResponseMonthTotalSales, *response.ErrorResponse)
	HandleYearlyTotalSalesByIdError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.CashierResponseYearTotalSales, *response.ErrorResponse)
	HandleMonthlySalesByIdError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.CashierResponseMonthSales, *response.ErrorResponse)
	HandleYearlySalesByIdError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.CashierResponseYearSales, *response.ErrorResponse)
}

type CashierStatsByMerchantError interface {
	HandleMonthlyTotalSalesByMerchantError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.CashierResponseMonthTotalSales, *response.ErrorResponse)
	HandleYearlyTotalSalesByMerchantError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.CashierResponseYearTotalSales, *response.ErrorResponse)
	HandleMonthlySalesByMerchantError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.CashierResponseMonthSales, *response.ErrorResponse)
	HandleYearlySalesByMerchantError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.CashierResponseYearSales, *response.ErrorResponse)
}
