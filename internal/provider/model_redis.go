package provider

import "context"

// cacheParameters are the keys used for caching the flag status.
type cacheParameters struct {
	ProjectKey      string
	EnvironmentName string
	FlagName        string
}

// ProviderCacheRepository provides the interface for acessing the cache for
// storing flag statuses.
type ProviderCacheRepository interface {
	// IsStatusCached checks if the flag status is already cached
	IsFlagStatusCached(
		ctx context.Context,
		params *cacheParameters,
	) (bool, error)

	// GetFlagStatus gets the currently cached value of the flag status.
	GetFlagStatus(ctx context.Context, params *cacheParameters,) (bool, error)

	// CacheFlagStatus caches the value of the flag status.
	CacheFlagStatus(
		ctx context.Context,
		params *cacheParameters,
		status bool,
	) error
}
