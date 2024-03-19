package provider

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

// The TTL in seconds of the keys stored in the Redis cache.
var cacheTTL, _ = time.ParseDuration(os.Getenv("FLAGGER_CACHE_TTL"))

type providerCacheRepository struct {
	rdb *redis.Client
}

func (r *providerCacheRepository) IsFlagStatusCached(
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

func (r *providerCacheRepository) GetFlagStatus(
	ctx context.Context,
	params *cacheParameters,
) (bool, error) {
	cacheKey := genFlagStatusCacheKey(params)

	return r.rdb.Get(ctx, cacheKey).Bool()
}

func (r *providerCacheRepository) CacheFlagStatus(
	ctx context.Context,
	params *cacheParameters,
	status bool,
) error {
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

func NewProviderCacheRepository(rdb *redis.Client) *providerCacheRepository {
	return &providerCacheRepository{
		rdb: rdb,
	}
}
