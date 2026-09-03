package repository

import (
	"context"

	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-shared/convert"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/merchant_errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/pb"
)

type merchantQueryRepository struct {
	client pb.MerchantServiceClient
}

func NewMerchantQueryRepository(client pb.MerchantServiceClient) MerchantQueryRepository {
	return &merchantQueryRepository{
		client: client,
	}
}

func (r *merchantQueryRepository) FindById(ctx context.Context, merchantID int) (*db.Merchant, error) {
	resp, err := r.client.FindById(ctx, &pb.FindByIdMerchantRequest{
		Id: int32(merchantID),
	})
	if err != nil {
		return nil, merchant_errors.ErrFindById
	}

	if resp == nil || resp.Data == nil {
		return nil, merchant_errors.ErrFindById
	}

	m := resp.Data
	res := &db.Merchant{
		MerchantID:   m.Id,
		UserID:       m.UserId,
		Name:         m.Name,
		Description:  convert.NullableString(m.Description),
		Address:      convert.NullableString(m.Address),
		ContactEmail: convert.NullableString(m.ContactEmail),
		ContactPhone: convert.NullableString(m.ContactPhone),
		Status:       m.Status,
		CreatedAt:    convert.PgTimestamp(m.CreatedAt),
		UpdatedAt:    convert.PgTimestamp(m.UpdatedAt),
	}

	return res, nil
}
