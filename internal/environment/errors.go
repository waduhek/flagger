package environment

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// EEnvironmentNameTaken is a GRPC error that is returned when creating a new
// environment and the provided environment name is already taken.
var EEnvironmentNameTaken = status.Error(
	codes.AlreadyExists,
	"an environment with that name already exists",
)

// EEnvironmentSave is a GRPC error that is returned when an unknown error
// occurs while creating a new environment.
var EEnvironmentSave = status.Error(
	codes.Internal,
	"could not create a new environment",
)

// EEnvironmentNotFound is a GRPC error that is returned when no environments
// are found for the provided query.
var EEnvironmentNotFound = status.Error(codes.NotFound, "environment not found")

// EEnvironmentFetch is a GRPC error that is returned when unknown error occurs
// while fetching an environment.
var EEnvironmentFetch = status.Error(
	codes.Internal,
	"error occurred while fetching environment",
)

// EEnvironmentTxn is a GRPC error that is returned when a new session for
// starting a transaction could not be created.
var EEnvironmentTxn = status.Error(
	codes.Internal,
	"could not create a transaction session",
)

// EEnvironmentIDCast is a GRPC error that is returned when the ID of the
// environment could not be cast.
var EEnvironmentIDCast = status.Error(
	codes.Internal,
	"could not create a new environment",
)
