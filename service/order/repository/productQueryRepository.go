package repository

import (
	"context"

	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-shared/convert"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/product_errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/pb"
)

type productQueryRepository struct {
	client pb.ProductServiceClient
}

func NewProductQueryRepository(client pb.ProductServiceClient) ProductQueryRepository {
	return &productQueryRepository{
		client: client,
	}
}

func parseNullableInt32(i int32) *int32 {
	return &i
}

func (r *productQueryRepository) FindById(ctx context.Context, product_id int) (*db.Product, error) {
	resp, err := r.client.FindById(ctx, &pb.FindByIdProductRequest{
		Id: int32(product_id),
	})
	if err != nil {
		return nil, product_errors.ErrFindById
	}

	if resp == nil || resp.Data == nil {
		return nil, product_errors.ErrFindById
	}

	p := resp.Data
	res := &db.Product{
		ProductID:    p.Id,
		MerchantID:   p.MerchantId,
		CategoryID:   p.CategoryId,
		Name:         p.Name,
		Description:  convert.NullableString(p.Description),
		Price:        p.Price,
		CountInStock: p.CountInStock,
		Brand:        convert.NullableString(p.Brand),
		Weight:       parseNullableInt32(p.Weight),
		SlugProduct:  convert.NullableString(p.SlugProduct),
		ImageProduct: convert.NullableString(p.ImageProduct),
		Barcode:      convert.NullableString(p.Barcode),
		CreatedAt:    convert.PgTimestamp(p.CreatedAt),
		UpdatedAt:    convert.PgTimestamp(p.UpdatedAt),
	}

	return res, nil
}
