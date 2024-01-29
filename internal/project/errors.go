package project

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// EProjectNameTaken is a GRPC error that is returned when attempting to save a
// new project and the project name is already taken.
var EProjectNameTaken = status.Error(
	codes.AlreadyExists,
	"project already exists",
)

// EProjectSave is a GRPC error that is returned when an unknown error occurs
// while saving the project.
var EProjectSave = status.Error(
	codes.Internal,
	"could not create a new project",
)

// EProjectNotFound is a GRPC error that is returned when no projects were
// found.
var EProjectNotFound = status.Error(
	codes.NotFound,
	"no project was found with that name",
)

// EProjectFetch is a GRPC error that is returned when an unknown error occurs
// while fetching the project.
var EProjectFetch = status.Error(
	codes.Internal,
	"error occurred while fetching project",
)

// EProjectAddEnvironment is a GRPC error that is returned when an error occurs
// while updating the environments array for a project.
var EProjectAddEnvironment = status.Error(
	codes.Internal,
	"could not create a new environment",
)

// EProjectAddFlag is a GRPC error that is returned when an error occurs while
// updating the flags array for a project.
var EProjectAddFlag = status.Error(
	codes.Internal,
	"could not update project with the flag",
)

// EProjectAddFlagSetting is a GRPC error that is returned when an error occurs
// while updating the flag settings array of a project.
var EProjectAddFlagSetting = status.Error(
	codes.Internal,
	"error while updating project with flag settings",
)

// EMetadataNotFound is a GRPC error that is returned when trying to
// authenticate a request using the project key and the required request
// metadata was not found.
var EMetadataNotFound = status.Error(
	codes.Internal,
	"request metadata was not found",
)

// EProjectKeyNotFound is a GRPC error that is returned when the project key
// was not found in the request metadata.
var EProjectKeyNotFound = status.Error(
	codes.Unauthenticated,
	"project key not found in request metadata",
)

// EKeyMetadataLength is a GRPC error that is returned when the length of the
// project key metadata does not match the expected length i.e. 1.
var EKeyMetadataLength = status.Error(
	codes.InvalidArgument,
	"invalid length of project key metadata",
)
