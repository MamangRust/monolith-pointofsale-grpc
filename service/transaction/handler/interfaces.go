package handler

import "github.com/MamangRust/monolith-point-of-sale-shared/pb"

type TransactionHandleGrpc interface {
	pb.TransactionServiceServer
}
