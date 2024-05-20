package project

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Project struct {
	ID           string
	Key          string
	Name         string
	Environments []string
	Flags        []string
	FlagSettings []string
	CreatedBy    string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// DataRepository is an interface to the operations that can be performed on
// the projects collection.
type DataRepository interface {
	// Save creates a new `Project` and returns the ID of the created project.
	Save(ctx context.Context, project *Project) (string, error)

	// GetByNameAndUserID gets the project by the name and the ID of the user
	// that created the project.
	GetByNameAndUserID(
		ctx context.Context,
		projectName string,
		userID string,
	) (*Project, error)

	// AddEnvironment adds a new environment ID to the project by the project
	// ID.
	AddEnvironment(
		ctx context.Context,
		projectID string,
		environmentID primitive.ObjectID,
	) (uint, error)

	// AddFlag adds a new flag ID to the project by the project ID.
	AddFlag(
		ctx context.Context,
		projectID string,
		flagID primitive.ObjectID,
	) (uint, error)

	// AddFlagSettings adds new flag setting IDs to the project.
	AddFlagSettings(
		ctx context.Context,
		projectID primitive.ObjectID,
		flagSettingIDs ...primitive.ObjectID,
	) (uint, error)
}
