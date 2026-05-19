package recordmapper

import (
	db "github.com/MamangRust/monolith-point-of-sale-pkg/database/schema"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/record"
)

type merchantDocumentRecordMapper struct {
}

func NewMerchantDocumentRecordMapper() *merchantDocumentRecordMapper {
	return &merchantDocumentRecordMapper{}
}

func (m *merchantDocumentRecordMapper) ToGetMerchantDocument(doc *db.MerchantDocument) *record.MerchantDocumentRecord {
	return &record.MerchantDocumentRecord{
		ID:           int(doc.DocumentID),
		MerchantID:   int(doc.MerchantID),
		DocumentType: doc.DocumentType,
		DocumentURL:  doc.DocumentUrl,
		Status:       doc.Status,
		Note:         doc.Note.String,
		CreatedAt:    doc.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt:    doc.UpdatedAt.Time.Format("2006-01-02"),
	}
}

func (m *merchantDocumentRecordMapper) ToMerchantDocumentRecord(doc *db.GetMerchantDocumentsRow) *record.MerchantDocumentRecord {
	var deletedAt *string

	if doc.DeletedAt.Valid {
		formattedDeletedAt := doc.DeletedAt.Time.Format("2006-01-02")
		deletedAt = &formattedDeletedAt
	}

	return &record.MerchantDocumentRecord{
		ID:           int(doc.DocumentID),
		MerchantID:   int(doc.MerchantID),
		DocumentType: doc.DocumentType,
		DocumentURL:  doc.DocumentUrl,
		Status:       doc.Status,
		Note:         doc.Note.String,
		CreatedAt:    doc.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt:    doc.UpdatedAt.Time.Format("2006-01-02"),
		DeletedAt:    deletedAt,
	}
}

func (m *merchantDocumentRecordMapper) ToMerchantDocumentsRecord(docs []*db.GetMerchantDocumentsRow) []*record.MerchantDocumentRecord {
	var records []*record.MerchantDocumentRecord
	for _, doc := range docs {
		records = append(records, m.ToMerchantDocumentRecord(doc))
	}
	return records
}

func (m *merchantDocumentRecordMapper) ToMerchantDocumentActiveRecord(doc *db.GetActiveMerchantDocumentsRow) *record.MerchantDocumentRecord {
	var deletedAt *string
	if doc.DeletedAt.Valid {
		formattedDeletedAt := doc.DeletedAt.Time.Format("2006-01-02")
		deletedAt = &formattedDeletedAt
	}
	return &record.MerchantDocumentRecord{
		ID:           int(doc.DocumentID),
		MerchantID:   int(doc.MerchantID),
		DocumentType: doc.DocumentType,
		DocumentURL:  doc.DocumentUrl,
		Status:       doc.Status,
		Note:         doc.Note.String,
		CreatedAt:    doc.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt:    doc.UpdatedAt.Time.Format("2006-01-02"),
		DeletedAt:    deletedAt,
	}
}

func (m *merchantDocumentRecordMapper) ToMerchantDocumentsActiveRecord(docs []*db.GetActiveMerchantDocumentsRow) []*record.MerchantDocumentRecord {
	var records []*record.MerchantDocumentRecord
	for _, doc := range docs {
		records = append(records, m.ToMerchantDocumentActiveRecord(doc))
	}
	return records
}

func (m *merchantDocumentRecordMapper) ToMerchantDocumentTrashedRecord(doc *db.GetTrashedMerchantDocumentsRow) *record.MerchantDocumentRecord {
	var deletedAt *string
	if doc.DeletedAt.Valid {
		formattedDeletedAt := doc.DeletedAt.Time.Format("2006-01-02")
		deletedAt = &formattedDeletedAt
	}

	return &record.MerchantDocumentRecord{
		ID:           int(doc.DocumentID),
		MerchantID:   int(doc.MerchantID),
		DocumentType: doc.DocumentType,
		DocumentURL:  doc.DocumentUrl,
		Status:       doc.Status,
		Note:         doc.Note.String,
		CreatedAt:    doc.CreatedAt.Time.Format("2006-01-02"),
		UpdatedAt:    doc.UpdatedAt.Time.Format("2006-01-02"),
		DeletedAt:    deletedAt,
	}
}

func (m *merchantDocumentRecordMapper) ToMerchantDocumentsTrashedRecord(docs []*db.GetTrashedMerchantDocumentsRow) []*record.MerchantDocumentRecord {
	var records []*record.MerchantDocumentRecord
	for _, doc := range docs {
		records = append(records, m.ToMerchantDocumentTrashedRecord(doc))
	}
	return records
}
