package handler

import (
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-product/service"
	"github.com/MamangRust/monolith-point-of-sale-shared/pb"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Product pb.ProductServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Product: NewProductHandleGrpc(
			deps.Service,
			deps.Logger,
		),
	}
}
