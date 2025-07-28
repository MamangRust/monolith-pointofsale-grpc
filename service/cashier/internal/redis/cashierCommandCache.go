package mencache

import (
	"context"
	"fmt"
)

type cashierCommandCache struct {
	store *CacheStore
}

func NewCashierCommandCache(store *CacheStore) *cashierCommandCache {
	return &cashierCommandCache{store: store}
}

func (c *cashierCommandCache) DeleteCashierCache(ctx context.Context, id int) {
	key := fmt.Sprintf(cashierByIdCacheKey, id)

	DeleteFromCache(ctx, c.store, key)
}
