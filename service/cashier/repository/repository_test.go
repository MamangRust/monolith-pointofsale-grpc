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
	repos := NewRepositories(queries, nil, nil)

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

func seedCashierDeps(t *testing.T, ctx context.Context, pool *pgxpool.Pool) (int, int) {
	t.Helper()

	var userID int
	err := pool.QueryRow(ctx, `INSERT INTO users (firstname, lastname, email, password, verification_code) 
		VALUES ('Cashier', 'Deps', 'cashierdeps@example.com', 'password', 'CODE') RETURNING user_id`).Scan(&userID)
	require.NoError(t, err)

	var merchantID int
	err = pool.QueryRow(ctx, `INSERT INTO merchants (user_id, name, status) VALUES ($1, 'Cashier Merchant', 'active') RETURNING merchant_id`, userID).Scan(&merchantID)
	require.NoError(t, err)

	return userID, merchantID
}

func TestCashierCommand_Create(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	userID, merchantID := seedCashierDeps(t, ctx, pool)

	cashier, err := repos.CashierCommand.CreateCashier(ctx, &requests.CreateCashierRequest{
		MerchantID: merchantID,
		UserID:     userID,
		Name:       "John Cashier",
	})
	require.NoError(t, err)
	require.NotNil(t, cashier)
	assert.Equal(t, "John Cashier", cashier.Name)
	assert.Equal(t, int32(merchantID), cashier.MerchantID)
	assert.Equal(t, int32(userID), cashier.UserID)
	assert.NotZero(t, cashier.CashierID)
}

func TestCashierQuery_FindByID(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	userID, merchantID := seedCashierDeps(t, ctx, pool)

	created, err := repos.CashierCommand.CreateCashier(ctx, &requests.CreateCashierRequest{
		MerchantID: merchantID,
		UserID:     userID,
		Name:       "Findable Cashier",
	})
	require.NoError(t, err)

	found, err := repos.CashierQuery.FindById(ctx, int(created.CashierID))
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, created.CashierID, found.CashierID)
	assert.Equal(t, "Findable Cashier", found.Name)
}

