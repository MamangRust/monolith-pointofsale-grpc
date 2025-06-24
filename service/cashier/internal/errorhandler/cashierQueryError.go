package errorhandler

import (
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/cashier_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type cashierQueryError struct {
	logger logger.LoggerInterface
}

func NewcashierQueryError(logger logger.LoggerInterface) *cashierQueryError {
	return &cashierQueryError{logger: logger}
}

func (t *cashierQueryError) HandleRepositoryPaginationError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.CashierResponse, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.CashierResponse](t.logger, err, method, tracePrefix, span, status, cashier_errors.ErrFailedFindAllCashiers, fields...)
}

func (t *cashierQueryError) HandleRepositoryPaginationDeleteAtError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) ([]*response.CashierResponseDeleteAt, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.CashierResponseDeleteAt](t.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

func (t *cashierQueryError) HandleRepositorySingleError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) (*response.CashierResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.CashierResponse](t.logger, err, method, tracePrefix, span, status, errResp, fields...)
}
