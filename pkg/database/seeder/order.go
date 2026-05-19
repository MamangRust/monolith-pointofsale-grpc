package seeder

import (
	"context"
	"database/sql"

	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"

	"go.uber.org/zap"
	"golang.org/x/exp/rand"
)

type orderSeeder struct {
	db     *db.Queries
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewOrderSeeder(db *db.Queries, ctx context.Context, logger logger.LoggerInterface) *orderSeeder {
	return &orderSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *orderSeeder) Seed() error {
	merchants, err := r.db.GetMerchants(r.ctx, db.GetMerchantsParams{
		Column1: "",
		Limit:   20,
		Offset:  0,
	})
	if err != nil {
		r.logger.Error("Failed to get merchants", zap.Error(err))
		return err
	}

	cashiers, err := r.db.GetCashiers(r.ctx, db.GetCashiersParams{
		Column1: "",
		Limit:   20,
		Offset:  0,
	})
	if err != nil {
		r.logger.Error("Failed to get cashiers", zap.Error(err))
		return err
	}

	if len(merchants) == 0 || len(cashiers) == 0 {
		r.logger.Error("No merchants or cashiers found, skipping order seeding")
		return nil
	}

	for i := 0; i < 10; i++ {
		merchant := merchants[rand.Intn(len(merchants))]
		cashier := cashiers[rand.Intn(len(cashiers))]
		totalPrice := int32(rand.Intn(500000) + 50000)

		order, err := r.db.CreateOrder(r.ctx, db.CreateOrderParams{
			MerchantID: merchant.MerchantID,
			CashierID:  cashier.CashierID,
			TotalPrice: int64(totalPrice),
		})
		if err != nil {
			r.logger.Error("Failed to create order", zap.Error(err))
			return err
		}

		orderID := order.OrderID

		products, err := r.db.GetProductsByMerchant(r.ctx, db.GetProductsByMerchantParams{
			MerchantID: merchant.MerchantID,
			Column2:    sql.NullString{},
			Column3:    0,
			Column4:    0,
			Column5:    0,
			Limit:      10,
			Offset:     0,
		})
		if err != nil {
			r.logger.Error("Failed to get products", zap.Error(err))
			return err
		}

		if len(products) == 0 {
			r.logger.Debug("No products found for merchant", zap.Int32("merchant_id", merchant.MerchantID))
			continue
		}

		for j := 0; j < rand.Intn(5)+1; j++ {
			product := products[rand.Intn(len(products))]
			quantity := int32(rand.Intn(5) + 1)
			price := product.Price * quantity

			_, err := r.db.CreateOrderItem(r.ctx, db.CreateOrderItemParams{
				OrderID:   orderID,
				ProductID: product.ProductID,
				Quantity:  quantity,
				Price:     price,
			})
			if err != nil {
				r.logger.Error("Failed to create order item", zap.Error(err))
				return err
			}
		}
	}

	r.logger.Info("Order seeding completed successfully.", zap.Int("count", 10))
	return nil
}
