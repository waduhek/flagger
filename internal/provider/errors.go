package provider

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrFetchFlagDetails is a GRPC error that is returned when an unknown error
// occurs while fetching flag details.
var ErrFetchFlagDetails = status.Error(
	codes.Internal,
	"error occurred while fetching flag details",
)

// ErrIncorrectFlagDetailCount is a GRPC error that is returned when the number
// of flag details does not match the expected count.
var ErrIncorrectFlagDetailCount = status.Error(
	codes.Internal,
	"an unexpected number of flags were found",
)

// ErrStatusCache is a GRPC error that is returned when an error has occurred
// while checking the cache for the flag status.
var ErrStatusCache = status.Error(
	codes.Internal,
	"error occurred while checking flag status cache",
)
