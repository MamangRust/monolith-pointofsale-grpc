package refreshtoken_errors

import "errors"

var (
	ErrTokenNotFound      = errors.New("refresh token not found")
	ErrFindByToken        = errors.New("failed to find refresh token by token")
	ErrFindByUserID       = errors.New("failed to find refresh token by user ID")
	ErrCreateRefreshToken = errors.New("failed to create refresh token")
	ErrUpdateRefreshToken = errors.New("failed to update refresh token")
	ErrDeleteRefreshToken = errors.New("failed to delete refresh token")
	ErrDeleteByUserID     = errors.New("failed to delete refresh token by user ID")
	ErrParseDate          = errors.New("failed to parse expiration date")
)
