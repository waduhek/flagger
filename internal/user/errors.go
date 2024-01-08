package user

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// EUserNotFound is a GRPC error that occurs when the user was not found.
var EUserNotFound = status.Error(
	codes.NotFound,
	"no user with that username exists",
)

// ECouldNotFetchUser is a GRPC error that occurs when an unknown error occurs
// while fetching a user.
var ECouldNotFetchUser = status.Error(
	codes.Internal,
	"error occurred while fetching the user",
)

// EUsernameTaken is a GRPC error that occurs when attempting to save a new user
// and the provided username is taken.
var EUsernameTaken = status.Error(codes.AlreadyExists, "username is taken")

// EUserNotSaved is a GRPC error that occurs when an unknown error occurs while
// saving a user.
var EUserNotSaved = status.Error(
	codes.Internal,
	"could not save the user details",
)

// EPasswordUpdate is a GRPC error that occurs when the password could not be
// updated.
var EPasswordUpdate = status.Error(codes.Internal, "could not save password")
