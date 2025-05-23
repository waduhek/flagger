package auth

import (
	"context"
	"strings"

	"google.golang.org/grpc"

	"github.com/waduhek/flagger/proto/authpb"

	"github.com/waduhek/flagger/internal/logger"
)

// UnaryServerInterceptor intercepts the requests coming to the
// authentication service.
func UnaryServerInterceptor(logger logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if info.FullMethod == authpb.Auth_ChangePassword_FullMethodName {
			newCtx, err := AuthoriseJWT(ctx, logger)
			if err != nil {
				return nil, err
			}

			return handler(newCtx, req)
		}

		return handler(ctx, req)
	}
}

// AuthoriseRequestInterceptor checks if the provided serverPath is the prefix
// in the intercepted method and then validates the JWT in the request metadata.
func AuthoriseRequestInterceptor(
	logger logger.Logger,
	serverPath string,
) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if !strings.HasPrefix(info.FullMethod, serverPath) {
			return handler(ctx, req)
		}

		newCtx, err := AuthoriseJWT(ctx, logger)
		if err != nil {
			return nil, err
		}

		return handler(newCtx, req)
	}
}
