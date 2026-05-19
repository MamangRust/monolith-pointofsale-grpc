package recordmapper

type RecordMapper struct {
	UserRecordMapper             UserRecordMapping
	RoleRecordMapper             RoleRecordMapping
	UserRoleRecordMapper         UserRoleRecordMapping
	RefreshTokenRecordMapper     RefreshTokenRecordMapping
	CategoryRecordMapper         CategoryRecordMapper
	CashierRecordMapper          CashierRecordMapping
	MerchantRecordMapper         MerchantRecordMapping
	OrderItemRecordMapper        OrderItemRecordMapping
	OrderRecordMapper            OrderRecordMapping
	ProductRecordMapper          ProductRecordMapping
	TransactionRecordMapper      TransactionRecordMapping
	ResetTokenRecordMapper       ResetTokenRecordMapping
	MerchantDocumentRecordMapper MerchantDocumentMapping
}

func NewRecordMapper() *RecordMapper {
	return &RecordMapper{
		UserRecordMapper:             NewUserRecordMapper(),
		RoleRecordMapper:             NewRoleRecordMapper(),
		UserRoleRecordMapper:         NewUserRoleRecordMapper(),
		RefreshTokenRecordMapper:     NewRefreshTokenRecordMapper(),
		ResetTokenRecordMapper:       NewResetTokenRecordMapper(),
		CategoryRecordMapper:         NewCategoryRecordMapper(),
		CashierRecordMapper:          NewCashierRecordMapper(),
		MerchantRecordMapper:         NewMerchantRecordMapper(),
		OrderItemRecordMapper:        NewOrderItemRecordMapper(),
		OrderRecordMapper:            NewOrderRecordMapper(),
		ProductRecordMapper:          NewProductRecordMapper(),
		TransactionRecordMapper:      NewTransactionRecordMapper(),
		MerchantDocumentRecordMapper: NewMerchantDocumentRecordMapper(),
	}
}
