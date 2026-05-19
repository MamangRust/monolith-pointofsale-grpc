package seeder

import (
	"context"
	"fmt"

	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"

	"go.uber.org/zap"
	"golang.org/x/exp/rand"
)

type cashierSeeder struct {
	db     *db.Queries
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewCashierSeeder(db *db.Queries, ctx context.Context, logger logger.LoggerInterface) *cashierSeeder {
	return &cashierSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *cashierSeeder) Seed() error {
	merchants, err := r.db.GetMerchants(r.ctx, db.GetMerchantsParams{
		Column1: "",
		Limit:   int32(20),
		Offset:  0,
	})
	if err != nil {
		r.logger.Error("Failed to fetch merchants:", zap.Any("err", err))
		return err
	}

	users, err := r.db.GetUsers(r.ctx, db.GetUsersParams{
		Column1: "",
		Limit:   int32(20),
		Offset:  0,
	})
	if err != nil {
		r.logger.Error("Failed to fetch users:", zap.Any("error", err))
		return err
	}

	if len(merchants) == 0 || len(users) == 0 {
		r.logger.Error("Merchants or Users not found. Seed operation aborted.")
		return fmt.Errorf("no merchants or users found")
	}

	for i := 1; i <= 10; i++ {
		merchant := merchants[rand.Intn(len(merchants))]
		user := users[rand.Intn(len(users))]

		cashierName := fmt.Sprintf("Cashier %d", i)
		_, err := r.db.CreateCashier(r.ctx, db.CreateCashierParams{
			MerchantID: merchant.MerchantID,
			UserID:     user.UserID,
			Name:       cashierName,
		})
		if err != nil {
			r.logger.Error("Failed to create cashier:", zap.Any("error", err))
			return err
		}

	}

	r.logger.Info("Cashier seeding completed successfully.", zap.Int("count", 10))
	return nil
}
