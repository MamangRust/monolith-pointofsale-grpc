package handler

import "github.com/MamangRust/monolith-point-of-sale-product/internal/service"

type Deps struct {
	Service *service.Service
}

type Handler struct {
	Product ProductHandleGrpc
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Product: NewProductHandleGrpc(deps.Service),
	}
}
