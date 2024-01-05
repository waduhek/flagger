package flagsetting

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// EFlagSettingSave is a GRPC error that is returned when an unknown error
// occurs while saving flag settings.
var EFlagSettingSave = status.Error(
	codes.Internal,
	"could not save flag settings",
)

// EFlagSettingStatusUpdate is a GRPC error that is returned when an error
// occurs while updating the flag settings status.
var EFlagSettingStatusUpdate = status.Error(
	codes.Internal,
	"error occurred while updating the flag setting",
)
