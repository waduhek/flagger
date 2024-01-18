package project

import (
	"context"
	"strings"

	"google.golang.org/grpc"
)

// ProjectTokenUnaryInterceptor intercepts an incoming request to the provided
// server path and ensures that the request contains the project token in the
// metadata.
func ProjectTokenUnaryInterceptor(serverPath string) grpc.UnaryServerInterceptor {
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
