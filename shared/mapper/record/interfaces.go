package recordmapper

import (
	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"

	"github.com/MamangRust/monolith-point-of-sale-shared/domain/record"
)

type UserRecordMapping interface {
	ToUserRecord(user *db.User) *record.UserRecord
	ToUserRecordPagination(user *db.GetUsersRow) *record.UserRecord
	ToUsersRecordPagination(users []*db.GetUsersRow) []*record.UserRecord

	ToUserRecordActivePagination(user *db.GetUsersActiveRow) *record.UserRecord
	ToUsersRecordActivePagination(users []*db.GetUsersActiveRow) []*record.UserRecord
	ToUserRecordTrashedPagination(user *db.GetUserTrashedRow) *record.UserRecord
	ToUsersRecordTrashedPagination(users []*db.GetUserTrashedRow) []*record.UserRecord
}

type RoleRecordMapping interface {
	ToRoleRecord(role *db.Role) *record.RoleRecord
	ToRolesRecord(roles []*db.Role) []*record.RoleRecord

	ToRoleRecordAll(role *db.GetRolesRow) *record.RoleRecord
	ToRolesRecordAll(roles []*db.GetRolesRow) []*record.RoleRecord

	ToRoleRecordActive(role *db.GetActiveRolesRow) *record.RoleRecord
	ToRolesRecordActive(roles []*db.GetActiveRolesRow) []*record.RoleRecord
	ToRoleRecordTrashed(role *db.GetTrashedRolesRow) *record.RoleRecord
	ToRolesRecordTrashed(roles []*db.GetTrashedRolesRow) []*record.RoleRecord
}

type ResetTokenRecordMapping interface {
	ToResetTokenRecord(resetToken *db.ResetToken) *record.ResetTokenRecord
}

type UserRoleRecordMapping interface {
	ToUserRoleRecord(userRole *db.UserRole) *record.UserRoleRecord
}

type RefreshTokenRecordMapping interface {
	ToRefreshTokenRecord(refreshToken *db.RefreshToken) *record.RefreshTokenRecord
	ToRefreshTokensRecord(refreshTokens []*db.RefreshToken) []*record.RefreshTokenRecord
}

