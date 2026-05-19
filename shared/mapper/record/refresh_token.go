package recordmapper

import (
	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"

	"github.com/MamangRust/monolith-point-of-sale-shared/domain/record"
)

type refreshTokenRecordMapper struct {
}

func NewRefreshTokenRecordMapper() *refreshTokenRecordMapper {
	return &refreshTokenRecordMapper{}
}

func (r *refreshTokenRecordMapper) ToRefreshTokenRecord(refreshToken *db.RefreshToken) *record.RefreshTokenRecord {
	return &record.RefreshTokenRecord{
		ID:        int(refreshToken.RefreshTokenID),
		UserID:    int(refreshToken.UserID),
		Token:     refreshToken.Token,
		ExpiredAt: refreshToken.Expiration.String(),
		CreatedAt: refreshToken.CreatedAt.Time.String(),
		UpdatedAt: refreshToken.UpdatedAt.Time.String(),
	}
}

func (r *refreshTokenRecordMapper) ToRefreshTokensRecord(refreshTokens []*db.RefreshToken) []*record.RefreshTokenRecord {
	var refreshTokenRecords []*record.RefreshTokenRecord
	for _, refreshToken := range refreshTokens {
		refreshTokenRecords = append(refreshTokenRecords, r.ToRefreshTokenRecord(refreshToken))
	}
	return refreshTokenRecords
}
