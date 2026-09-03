package repository

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

func setupTestDB(t *testing.T) (*db.Queries, *Repositories, *pgxpool.Pool, func()) {
	t.Helper()

	ctx := context.Background()

	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
	)
	require.NoError(t, err)

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err)

	runMigrations(t, ctx, pool, "../../../pkg/database/migrations")

	queries := db.New(pool)
	repos := NewRepositories(queries, nil, nil, nil, nil)

	cleanup := func() {
		pool.Close()
		pgContainer.Terminate(ctx)
	}

	return queries, repos, pool, cleanup
}

func runMigrations(t *testing.T, ctx context.Context, pool *pgxpool.Pool, migrationDir string) {
	t.Helper()

	migrations := []string{
		"20250204052452_add_pgcrypto_and_uuid_columns.sql",
		"20250204052518_create_users.sql",
		"20250204052555_create_role.sql",
		"20250204052605_create_user_role.sql",
		"20250204052613_create_refreshtoken.sql",
		"20250204052614_create_reset_token.sql",
		"20250206050402_create_merchants.sql",
		"20250206050403_create_merchant_documents.sql",
		"20250206050448_create_cashiers.sql",
		"20250206050455_create_categories.sql",
		"20250206050505_create_products.sql",
		"20250206050519_create_orders.sql",
		"20250206050525_create_order_items.sql",
		"20250206050554_create_transactions.sql",
	}

	for _, file := range migrations {
		path := migrationDir + "/" + file
		sql, err := os.ReadFile(path)
		require.NoError(t, err)

		_, err = pool.Exec(ctx, string(sql))
		if err != nil {
			t.Logf("Warning executing %s: %v", file, err)
		}
	}
}

func ptr[T any](v T) *T {
	return &v
}

func seedOrderDeps(t *testing.T, ctx context.Context, pool *pgxpool.Pool) (int, int) {
	t.Helper()

	var userID int
	err := pool.QueryRow(ctx, `INSERT INTO users (firstname, lastname, email, password, verification_code) 
		VALUES ('Order', 'User', 'orderuser@example.com', 'password', 'CODE') RETURNING user_id`).Scan(&userID)
	require.NoError(t, err)

	var merchantID int
	err = pool.QueryRow(ctx, `INSERT INTO merchants (user_id, name, status) VALUES ($1, 'Order Merchant', 'active') RETURNING merchant_id`, userID).Scan(&merchantID)
	require.NoError(t, err)

	var cashierID int
	err = pool.QueryRow(ctx, `INSERT INTO cashiers (merchant_id, user_id, name) VALUES ($1, $2, 'Order Cashier') RETURNING cashier_id`, merchantID, userID).Scan(&cashierID)
	require.NoError(t, err)

	return merchantID, cashierID
}

func TestOrderCommand_Create(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	merchantID, cashierID := seedOrderDeps(t, ctx, pool)

	order, err := repos.OrderCommand.CreateOrder(ctx, &requests.CreateOrderRecordRequest{
		MerchantID: merchantID,
		CashierID:  cashierID,
		TotalPrice: 50000,
	})
	require.NoError(t, err)
	require.NotNil(t, order)
	assert.Equal(t, int32(merchantID), order.MerchantID)
	assert.Equal(t, int32(cashierID), order.CashierID)
	assert.Equal(t, int64(50000), order.TotalPrice)
	assert.NotZero(t, order.OrderID)
}

func TestOrderQuery_FindByID(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	merchantID, cashierID := seedOrderDeps(t, ctx, pool)

	created, err := repos.OrderCommand.CreateOrder(ctx, &requests.CreateOrderRecordRequest{
		MerchantID: merchantID, CashierID: cashierID, TotalPrice: 75000,
	})
	require.NoError(t, err)

	found, err := repos.OrderQuery.FindById(ctx, int(created.OrderID))
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, created.OrderID, found.OrderID)
	assert.Equal(t, int64(75000), found.TotalPrice)
}

