package provider

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCacheRepository struct {
	rdb *redis.Client
}

func (r *RedisCacheRepository) IsFlagStatusCached(
	ctx context.Context,
	params *cacheParameters,
) (bool, error) {
	cacheKey := genFlagStatusCacheKey(params)

	result, err := r.rdb.Exists(ctx, cacheKey).Result()
	if err != nil {
		return false, err
	}

	return result >= 1, nil
}

func (r *RedisCacheRepository) GetFlagStatus(
	ctx context.Context,
	params *cacheParameters,
) (bool, error) {
	cacheKey := genFlagStatusCacheKey(params)

	return r.rdb.Get(ctx, cacheKey).Bool()
}

func (r *RedisCacheRepository) CacheFlagStatus(
	ctx context.Context,
	params *cacheParameters,
	status bool,
) error {
	// The TTL in seconds of the keys stored in the Redis cache.
	cacheTTL, _ := time.ParseDuration(os.Getenv("FLAGGER_CACHE_TTL"))

	cacheKey := genFlagStatusCacheKey(params)

	return r.rdb.SetEx(
		ctx,
		cacheKey,
		status,
		cacheTTL,
	).Err()
}

// getFlagStatusCacheKey generates the cache key used for caching the flag
// status.
func genFlagStatusCacheKey(params *cacheParameters) string {
	return fmt.Sprintf(
		"status:%v:%v:%v",
		params.ProjectKey,
		params.EnvironmentName,
		params.FlagName,
	)
}

func NewProviderCacheRepository(rdb *redis.Client) *RedisCacheRepository {
	return &RedisCacheRepository{
		rdb: rdb,
	}
}
