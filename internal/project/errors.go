package project

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrNameTaken is a GRPC error that is returned when attempting to save a new
// project and the project name is already taken.
var ErrNameTaken = status.Error(
	codes.AlreadyExists,
	"project already exists",
)

// ErrCouldNotSave is a GRPC error that is returned when an unknown error occurs
// while saving the project.
var ErrCouldNotSave = status.Error(
	codes.Internal,
	"could not create a new project",
)

// ErrNotFound is a GRPC error that is returned when no projects were found.
var ErrNotFound = status.Error(
	codes.NotFound,
	"no project was found with that name",
)

// ErrCouldNotFetch is a GRPC error that is returned when an unknown error
// occurs while fetching the project.
var ErrCouldNotFetch = status.Error(
	codes.Internal,
	"error occurred while fetching project",
)

// ErrAddEnvironment is a GRPC error that is returned when an error occurs while
// updating the environments array for a project.
var ErrAddEnvironment = status.Error(
	codes.Internal,
	"could not create a new environment",
)

// ErrAddFlag is a GRPC error that is returned when an error occurs while
// updating the flags array for a project.
var ErrAddFlag = status.Error(
	codes.Internal,
	"could not update project with the flag",
)

// ErrAddFlagSetting is a GRPC error that is returned when an error occurs while
// updating the flag settings array of a project.
var ErrAddFlagSetting = status.Error(
	codes.Internal,
	"error while updating project with flag settings",
)

// ErrMetadataNotFound is a GRPC error that is returned when trying to
// authenticate a request using the project key and the required request
// metadata was not found.
var ErrMetadataNotFound = status.Error(
	codes.Internal,
	"request metadata was not found",
)

// ErrProjectKeyNotFound is a GRPC error that is returned when the project key
// was not found in the request metadata.
var ErrProjectKeyNotFound = status.Error(
	codes.Unauthenticated,
	"project key not found in request metadata",
)

// ErrKeyMetadataLength is a GRPC error that is returned when the length of the
// project key metadata does not match the expected length i.e. 1.
var ErrKeyMetadataLength = status.Error(
	codes.InvalidArgument,
	"invalid length of project key metadata",
)
