package refreshtoken_errors

import (
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"

	"google.golang.org/grpc/codes"
)

var ErrGrpcRefreshToken = response.NewGrpcError("error", "refresh token failed", int(codes.Unauthenticated))
