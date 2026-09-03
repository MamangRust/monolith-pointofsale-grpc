package handler

import (
	"github.com/MamangRust/monolith-point-of-sale-cashier/service"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/pb"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Cashier pb.CashierServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Cashier: NewCashierHandleGrpc(
			deps.Service,
			deps.Logger,
		),
	}
}
