package role_errors

import "errors"

var (
	ErrRoleNotFound     = errors.New("role not found")
	ErrFindAllRoles     = errors.New("failed to find all Roles")
	ErrFindActiveRoles  = errors.New("failed to find active Roles")
	ErrFindTrashedRoles = errors.New("failed to find trashed Roles")
	ErrRoleConflict     = errors.New("failed Role already exists")

	ErrCreateRole = errors.New("failed to create Role")
	ErrUpdateRole = errors.New("failed to update Role")

	ErrTrashedRole         = errors.New("failed to move Role to trash")
	ErrRestoreRole         = errors.New("failed to restore Role from trash")
	ErrDeleteRolePermanent = errors.New("failed to permanently delete Role")

	ErrRestoreAllRoles = errors.New("failed to restore all Roles")
	ErrDeleteAllRoles  = errors.New("failed to permanently delete all Roles")
)
