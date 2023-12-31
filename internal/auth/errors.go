package auth

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ENoTokenClaims is a GRPC error that occurs when the claims in the provided
// token were not found.
var ENoTokenClaims = status.Error(
	codes.Internal,
	"could not find token claims",
)
