package mencache

import (
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/redis/go-redis/v9"
)

type Mencache struct {
	OrderItemQueryCache OrderItemQueryCache
}

type Deps struct {
	Redis  *redis.Client
	Logger logger.LoggerInterface
}

func NewMencache(deps *Deps) *Mencache {
	cacheStore := NewCacheStore(deps.Redis, deps.Logger)

	return &Mencache{
		OrderItemQueryCache: NewOrderItemQueryCache(cacheStore),
	}
}
