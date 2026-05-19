package order_errors

import "errors"

var (
	ErrGetMonthlyTotalRevenue           = errors.New("failed to get monthly total revenue")
	ErrGetYearlyTotalRevenue            = errors.New("failed to get yearly total revenue")
	ErrGetMonthlyTotalRevenueById       = errors.New("failed to get monthly total revenue by order ID")
	ErrGetYearlyTotalRevenueById        = errors.New("failed to get yearly total revenue by order ID")
	ErrGetMonthlyTotalRevenueByMerchant = errors.New("failed to get monthly total revenue by merchant")
	ErrGetYearlyTotalRevenueByMerchant  = errors.New("failed to get yearly total revenue by merchant")

	ErrGetMonthlyOrder           = errors.New("failed to get monthly orders")
	ErrGetYearlyOrder            = errors.New("failed to get yearly orders")
	ErrGetMonthlyOrderByMerchant = errors.New("failed to get monthly orders by merchant")
	ErrGetYearlyOrderByMerchant  = errors.New("failed to get yearly orders by merchant")

	ErrFindAllOrders           = errors.New("failed to find all orders")
	ErrFindByActive            = errors.New("failed to find active orders")
	ErrFindByTrashed           = errors.New("failed to find trashed orders")
	ErrFindByMerchant          = errors.New("failed to find orders by merchant")
	ErrFindById                = errors.New("failed to find order by ID")
	ErrCreateOrder             = errors.New("failed to create order")
	ErrUpdateOrder             = errors.New("failed to update order")
	ErrTrashedOrder            = errors.New("failed to move order to trash")
	ErrRestoreOrder            = errors.New("failed to restore order from trash")
	ErrDeleteOrderPermanent    = errors.New("failed to permanently delete order")
	ErrRestoreAllOrder         = errors.New("failed to restore all trashed orders")
	ErrDeleteAllOrderPermanent = errors.New("failed to permanently delete all trashed orders")
)
