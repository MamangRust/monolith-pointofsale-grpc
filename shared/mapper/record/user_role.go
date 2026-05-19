package recordmapper

import (
	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"

	"github.com/MamangRust/monolith-point-of-sale-shared/domain/record"
)

type userRoleRecordMapper struct {
}

func NewUserRoleRecordMapper() *userRoleRecordMapper {
	return &userRoleRecordMapper{}
}

func (t *userRoleRecordMapper) ToUserRoleRecord(userRole *db.UserRole) *record.UserRoleRecord {
	return &record.UserRoleRecord{
		UserRoleID: int32(userRole.UserRoleID),
		UserID:     int32(userRole.UserID),
		RoleID:     int32(userRole.RoleID),
		CreatedAt:  userRole.CreatedAt.Time,
		UpdatedAt:  userRole.UpdatedAt.Time,
	}
}
