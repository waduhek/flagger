package auth

import (
	"context"
	"regexp"
	"strings"

	"google.golang.org/grpc/metadata"

	"github.com/waduhek/flagger/internal/logger"
)

// AuthoriseJWT takes a GRPC context and validates that the current request has
// the "authorization" header and is a valid JWT. If the token is present and
// valid, adds the claims to the context and returns a new context for the
// handler. All errors from this middleware will be GRPC compliant.
func AuthoriseJWT(ctx context.Context, logger logger.Logger) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Error("could not find incoming request metadata")
		return nil, ErrMetadataNotFound
	}

	authHeader, ok := md["authorization"]
	if !ok {
		logger.Error("could not find authorization header")
		return nil, ErrAuthMetadataNotFound
	}

	if len(authHeader) != 1 {
		logger.Error(
			"authorization header was found to be of len %d which is not expected",
			len(authHeader),
		)
		return nil, ErrAuthMetadataLength
	}

	bearerToken := authHeader[0]

	claims, err := validateJWT(logger, bearerToken)
	if err != nil {
		return nil, err
	}

	claimCtx := InjectClaimsIntoContext(ctx, claims)

	return claimCtx, nil
}

// validateJWT accepts the token header value of the "authorization" header and
// validates it. If the token is valid, returns the claims from the body of the
// token. If an error occurs, will always return a GRPC compliant error.
func validateJWT(logger logger.Logger, token string) (*FlaggerJWTClaims, error) {
	bearerTokenRegEx := regexp.MustCompile(
		`^[b|B]earer [a-zA-Z0-9-_]+\.[a-zA-Z0-9-_]+\.[a-zA-Z0-9-_]+$`,
	)
	if !bearerTokenRegEx.MatchString(token) {
		logger.Error("token header format does not match")
		return nil, ErrInvalidTokenFormat
	}

	headerJWT := strings.Split(token, " ")[1]

	return VerifyJWT(logger, headerJWT)
}
