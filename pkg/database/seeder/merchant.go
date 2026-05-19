package seeder

import (
	"context"
	"database/sql"
	"fmt"

	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"

	"go.uber.org/zap"
)

type merchantSeeder struct {
	db     *db.Queries
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewMerchantSeeder(db *db.Queries, ctx context.Context, logger logger.LoggerInterface) *merchantSeeder {
	return &merchantSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *merchantSeeder) Seed() error {
	users, err := r.db.GetUsers(r.ctx, db.GetUsersParams{
		Column1: "",
		Limit:   int32(20),
		Offset:  0,
	})
	if err != nil {
		r.logger.Error("Failed to fetch merchants:", zap.Any("error", err))
		return err
	}

	for i := 1; i <= 10; i++ {
		userID := users[i%len(users)].UserID

		merchant := db.CreateMerchantParams{
			UserID:       userID,
			Name:         fmt.Sprintf("Toko %d", i),
			Description:  sql.NullString{String: fmt.Sprintf("Deskripsi untuk Toko %d", i), Valid: true},
			Address:      sql.NullString{String: fmt.Sprintf("Jl. Toko %d", i), Valid: true},
			ContactEmail: sql.NullString{String: fmt.Sprintf("toko%d@example.com", i), Valid: true},
			ContactPhone: sql.NullString{String: fmt.Sprintf("0812345678%d", i), Valid: true},
			Status:       "active",
		}

		_, err = r.db.CreateMerchant(r.ctx, merchant)
		if err != nil {
			r.logger.Error("Failed to create merchant:", zap.Error(err))
			return err
		}
	}

	r.logger.Info("Merchant seeding completed successfully.", zap.Int("count", 10))
	return nil
}
