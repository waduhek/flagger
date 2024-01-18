package provider

import (
	"context"

	"github.com/waduhek/flagger/proto/providerpb"
)

// FlagProviderServer is an implementation of the FlagProvider service.
type FlagProviderServer struct {
	providerpb.UnimplementedFlagProviderServer
}

func (s *FlagProviderServer) GetFlag(
	ctx context.Context,
	req *providerpb.GetFlagRequest,
) (*providerpb.GetFlagResponse, error) {
	return &providerpb.GetFlagResponse{Status: true}, nil
}

func NewFlagProviderServer() *FlagProviderServer {
	return &FlagProviderServer{}
}