func TestCashierCommand_Update(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	userID, merchantID := seedCashierDeps(t, ctx, pool)

	created, err := repos.CashierCommand.CreateCashier(ctx, &requests.CreateCashierRequest{
		MerchantID: merchantID,
		UserID:     userID,
		Name:       "Old Name",
	})
	require.NoError(t, err)

	updated, err := repos.CashierCommand.UpdateCashier(ctx, &requests.UpdateCashierRequest{
		CashierID: ptr(int(created.CashierID)),
		Name:      "New Name",
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, "New Name", updated.Name)
}

func TestCashierQuery_FindAll(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	userID, merchantID := seedCashierDeps(t, ctx, pool)

	_, err := repos.CashierCommand.CreateCashier(ctx, &requests.CreateCashierRequest{
		MerchantID: merchantID, UserID: userID, Name: "Cashier1",
	})
	require.NoError(t, err)
	_, err = repos.CashierCommand.CreateCashier(ctx, &requests.CreateCashierRequest{
		MerchantID: merchantID, UserID: userID, Name: "Cashier2",
	})
	require.NoError(t, err)

	results, total, err := repos.CashierQuery.FindAllCashiers(ctx, &requests.FindAllCashiers{
		Search:   "",
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	require.NotNil(t, total)
	assert.GreaterOrEqual(t, *total, 2)
	assert.GreaterOrEqual(t, len(results), 2)
}

func TestCashierCommand_Trash(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	userID, merchantID := seedCashierDeps(t, ctx, pool)
	created, err := repos.CashierCommand.CreateCashier(ctx, &requests.CreateCashierRequest{
		MerchantID: merchantID, UserID: userID, Name: "TrashCashier",
	})
	require.NoError(t, err)

	trashed, err := repos.CashierCommand.TrashedCashier(ctx, int(created.CashierID))
	require.NoError(t, err)
	require.NotNil(t, trashed)
	assert.True(t, trashed.DeletedAt.Valid)
}

func TestCashierQuery_FindByTrashed(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	userID, merchantID := seedCashierDeps(t, ctx, pool)
	created, err := repos.CashierCommand.CreateCashier(ctx, &requests.CreateCashierRequest{
		MerchantID: merchantID, UserID: userID, Name: "TrashFind",
	})
	require.NoError(t, err)

	_, err = repos.CashierCommand.TrashedCashier(ctx, int(created.CashierID))
	require.NoError(t, err)

	results, total, err := repos.CashierQuery.FindByTrashed(ctx, &requests.FindAllCashiers{
		Search:   "",
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	require.NotNil(t, total)
	assert.GreaterOrEqual(t, *total, 1)
	assert.GreaterOrEqual(t, len(results), 1)
}

func TestCashierCommand_Restore(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	userID, merchantID := seedCashierDeps(t, ctx, pool)
	created, err := repos.CashierCommand.CreateCashier(ctx, &requests.CreateCashierRequest{
		MerchantID: merchantID, UserID: userID, Name: "RestoreCashier",
	})
	require.NoError(t, err)

	_, err = repos.CashierCommand.TrashedCashier(ctx, int(created.CashierID))
	require.NoError(t, err)

	restored, err := repos.CashierCommand.RestoreCashier(ctx, int(created.CashierID))
	require.NoError(t, err)
	require.NotNil(t, restored)
	assert.False(t, restored.DeletedAt.Valid)
}

func TestCashierQuery_FindByActive(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	userID, merchantID := seedCashierDeps(t, ctx, pool)
	created, err := repos.CashierCommand.CreateCashier(ctx, &requests.CreateCashierRequest{
		MerchantID: merchantID, UserID: userID, Name: "ActiveCashier",
	})
	require.NoError(t, err)

	results, total, err := repos.CashierQuery.FindByActive(ctx, &requests.FindAllCashiers{
		Search:   "",
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	require.NotNil(t, total)
	assert.GreaterOrEqual(t, *total, 1)

	found := false
	for _, r := range results {
		if int(r.CashierID) == int(created.CashierID) {
			found = true
			break
		}
	}
	assert.True(t, found)
}

func TestCashierCommand_DeletePermanent(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	userID, merchantID := seedCashierDeps(t, ctx, pool)
	created, err := repos.CashierCommand.CreateCashier(ctx, &requests.CreateCashierRequest{
		MerchantID: merchantID, UserID: userID, Name: "PermDelete",
	})
	require.NoError(t, err)

	_, err = repos.CashierCommand.TrashedCashier(ctx, int(created.CashierID))
	require.NoError(t, err)

	deleted, err := repos.CashierCommand.DeleteCashierPermanent(ctx, int(created.CashierID))
	require.NoError(t, err)
	assert.True(t, deleted)

	_, err = repos.CashierQuery.FindById(ctx, int(created.CashierID))
	assert.Error(t, err)
}

func TestCashierQuery_FindByMerchant(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	userID, merchantID := seedCashierDeps(t, ctx, pool)

	_, err := repos.CashierCommand.CreateCashier(ctx, &requests.CreateCashierRequest{
		MerchantID: merchantID, UserID: userID, Name: "MCashier",
	})
	require.NoError(t, err)

	results, total, err := repos.CashierQuery.FindByMerchant(ctx, &requests.FindAllCashierMerchant{
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

func TestCashierStats_GetMonthlyTotalSales(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	userID, merchantID := seedCashierDeps(t, ctx, pool)

	_, err := repos.CashierCommand.CreateCashier(ctx, &requests.CreateCashierRequest{
		MerchantID: merchantID, UserID: userID, Name: "StatsCashier",
	})
	require.NoError(t, err)

	results, err := repos.CashierStats.GetMonthlyTotalSales(ctx, &requests.MonthTotalSales{
		Year:  2026,
		Month: 1,
	})
	require.NoError(t, err)
	require.NotNil(t, results)
}

func TestCashierStats_GetYearlyTotalSales(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	results, err := repos.CashierStats.GetYearlyTotalSales(ctx, 2026)
	require.NoError(t, err)
	require.NotNil(t, results)
}

func TestCashierStats_GetMonthlyCashier(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	results, err := repos.CashierStats.GetMonthyCashier(ctx, 2026)
	require.NoError(t, err)
	require.NotNil(t, results)
}

func TestCashierStats_GetYearlyCashier(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	results, err := repos.CashierStats.GetYearlyCashier(ctx, 2026)
	require.NoError(t, err)
	require.NotNil(t, results)
}
