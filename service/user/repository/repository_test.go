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

func TestUserCommand_Create(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	user, err := repos.UserCommand.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName:       "Alice",
		LastName:        "Johnson",
		Email:           "alice@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	})
	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, "Alice", user.Firstname)
	assert.Equal(t, "Johnson", user.Lastname)
	assert.Equal(t, "alice@example.com", user.Email)
	assert.NotZero(t, user.UserID)
}

func TestUserQuery_FindByID(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	created, err := repos.UserCommand.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName:       "Bob",
		LastName:        "Smith",
		Email:           "bob@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	})
	require.NoError(t, err)

	found, err := repos.UserQuery.FindById(ctx, int(created.UserID))
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, created.UserID, found.UserID)
	assert.Equal(t, "Bob", found.Firstname)
}

func TestUserQuery_FindByEmail(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	_, err := repos.UserCommand.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName:       "Email",
		LastName:        "Test",
		Email:           "emailtest@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	})
	require.NoError(t, err)

	found, err := repos.UserQuery.FindByEmail(ctx, "emailtest@example.com")
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, "emailtest@example.com", found.Email)
}

func TestUserCommand_Update(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	created, err := repos.UserCommand.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName:       "Update",
		LastName:        "Me",
		Email:           "updateme@example.com",
		Password:        "oldpassword",
		ConfirmPassword: "oldpassword",
	})
	require.NoError(t, err)

	updated, err := repos.UserCommand.UpdateUser(ctx, &requests.UpdateUserRequest{
		UserID:          ptr(int(created.UserID)),
		FirstName:       "Updated",
		LastName:        "User",
		Email:           "updated@example.com",
		Password:        "newpassword",
		ConfirmPassword: "newpassword",
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, "Updated", updated.Firstname)
	assert.Equal(t, "updated@example.com", updated.Email)
}

func TestUserQuery_FindAll(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	_, err := repos.UserCommand.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "User1", LastName: "Test1", Email: "user1@example.com", Password: "password123", ConfirmPassword: "password123",
	})
	require.NoError(t, err)
	_, err = repos.UserCommand.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "User2", LastName: "Test2", Email: "user2@example.com", Password: "password123", ConfirmPassword: "password123",
	})
	require.NoError(t, err)

	results, total, err := repos.UserQuery.FindAllUsers(ctx, &requests.FindAllUsers{
		Search:   "",
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	require.NotNil(t, total)
	assert.GreaterOrEqual(t, *total, 2)
	assert.GreaterOrEqual(t, len(results), 2)
}

func TestUserCommand_Trash(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	created, err := repos.UserCommand.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "Trash", LastName: "User", Email: "trash@example.com", Password: "password123", ConfirmPassword: "password123",
	})
	require.NoError(t, err)

	trashed, err := repos.UserCommand.TrashedUser(ctx, int(created.UserID))
	require.NoError(t, err)
	require.NotNil(t, trashed)
	assert.True(t, trashed.DeletedAt.Valid)
}

func TestUserQuery_FindByTrashed(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	created, err := repos.UserCommand.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "TrashFind", LastName: "User", Email: "trashfind@example.com", Password: "password123", ConfirmPassword: "password123",
	})
	require.NoError(t, err)

	_, err = repos.UserCommand.TrashedUser(ctx, int(created.UserID))
	require.NoError(t, err)

	results, total, err := repos.UserQuery.FindByTrashed(ctx, &requests.FindAllUsers{
		Search:   "",
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	require.NotNil(t, total)
	assert.GreaterOrEqual(t, *total, 1)
	assert.GreaterOrEqual(t, len(results), 1)
}

func TestUserCommand_Restore(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	created, err := repos.UserCommand.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "Restore", LastName: "User", Email: "restore@example.com", Password: "password123", ConfirmPassword: "password123",
	})
	require.NoError(t, err)

	_, err = repos.UserCommand.TrashedUser(ctx, int(created.UserID))
	require.NoError(t, err)

	restored, err := repos.UserCommand.RestoreUser(ctx, int(created.UserID))
	require.NoError(t, err)
	require.NotNil(t, restored)
	assert.False(t, restored.DeletedAt.Valid)
}

func TestUserQuery_FindByActive(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	created, err := repos.UserCommand.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "Active", LastName: "User", Email: "active@example.com", Password: "password123", ConfirmPassword: "password123",
	})
	require.NoError(t, err)

	results, total, err := repos.UserQuery.FindByActive(ctx, &requests.FindAllUsers{
		Search:   "",
		Page:     1,
		PageSize: 10,
	})
	require.NoError(t, err)
	require.NotNil(t, total)
	assert.GreaterOrEqual(t, *total, 1)

	found := false
	for _, r := range results {
		if int(r.UserID) == int(created.UserID) {
			found = true
			break
		}
	}
	assert.True(t, found)
}

func TestUserCommand_DeletePermanent(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	created, err := repos.UserCommand.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "PermDelete", LastName: "User", Email: "permdelete@example.com", Password: "password123", ConfirmPassword: "password123",
	})
	require.NoError(t, err)

	_, err = repos.UserCommand.TrashedUser(ctx, int(created.UserID))
	require.NoError(t, err)

	deleted, err := repos.UserCommand.DeleteUserPermanent(ctx, int(created.UserID))
	require.NoError(t, err)
	assert.True(t, deleted)

	_, err = repos.UserQuery.FindById(ctx, int(created.UserID))
	assert.Error(t, err)
}

func TestUserRole_FindByName(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	_, err := pool.Exec(ctx, `INSERT INTO roles (role_name) VALUES ('admin')`)
	require.NoError(t, err)

	found, err := repos.Role.FindByName(ctx, "admin")
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, "admin", found.RoleName)
}
