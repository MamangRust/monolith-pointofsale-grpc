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
	repos := NewRepositories(queries)

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

func seedOrderItemData(t *testing.T, ctx context.Context, pool *pgxpool.Pool) (int, int) {
	t.Helper()

	// Seed user
	var userID int
	err := pool.QueryRow(ctx, `INSERT INTO users (firstname, lastname, email, password, verification_code) 
		VALUES ('Order', 'User', 'orderuser@example.com', 'password', 'CODE') RETURNING user_id`).Scan(&userID)
	require.NoError(t, err)

	// Seed merchant
	var merchantID int
	err = pool.QueryRow(ctx, `INSERT INTO merchants (user_id, name, status) VALUES ($1, 'Test Merchant', 'active') RETURNING merchant_id`, userID).Scan(&merchantID)
	require.NoError(t, err)

	// Seed category
	var categoryID int
	err = pool.QueryRow(ctx, `INSERT INTO categories (name) VALUES ('Test Category') RETURNING category_id`).Scan(&categoryID)
	require.NoError(t, err)

	// Seed product
	var productID int
	err = pool.QueryRow(ctx, `INSERT INTO products (merchant_id, category_id, name, price, count_in_stock) 
		VALUES ($1, $2, 'Test Product', 10000, 50) RETURNING product_id`, merchantID, categoryID).Scan(&productID)
	require.NoError(t, err)

	// Seed cashier
	var cashierID int
	err = pool.QueryRow(ctx, `INSERT INTO cashiers (merchant_id, user_id, name) 
		VALUES ($1, $2, 'Test Cashier') RETURNING cashier_id`, merchantID, userID).Scan(&cashierID)
	require.NoError(t, err)

	// Seed order
	var orderID int
	err = pool.QueryRow(ctx, `INSERT INTO orders (merchant_id, cashier_id, total_price) 
		VALUES ($1, $2, 10000) RETURNING order_id`, merchantID, cashierID).Scan(&orderID)
	require.NoError(t, err)

	// Seed order item
	_, err = pool.Exec(ctx, `INSERT INTO order_items (order_id, product_id, quantity, price) VALUES ($1, $2, 2, 10000)`, orderID, productID)
	require.NoError(t, err)

	return orderID, productID
}

func TestOrderItemQuery_FindAll(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	seedOrderItemData(t, ctx, pool)

	results, total, err := repos.OrderItemQuery.FindAllOrderItems(ctx, &requests.FindAllOrderItems{
		Search:   "",
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	require.NotNil(t, total)
	assert.GreaterOrEqual(t, *total, 1)
	assert.GreaterOrEqual(t, len(results), 1)
}

func TestOrderItemQuery_FindByActive(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	seedOrderItemData(t, ctx, pool)

	results, total, err := repos.OrderItemQuery.FindByActive(ctx, &requests.FindAllOrderItems{
		Search:   "",
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	require.NotNil(t, total)
	assert.GreaterOrEqual(t, *total, 1)
	assert.GreaterOrEqual(t, len(results), 1)
}

func TestOrderItemQuery_FindByTrashed(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	orderID, _ := seedOrderItemData(t, ctx, pool)

	// First trash the order item
	_, err := pool.Exec(ctx, `UPDATE order_items SET deleted_at = NOW() WHERE order_id = $1`, orderID)
	require.NoError(t, err)

	results, total, err := repos.OrderItemQuery.FindByTrashed(ctx, &requests.FindAllOrderItems{
		Search:   "",
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	require.NotNil(t, total)
	assert.GreaterOrEqual(t, *total, 1)
	assert.GreaterOrEqual(t, len(results), 1)
}

func TestOrderItemQuery_FindOrderItemByOrder(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	orderID, _ := seedOrderItemData(t, ctx, pool)

	items, err := repos.OrderItemQuery.FindOrderItemByOrder(ctx, orderID)
	require.NoError(t, err)
	require.NotNil(t, items)
	assert.GreaterOrEqual(t, len(items), 1)
	assert.Equal(t, int32(orderID), items[0].OrderID)
}
