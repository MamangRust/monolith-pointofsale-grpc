package errorhandler

import "github.com/MamangRust/monolith-point-of-sale-pkg/logger"

type ErrorHandler struct {
	OrderCommandError    OrderCommandError
	OrderQueryError      OrderQueryError
	OrderStats           OrderStatsError
	OrderStatsByMerchant OrderStatsByMerchantError
}

func NewErrorHandler(logger logger.LoggerInterface) *ErrorHandler {
	return &ErrorHandler{
		OrderCommandError:    NewOrderCommandError(logger),
		OrderQueryError:      NewOrderQueryError(logger),
		OrderStats:           NewOrderStatsError(logger),
		OrderStatsByMerchant: NewOrderStatsByMerchantError(logger),
	}
}
