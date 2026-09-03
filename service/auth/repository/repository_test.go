package repository

import (
	"context"
	"os"
	"testing"
	"time"

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

	// Run migrations
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

	// Order matters: create extensions first, then users, then dependent tables
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

func seedUser(t *testing.T, ctx context.Context, pool *pgxpool.Pool) int {
	t.Helper()
	var userID int
	err := pool.QueryRow(ctx, `INSERT INTO users (firstname, lastname, email, password, verification_code, is_verified) 
		VALUES ('John', 'Doe', 'john@example.com', 'hashedpassword', 'VERIF123', false) 
		RETURNING user_id`).Scan(&userID)
	require.NoError(t, err)
	return userID
}

func seedRole(t *testing.T, ctx context.Context, pool *pgxpool.Pool) int {
	t.Helper()
	var roleID int
	err := pool.QueryRow(ctx, `INSERT INTO roles (role_name) VALUES ('admin') RETURNING role_id`).Scan(&roleID)
	require.NoError(t, err)
	return roleID
}

// ----- User Repository Tests -----

func TestAuthUser_CreateUser(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	req := &requests.RegisterRequest{
		FirstName:       "Jane",
		LastName:        "Smith",
		Email:           "jane@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
		VerifiedCode:    "CODE001",
		IsVerified:      false,
	}

	user, err := repos.User.CreateUser(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, "Jane", user.Firstname)
	assert.Equal(t, "Smith", user.Lastname)
	assert.Equal(t, "jane@example.com", user.Email)
	assert.NotEmpty(t, user.Password)
	assert.Equal(t, "CODE001", user.VerificationCode)
	assert.NotNil(t, user.UserID)
}

func TestAuthUser_FindByEmail(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	req := &requests.RegisterRequest{
		FirstName:       "Find",
		LastName:        "ByEmail",
		Email:           "findbyemail@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
		VerifiedCode:    "CODE002",
		IsVerified:      true,
	}

	created, err := repos.User.CreateUser(ctx, req)
	require.NoError(t, err)

	found, err := repos.User.FindByEmail(ctx, created.Email)
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, created.UserID, found.UserID)
	assert.Equal(t, "Find", found.Firstname)
}

func TestAuthUser_FindById(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	req := &requests.RegisterRequest{
		FirstName:       "Find",
		LastName:        "ById",
		Email:           "findbyid@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
		VerifiedCode:    "CODE003",
		IsVerified:      false,
	}

	created, err := repos.User.CreateUser(ctx, req)
	require.NoError(t, err)

	found, err := repos.User.FindById(ctx, int(created.UserID))
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, created.UserID, found.UserID)
}

func TestAuthUser_UpdateUserIsVerified(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	req := &requests.RegisterRequest{
		FirstName:       "Verify",
		LastName:        "Test",
		Email:           "verify@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
		VerifiedCode:    "CODE004",
		IsVerified:      false,
	}

	created, err := repos.User.CreateUser(ctx, req)
	require.NoError(t, err)

	updated, err := repos.User.UpdateUserIsVerified(ctx, int(created.UserID), true)
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.NotNil(t, updated.IsVerified)
	assert.True(t, *updated.IsVerified)
}

func TestAuthUser_UpdateUserPassword(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	req := &requests.RegisterRequest{
		FirstName:       "Pass",
		LastName:        "Test",
		Email:           "pass@example.com",
		Password:        "oldpassword",
		ConfirmPassword: "oldpassword",
		VerifiedCode:    "CODE005",
		IsVerified:      false,
	}

	created, err := repos.User.CreateUser(ctx, req)
	require.NoError(t, err)

	updated, err := repos.User.UpdateUserPassword(ctx, int(created.UserID), "newpassword")
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, "newpassword", updated.Password)
}

func TestAuthUser_FindByVerificationCode(t *testing.T) {
	_, repos, _, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	req := &requests.RegisterRequest{
		FirstName:       "Code",
		LastName:        "Test",
		Email:           "code@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
		VerifiedCode:    "UNIQUE_VERIF_CODE",
		IsVerified:      false,
	}

	created, err := repos.User.CreateUser(ctx, req)
	require.NoError(t, err)

	found, err := repos.User.FindByVerificationCode(ctx, created.VerificationCode)
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, created.UserID, found.UserID)
}

// ----- Role Repository Tests (read-only) -----

func TestAuthRole_FindById(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	roleID := seedRole(t, ctx, pool)

	found, err := repos.Role.FindById(ctx, roleID)
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, "admin", found.RoleName)
}

func TestAuthRole_FindByName(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	seedRole(t, ctx, pool)

	found, err := repos.Role.FindByName(ctx, "admin")
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, "admin", found.RoleName)
}

