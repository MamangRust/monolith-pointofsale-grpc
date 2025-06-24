package mencache

import "fmt"

type cashierCommandCache struct {
	store *CacheStore
}

func NewCashierCommandCache(store *CacheStore) *cashierCommandCache {
	return &cashierCommandCache{store: store}
}

func (c *cashierCommandCache) DeleteCashierCache(id int) {
	key := fmt.Sprintf(cashierByIdCacheKey, id)

	DeleteFromCache(c.store, key)
}
