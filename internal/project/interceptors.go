package project

import (
	"context"
	"strings"

	"google.golang.org/grpc"
)

// ProjectKeyUnaryInterceptor intercepts an incoming request to the provided
// server path and ensures that the request contains the project key in the
// metadata.
func ProjectKeyUnaryInterceptor(serverPath string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if !strings.HasPrefix(info.FullMethod, serverPath) {
			return handler(ctx, req)
		}

		newCtx, err := AuthoriseProject(ctx)
		if err != nil {
			return nil, err
		}

		return handler(newCtx, req)
	}
}
