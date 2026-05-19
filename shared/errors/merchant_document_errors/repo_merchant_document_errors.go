package merchantdocument_errors

import "errors"

var (
	ErrFindAllMerchantDocumentsFailed     = errors.New("failed to find all merchant documents")
	ErrFindActiveMerchantDocumentsFailed  = errors.New("failed to find active merchant documents")
	ErrFindTrashedMerchantDocumentsFailed = errors.New("failed to find trashed merchant documents")
	ErrFindMerchantDocumentByIdFailed     = errors.New("failed to find merchant document by ID")

	ErrCreateMerchantDocumentFailed       = errors.New("failed to create merchant document")
	ErrUpdateMerchantDocumentFailed       = errors.New("failed to update merchant document")
	ErrUpdateMerchantDocumentStatusFailed = errors.New("failed to update merchant document status")

	ErrTrashedMerchantDocumentFailed             = errors.New("failed to soft-delete (trash) merchant document")
	ErrRestoreMerchantDocumentFailed             = errors.New("failed to restore merchant document")
	ErrDeleteMerchantDocumentPermanentFailed     = errors.New("failed to permanently delete merchant document")
	ErrRestoreAllMerchantDocumentsFailed         = errors.New("failed to restore all merchant documents")
	ErrDeleteAllMerchantDocumentsPermanentFailed = errors.New("failed to permanently delete all merchant documents")
)
