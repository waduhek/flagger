package interceptors

import (
	"context"
	"strings"

	"google.golang.org/grpc"

	"github.com/waduhek/flagger/internal/middleware"
)

// ProjectServerUnaryInterceptor intercepts the requests coming to the project
// service.
func ProjectServerUnaryInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	if !strings.HasPrefix(info.FullMethod, "/projectpb.Project/") {
		return handler(ctx, req)
	}

	newCtx, err := middleware.AuthoriseJWT(ctx)
	if err != nil {
		return nil, err
	}

	return handler(newCtx, req)
}
