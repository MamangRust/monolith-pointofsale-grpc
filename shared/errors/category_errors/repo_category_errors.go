package category_errors

import "errors"

var (
	ErrGetMonthlyTotalPrice           = errors.New("failed to get monthly total price")
	ErrGetYearlyTotalPrices           = errors.New("failed to get yearly total prices")
	ErrGetMonthlyTotalPriceById       = errors.New("failed to get monthly total price by category ID")
	ErrGetYearlyTotalPricesById       = errors.New("failed to get yearly total prices by category ID")
	ErrGetMonthlyTotalPriceByMerchant = errors.New("failed to get monthly total price by merchant")
	ErrGetYearlyTotalPricesByMerchant = errors.New("failed to get yearly total prices by merchant")

	ErrGetMonthPrice           = errors.New("failed to get month price")
	ErrGetYearPrice            = errors.New("failed to get year price")
	ErrGetMonthPriceByMerchant = errors.New("failed to get month price by merchant")
	ErrGetYearPriceByMerchant  = errors.New("failed to get year price by merchant")
	ErrGetMonthPriceById       = errors.New("failed to get month price by category ID")
	ErrGetYearPriceById        = errors.New("failed to get year price by category ID")

	ErrFindAllCategory = errors.New("failed to find all categories")
	ErrFindById        = errors.New("failed to find category by ID")
	ErrFindByNameAndId = errors.New("failed to find category by name and ID")
	ErrFindByName      = errors.New("failed to find category by name")
	ErrFindByActive    = errors.New("failed to find active categories")
	ErrFindByTrashed   = errors.New("failed to find trashed categories")

	ErrCreateCategory               = errors.New("failed to create category")
	ErrUpdateCategory               = errors.New("failed to update category")
	ErrTrashedCategory              = errors.New("failed to trash category")
	ErrRestoreCategory              = errors.New("failed to restore category")
	ErrDeleteCategoryPermanently    = errors.New("failed to permanently delete category")
	ErrRestoreAllCategories         = errors.New("failed to restore all categories")
	ErrDeleteAllPermanentCategories = errors.New("failed to permanently delete all trashed categories")
)
