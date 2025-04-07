package provider

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"

	"github.com/waduhek/flagger/proto/providerpb"

	"github.com/waduhek/flagger/internal/logger"
	"github.com/waduhek/flagger/internal/project"
)

// FlagProviderServer is an implementation of the FlagProvider service.
type FlagProviderServer struct {
	providerpb.UnimplementedFlagProviderServer
	providerDataRepo  DataRepository
	providerCacheRepo CacheRepository
	logger            logger.Logger
}

func (s *FlagProviderServer) GetFlag(
	ctx context.Context,
	req *providerpb.GetFlagRequest,
) (*providerpb.GetFlagResponse, error) {
	projectKey, ok := project.KeyFromContext(ctx)
	if !ok {
		s.logger.Error("could not find project key in request")
		return nil, project.ErrProjectKeyNotFound
	}

	// Check if the flag status has been cached previously and return it.
	isFlagStatusCached, isCachedErr :=
		s.checkIfFlagStatusIsCached(ctx, projectKey, req)
	if isCachedErr != nil {
		return nil, isCachedErr
	}

	if isFlagStatusCached {
		s.logger.Info("found cached value for flag status")

		cachedStatus, cachedStatusErr := s.getCachedFlagStatus(ctx, projectKey, req)
		if cachedStatusErr != nil {
			return nil, cachedStatusErr
		}

		response := &providerpb.GetFlagResponse{
			Status: cachedStatus,
		}

		return response, nil
	}

	environmentName := req.GetEnvironment()
	flagName := req.GetFlagName()

	flagDetails, err := s.providerDataRepo.GetFlagDetailsByProjectKey(
		ctx,
		projectKey,
		environmentName,
		flagName,
	)
	if err != nil {
		s.logger.Error("error while fetching details of the flag: %v", err)
		return nil, ErrFetchFlagDetails
	}

	if len(flagDetails) != 1 {
		s.logger.Error("found %d responses of flag details", len(flagDetails))
		return nil, ErrIncorrectFlagDetailCount
	}

	flagDetail := flagDetails[0]
	status := flagDetail.FlagSetting.IsActive

	// Cache the result of this flag status for the next time.
	cacheParams := cacheParameters{
		ProjectKey:      projectKey,
		EnvironmentName: environmentName,
		FlagName:        flagName,
	}

	cacheErr := s.providerCacheRepo.CacheFlagStatus(ctx, &cacheParams, status)
	if cacheErr != nil {
		s.logger.Warn("could not cache flag status: %v. ignoring error", cacheErr)
	}

	response := &providerpb.GetFlagResponse{
		Status: status,
	}

	return response, nil
}

// checkIfFlagStatusIsCached checks if the requested flag's status has already
// been cached.
func (s *FlagProviderServer) checkIfFlagStatusIsCached(
	ctx context.Context,
	projectKey string,
	req *providerpb.GetFlagRequest,
) (bool, error) {
	environmentName := req.GetEnvironment()
	flagName := req.GetFlagName()

	cacheParams := cacheParameters{
		ProjectKey:      projectKey,
		EnvironmentName: environmentName,
		FlagName:        flagName,
	}

	statusExists, err := s.providerCacheRepo.IsFlagStatusCached(ctx, &cacheParams)
	if err != nil {
		s.logger.Error(
			"error occurred while checking if flag status is cached: %v",
			err,
		)
		return false, ErrStatusCache
	}

	return statusExists, nil
}

// getCachedFlagStatus gets the value of the cached flag status.
func (s *FlagProviderServer) getCachedFlagStatus(
	ctx context.Context,
	projectKey string,
	req *providerpb.GetFlagRequest,
) (bool, error) {
	environmentName := req.GetEnvironment()
	flagName := req.GetFlagName()

	cacheParams := cacheParameters{
		ProjectKey:      projectKey,
		EnvironmentName: environmentName,
		FlagName:        flagName,
	}

	cachedStatus, err := s.providerCacheRepo.GetFlagStatus(ctx, &cacheParams)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			s.logger.Error("Could not find cached flag status")
			return false, nil
		}

		return false, err
	}

	return cachedStatus, nil
}

func NewFlagProviderServer(
	providerDataRepo DataRepository,
	providerCacheRepo CacheRepository,
	logger logger.Logger,
) *FlagProviderServer {
	return &FlagProviderServer{
		providerDataRepo:  providerDataRepo,
		providerCacheRepo: providerCacheRepo,
		logger:            logger,
	}
}
