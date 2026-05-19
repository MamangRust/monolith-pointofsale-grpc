package transaction_errors

import "errors"

var (
	ErrGetMonthlyAmountSuccess = errors.New("failed to get monthly amount success")
	ErrGetYearlyAmountSuccess  = errors.New("failed to get yearly amount success")
	ErrGetMonthlyAmountFailed  = errors.New("failed to get monthly amount failed")
	ErrGetYearlyAmountFailed   = errors.New("failed to get yearly amount failed")

	ErrGetMonthlyAmountSuccessByMerchant = errors.New("failed to get monthly amount success by merchant")
	ErrGetYearlyAmountSuccessByMerchant  = errors.New("failed to get yearly amount success by merchant")
	ErrGetMonthlyAmountFailedByMerchant  = errors.New("failed to get monthly amount failed by merchant")
	ErrGetYearlyAmountFailedByMerchant   = errors.New("failed to get yearly amount failed by merchant")

	ErrGetMonthlyTransactionMethod           = errors.New("failed to get monthly transaction method")
	ErrGetYearlyTransactionMethod            = errors.New("failed to get yearly transaction method")
	ErrGetMonthlyTransactionMethodByMerchant = errors.New("failed to get monthly transaction method by merchant")
	ErrGetYearlyTransactionMethodByMerchant  = errors.New("failed to get yearly transaction method by merchant")

	ErrFindAllTransactions = errors.New("failed to find all transactions")
	ErrFindByActive        = errors.New("failed to find active transactions")
	ErrFindByTrashed       = errors.New("failed to find trashed transactions")
	ErrFindByMerchant      = errors.New("failed to find transactions by merchant")
	ErrFindById            = errors.New("failed to find transaction by ID")
	ErrFindByOrderId       = errors.New("failed to find transaction by order ID")

	ErrCreateTransaction             = errors.New("failed to create transaction")
	ErrUpdateTransaction             = errors.New("failed to update transaction")
	ErrTrashTransaction              = errors.New("failed to move transaction to trash")
	ErrRestoreTransaction            = errors.New("failed to restore transaction")
	ErrDeleteTransactionPermanently  = errors.New("failed to permanently delete transaction")
	ErrRestoreAllTransactions        = errors.New("failed to restore all transactions")
	ErrDeleteAllTransactionPermanent = errors.New("failed to permanently delete all transactions")
)
