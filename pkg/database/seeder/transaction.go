package seeder

import (
	"context"
	"database/sql"
	"math/rand"

	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"

	"go.uber.org/zap"
)

type transactionSeeder struct {
	db     *db.Queries
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewTransactionSeeder(db *db.Queries, ctx context.Context, logger logger.LoggerInterface) *transactionSeeder {
	return &transactionSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *transactionSeeder) Seed() error {
	orders, err := r.db.GetOrders(r.ctx, db.GetOrdersParams{
		Column1: "",
		Limit:   20,
		Offset:  0,
	})

	if err != nil {
		r.logger.Error("Failed to get transactions:", zap.Any("error", err))
		return err
	}

	merchants, err := r.db.GetMerchants(r.ctx, db.GetMerchantsParams{
		Column1: "",
		Limit:   20,
		Offset:  0,
	})

	if err != nil {
		r.logger.Error("Failed to get transactions:", zap.Any("error", err))
		return err
	}

	for i := 0; i < 10; i++ {
		selectedMerchantId := merchants[rand.Intn(len(merchants))]
		selectedOrderId := orders[rand.Intn(len(orders))]

		var paymentMethod string
		var amount, changeAmount float64
		var paymentStatus string

		paymentMethod = "Credit Card"
		amount = float64(100 + i)
		changeAmount = float64(5 + i)
		paymentStatus = "Completed"

		_, err := r.db.CreateTransaction(r.ctx, db.CreateTransactionParams{
			OrderID:       selectedOrderId.OrderID,
			PaymentMethod: paymentMethod,
			Amount:        int32(amount),
			ChangeAmount: sql.NullInt32{
				Int32: int32(changeAmount),
				Valid: true,
			},
			PaymentStatus: paymentStatus,
			MerchantID:    selectedMerchantId.MerchantID,
		})
		if err != nil {
			r.logger.Error("Failed to create transaction:", zap.Any("error", err))
			return err
		}
	}

	r.logger.Info("Successfully seeded 10 transactions.", zap.Int("count", 10))
	return nil
}
