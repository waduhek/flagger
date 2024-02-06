package provider

import (
	"context"
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
	projectKey, ok := project.ProjectKeyFromContext(ctx)
	if !ok {
		log.Printf("could not find project key in request")
		return nil, project.EProjectKeyNotFound
	}

	// Check if the flag status has been cached previously and return it.
	isFlagStatusCached, err := s.checkIfFlagStatusIsCached(ctx, projectKey, req)
	if err != nil {
		return nil, err
	}

	if isFlagStatusCached {
		log.Printf("found cached value for flag status")

		cachedStatus, err := s.getCachedFlagStatus(ctx, projectKey, req)
		if err != nil {
			return nil, err
		}

		response := &providerpb.GetFlagResponse{
			Status: cachedStatus,
		}

		return response, nil
	}

	flagDetails, err := s.providerRepo.GetFlagDetailsByProjectKey(
		ctx,
		projectKey,
		req.Environment,
		req.FlagName,
	)
	if err != nil {
		log.Printf("error while fetching details of the flag: %v", err)
		return nil, EFetchFlagDetails
	}

	if len(flagDetails) != 1 {
		log.Printf("found %d responses of flag details", len(flagDetails))
		return nil, EIncorrectFlagDetailCount
	}

	flagDetail := flagDetails[0]
	status := flagDetail.FlagSetting.IsActive

	// Cache the result of this flag status for the next time.
	cacheParams := cacheParameters{
		ProjectKey: projectKey,
		EnvironmentName: req.Environment,
		FlagName: req.FlagName,
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
	cacheParams := cacheParameters{
		ProjectKey: projectKey,
		EnvironmentName: req.Environment,
		FlagName: req.FlagName,
	}

	statusExists, err := s.cacheRepo.IsFlagStatusCached(ctx, &cacheParams)
	if err != nil {
		log.Printf(
			"error occurred while checking if flag status is cached: %v",
			err,
		)
		return false, EStatusCache
	}

	return statusExists, nil
}

// getCachedFlagStatus gets the value of the cached flag status.
func (s *FlagProviderServer) getCachedFlagStatus(
	ctx context.Context,
	projectKey string,
	req *providerpb.GetFlagRequest,
) (bool, error) {
	cacheParams := cacheParameters{
		ProjectKey: projectKey,
		EnvironmentName: req.Environment,
		FlagName: req.FlagName,
	}

	cachedStatus, err := s.cacheRepo.GetFlagStatus(ctx, &cacheParams)
	if err != nil {
		if err == redis.Nil {
			log.Printf("")
		}
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
