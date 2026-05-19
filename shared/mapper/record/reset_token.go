package recordmapper

import (
	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/record"
)

type resetTokenRecordMapper struct {
}

func NewResetTokenRecordMapper() *resetTokenRecordMapper {
	return &resetTokenRecordMapper{}
}

func (r *resetTokenRecordMapper) ToResetTokenRecord(resetToken *db.ResetToken) *record.ResetTokenRecord {
	return &record.ResetTokenRecord{
		ID:        int64(resetToken.ID),
		UserID:    resetToken.UserID,
		Token:     resetToken.Token,
		ExpiredAt: resetToken.ExpiryDate.String(),
	}
}
