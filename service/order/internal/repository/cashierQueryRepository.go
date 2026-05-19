package repository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/cashier_errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/pb"
	"github.com/jackc/pgx/v5/pgtype"
)

type cashierQueryRepository struct {
	client pb.CashierServiceClient
}

func NewCashierQueryRepository(client pb.CashierServiceClient) CashierQueryRepository {
	return &cashierQueryRepository{
		client: client,
	}
}

func parsePgTimestamp(s string) pgtype.Timestamp {
	if s == "" {
		return pgtype.Timestamp{}
	}
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		t, err = time.Parse(time.RFC3339, s)
		if err != nil {
			return pgtype.Timestamp{}
		}
	}
	return pgtype.Timestamp{Time: t, Valid: true}
}

func (r *cashierQueryRepository) FindById(ctx context.Context, cashierID int) (*db.Cashier, error) {
	resp, err := r.client.FindById(ctx, &pb.FindByIdCashierRequest{
		Id: int32(cashierID),
	})
	if err != nil {
		return nil, cashier_errors.ErrFindCashierById
	}

	if resp == nil || resp.Data == nil {
		return nil, cashier_errors.ErrFindCashierById
	}

	c := resp.Data
	res := &db.Cashier{
		CashierID:  c.Id,
		MerchantID: c.MerchantId,
		Name:       c.Name,
		CreatedAt:  parsePgTimestamp(c.CreatedAt),
		UpdatedAt:  parsePgTimestamp(c.UpdatedAt),
	}

	return res, nil
}
