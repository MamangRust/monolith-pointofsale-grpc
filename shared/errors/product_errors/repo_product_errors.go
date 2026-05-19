package product_errors

import "errors"

var (
	ErrFindAllProducts           = errors.New("failed to find all products")
	ErrFindByActive              = errors.New("failed to find active products")
	ErrFindByTrashed             = errors.New("failed to find trashed products")
	ErrFindByMerchant            = errors.New("failed to find products by merchant")
	ErrFindByCategory            = errors.New("failed to find products by category")
	ErrFindById                  = errors.New("failed to find product by ID")
	ErrFindByIdTrashed           = errors.New("failed to find trashed product by ID")
	ErrCreateProduct             = errors.New("failed to create product")
	ErrUpdateProduct             = errors.New("failed to update product")
	ErrUpdateProductCountStock   = errors.New("failed to update product stock count")
	ErrTrashedProduct            = errors.New("failed to move product to trash")
	ErrRestoreProduct            = errors.New("failed to restore product")
	ErrDeleteProductPermanent    = errors.New("failed to permanently delete product")
	ErrRestoreAllProducts        = errors.New("failed to restore all products")
	ErrDeleteAllProductPermanent = errors.New("failed to permanently delete all products")
)
