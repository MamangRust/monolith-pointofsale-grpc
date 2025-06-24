package errorhandler

import (
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/cashier_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type cashierStatsByMerchantError struct {
	logger logger.LoggerInterface
}

func NewcashierStatsByMerchantError(logger logger.LoggerInterface) *cashierStatsByMerchantError {
	return &cashierStatsByMerchantError{logger: logger}
}

func (c *cashierStatsByMerchantError) HandleMonthlyTotalSalesByMerchantError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.CashierResponseMonthTotalSales, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CashierResponseMonthTotalSales](c.logger, err, method, tracePrefix, span, status, cashier_errors.ErrFailedFindMonthlyTotalSalesByMerchant, fields...)
}

func (c *cashierStatsByMerchantError) HandleYearlyTotalSalesByMerchantError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.CashierResponseYearTotalSales, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CashierResponseYearTotalSales](c.logger, err, method, tracePrefix, span, status, cashier_errors.ErrFailedFindYearlyTotalSalesByMerchant, fields...)
}

func (c *cashierStatsByMerchantError) HandleMonthlySalesByMerchantError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.CashierResponseMonthSales, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CashierResponseMonthSales](c.logger, err, method, tracePrefix, span, status, cashier_errors.ErrFailedFindMonthlyTotalSalesByMerchant, fields...)
}

func (c *cashierStatsByMerchantError) HandleYearlySalesByMerchantError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) ([]*response.CashierResponseYearSales, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.CashierResponseYearSales](c.logger, err, method, tracePrefix, span, status, cashier_errors.ErrFailedFindYearlyTotalSalesByMerchant, fields...)
}
