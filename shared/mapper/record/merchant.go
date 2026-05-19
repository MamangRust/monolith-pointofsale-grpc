package recordmapper

import (
	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"

	"github.com/MamangRust/monolith-point-of-sale-shared/domain/record"
)

type merchantRecordMapper struct {
}

func NewMerchantRecordMapper() *merchantRecordMapper {
	return &merchantRecordMapper{}
}

func (s *merchantRecordMapper) ToMerchantRecord(Merchant *db.Merchant) *record.MerchantRecord {
	var deletedAt *string
	if Merchant.DeletedAt.Valid {
		deletedAtStr := Merchant.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.MerchantRecord{
		ID:           int(Merchant.MerchantID),
		UserID:       int(Merchant.UserID),
		Name:         Merchant.Name,
		Description:  Merchant.Description.String,
		Address:      Merchant.Address.String,
		ContactEmail: Merchant.ContactEmail.String,
		ContactPhone: Merchant.ContactPhone.String,
		Status:       Merchant.Status,
		CreatedAt:    Merchant.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:    Merchant.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:    deletedAt,
	}
}

func (s *merchantRecordMapper) ToMerchantRecordPagination(Merchant *db.GetMerchantsRow) *record.MerchantRecord {
	var deletedAt *string
	if Merchant.DeletedAt.Valid {
		deletedAtStr := Merchant.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.MerchantRecord{
		ID:           int(Merchant.MerchantID),
		UserID:       int(Merchant.UserID),
		Name:         Merchant.Name,
		Description:  Merchant.Description.String,
		Address:      Merchant.Address.String,
		ContactEmail: Merchant.ContactEmail.String,
		ContactPhone: Merchant.ContactPhone.String,
		Status:       Merchant.Status,
		CreatedAt:    Merchant.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:    Merchant.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:    deletedAt,
	}
}

func (s *merchantRecordMapper) ToMerchantsRecordPagination(Merchants []*db.GetMerchantsRow) []*record.MerchantRecord {
	var result []*record.MerchantRecord

	for _, Merchant := range Merchants {
		result = append(result, s.ToMerchantRecordPagination(Merchant))
	}

	return result
}

func (s *merchantRecordMapper) ToMerchantRecordActivePagination(Merchant *db.GetMerchantsActiveRow) *record.MerchantRecord {
	var deletedAt *string
	if Merchant.DeletedAt.Valid {
		deletedAtStr := Merchant.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.MerchantRecord{
		ID:           int(Merchant.MerchantID),
		UserID:       int(Merchant.UserID),
		Name:         Merchant.Name,
		Description:  Merchant.Description.String,
		Address:      Merchant.Address.String,
		ContactEmail: Merchant.ContactEmail.String,
		ContactPhone: Merchant.ContactPhone.String,
		Status:       Merchant.Status,
		CreatedAt:    Merchant.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:    Merchant.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:    deletedAt,
	}
}

func (s *merchantRecordMapper) ToMerchantsRecordActivePagination(Merchants []*db.GetMerchantsActiveRow) []*record.MerchantRecord {
	var result []*record.MerchantRecord

	for _, Merchant := range Merchants {
		result = append(result, s.ToMerchantRecordActivePagination(Merchant))
	}

	return result
}

func (s *merchantRecordMapper) ToMerchantRecordTrashedPagination(Merchant *db.GetMerchantsTrashedRow) *record.MerchantRecord {
	var deletedAt *string
	if Merchant.DeletedAt.Valid {
		deletedAtStr := Merchant.DeletedAt.Time.Format("2006-01-02 15:04:05.000")
		deletedAt = &deletedAtStr
	}

	return &record.MerchantRecord{
		ID:           int(Merchant.MerchantID),
		UserID:       int(Merchant.UserID),
		Name:         Merchant.Name,
		Description:  Merchant.Description.String,
		Address:      Merchant.Address.String,
		ContactEmail: Merchant.ContactEmail.String,
		ContactPhone: Merchant.ContactPhone.String,
		Status:       Merchant.Status,
		CreatedAt:    Merchant.CreatedAt.Time.Format("2006-01-02 15:04:05.000"),
		UpdatedAt:    Merchant.UpdatedAt.Time.Format("2006-01-02 15:04:05.000"),
		DeletedAt:    deletedAt,
	}
}

func (s *merchantRecordMapper) ToMerchantsRecordTrashedPagination(Merchants []*db.GetMerchantsTrashedRow) []*record.MerchantRecord {
	var result []*record.MerchantRecord

	for _, Merchant := range Merchants {
		result = append(result, s.ToMerchantRecordTrashedPagination(Merchant))
	}

	return result
}
