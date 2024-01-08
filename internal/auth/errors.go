package auth

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ENoTokenClaims is a GRPC error that is returned when the claims in the
// provided token were not found.
var ENoTokenClaims = status.Error(
	codes.Internal,
	"could not find token claims",
)

// EIncorrectUsernameOrPassword is a GRPC error that is returned when the
// provided usernam or password is incorrect.
var EIncorrectUsernameOrPassword = status.Error(
	codes.Unauthenticated,
	"incorrect username or password",
)

// EJWTSign is a GRPC error that is returned when an error occurs while signing
// a JWT.
var EJWTSign = status.Error(codes.Internal, "error while signing JWT")

// EInvalidJWT is a GRPC error that is returned when the provided JWT doesn't
// complete validation.
var EInvalidJWT = status.Error(codes.Unauthenticated, "invalid jwt")

// EMetadataNotFound is a GRPC error that is returned when the request metadata
// could not be found.
var EMetadataNotFound = status.Error(
	codes.InvalidArgument,
	"could not find incoming request metadata",
)

// EAuthMetadataNotFound is a GRPC error that is returned when "authorization"
// metadata was not found.
var EAuthMetadataNotFound = status.Error(
	codes.Unauthenticated,
	"authorization metdata not found",
)

// EAuthMetadataLength is a GRPC error that is returned when the length of
// "authorization" metadata is of incorrect length.
var EAuthMetadataLength = status.Error(
	codes.InvalidArgument,
	"invalid authorization metadata value",
)

// EInvalidTokenFormat is a GRPC error that is returned when the format of the
// provided token does match the expected format.
var EInvalidTokenFormat = status.Error(
	codes.InvalidArgument,
	"invalid bearer header format",
)
