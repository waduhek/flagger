package interceptors

import (
	"context"

	"google.golang.org/grpc"

	"github.com/waduhek/flagger/proto/authpb"

	"github.com/waduhek/flagger/internal/middleware"
)

// AuthServerUnaryInterceptor intercepts the requests coming to the
// authentication service.
func AuthServerUnaryInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	if info.FullMethod == authpb.Auth_ChangePassword_FullMethodName {
		newCtx, err := middleware.AuthoriseJWT(ctx)
		if err != nil {
			return nil, err
		}

		return handler(newCtx, req)
	}

	return handler(ctx, req)
}
