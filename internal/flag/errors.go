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

// ENoEnvironments is a GRPC error that is returned when no environments have
// been configured and the user tries to create a new flag.
var ENoEnvironments = status.Error(
	codes.FailedPrecondition,
	"no environments have been configured for the project",
)

// EFlagSaveTxn is a GRPC error that is returned when a new session for starting
// a transaction could not be created.
var EFlagSaveTxn = status.Error(
	codes.Internal,
	"could start the transaction to save flag",
)

// EUpdateFlagStatus is a GRPC error that is returned when the status of a flag
// could not be updated.
var EUpdateFlagStatus = status.Error(
	codes.Internal,
	"no flag settings could be updated",
)
