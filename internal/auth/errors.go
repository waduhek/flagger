package auth

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrNoTokenClaims is a GRPC error that is returned when the claims in the
// provided token were not found.
var ErrNoTokenClaims = status.Error(
	codes.Internal,
	"could not find token claims",
)

// ErrIncorrectUsernameOrPassword is a GRPC error that is returned when the
// provided usernam or password is incorrect.
var ErrIncorrectUsernameOrPassword = status.Error(
	codes.Unauthenticated,
	"incorrect username or password",
)

// ErrJWTSign is a GRPC error that is returned when an error occurs while signing
// a JWT.
var ErrJWTSign = status.Error(codes.Internal, "error while signing JWT")

// ErrInvalidJWT is a GRPC error that is returned when the provided JWT doesn't
// complete validation.
var ErrInvalidJWT = status.Error(codes.Unauthenticated, "invalid jwt")

// ErrMetadataNotFound is a GRPC error that is returned when the request
// metadata could not be found.
var ErrMetadataNotFound = status.Error(
	codes.InvalidArgument,
	"could not find incoming request metadata",
)

// ErrAuthMetadataNotFound is a GRPC error that is returned when "authorization"
// metadata was not found.
var ErrAuthMetadataNotFound = status.Error(
	codes.Unauthenticated,
	"authorization metdata not found",
)

// ErrAuthMetadataLength is a GRPC error that is returned when the length of
// "authorization" metadata is of incorrect length.
var ErrAuthMetadataLength = status.Error(
	codes.InvalidArgument,
	"invalid authorization metadata value",
)

// ErrInvalidTokenFormat is a GRPC error that is returned when the format of the
// provided token does match the expected format.
var ErrInvalidTokenFormat = status.Error(
	codes.InvalidArgument,
	"invalid bearer header format",
)
