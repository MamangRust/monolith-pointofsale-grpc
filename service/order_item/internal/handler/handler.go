package handler

import (
	"github.com/MamangRust/monolith-point-of-sale-order-item/internal/service"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/pb"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	OrderItem pb.OrderItemServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		OrderItem: NewOrderItemHandleGrpc(
			deps.Service.OrderItemQuery,
			deps.Logger,
		),
	}
}
