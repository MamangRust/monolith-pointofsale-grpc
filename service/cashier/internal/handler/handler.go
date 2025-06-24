package handler

import "github.com/MamangRust/monolith-point-of-sale-cashier/internal/service"

type Deps struct {
	Service *service.Service
}

type Handler struct {
	Cashier CashierHandleGrpc
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Cashier: NewCashierHandleGrpc(deps.Service),
	}
}
