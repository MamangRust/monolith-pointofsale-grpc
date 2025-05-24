package handler

import "github.com/MamangRust/monolith-point-of-sale-role/internal/service"

type Deps struct {
	Service service.Service
}

type Handler struct {
	Role RoleHandleGrpc
}

func NewHandler(deps Deps) *Handler {
	return &Handler{
		Role: NewRoleHandleGrpc(deps.Service),
	}
}