type CategoryRecordMapper interface {
	ToCategoryMonthlyTotalPrice(c *db.GetMonthlyTotalPriceRow) *record.CategoriesMonthlyTotalPriceRecord
	ToCategoryMonthlyTotalPrices(c []*db.GetMonthlyTotalPriceRow) []*record.CategoriesMonthlyTotalPriceRecord
	ToCategoryYearlyTotalPrice(c *db.GetYearlyTotalPriceRow) *record.CategoriesYearlyTotalPriceRecord
	ToCategoryYearlyTotalPrices(c []*db.GetYearlyTotalPriceRow) []*record.CategoriesYearlyTotalPriceRecord
	ToCategoryMonthlyTotalPriceById(c *db.GetMonthlyTotalPriceByIdRow) *record.CategoriesMonthlyTotalPriceRecord
	ToCategoryMonthlyTotalPricesById(c []*db.GetMonthlyTotalPriceByIdRow) []*record.CategoriesMonthlyTotalPriceRecord
	ToCategoryYearlyTotalPriceById(c *db.GetYearlyTotalPriceByIdRow) *record.CategoriesYearlyTotalPriceRecord
	ToCategoryYearlyTotalPricesById(c []*db.GetYearlyTotalPriceByIdRow) []*record.CategoriesYearlyTotalPriceRecord

	ToCategoryMonthlyTotalPriceByMerchant(c *db.GetMonthlyTotalPriceByMerchantRow) *record.CategoriesMonthlyTotalPriceRecord
	ToCategoryMonthlyTotalPricesByMerchant(c []*db.GetMonthlyTotalPriceByMerchantRow) []*record.CategoriesMonthlyTotalPriceRecord
	ToCategoryYearlyTotalPriceByMerchant(c *db.GetYearlyTotalPriceByMerchantRow) *record.CategoriesYearlyTotalPriceRecord
	ToCategoryYearlyTotalPricesByMerchant(c []*db.GetYearlyTotalPriceByMerchantRow) []*record.CategoriesYearlyTotalPriceRecord

	ToCategoryMonthlyPriceById(category *db.GetMonthlyCategoryByIdRow) *record.CategoriesMonthPriceRecord
	ToCategoryMonthlyPricesById(c []*db.GetMonthlyCategoryByIdRow) []*record.CategoriesMonthPriceRecord
	ToCategoryYearlyPriceById(category *db.GetYearlyCategoryByIdRow) *record.CategoriesYearPriceRecord
	ToCategoryYearlyPricesById(c []*db.GetYearlyCategoryByIdRow) []*record.CategoriesYearPriceRecord

	ToCategoryMonthlyPrice(category *db.GetMonthlyCategoryRow) *record.CategoriesMonthPriceRecord
	ToCategoryMonthlyPrices(c []*db.GetMonthlyCategoryRow) []*record.CategoriesMonthPriceRecord
	ToCategoryYearlyPrice(category *db.GetYearlyCategoryRow) *record.CategoriesYearPriceRecord
	ToCategoryYearlyPrices(c []*db.GetYearlyCategoryRow) []*record.CategoriesYearPriceRecord
	ToCategoryMonthlyPriceByMerchant(category *db.GetMonthlyCategoryByMerchantRow) *record.CategoriesMonthPriceRecord
	ToCategoryMonthlyPricesByMerchant(c []*db.GetMonthlyCategoryByMerchantRow) []*record.CategoriesMonthPriceRecord
	ToCategoryYearlyPriceByMerchant(category *db.GetYearlyCategoryByMerchantRow) *record.CategoriesYearPriceRecord
	ToCategoryYearlyPricesByMerchant(c []*db.GetYearlyCategoryByMerchantRow) []*record.CategoriesYearPriceRecord

	ToCategoryRecord(category *db.Category) *record.CategoriesRecord
	ToCategoryRecordPagination(category *db.GetCategoriesRow) *record.CategoriesRecord
	ToCategoriesRecordPagination(categories []*db.GetCategoriesRow) []*record.CategoriesRecord
	ToCategoryRecordActivePagination(category *db.GetCategoriesActiveRow) *record.CategoriesRecord
	ToCategoriesRecordActivePagination(categories []*db.GetCategoriesActiveRow) []*record.CategoriesRecord
	ToCategoryRecordTrashedPagination(category *db.GetCategoriesTrashedRow) *record.CategoriesRecord
	ToCategoriesRecordTrashedPagination(categories []*db.GetCategoriesTrashedRow) []*record.CategoriesRecord
}

