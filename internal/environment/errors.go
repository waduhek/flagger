package environment

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrNameTaken is a GRPC error that is returned when creating a new environment
// and the provided environment name is already taken.
var ErrNameTaken = status.Error(
	codes.AlreadyExists,
	"an environment with that name already exists",
)

// ErrCouldNotSave is a GRPC error that is returned when an unknown error occurs
// while creating a new environment.
var ErrCouldNotSave = status.Error(
	codes.Internal,
	"could not create a new environment",
)

// ErrNotFound is a GRPC error that is returned when no environments are found
// for the provided query.
var ErrNotFound = status.Error(codes.NotFound, "environment not found")

// ErrCouldNotFetch is a GRPC error that is returned when unknown error occurs
// while fetching an environment.
var ErrCouldNotFetch = status.Error(
	codes.Internal,
	"error occurred while fetching environment",
)

// ErrTxnSession is a GRPC error that is returned when a new session for
// starting a transaction could not be created.
var ErrTxnSession = status.Error(
	codes.Internal,
	"could not create a transaction session",
)

// ErrEnvironmentIDCast is a GRPC error that is returned when the ID of the
// environment could not be cast.
var ErrEnvironmentIDCast = status.Error(
	codes.Internal,
	"could not create a new environment",
)
