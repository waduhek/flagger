package provider

import (
	"context"
	"log"

	"github.com/waduhek/flagger/proto/providerpb"

	"github.com/waduhek/flagger/internal/project"
)

// FlagProviderServer is an implementation of the FlagProvider service.
type FlagProviderServer struct {
	providerpb.UnimplementedFlagProviderServer
	providerRepo ProviderRepository
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

	response := &providerpb.GetFlagResponse{
		Status: flagDetail.FlagSetting.IsActive,
	}

	return response, nil
}

func NewFlagProviderServer(providerRepo ProviderRepository) *FlagProviderServer {
	return &FlagProviderServer{
		providerRepo: providerRepo,
	}
}
