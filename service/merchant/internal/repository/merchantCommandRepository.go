package repository

import (
	"context"

	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-point-of-sale-shared/errors"
)

type merchantCommandRepository struct {
	db *db.Queries
}

func NewMerchantCommandRepository(db *db.Queries) MerchantCommandRepository {
	return &merchantCommandRepository{
		db: db,
	}
}

func (r *merchantCommandRepository) CreateMerchant(ctx context.Context, request *requests.CreateMerchantRequest) (*db.Merchant, error) {
	req := db.CreateMerchantParams{
		UserID:       int32(request.UserID),
		Name:         request.Name,
		Description:  &request.Description,
		Address:      &request.Address,
		ContactEmail: &request.ContactEmail,
		ContactPhone: &request.ContactPhone,
		Status:       "inactive",
	}

	merchant, err := r.db.CreateMerchant(ctx, req)
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return merchant, nil
}

func (r *merchantCommandRepository) UpdateMerchant(ctx context.Context, request *requests.UpdateMerchantRequest) (*db.Merchant, error) {
	req := db.UpdateMerchantParams{
		MerchantID:   int32(*request.MerchantID),
		Name:         request.Name,
		Description:  &request.Description,
		Address:      &request.Address,
		ContactEmail: &request.ContactEmail,
		ContactPhone: &request.ContactPhone,
		Status:       request.Status,
	}

	res, err := r.db.UpdateMerchant(ctx, req)
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return res, nil
}

func (r *merchantCommandRepository) UpdateMerchantStatus(ctx context.Context, request *requests.UpdateMerchantStatusRequest) (*db.Merchant, error) {
	req := db.UpdateMerchantStatusParams{
		MerchantID: int32(*request.MerchantID),
		Status:     request.Status,
	}

	res, err := r.db.UpdateMerchantStatus(ctx, req)
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return res, nil
}

func (r *merchantCommandRepository) TrashedMerchant(ctx context.Context, merchant_id int) (*db.Merchant, error) {
	res, err := r.db.TrashMerchant(ctx, int32(merchant_id))
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return res, nil
}

func (r *merchantCommandRepository) RestoreMerchant(ctx context.Context, merchant_id int) (*db.Merchant, error) {
	res, err := r.db.RestoreMerchant(ctx, int32(merchant_id))
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return res, nil
}

func (r *merchantCommandRepository) DeleteMerchantPermanent(ctx context.Context, merchant_id int) (bool, error) {
	err := r.db.DeleteMerchantPermanently(ctx, int32(merchant_id))
	if err != nil {
		return false, sharedErrors.ErrInternal.WithInternal(err)
	}

	return true, nil
}

func (r *merchantCommandRepository) RestoreAllMerchant(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllMerchants(ctx)
	if err != nil {
		return false, sharedErrors.ErrInternal.WithInternal(err)
	}
	return true, nil
}

func (r *merchantCommandRepository) DeleteAllMerchantPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentMerchants(ctx)
	if err != nil {
		return false, sharedErrors.ErrInternal.WithInternal(err)
	}
	return true, nil
}
