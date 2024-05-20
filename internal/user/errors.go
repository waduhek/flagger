package user

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrNotFound is a GRPC error that occurs when the user was not found.
var ErrNotFound = status.Error(
	codes.NotFound,
	"no user with that username exists",
)

// ErrCouldNotFetch is a GRPC error that occurs when an unknown error occurs
// while fetching a user.
var ErrCouldNotFetch = status.Error(
	codes.Internal,
	"error occurred while fetching the user",
)

// ErrUsernameTaken is a GRPC error that occurs when attempting to save a new
// user and the provided username is taken.
var ErrUsernameTaken = status.Error(codes.AlreadyExists, "username is taken")

// ErrNotSaved is a GRPC error that occurs when an unknown error occurs while
// saving a user.
var ErrNotSaved = status.Error(
	codes.Internal,
	"could not save the user details",
)

// ErrPasswordUpdate is a GRPC error that occurs when the password could not be
// updated.
var ErrPasswordUpdate = status.Error(codes.Internal, "could not save password")

// ErrUserIDConvert is a GRPC error that is returned when the user ID could not
// be converted to the expected type.
var ErrUserIDConvert = status.Error(
	codes.Internal,
	"could not convert user ID to expected format",
)
