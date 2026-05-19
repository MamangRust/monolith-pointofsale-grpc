package mencache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-point-of-sale-shared/cache"
)

type roleCommandCache struct {
	store *cache.CacheStore
}

func NewRoleCommandCache(store *cache.CacheStore) RoleCommandCache {
	return &roleCommandCache{store: store}
}

func (s *roleCommandCache) DeleteCachedRole(ctx context.Context, id int) {
	key := fmt.Sprintf(roleByIdCacheKey, id)

	cache.DeleteFromCache(ctx, s.store, key)
}
