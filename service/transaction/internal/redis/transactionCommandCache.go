package mencache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-point-of-sale-shared/cache"
)

type transactionCommandCache struct {
	store *cache.CacheStore
}

func NewTransactionCommandCache(store *cache.CacheStore) TransactionCommandCache {
	return &transactionCommandCache{store: store}
}

func (t *transactionCommandCache) DeleteTransactionCache(ctx context.Context, transactionID int) {
	key := fmt.Sprintf(transactionByIdCacheKey, transactionID)

	cache.DeleteFromCache(ctx, t.store, key)
}
