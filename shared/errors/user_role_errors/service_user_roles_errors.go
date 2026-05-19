package userrole_errors

import (
	"github.com/MamangRust/monolith-point-of-sale-shared/errors"
)


var (
	ErrFailedAssignRoleToUser = errors.ErrInternal.WithMessage("Failed to assign role to user")
	ErrFailedRemoveRole = errors.ErrInternal.WithMessage("Failed to remove role from user")
)
