package user_errors

import "errors"

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrFindAllUsers     = errors.New("failed to find all users")
	ErrFindActiveUsers  = errors.New("failed to find active users")
	ErrFindTrashedUsers = errors.New("failed to find trashed users")
	ErrUserConflict     = errors.New("failed user already exists")

	ErrCreateUser                 = errors.New("failed to create user")
	ErrUpdateUser                 = errors.New("failed to update user")
	ErrUpdateUserVerificationCode = errors.New("failed to update user verification code")
	ErrUpdateUserPassword         = errors.New("failed to update user password")

	ErrTrashedUser         = errors.New("failed to move user to trash")
	ErrRestoreUser         = errors.New("failed to restore user from trash")
	ErrDeleteUserPermanent = errors.New("failed to permanently delete user")

	ErrRestoreAllUsers = errors.New("failed to restore all users")
	ErrDeleteAllUsers  = errors.New("failed to permanently delete all users")
)
