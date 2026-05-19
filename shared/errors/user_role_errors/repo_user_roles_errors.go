package userrole_errors

import "errors"

var (
	ErrAssignRoleToUser = errors.New("failed to assign role to user")
	ErrRemoveRole       = errors.New("failed to remove role from user")
)