// ----- UserRole Repository Tests -----

func TestAuthUserRole_AssignAndRemoveRole(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	userID := seedUser(t, ctx, pool)
	roleID := seedRole(t, ctx, pool)

	assigned, err := repos.UserRole.AssignRoleToUser(ctx, &requests.CreateUserRoleRequest{
		UserId: userID,
		RoleId: roleID,
	})
	require.NoError(t, err)
	require.NotNil(t, assigned)
	assert.Equal(t, int32(userID), assigned.UserID)
	assert.Equal(t, int32(roleID), assigned.RoleID)

	err = repos.UserRole.RemoveRoleFromUser(ctx, &requests.RemoveUserRoleRequest{
		UserId: userID,
		RoleId: roleID,
	})
	require.NoError(t, err)
}

// ----- RefreshToken Repository Tests -----

func TestAuthRefreshToken_CreateAndFind(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	userID := seedUser(t, ctx, pool)

	expiresAt := time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05")
	created, err := repos.RefreshToken.CreateRefreshToken(ctx, &requests.CreateRefreshToken{
		UserId:    userID,
		Token:     "test-refresh-token-001",
		ExpiresAt: expiresAt,
	})
	require.NoError(t, err)
	require.NotNil(t, created)
	assert.Equal(t, int32(userID), created.UserID)
	assert.Equal(t, "test-refresh-token-001", created.Token)

	// Find by token
	foundByToken, err := repos.RefreshToken.FindByToken(ctx, created.Token)
	require.NoError(t, err)
	require.NotNil(t, foundByToken)
	assert.Equal(t, created.RefreshTokenID, foundByToken.RefreshTokenID)

	// Find by user ID
	foundByUser, err := repos.RefreshToken.FindByUserId(ctx, userID)
	require.NoError(t, err)
	require.NotNil(t, foundByUser)
	assert.Equal(t, created.RefreshTokenID, foundByUser.RefreshTokenID)
}

func TestAuthRefreshToken_UpdateAndDelete(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	userID := seedUser(t, ctx, pool)

	expiresAt := time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05")
	_, err := repos.RefreshToken.CreateRefreshToken(ctx, &requests.CreateRefreshToken{
		UserId:    userID,
		Token:     "test-refresh-token-002",
		ExpiresAt: expiresAt,
	})
	require.NoError(t, err)

	newExpiresAt := time.Now().Add(48 * time.Hour).Format("2006-01-02 15:04:05")
	updated, err := repos.RefreshToken.UpdateRefreshToken(ctx, &requests.UpdateRefreshToken{
		UserId:    userID,
		Token:     "updated-refresh-token",
		ExpiresAt: newExpiresAt,
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, "updated-refresh-token", updated.Token)

	// Delete by token
	err = repos.RefreshToken.DeleteRefreshToken(ctx, updated.Token)
	require.NoError(t, err)
}

func TestAuthRefreshToken_DeleteByUserId(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	userID := seedUser(t, ctx, pool)

	expiresAt := time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05")
	_, err := repos.RefreshToken.CreateRefreshToken(ctx, &requests.CreateRefreshToken{
		UserId:    userID,
		Token:     "test-refresh-token-003",
		ExpiresAt: expiresAt,
	})
	require.NoError(t, err)

	err = repos.RefreshToken.DeleteRefreshTokenByUserId(ctx, userID)
	require.NoError(t, err)
}

// ----- ResetToken Repository Tests -----

func TestAuthResetToken_CreateAndFind(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	userID := seedUser(t, ctx, pool)

	expiredAt := time.Now().Add(1 * time.Hour).Format("2006-01-02 15:04:05")
	created, err := repos.ResetToken.CreateResetToken(ctx, &requests.CreateResetTokenRequest{
		UserID:     userID,
		ResetToken: "reset-token-001",
		ExpiredAt:  expiredAt,
	})
	require.NoError(t, err)
	require.NotNil(t, created)
	assert.Equal(t, int64(userID), created.UserID)
	assert.Equal(t, "reset-token-001", created.Token)

	found, err := repos.ResetToken.FindByToken(ctx, created.Token)
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, created.ID, found.ID)
}

func TestAuthResetToken_Delete(t *testing.T) {
	_, repos, pool, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	userID := seedUser(t, ctx, pool)

	expiredAt := time.Now().Add(1 * time.Hour).Format("2006-01-02 15:04:05")
	_, err := repos.ResetToken.CreateResetToken(ctx, &requests.CreateResetTokenRequest{
		UserID:     userID,
		ResetToken: "reset-token-002",
		ExpiredAt:  expiredAt,
	})
	require.NoError(t, err)

	err = repos.ResetToken.DeleteResetToken(ctx, userID)
	require.NoError(t, err)
}