type CashierRecordMapping interface {
	ToCashierMonthlyTotalSale(c *db.GetMonthlyTotalSalesCashierRow) *record.CashierRecordMonthTotalSales
	ToCashierMonthlyTotalSales(c []*db.GetMonthlyTotalSalesCashierRow) []*record.CashierRecordMonthTotalSales
	ToCashierYearlyTotalSale(c *db.GetYearlyTotalSalesCashierRow) *record.CashierRecordYearTotalSales
	ToCashierYearlyTotalSales(c []*db.GetYearlyTotalSalesCashierRow) []*record.CashierRecordYearTotalSales
	ToCashierMonthlyTotalSaleById(c *db.GetMonthlyTotalSalesByIdRow) *record.CashierRecordMonthTotalSales
	ToCashierMonthlyTotalSalesById(c []*db.GetMonthlyTotalSalesByIdRow) []*record.CashierRecordMonthTotalSales
	ToCashierYearlyTotalSaleById(c *db.GetYearlyTotalSalesByIdRow) *record.CashierRecordYearTotalSales
	ToCashierYearlyTotalSalesById(c []*db.GetYearlyTotalSalesByIdRow) []*record.CashierRecordYearTotalSales
	ToCashierMonthlyTotalSaleByMerchant(c *db.GetMonthlyTotalSalesByMerchantRow) *record.CashierRecordMonthTotalSales
	ToCashierMonthlyTotalSalesByMerchant(c []*db.GetMonthlyTotalSalesByMerchantRow) []*record.CashierRecordMonthTotalSales
	ToCashierYearlyTotalSaleByMerchant(c *db.GetYearlyTotalSalesByMerchantRow) *record.CashierRecordYearTotalSales
	ToCashierYearlyTotalSalesByMerchant(c []*db.GetYearlyTotalSalesByMerchantRow) []*record.CashierRecordYearTotalSales

	ToCashierMonthlySaleById(cashier *db.GetMonthlyCashierByCashierIdRow) *record.CashierRecordMonthSales
	ToCashierMonthlySalesById(c []*db.GetMonthlyCashierByCashierIdRow) []*record.CashierRecordMonthSales
	ToCashierYearlySaleById(cashier *db.GetYearlyCashierByCashierIdRow) *record.CashierRecordYearSales
	ToCashierYearlySalesById(c []*db.GetYearlyCashierByCashierIdRow) []*record.CashierRecordYearSales

	ToCashierMonthlySale(cashier *db.GetMonthlyCashierRow) *record.CashierRecordMonthSales
	ToCashierMonthlySales(c []*db.GetMonthlyCashierRow) []*record.CashierRecordMonthSales
	ToCashierYearlySale(cashier *db.GetYearlyCashierRow) *record.CashierRecordYearSales
	ToCashierYearlySales(c []*db.GetYearlyCashierRow) []*record.CashierRecordYearSales

	ToCashierMonthlySaleByMerchant(cashier *db.GetMonthlyCashierByMerchantRow) *record.CashierRecordMonthSales
	ToCashierMonthlySalesByMerchant(c []*db.GetMonthlyCashierByMerchantRow) []*record.CashierRecordMonthSales
	ToCashierYearlySaleByMerchant(cashier *db.GetYearlyCashierByMerchantRow) *record.CashierRecordYearSales
	ToCashierYearlySalesByMerchant(c []*db.GetYearlyCashierByMerchantRow) []*record.CashierRecordYearSales

	ToCashierRecord(Cashier *db.Cashier) *record.CashierRecord

	ToCashierRecordPagination(Cashier *db.GetCashiersRow) *record.CashierRecord
	ToCashiersRecordPagination(Cashiers []*db.GetCashiersRow) []*record.CashierRecord

	ToCashierRecordActivePagination(Cashier *db.GetCashiersActiveRow) *record.CashierRecord
	ToCashiersRecordActivePagination(Cashiers []*db.GetCashiersActiveRow) []*record.CashierRecord
	ToCashierRecordTrashedPagination(Cashier *db.GetCashiersTrashedRow) *record.CashierRecord
	ToCashiersRecordTrashedPagination(Cashiers []*db.GetCashiersTrashedRow) []*record.CashierRecord

	ToCashierMerchantRecordPagination(cashier *db.GetCashiersByMerchantRow) *record.CashierRecord
	ToCashiersMerchantRecordPagination(cashiers []*db.GetCashiersByMerchantRow) []*record.CashierRecord
}

type MerchantRecordMapping interface {
	ToMerchantRecord(Merchant *db.Merchant) *record.MerchantRecord
	ToMerchantRecordPagination(Merchant *db.GetMerchantsRow) *record.MerchantRecord
	ToMerchantsRecordPagination(Merchants []*db.GetMerchantsRow) []*record.MerchantRecord

	ToMerchantRecordActivePagination(Merchant *db.GetMerchantsActiveRow) *record.MerchantRecord
	ToMerchantsRecordActivePagination(Merchants []*db.GetMerchantsActiveRow) []*record.MerchantRecord
	ToMerchantRecordTrashedPagination(Merchant *db.GetMerchantsTrashedRow) *record.MerchantRecord
	ToMerchantsRecordTrashedPagination(Merchants []*db.GetMerchantsTrashedRow) []*record.MerchantRecord
}

type MerchantDocumentMapping interface {
	ToGetMerchantDocument(doc *db.MerchantDocument) *record.MerchantDocumentRecord

	ToMerchantDocumentRecord(doc *db.GetMerchantDocumentsRow) *record.MerchantDocumentRecord
	ToMerchantDocumentsRecord(docs []*db.GetMerchantDocumentsRow) []*record.MerchantDocumentRecord

	ToMerchantDocumentActiveRecord(doc *db.GetActiveMerchantDocumentsRow) *record.MerchantDocumentRecord
	ToMerchantDocumentsActiveRecord(docs []*db.GetActiveMerchantDocumentsRow) []*record.MerchantDocumentRecord

	ToMerchantDocumentTrashedRecord(doc *db.GetTrashedMerchantDocumentsRow) *record.MerchantDocumentRecord
	ToMerchantDocumentsTrashedRecord(docs []*db.GetTrashedMerchantDocumentsRow) []*record.MerchantDocumentRecord
}

