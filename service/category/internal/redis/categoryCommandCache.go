package mencache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-point-of-sale-shared/cache"
)

type categoryCommandCache struct {
	store *cache.CacheStore
}

func NewCategoryCommandCache(store *cache.CacheStore) CategoryCommandCache {
	return &categoryCommandCache{store: store}
}

func (c *categoryCommandCache) DeleteCachedCategoryCache(ctx context.Context, id int) {
	key := fmt.Sprintf(categoryByIdCacheKey, id)
	cache.DeleteFromCache(ctx, c.store, key)
}
