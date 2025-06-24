package errorhandler

import "github.com/MamangRust/monolith-point-of-sale-pkg/logger"

type ErrorHandler struct {
	MerchantQueryError           MerchantQueryErrorHandler
	MerchantCommandError         MerchantCommandErrorHandler
	MerchantDocumentQueryError   MerchantDocumentQueryErrorHandler
	MerchantDocumentCommandError MerchantDocumentCommandErrorHandler
}

func NewErrorHandler(logger logger.LoggerInterface) *ErrorHandler {
	return &ErrorHandler{
		MerchantQueryError: NewMerchantQueryError(logger),
		MerchantCommandError: NewMerchantCommandError(
			logger,
		),
		MerchantDocumentQueryError: NewMerchantDocumentQueryError(
			logger,
		),
		MerchantDocumentCommandError: NewMerchantDocumentCommandError(
			logger),
	}
}