type OrderItemRecordMapping interface {
	ToOrderItemRecord(orderItems *db.OrderItem) *record.OrderItemRecord
	ToOrderItemsRecord(orders []*db.OrderItem) []*record.OrderItemRecord

	ToOrderItemRecordPagination(OrderItem *db.GetOrderItemsRow) *record.OrderItemRecord
	ToOrderItemsRecordPagination(OrderItem []*db.GetOrderItemsRow) []*record.OrderItemRecord

	ToOrderItemRecordActivePagination(OrderItem *db.GetOrderItemsActiveRow) *record.OrderItemRecord
	ToOrderItemsRecordActivePagination(OrderItem []*db.GetOrderItemsActiveRow) []*record.OrderItemRecord
	ToOrderItemRecordTrashedPagination(OrderItem *db.GetOrderItemsTrashedRow) *record.OrderItemRecord
	ToOrderItemsRecordTrashedPagination(OrderItem []*db.GetOrderItemsTrashedRow) []*record.OrderItemRecord
}

type OrderRecordMapping interface {
	ToOrderMonthlyTotalRevenue(c *db.GetMonthlyTotalRevenueRow) *record.OrderMonthlyTotalRevenueRecord
	ToOrderMonthlyTotalRevenues(c []*db.GetMonthlyTotalRevenueRow) []*record.OrderMonthlyTotalRevenueRecord
	ToOrderYearlyTotalRevenue(c *db.GetYearlyTotalRevenueRow) *record.OrderYearlyTotalRevenueRecord
	ToOrderYearlyTotalRevenues(c []*db.GetYearlyTotalRevenueRow) []*record.OrderYearlyTotalRevenueRecord
	ToOrderMonthlyTotalRevenueById(c *db.GetMonthlyTotalRevenueByIdRow) *record.OrderMonthlyTotalRevenueRecord
	ToOrderMonthlyTotalRevenuesById(c []*db.GetMonthlyTotalRevenueByIdRow) []*record.OrderMonthlyTotalRevenueRecord
	ToOrderYearlyTotalRevenueById(c *db.GetYearlyTotalRevenueByIdRow) *record.OrderYearlyTotalRevenueRecord
	ToOrderYearlyTotalRevenuesById(c []*db.GetYearlyTotalRevenueByIdRow) []*record.OrderYearlyTotalRevenueRecord
	ToOrderMonthlyTotalRevenueByMerchant(c *db.GetMonthlyTotalRevenueByMerchantRow) *record.OrderMonthlyTotalRevenueRecord
	ToOrderMonthlyTotalRevenuesByMerchant(c []*db.GetMonthlyTotalRevenueByMerchantRow) []*record.OrderMonthlyTotalRevenueRecord
	ToOrderYearlyTotalRevenueByMerchant(c *db.GetYearlyTotalRevenueByMerchantRow) *record.OrderYearlyTotalRevenueRecord
	ToOrderYearlyTotalRevenuesByMerchant(c []*db.GetYearlyTotalRevenueByMerchantRow) []*record.OrderYearlyTotalRevenueRecord

	ToOrderMonthlyPrice(category *db.GetMonthlyOrderRow) *record.OrderMonthlyRecord
	ToOrderMonthlyPrices(c []*db.GetMonthlyOrderRow) []*record.OrderMonthlyRecord
	ToOrderYearlyPrice(category *db.GetYearlyOrderRow) *record.OrderYearlyRecord
	ToOrderYearlyPrices(c []*db.GetYearlyOrderRow) []*record.OrderYearlyRecord
	ToOrderMonthlyPriceByMerchant(category *db.GetMonthlyOrderByMerchantRow) *record.OrderMonthlyRecord
	ToOrderMonthlyPricesByMerchant(c []*db.GetMonthlyOrderByMerchantRow) []*record.OrderMonthlyRecord
	ToOrderYearlyPriceByMerchant(category *db.GetYearlyOrderByMerchantRow) *record.OrderYearlyRecord
	ToOrderYearlyPricesByMerchant(c []*db.GetYearlyOrderByMerchantRow) []*record.OrderYearlyRecord

	ToOrderRecord(order *db.Order) *record.OrderRecord
	ToOrdersRecord(orders []*db.Order) []*record.OrderRecord
	ToOrderRecordPagination(order *db.GetOrdersRow) *record.OrderRecord
	ToOrdersRecordPagination(orders []*db.GetOrdersRow) []*record.OrderRecord
	ToOrderRecordActivePagination(order *db.GetOrdersActiveRow) *record.OrderRecord
	ToOrdersRecordActivePagination(orders []*db.GetOrdersActiveRow) []*record.OrderRecord
	ToOrderRecordTrashedPagination(order *db.GetOrdersTrashedRow) *record.OrderRecord
	ToOrdersRecordTrashedPagination(orders []*db.GetOrdersTrashedRow) []*record.OrderRecord

	ToOrderRecordByMerchantPagination(order *db.GetOrdersByMerchantRow) *record.OrderRecord
	ToOrdersRecordByMerchantPagination(orders []*db.GetOrdersByMerchantRow) []*record.OrderRecord
}