func TestOrderCommand_Update(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	merchantID, cashierID := seedOrderDeps(t, ctx, pool)

	created, err := repos.OrderCommand.CreateOrder(ctx, &requests.CreateOrderRecordRequest{
		MerchantID: merchantID, CashierID: cashierID, TotalPrice: 30000,
	})
	require.NoError(t, err)

	updated, err := repos.OrderCommand.UpdateOrder(ctx, &requests.UpdateOrderRecordRequest{
		OrderID:    int(created.OrderID),
		TotalPrice: 60000,
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, int64(60000), updated.TotalPrice)
}

func TestOrderQuery_FindAll(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	merchantID, cashierID := seedOrderDeps(t, ctx, pool)

	_, err := repos.OrderCommand.CreateOrder(ctx, &requests.CreateOrderRecordRequest{
		MerchantID: merchantID, CashierID: cashierID, TotalPrice: 10000,
	})
	require.NoError(t, err)
	_, err = repos.OrderCommand.CreateOrder(ctx, &requests.CreateOrderRecordRequest{
		MerchantID: merchantID, CashierID: cashierID, TotalPrice: 20000,
	})
	require.NoError(t, err)

	results, total, err := repos.OrderQuery.FindAllOrders(ctx, &requests.FindAllOrders{
		Search:   "",
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	require.NotNil(t, total)
	assert.GreaterOrEqual(t, *total, 2)
	assert.GreaterOrEqual(t, len(results), 2)
}

func TestOrderCommand_Trash(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	merchantID, cashierID := seedOrderDeps(t, ctx, pool)
	created, err := repos.OrderCommand.CreateOrder(ctx, &requests.CreateOrderRecordRequest{
		MerchantID: merchantID, CashierID: cashierID, TotalPrice: 50000,
	})
	require.NoError(t, err)

	trashed, err := repos.OrderCommand.TrashedOrder(ctx, int(created.OrderID))
	require.NoError(t, err)
	require.NotNil(t, trashed)
	assert.True(t, trashed.DeletedAt.Valid)
}

func TestOrderQuery_FindByTrashed(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	merchantID, cashierID := seedOrderDeps(t, ctx, pool)
	created, err := repos.OrderCommand.CreateOrder(ctx, &requests.CreateOrderRecordRequest{
		MerchantID: merchantID, CashierID: cashierID, TotalPrice: 50000,
	})
	require.NoError(t, err)

	_, err = repos.OrderCommand.TrashedOrder(ctx, int(created.OrderID))
	require.NoError(t, err)

	results, total, err := repos.OrderQuery.FindByTrashed(ctx, &requests.FindAllOrders{
		Search:   "",
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	require.NotNil(t, total)
	assert.GreaterOrEqual(t, *total, 1)
	assert.GreaterOrEqual(t, len(results), 1)
}

func TestOrderCommand_Restore(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	merchantID, cashierID := seedOrderDeps(t, ctx, pool)
	created, err := repos.OrderCommand.CreateOrder(ctx, &requests.CreateOrderRecordRequest{
		MerchantID: merchantID, CashierID: cashierID, TotalPrice: 50000,
	})
	require.NoError(t, err)

	_, err = repos.OrderCommand.TrashedOrder(ctx, int(created.OrderID))
	require.NoError(t, err)

	restored, err := repos.OrderCommand.RestoreOrder(ctx, int(created.OrderID))
	require.NoError(t, err)
	require.NotNil(t, restored)
	assert.False(t, restored.DeletedAt.Valid)
}

func TestOrderQuery_FindByActive(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	merchantID, cashierID := seedOrderDeps(t, ctx, pool)
	created, err := repos.OrderCommand.CreateOrder(ctx, &requests.CreateOrderRecordRequest{
		MerchantID: merchantID, CashierID: cashierID, TotalPrice: 50000,
	})
	require.NoError(t, err)

	results, total, err := repos.OrderQuery.FindByActive(ctx, &requests.FindAllOrders{
		Search:   "",
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	require.NotNil(t, total)
	assert.GreaterOrEqual(t, *total, 1)

	found := false
	for _, r := range results {
		if int(r.OrderID) == int(created.OrderID) {
			found = true
			break
		}
	}
	assert.True(t, found)
}

func TestOrderCommand_DeletePermanent(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	merchantID, cashierID := seedOrderDeps(t, ctx, pool)
	created, err := repos.OrderCommand.CreateOrder(ctx, &requests.CreateOrderRecordRequest{
		MerchantID: merchantID, CashierID: cashierID, TotalPrice: 50000,
	})
	require.NoError(t, err)

	_, err = repos.OrderCommand.TrashedOrder(ctx, int(created.OrderID))
	require.NoError(t, err)

	deleted, err := repos.OrderCommand.DeleteOrderPermanent(ctx, int(created.OrderID))
	require.NoError(t, err)
	assert.True(t, deleted)

	_, err = repos.OrderQuery.FindById(ctx, int(created.OrderID))
	assert.Error(t, err)
}

func TestOrderQuery_FindByMerchant(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	merchantID, cashierID := seedOrderDeps(t, ctx, pool)
	_, err := repos.OrderCommand.CreateOrder(ctx, &requests.CreateOrderRecordRequest{
		MerchantID: merchantID, CashierID: cashierID, TotalPrice: 50000,
	})
	require.NoError(t, err)

	results, total, err := repos.OrderQuery.FindByMerchant(ctx, &requests.FindAllOrderMerchant{
		MerchantID: merchantID,
		Search:     "",
		Page:       1,
		PageSize:   10,
	})
	require.NoError(t, err)
	require.NotNil(t, total)
	assert.GreaterOrEqual(t, *total, 1)
	assert.GreaterOrEqual(t, len(results), 1)
}

// ----- OrderItemCommand Tests (DB-backed in order service) -----

func TestOrderItemCommand_Create(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	merchantID, cashierID := seedOrderDeps(t, ctx, pool)

	// Seed product dependencies
	var categoryID int
	err := pool.QueryRow(ctx, `INSERT INTO categories (name) VALUES ('Test Cat') RETURNING category_id`).Scan(&categoryID)
	require.NoError(t, err)

	var productID int
	err = pool.QueryRow(ctx, `INSERT INTO products (merchant_id, category_id, name, price, count_in_stock) 
		VALUES ($1, $2, 'Test Product', 10000, 50) RETURNING product_id`, merchantID, categoryID).Scan(&productID)
	require.NoError(t, err)

	created, err := repos.OrderCommand.CreateOrder(ctx, &requests.CreateOrderRecordRequest{
		MerchantID: merchantID, CashierID: cashierID, TotalPrice: 20000,
	})
	require.NoError(t, err)

	item, err := repos.OrderItemCommand.CreateOrderItem(ctx, &requests.CreateOrderItemRecordRequest{
		OrderID:   int(created.OrderID),
		ProductID: productID,
		Quantity:  2,
		Price:     10000,
	})
	require.NoError(t, err)
	require.NotNil(t, item)
	assert.Equal(t, int32(created.OrderID), item.OrderID)
	assert.Equal(t, int32(productID), item.ProductID)
	assert.Equal(t, int32(2), item.Quantity)
	assert.NotZero(t, item.OrderItemID)
}

func TestOrderItemCommand_Update(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	merchantID, cashierID := seedOrderDeps(t, ctx, pool)

	var categoryID int
	err := pool.QueryRow(ctx, `INSERT INTO categories (name) VALUES ('Cat') RETURNING category_id`).Scan(&categoryID)
	require.NoError(t, err)

	var productID int
	err = pool.QueryRow(ctx, `INSERT INTO products (merchant_id, category_id, name, price, count_in_stock) 
		VALUES ($1, $2, 'P', 5000, 10) RETURNING product_id`, merchantID, categoryID).Scan(&productID)
	require.NoError(t, err)

	created, err := repos.OrderCommand.CreateOrder(ctx, &requests.CreateOrderRecordRequest{
		MerchantID: merchantID, CashierID: cashierID, TotalPrice: 10000,
	})
	require.NoError(t, err)

	item, err := repos.OrderItemCommand.CreateOrderItem(ctx, &requests.CreateOrderItemRecordRequest{
		OrderID: int(created.OrderID), ProductID: productID, Quantity: 1, Price: 5000,
	})
	require.NoError(t, err)

	updated, err := repos.OrderItemCommand.UpdateOrderItem(ctx, &requests.UpdateOrderItemRecordRequest{
		OrderItemID: int(item.OrderItemID),
		OrderID:     int(created.OrderID),
		ProductID:   productID,
		Quantity:    3,
		Price:       6000,
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, int32(3), updated.Quantity)
	assert.Equal(t, int32(6000), updated.Price)
}

// ----- OrderStats Tests -----

func TestOrderStats_GetMonthlyTotalRevenue(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	results, err := repos.OrderStats.GetMonthlyTotalRevenue(ctx, &requests.MonthTotalRevenue{
		Year:  2026,
		Month: 1,
	})
	require.NoError(t, err)
	require.NotNil(t, results)
}

func TestOrderStats_GetYearlyTotalRevenue(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	results, err := repos.OrderStats.GetYearlyTotalRevenue(ctx, 2026)
	require.NoError(t, err)
	require.NotNil(t, results)
}
