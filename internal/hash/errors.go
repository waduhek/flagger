package hash

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// EHashGenPasswordHash is a GRPC error that is returned when the hash of a password
// could not be generated.
var EHashGenPasswordHash = status.Error(
	codes.Internal,
	"could not hash the password",
)
