package cashier_errors

import "errors"

var (
	ErrGetMonthlyTotalSales           = errors.New("failed to get monthly total sales")
	ErrGetYearlyTotalSales            = errors.New("failed to get yearly total sales")
	ErrGetMonthlyTotalSalesById       = errors.New("failed to get monthly total sales by cashier ID")
	ErrGetYearlyTotalSalesById        = errors.New("failed to get yearly total sales by cashier ID")
	ErrGetMonthlyTotalSalesByMerchant = errors.New("failed to get monthly total sales by merchant")
	ErrGetYearlyTotalSalesByMerchant  = errors.New("failed to get yearly total sales by merchant")

	ErrGetMonthlyCashier           = errors.New("failed to get monthly cashier sales")
	ErrGetYearlyCashier            = errors.New("failed to get yearly cashier sales")
	ErrGetMonthlyCashierByMerchant = errors.New("failed to get monthly cashier sales by merchant")
	ErrGetYearlyCashierByMerchant  = errors.New("failed to get yearly cashier sales by merchant")
	ErrGetMonthlyCashierById       = errors.New("failed to get monthly cashier sales by cashier ID")
	ErrGetYearlyCashierById        = errors.New("failed to get yearly cashier sales by cashier ID")

	ErrFindAllCashiers        = errors.New("failed to find all cashiers")
	ErrFindCashierById        = errors.New("failed to find cashier by ID")
	ErrFindActiveCashiers     = errors.New("failed to find active cashiers")
	ErrFindTrashedCashiers    = errors.New("failed to find trashed cashiers")
	ErrFindCashiersByMerchant = errors.New("failed to find cashiers by merchant")

	ErrCreateCashier              = errors.New("failed to create cashier")
	ErrUpdateCashier              = errors.New("failed to update cashier")
	ErrTrashedCashier             = errors.New("failed to move cashier to trash")
	ErrRestoreCashier             = errors.New("failed to restore cashier from trash")
	ErrDeleteCashierPermanent     = errors.New("failed to permanently delete cashier")
	ErrRestoreAllCashiers         = errors.New("failed to restore all cashiers")
	ErrDeleteAllCashiersPermanent = errors.New("failed to permanently delete all trashed cashiers")
)
