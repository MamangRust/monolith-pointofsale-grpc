package recordmapper

import (
	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"

	"github.com/MamangRust/monolith-point-of-sale-shared/domain/record"
)

type roleRecordMapper struct {
}

func NewRoleRecordMapper() *roleRecordMapper {
	return &roleRecordMapper{}
}

func (s *roleRecordMapper) ToRoleRecord(role *db.Role) *record.RoleRecord {
	deletedAt := role.DeletedAt.Time.Format("2006-01-02 15:04:05.000")

	return &record.RoleRecord{
		ID:        int(role.RoleID),
		Name:      role.RoleName,
		CreatedAt: role.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: role.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt: &deletedAt,
	}
}

func (s *roleRecordMapper) ToRolesRecord(roles []*db.Role) []*record.RoleRecord {
	var result []*record.RoleRecord

	for _, role := range roles {
		result = append(result, s.ToRoleRecord(role))
	}

	return result
}

func (s *roleRecordMapper) ToRoleRecordAll(role *db.GetRolesRow) *record.RoleRecord {
	deletedAt := role.DeletedAt.Time.Format("2006-01-02 15:04:05.000")

	return &record.RoleRecord{
		ID:        int(role.RoleID),
		Name:      role.RoleName,
		CreatedAt: role.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: role.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt: &deletedAt,
	}
}

func (s *roleRecordMapper) ToRolesRecordAll(roles []*db.GetRolesRow) []*record.RoleRecord {
	var result []*record.RoleRecord

	for _, role := range roles {
		result = append(result, s.ToRoleRecordAll(role))
	}

	return result
}

func (s *roleRecordMapper) ToRoleRecordActive(role *db.GetActiveRolesRow) *record.RoleRecord {
	deletedAt := role.DeletedAt.Time.Format("2006-01-02 15:04:05.000")

	return &record.RoleRecord{
		ID:        int(role.RoleID),
		Name:      role.RoleName,
		CreatedAt: role.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: role.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt: &deletedAt,
	}
}

func (s *roleRecordMapper) ToRolesRecordActive(roles []*db.GetActiveRolesRow) []*record.RoleRecord {
	var result []*record.RoleRecord

	for _, role := range roles {
		result = append(result, s.ToRoleRecordActive(role))
	}

	return result
}

func (s *roleRecordMapper) ToRoleRecordTrashed(role *db.GetTrashedRolesRow) *record.RoleRecord {
	deletedAt := role.DeletedAt.Time.Format("2006-01-02 15:04:05.000")

	return &record.RoleRecord{
		ID:        int(role.RoleID),
		Name:      role.RoleName,
		CreatedAt: role.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt: role.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt: &deletedAt,
	}
}

func (s *roleRecordMapper) ToRolesRecordTrashed(roles []*db.GetTrashedRolesRow) []*record.RoleRecord {
	var result []*record.RoleRecord

	for _, role := range roles {
		result = append(result, s.ToRoleRecordTrashed(role))
	}

	return result
}
