package provider

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// EFetchFlagDetails is a GRPC error that is returned when an unknown error
// occurs while fetching flag details.
var EFetchFlagDetails = status.Error(
	codes.Internal,
	"error occurred while fetching flag details",
)

// EIncorrectFlagDetailCount is a GRPC error that is returned when the number of
// flag details does not match the expected count.
var EIncorrectFlagDetailCount = status.Error(
	codes.Internal,
	"an unexpected number of flags were found",
)

// EStatusCache is a GRPC error that is returned when an error has occurred
// while checking the cache for the flag status.
var EStatusCache = status.Error(
	codes.Internal,
	"error occurred while checking flag status cache",
)
