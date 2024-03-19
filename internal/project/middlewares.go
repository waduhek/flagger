package project

import (
	"context"
	"log"

	"google.golang.org/grpc/metadata"
)

// projectTokenMetadataKey is the metadata key for the incoming request
// metadata to find the project token.
//
//nolint:gosec // This isn't a secret but a header key.
const projectTokenMetadataKey = "x-flagger-token"

// AuthoriseProject takes an incoming GRPC context and checks if the project
// token is present in the incoming context. If the project token exists, it
// adds the token to the returned context. All errors returned by this function
// are GRPC compliant.
func AuthoriseProject(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("could not find message metadata")
		return nil, EMetadataNotFound
	}

	projectTokens, ok := md[projectTokenMetadataKey]
	if !ok {
		log.Println("could not find the project token")
		return nil, EProjectKeyNotFound
	}

	if len(projectTokens) != 1 {
		log.Println("multiple project tokens found in metadata")
		return nil, EKeyMetadataLength
	}

	projectToken := projectTokens[0]

	return injectProjectKeyIntoContext(ctx, projectToken), nil
}
