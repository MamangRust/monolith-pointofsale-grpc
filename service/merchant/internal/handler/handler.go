package handler

import (
	"github.com/MamangRust/monolith-point-of-sale-merchant/internal/service"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/pb"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Merchant         pb.MerchantServiceServer
	MerchantDocument pb.MerchantDocumentServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Merchant: NewMerchantHandleGrpc(
			deps.Service,
			deps.Logger,
		),
		MerchantDocument: NewMerchantDocumentHandleGrpc(
			deps.Service,
			deps.Logger,
		),
	}
}
