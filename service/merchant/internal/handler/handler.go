package handler

import (
	"github.com/MamangRust/monolith-point-of-sale-merchant/internal/service"
	protomapper "github.com/MamangRust/monolith-point-of-sale-shared/mapper/proto"
)

type Deps struct {
	Service service.Service
}

type Handler struct {
	Merchant         MerchantHandleGrpc
	MerchantDocument MerchantDocumentHandleGrpc
}

func NewHandler(deps Deps) *Handler {
	merchantProto := protomapper.NewMerchantProtoMaper()
	merchantDocumentProto := protomapper.NewMerchantDocumentProtoMapper()

	return &Handler{
		Merchant:         NewMerchantHandleGrpc(deps.Service, merchantProto),
		MerchantDocument: NewMerchantDocumentHandleGrpc(deps.Service, merchantDocumentProto),
	}
}
