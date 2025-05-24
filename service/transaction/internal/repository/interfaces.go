package repository

import (
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/record"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type CashierQueryRepository interface {
	FindById(id int) (*record.CashierRecord, error)
}

type MerchantQueryRepository interface {
	FindById(id int) (*record.MerchantRecord, error)
}

type OrderItemQueryRepository interface {
	FindOrderItemByOrder(order_id int) ([]*record.OrderItemRecord, error)
}
type OrderQueryRepository interface {
	FindById(id int) (*record.OrderRecord, error)
}

type TransactionStatsRepository interface {
	GetMonthlyAmountSuccess(req *requests.MonthAmountTransaction) ([]*record.TransactionMonthlyAmountSuccessRecord, error)
	GetYearlyAmountSuccess(year int) ([]*record.TransactionYearlyAmountSuccessRecord, error)
	GetMonthlyAmountFailed(req *requests.MonthAmountTransaction) ([]*record.TransactionMonthlyAmountFailedRecord, error)
	GetYearlyAmountFailed(year int) ([]*record.TransactionYearlyAmountFailedRecord, error)

	GetMonthlyTransactionMethodSuccess(req *requests.MonthMethodTransaction) ([]*record.TransactionMonthlyMethodRecord, error)
	GetYearlyTransactionMethodSuccess(year int) ([]*record.TransactionYearlyMethodRecord, error)
	GetMonthlyTransactionMethodFailed(req *requests.MonthMethodTransaction) ([]*record.TransactionMonthlyMethodRecord, error)
	GetYearlyTransactionMethodFailed(year int) ([]*record.TransactionYearlyMethodRecord, error)
}

type TransactionStatsByMerchantRepository interface {
	GetMonthlyAmountSuccessByMerchant(req *requests.MonthAmountTransactionMerchant) ([]*record.TransactionMonthlyAmountSuccessRecord, error)
	GetYearlyAmountSuccessByMerchant(req *requests.YearAmountTransactionMerchant) ([]*record.TransactionYearlyAmountSuccessRecord, error)
	GetMonthlyAmountFailedByMerchant(req *requests.MonthAmountTransactionMerchant) ([]*record.TransactionMonthlyAmountFailedRecord, error)
	GetYearlyAmountFailedByMerchant(req *requests.YearAmountTransactionMerchant) ([]*record.TransactionYearlyAmountFailedRecord, error)

	GetMonthlyTransactionMethodByMerchantSuccess(req *requests.MonthMethodTransactionMerchant) ([]*record.TransactionMonthlyMethodRecord, error)
	GetYearlyTransactionMethodByMerchantSuccess(req *requests.YearMethodTransactionMerchant) ([]*record.TransactionYearlyMethodRecord, error)
	GetMonthlyTransactionMethodByMerchantFailed(req *requests.MonthMethodTransactionMerchant) ([]*record.TransactionMonthlyMethodRecord, error)
	GetYearlyTransactionMethodByMerchantFailed(req *requests.YearMethodTransactionMerchant) ([]*record.TransactionYearlyMethodRecord, error)
}

type TransactionQueryRepository interface {
	FindAllTransactions(req *requests.FindAllTransaction) ([]*record.TransactionRecord, *int, error)
	FindByActive(req *requests.FindAllTransaction) ([]*record.TransactionRecord, *int, error)
	FindByTrashed(req *requests.FindAllTransaction) ([]*record.TransactionRecord, *int, error)
	FindByMerchant(req *requests.FindAllTransactionByMerchant) ([]*record.TransactionRecord, *int, error)
	FindById(transaction_id int) (*record.TransactionRecord, error)
	FindByOrderId(order_id int) (*record.TransactionRecord, error)
}

type TransactionCommandRepository interface {
	CreateTransaction(request *requests.CreateTransactionRequest) (*record.TransactionRecord, error)
	UpdateTransaction(request *requests.UpdateTransactionRequest) (*record.TransactionRecord, error)
	TrashTransaction(transaction_id int) (*record.TransactionRecord, error)
	RestoreTransaction(transaction_id int) (*record.TransactionRecord, error)
	DeleteTransactionPermanently(transaction_id int) (bool, error)
	RestoreAllTransactions() (bool, error)
	DeleteAllTransactionPermanent() (bool, error)
}
