package errorhandler

import "github.com/MamangRust/monolith-point-of-sale-pkg/logger"

type ErrorHandler struct {
	UserQueryError   UserQueryError
	UserCommandError UserCommandError
}

func NewErrorHandler(logger logger.LoggerInterface) *ErrorHandler {
	return &ErrorHandler{
		UserQueryError:   NewUserQueryError(logger),
		UserCommandError: NewUserCommandError(logger),
	}
}
