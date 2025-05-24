package handler

import "github.com/MamangRust/monolith-point-of-sale-shared/pb"

type ProductHandleGrpc interface {
	pb.ProductServiceServer
}