type ProductRecordMapping interface {
	ToProductRecord(product *db.Product) *record.ProductRecord
	ToProductsRecord(products []*db.Product) []*record.ProductRecord
	ToProductRecordPagination(product *db.GetProductsRow) *record.ProductRecord
	ToProductsRecordPagination(products []*db.GetProductsRow) []*record.ProductRecord
	ToProductRecordActivePagination(product *db.GetProductsActiveRow) *record.ProductRecord
	ToProductsRecordActivePagination(products []*db.GetProductsActiveRow) []*record.ProductRecord
	ToProductRecordTrashedPagination(product *db.GetProductsTrashedRow) *record.ProductRecord
	ToProductsRecordTrashedPagination(products []*db.GetProductsTrashedRow) []*record.ProductRecord

	ToProductRecordMerchantPagination(product *db.GetProductsByMerchantRow) *record.ProductRecord
	ToProductsRecordMerchantPagination(products []*db.GetProductsByMerchantRow) []*record.ProductRecord

	ToProductRecordCategoryPagination(product *db.GetProductsByCategoryNameRow) *record.ProductRecord
	ToProductsRecordCategoryPagination(products []*db.GetProductsByCategoryNameRow) []*record.ProductRecord
}

type TransactionRecordMapping interface {
	ToTransactionMonthAmountSuccess(row *db.GetMonthlyAmountTransactionSuccessRow) *record.TransactionMonthlyAmountSuccessRecord
	ToTransactionMonthlyAmountSuccess(rows []*db.GetMonthlyAmountTransactionSuccessRow) []*record.TransactionMonthlyAmountSuccessRecord
	ToTransactionYearAmountSuccess(row *db.GetYearlyAmountTransactionSuccessRow) *record.TransactionYearlyAmountSuccessRecord
	ToTransactionYearlyAmountSuccess(rows []*db.GetYearlyAmountTransactionSuccessRow) []*record.TransactionYearlyAmountSuccessRecord
	ToTransactionMonthAmountFailed(row *db.GetMonthlyAmountTransactionFailedRow) *record.TransactionMonthlyAmountFailedRecord
	ToTransactionMonthlyAmountFailed(rows []*db.GetMonthlyAmountTransactionFailedRow) []*record.TransactionMonthlyAmountFailedRecord
	ToTransactionYearAmountFailed(row *db.GetYearlyAmountTransactionFailedRow) *record.TransactionYearlyAmountFailedRecord
	ToTransactionYearlyAmountFailed(rows []*db.GetYearlyAmountTransactionFailedRow) []*record.TransactionYearlyAmountFailedRecord
	ToTransactionMonthAmountSuccessByMerchant(row *db.GetMonthlyAmountTransactionSuccessByMerchantRow) *record.TransactionMonthlyAmountSuccessRecord
	ToTransactionMonthlyAmountSuccessByMerchant(rows []*db.GetMonthlyAmountTransactionSuccessByMerchantRow) []*record.TransactionMonthlyAmountSuccessRecord
	ToTransactionYearAmountSuccessByMerchant(row *db.GetYearlyAmountTransactionSuccessByMerchantRow) *record.TransactionYearlyAmountSuccessRecord
	ToTransactionYearlyAmountSuccessByMerchant(rows []*db.GetYearlyAmountTransactionSuccessByMerchantRow) []*record.TransactionYearlyAmountSuccessRecord
	ToTransactionMonthAmountFailedByMerchant(row *db.GetMonthlyAmountTransactionFailedByMerchantRow) *record.TransactionMonthlyAmountFailedRecord
	ToTransactionMonthlyAmountFailedByMerchant(rows []*db.GetMonthlyAmountTransactionFailedByMerchantRow) []*record.TransactionMonthlyAmountFailedRecord
	ToTransactionYearAmountFailedByMerchant(row *db.GetYearlyAmountTransactionFailedByMerchantRow) *record.TransactionYearlyAmountFailedRecord
	ToTransactionYearlyAmountFailedByMerchant(rows []*db.GetYearlyAmountTransactionFailedByMerchantRow) []*record.TransactionYearlyAmountFailedRecord

	ToTransactionMonthMethodSuccess(row *db.GetMonthlyTransactionMethodsSuccessRow) *record.TransactionMonthlyMethodRecord
	ToTransactionMonthlyMethodSuccess(rows []*db.GetMonthlyTransactionMethodsSuccessRow) []*record.TransactionMonthlyMethodRecord
	ToTransactionYearMethodSuccess(row *db.GetYearlyTransactionMethodsSuccessRow) *record.TransactionYearlyMethodRecord
	ToTransactionYearlyMethodSuccess(rows []*db.GetYearlyTransactionMethodsSuccessRow) []*record.TransactionYearlyMethodRecord

	ToTransactionMonthMethodFailed(row *db.GetMonthlyTransactionMethodsFailedRow) *record.TransactionMonthlyMethodRecord
	ToTransactionMonthlyMethodFailed(rows []*db.GetMonthlyTransactionMethodsFailedRow) []*record.TransactionMonthlyMethodRecord
	ToTransactionYearMethodFailed(row *db.GetYearlyTransactionMethodsFailedRow) *record.TransactionYearlyMethodRecord
	ToTransactionYearlyMethodFailed(rows []*db.GetYearlyTransactionMethodsFailedRow) []*record.TransactionYearlyMethodRecord

	ToTransactionMonthMethodByMerchantSuccess(row *db.GetMonthlyTransactionMethodsByMerchantSuccessRow) *record.TransactionMonthlyMethodRecord
	ToTransactionMonthlyByMerchantMethodSuccess(rows []*db.GetMonthlyTransactionMethodsByMerchantSuccessRow) []*record.TransactionMonthlyMethodRecord
	ToTransactionYearMethodByMerchantSuccess(row *db.GetYearlyTransactionMethodsByMerchantSuccessRow) *record.TransactionYearlyMethodRecord
	ToTransactionYearlyMethodByMerchantSuccess(rows []*db.GetYearlyTransactionMethodsByMerchantSuccessRow) []*record.TransactionYearlyMethodRecord

	ToTransactionMonthMethodByMerchantFailed(row *db.GetMonthlyTransactionMethodsByMerchantFailedRow) *record.TransactionMonthlyMethodRecord
	ToTransactionMonthlyByMerchantMethodFailed(rows []*db.GetMonthlyTransactionMethodsByMerchantFailedRow) []*record.TransactionMonthlyMethodRecord
	ToTransactionYearMethodByMerchantFailed(row *db.GetYearlyTransactionMethodsByMerchantFailedRow) *record.TransactionYearlyMethodRecord
	ToTransactionYearlyMethodByMerchantFailed(rows []*db.GetYearlyTransactionMethodsByMerchantFailedRow) []*record.TransactionYearlyMethodRecord

	ToTransactionRecord(transaction *db.Transaction) *record.TransactionRecord
	ToTransactionsRecord(transactions []*db.Transaction) []*record.TransactionRecord
	ToTransactionRecordPagination(transaction *db.GetTransactionsRow) *record.TransactionRecord
	ToTransactionsRecordPagination(transaction []*db.GetTransactionsRow) []*record.TransactionRecord
	ToTransactionRecordActivePagination(transaction *db.GetTransactionsActiveRow) *record.TransactionRecord
	ToTransactionsRecordActivePagination(transactions []*db.GetTransactionsActiveRow) []*record.TransactionRecord
	ToTransactionRecordTrashedPagination(transaction *db.GetTransactionsTrashedRow) *record.TransactionRecord
	ToTransactionsRecordTrashedPagination(products []*db.GetTransactionsTrashedRow) []*record.TransactionRecord

	ToTransactionMerchantRecordPagination(transaction *db.GetTransactionByMerchantRow) *record.TransactionRecord
	ToTransactionMerchantsRecordPagination(products []*db.GetTransactionByMerchantRow) []*record.TransactionRecord
}
