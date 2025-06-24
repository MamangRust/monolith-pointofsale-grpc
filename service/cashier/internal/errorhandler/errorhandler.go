package errorhandler

import "github.com/MamangRust/monolith-point-of-sale-pkg/logger"

type ErrorHandler struct {
	CashierCommandError         CashierCommadError
	CashierQueryError           CashierQueryError
	CashierStatsError           CashierStatsError
	CashierStatsByIdError       CashierStatsByIdError
	CashierStatsByMerchantError CashierStatsByMerchantError
}

func NewErrorHandler(logger logger.LoggerInterface) *ErrorHandler {
	return &ErrorHandler{
		CashierCommandError:         NewCashierCommandError(logger),
		CashierQueryError:           NewcashierQueryError(logger),
		CashierStatsError:           NewCashierStatsError(logger),
		CashierStatsByIdError:       NewcashierStatsByIdError(logger),
		CashierStatsByMerchantError: NewcashierStatsByMerchantError(logger),
	}
}
