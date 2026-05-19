package seeder

import (
	"context"
	"database/sql"
	"fmt"

	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"

	"go.uber.org/zap"
	"golang.org/x/exp/rand"
)

type productSeeder struct {
	db     *db.Queries
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewProductSeeder(db *db.Queries, ctx context.Context, logger logger.LoggerInterface) *productSeeder {
	return &productSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *productSeeder) Seed() error {
	merchants, err := r.db.GetMerchants(r.ctx, db.GetMerchantsParams{
		Column1: "",
		Limit:   20,
		Offset:  0,
	})
	if err != nil {
		r.logger.Error("Failed to get merchants:", zap.Any("error", err))
		return err
	}

	categories, err := r.db.GetCategories(r.ctx, db.GetCategoriesParams{
		Column1: "",
		Limit:   20,
		Offset:  0,
	})
	if err != nil {
		r.logger.Error("Failed to get categories:", zap.Any("error", err))
		return err
	}

	if len(merchants) == 0 || len(categories) == 0 {
		r.logger.Error("No merchants or categories found, skipping seeding")
		return nil
	}

	productNames := []string{
		"Smartphone", "Laptop", "Wireless Earbuds", "Gaming Mouse", "Mechanical Keyboard",
		"Smartwatch", "Power Bank", "Bluetooth Speaker", "External Hard Drive", "Monitor",
	}
	brands := []string{"Samsung", "Apple", "Sony", "Logitech", "Razer", "Xiaomi", "HP", "Dell", "Acer", "Asus"}
	images := []string{
		"image1.jpg", "image2.jpg", "image3.jpg", "image4.jpg", "image5.jpg",
		"image6.jpg", "image7.jpg", "image8.jpg", "image9.jpg", "image10.jpg",
	}

	for i := 0; i < 10; i++ {
		merchant := merchants[rand.Intn(len(merchants))]
		category := categories[rand.Intn(len(categories))]
		name := productNames[rand.Intn(len(productNames))]
		brand := sql.NullString{String: brands[rand.Intn(len(brands))], Valid: true}
		price := int32(rand.Intn(5000000) + 50000)
		countInStock := int32(rand.Intn(100) + 1)
		weight := sql.NullInt32{Int32: int32(rand.Intn(5000) + 100), Valid: true}
		slug := sql.NullString{String: fmt.Sprintf("%s-%d", name, rand.Intn(1000)), Valid: true}
		image := sql.NullString{String: images[rand.Intn(len(images))], Valid: true}
		barcode := sql.NullString{String: fmt.Sprintf("BC-%d", rand.Intn(9999999)), Valid: true}

		_, err := r.db.CreateProduct(r.ctx, db.CreateProductParams{
			MerchantID:   merchant.MerchantID,
			CategoryID:   category.CategoryID,
			Name:         name,
			Description:  sql.NullString{String: fmt.Sprintf("Description for %s", name), Valid: true},
			Price:        price,
			CountInStock: countInStock,
			Brand:        brand,
			Weight:       weight,
			SlugProduct:  slug,
			ImageProduct: image,
			Barcode:      barcode,
		})

		if err != nil {
			r.logger.Error("Failed to create product:", zap.Error(err))
			return err
		}

	}

	r.logger.Info("Product seeding completed successfully.", zap.Int("count", 10))
	return nil
}
