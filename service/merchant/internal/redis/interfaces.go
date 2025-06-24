package mencache

import (
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
)

type MerchantQueryCache interface {
	GetCachedMerchants(req *requests.FindAllMerchants) ([]*response.MerchantResponse, *int, bool)
	SetCachedMerchants(req *requests.FindAllMerchants, data []*response.MerchantResponse, total *int)
	GetCachedMerchantActive(req *requests.FindAllMerchants) ([]*response.MerchantResponseDeleteAt, *int, bool)
	SetCachedMerchantActive(req *requests.FindAllMerchants, data []*response.MerchantResponseDeleteAt, total *int)
	GetCachedMerchantTrashed(req *requests.FindAllMerchants) ([]*response.MerchantResponseDeleteAt, *int, bool)
	SetCachedMerchantTrashed(req *requests.FindAllMerchants, data []*response.MerchantResponseDeleteAt, total *int)
	GetCachedMerchant(id int) (*response.MerchantResponse, bool)
	SetCachedMerchant(data *response.MerchantResponse)
	GetCachedMerchantsByUserId(id int) ([]*response.MerchantResponse, bool)
	SetCachedMerchantsByUserId(userId int, data []*response.MerchantResponse)
}

type MerchantDocumentQueryCache interface {
	GetCachedMerchantDocuments(req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponse, *int, bool)
	SetCachedMerchantDocuments(req *requests.FindAllMerchantDocuments, data []*response.MerchantDocumentResponse, total *int)
	SetCachedMerchantDocumentsTrashed(req *requests.FindAllMerchantDocuments, data []*response.MerchantDocumentResponseDeleteAt, total *int)
	GetCachedMerchantDocumentsActive(req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponseDeleteAt, *int, bool)
	SetCachedMerchantDocumentsActive(req *requests.FindAllMerchantDocuments, data []*response.MerchantDocumentResponseDeleteAt, total *int)
	GetCachedMerchantDocumentsTrashed(req *requests.FindAllMerchantDocuments) ([]*response.MerchantDocumentResponseDeleteAt, *int, bool)
	GetCachedMerchantDocument(id int) (*response.MerchantDocumentResponse, bool)
	SetCachedMerchantDocument(data *response.MerchantDocumentResponse)
}

type MerchantCommandCache interface {
	DeleteCachedMerchant(id int)
}

type MerchantDocumentCommandCache interface {
	DeleteCachedMerchantDocuments(id int)
}
