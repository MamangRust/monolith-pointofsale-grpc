package handler

import "github.com/MamangRust/monolith-point-of-sale-shared/pb"

type MerchantDocumentHandleGrpc interface {
	pb.MerchantDocumentServiceServer
}

type MerchantHandleGrpc interface {
	pb.MerchantServiceServer
}
