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
