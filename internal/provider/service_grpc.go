package provider

import (
	"context"
	"errors"
	"log"

	"github.com/redis/go-redis/v9"

	"github.com/waduhek/flagger/proto/providerpb"

	"github.com/waduhek/flagger/internal/project"
)

// FlagProviderServer is an implementation of the FlagProvider service.
type FlagProviderServer struct {
	providerpb.UnimplementedFlagProviderServer
	providerRepo ProviderRepository
	cacheRepo    ProviderCacheRepository
}

func (s *FlagProviderServer) GetFlag(
	ctx context.Context,
	req *providerpb.GetFlagRequest,
) (*providerpb.GetFlagResponse, error) {
	projectKey, ok := project.KeyFromContext(ctx)
	if !ok {
		log.Printf("could not find project key in request")
		return nil, project.ErrProjectKeyNotFound
	}

	// Check if the flag status has been cached previously and return it.
	isFlagStatusCached, isCachedErr :=
		s.checkIfFlagStatusIsCached(ctx, projectKey, req)
	if isCachedErr != nil {
		return nil, isCachedErr
	}

	if isFlagStatusCached {
		log.Printf("found cached value for flag status")

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

	flagDetails, err := s.providerRepo.GetFlagDetailsByProjectKey(
		ctx,
		projectKey,
		environmentName,
		flagName,
	)
	if err != nil {
		log.Printf("error while fetching details of the flag: %v", err)
		return nil, ErrFetchFlagDetails
	}

	if len(flagDetails) != 1 {
		log.Printf("found %d responses of flag details", len(flagDetails))
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

	cacheErr := s.cacheRepo.CacheFlagStatus(ctx, &cacheParams, status)
	if cacheErr != nil {
		log.Printf("could not cache flag status: %v. ignoring error", cacheErr)
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

	statusExists, err := s.cacheRepo.IsFlagStatusCached(ctx, &cacheParams)
	if err != nil {
		log.Printf(
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

	cachedStatus, err := s.cacheRepo.GetFlagStatus(ctx, &cacheParams)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			log.Print("Could not find cached flag status")
			return false, nil
		}

		return false, err
	}

	return cachedStatus, nil
}

func NewFlagProviderServer(
	providerRepo ProviderRepository,
	cacheRepo ProviderCacheRepository,
) *FlagProviderServer {
	return &FlagProviderServer{
		providerRepo: providerRepo,
		cacheRepo:    cacheRepo,
	}
}
