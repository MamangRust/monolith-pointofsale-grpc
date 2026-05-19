package seeder

import (
	"context"
	"database/sql"
	"fmt"

	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"

	"go.uber.org/zap"
)

type categorySeeder struct {
	db     *db.Queries
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewCategorySeeder(db *db.Queries, ctx context.Context, logger logger.LoggerInterface) *categorySeeder {
	return &categorySeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *categorySeeder) Seed() error {
	categoryNames := []string{
		"Electronics", "Clothing", "Groceries", "Toys", "Home & Kitchen",
		"Books", "Beauty & Health", "Sports & Outdoors", "Automotive", "Furniture",
	}

	categoryDescriptions := []string{
		"Best electronics products", "Latest fashion trends", "Fresh groceries",
		"Fun toys for kids", "Essentials for home & kitchen",
		"Books for all ages", "Beauty and health products",
		"Outdoor sports equipment", "Automotive accessories", "Stylish furniture",
	}

	for i := 0; i < 10; i++ {
		name := categoryNames[i%len(categoryNames)]
		description := sql.NullString{String: categoryDescriptions[i%len(categoryDescriptions)], Valid: true}
		slugCategory := sql.NullString{String: fmt.Sprintf("%s-%d", name, i+1), Valid: true}

		_, err := r.db.CreateCategory(r.ctx, db.CreateCategoryParams{
			Name:         name,
			Description:  description,
			SlugCategory: slugCategory,
		})
		if err != nil {
			r.logger.Error("Failed to create category:", zap.Error(err))
			return err
		}
		r.logger.Debug("Category created:", zap.String("slug", slugCategory.String))
	}

	r.logger.Info("Category seeding completed successfully.", zap.Int("count", 10))
	return nil
}
