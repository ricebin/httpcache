package httpcache

import (
	"net/http"
	"time"

	"github.com/bxcodec/gotcha"
	inmemcache "github.com/bxcodec/gotcha/cache"
	"github.com/bxcodec/httpcache/cache"
	"github.com/bxcodec/httpcache/cache/inmem"
	rediscache "github.com/bxcodec/httpcache/cache/redis"
	"github.com/redis/go-redis/v9"
)

// NewWithCustomStorageCache will initiate the httpcache with your defined cache storage
// To use your own cache storage handler, you need to implement the cache.Interactor interface
// And pass it to httpcache.
func NewWithCustomStorageCache(transport http.RoundTripper, rfcCompliance bool,
	cacheInteractor cache.ICacheInteractor) *CacheHandler {
	return newClient(transport, rfcCompliance, cacheInteractor)
}

func newClient(transport http.RoundTripper, rfcCompliance bool,
	cacheInteractor cache.ICacheInteractor) *CacheHandler {
	return NewCacheHandlerRoundtrip(transport, rfcCompliance, cacheInteractor)
}

const (
	MaxSizeCacheItem = 100
)

// NewWithInmemoryCache will create a complete cache-support of HTTP client with using inmemory cache.
// If the duration not set, the cache will use LFU algorithm
func NewWithInmemoryCache(transport http.RoundTripper, rfcCompliance bool, expiryTime time.Duration) *CacheHandler {
	c := gotcha.New(
		gotcha.NewOption().SetAlgorithm(inmemcache.LRUAlgorithm).
			SetExpiryTime(expiryTime).SetMaxSizeItem(MaxSizeCacheItem),
	)

	return newClient(transport, rfcCompliance, inmem.NewCache(c))
}

// NewWithRedisCache will create a complete cache-support of HTTP client with using redis cache.
// If the duration not set, the cache will use LFU algorithm
func NewWithRedisCache(transport http.RoundTripper, rfcCompliance bool, c *redis.Client, expiryTime time.Duration) *CacheHandler {
	return newClient(transport, rfcCompliance, rediscache.NewCache(c, expiryTime))
}
