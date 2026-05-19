package merchant_errors

import "errors"

var (
	ErrFindAllMerchants = errors.New("failed to find all merchants")
	ErrFindByActive     = errors.New("failed to find active merchants")
	ErrFindByTrashed    = errors.New("failed to find trashed merchants")
	ErrFindById         = errors.New("failed to find merchant by ID")

	ErrCreateMerchant             = errors.New("failed to create merchant")
	ErrUpdateMerchant             = errors.New("failed to update merchant")
	ErrUpdateMerchantStatusFailed = errors.New("failed to update merchant status")

	ErrTrashedMerchant            = errors.New("failed to move merchant to trash")
	ErrRestoreMerchant            = errors.New("failed to restore merchant from trash")
	ErrDeleteMerchantPermanent    = errors.New("failed to permanently delete merchant")
	ErrRestoreAllMerchant         = errors.New("failed to restore all trashed merchants")
	ErrDeleteAllMerchantPermanent = errors.New("failed to permanently delete all trashed merchants")
)
