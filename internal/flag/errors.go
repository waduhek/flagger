package flag

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrNameTaken is a GRPC error that is returned when attempting to save a new
// flag and the flag name has already been taken.
var ErrNameTaken = status.Error(
	codes.AlreadyExists,
	"a flag with that name already exists",
)

// ErrCouldNotSave is a GRPC error that is returned when an unknown error occurs
// while saving a flag.
var ErrCouldNotSave = status.Error(
	codes.Internal,
	"error occurred while saving the flag",
)

// ErrNotFound is a GRPC error that is returned when the requested flag was not
// found.
var ErrNotFound = status.Error(codes.NotFound, "flag not found")

// ErrCouldNotFetch is a GRPC error that is returned when an unknown error
// occurs while fetching a flag.
var ErrCouldNotFetch = status.Error(
	codes.Internal,
	"error occurred while fetching flag",
)

// ErrNoEnvironments is a GRPC error that is returned when no environments have
// been configured and the user tries to create a new flag.
var ErrNoEnvironments = status.Error(
	codes.FailedPrecondition,
	"no environments have been configured for the project",
)

// ErrTxnSession is a GRPC error that is returned when a new session for
// starting a transaction could not be created.
var ErrTxnSession = status.Error(
	codes.Internal,
	"could start the transaction to save flag",
)

// ErrUpdateStatus is a GRPC error that is returned when the status of a flag
// could not be updated.
var ErrUpdateStatus = status.Error(
	codes.Internal,
	"no flag settings could be updated",
)
