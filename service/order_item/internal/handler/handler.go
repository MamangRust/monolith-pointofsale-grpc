package handler

import (
	"github.com/MamangRust/monolith-point-of-sale-order-item/internal/service"
)

type Deps struct {
	Service *service.Service
}

type Handler struct {
	OrderItem OrderItemHandlerGrpc
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		OrderItem: NewOrderItemHandleGrpc(deps.Service.OrderItemQuery),
	}
}
