package hash

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrGenPasswordHash is a GRPC error that is returned when the hash of a password
// could not be generated.
var ErrGenPasswordHash = status.Error(
	codes.Internal,
	"could not hash the password",
)
