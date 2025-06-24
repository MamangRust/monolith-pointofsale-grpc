package errorhandler

import (
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/cashier_errors"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type cashierCommandError struct {
	logger logger.LoggerInterface
}

func NewCashierCommandError(logger logger.LoggerInterface) *cashierCommandError {
	return &cashierCommandError{logger: logger}
}

func (t *cashierCommandError) HandleRepositorySingleError(err error, method, tracePrefix string, span trace.Span, status *string, errResp *response.ErrorResponse, fields ...zap.Field) (*response.CashierResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.CashierResponse](t.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

func (t *cashierCommandError) HandleCreateCashierError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.CashierResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.CashierResponse](t.logger, err, method, tracePrefix, span, status, cashier_errors.ErrFailedCreateCashier, fields...)
}

func (t *cashierCommandError) HandleUpdateCashierError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.CashierResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.CashierResponse](t.logger, err, method, tracePrefix, span, status, cashier_errors.ErrFailedUpdateCashier, fields...)
}

func (t *cashierCommandError) HandleTrashedCashierError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.CashierResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.CashierResponseDeleteAt](t.logger, err, method, tracePrefix, span, status, cashier_errors.ErrFailedTrashedCashier, fields...)
}

func (t *cashierCommandError) HandleRestoreCashierError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.CashierResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.CashierResponseDeleteAt](t.logger, err, method, tracePrefix, span, status, cashier_errors.ErrFailedRestoreCashier, fields...)
}

func (t *cashierCommandError) HandleDeleteCashierPermanentError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](t.logger, err, method, tracePrefix, span, status, cashier_errors.ErrFailedDeleteCashierPermanent, fields...)
}

func (t *cashierCommandError) HandleRestoreAllCashierError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](t.logger, err, method, tracePrefix, span, status, cashier_errors.ErrFailedRestoreAllCashiers, fields...)
}

func (t *cashierCommandError) HandleDeleteAllCashierPermanentError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](t.logger, err, method, tracePrefix, span, status, cashier_errors.ErrFailedDeleteAllCashierPermanent, fields...)
}
