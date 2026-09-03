package mencache

import (
	"github.com/MamangRust/monolith-point-of-sale-shared/cache"
)

type Mencache interface {
	MerchantQueryCache
	MerchantCommandCache
	MerchantDocumentQueryCache
	MerchantDocumentCommandCache
}

type mencache struct {
	MerchantQueryCache
	MerchantCommandCache
	MerchantDocumentQueryCache
	MerchantDocumentCommandCache
}

func NewMencache(cacheStore *cache.CacheStore) Mencache {
	return &mencache{
		MerchantQueryCache:           NewMerchantQueryCache(cacheStore),
		MerchantCommandCache:         NewMerchantCommandCache(cacheStore),
		MerchantDocumentQueryCache:   NewMerchantDocumentQueryCache(cacheStore),
		MerchantDocumentCommandCache: NewMerchantDocumentCommandCache(cacheStore),
	}
}
