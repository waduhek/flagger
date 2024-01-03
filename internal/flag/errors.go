package flag

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// EFlagNameTaken is a GRPC error that is returned when attempting to save a new
// flag and the flag name has already been taken.
var EFlagNameTaken = status.Error(
	codes.AlreadyExists,
	"a flag with that name already exists",
)

// EFlagSave is a GRPC error that is returned when an unknown error occurs while
// saving a flag.
var EFlagSave = status.Error(
	codes.Internal,
	"error occurred while saving the flag",
)

// EFlagNotFound is a GRPC error that is returned when the requested flag was
// not found.
var EFlagNotFound = status.Error(codes.NotFound, "flag not found")

// EFlagFetch is a GRPC error that is returned when an unknown error occurs
// while fetching a flag.
var EFlagFetch = status.Error(
	codes.Internal,
	"error occurred while fetching flag",
)
