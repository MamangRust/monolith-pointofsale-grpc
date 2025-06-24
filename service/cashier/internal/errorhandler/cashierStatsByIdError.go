package errorhandler

import (
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/cashier_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type cashierStatsByIdError struct {
	logger logger.LoggerInterface
}

func NewcashierStatsByIdError(logger logger.LoggerInterface) *cashierStatsByIdError {
	return &cashierStatsByIdError{logger: logger}
}

func (c *cashierStatsByIdError) HandleMonthlyTotalSalesByIdError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.CashierResponseMonthTotalSales, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CashierResponseMonthTotalSales](c.logger, err, method, tracePrefix, span, status, cashier_errors.ErrFailedFindMonthlyTotalSalesById, fields...)
}

func (c *cashierStatsByIdError) HandleYearlyTotalSalesByIdError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.CashierResponseYearTotalSales, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CashierResponseYearTotalSales](c.logger, err, method, tracePrefix, span, status, cashier_errors.ErrFailedFindYearlyTotalSalesById, fields...)
}

func (c *cashierStatsByIdError) HandleMonthlySalesByIdError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.CashierResponseMonthSales, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CashierResponseMonthSales](c.logger, err, method, tracePrefix, span, status, cashier_errors.ErrFailedFindMonthlyTotalSalesById, fields...)
}

func (c *cashierStatsByIdError) HandleYearlySalesByIdError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.CashierResponseYearSales, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CashierResponseYearSales](c.logger, err, method, tracePrefix, span, status, cashier_errors.ErrFailedFindYearlyTotalSalesById, fields...)
}
