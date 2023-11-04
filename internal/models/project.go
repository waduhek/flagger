package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Project struct {
	ID           primitive.ObjectID   `bson:"_id,omitempty"`
	Name         string               `bson:"name"`
	Environments []primitive.ObjectID `bson:"environments,omitempty"`
	Flags        []primitive.ObjectID `bson:"flags,omitempty"`
	CreatedBy    primitive.ObjectID   `bson:"created_by"`
	CreatedAt    time.Time            `bson:"created_at"`
	UpdatedAt    time.Time            `bson:"updated_at"`
}

// ProjectRepository is an interface to the operations that can be performed on
// the projects collection.
type ProjectRepository interface {
	// Creates a new `Project` document
	Save(ctx context.Context, project *Project) (*mongo.InsertOneResult, error)

	// Gets the project by the name and the ID of the user that created the
	// project.
	GetByNameAndUserID(
		ctx context.Context,
		projectName string,
		userID primitive.ObjectID,
	) (*Project, error)

	// Adds a new environment to the project by the project name and the user
	// ID.
	AddEnvironment(
		ctx context.Context,
		projectName string,
		userID primitive.ObjectID,
		environmentID primitive.ObjectID,
	) (*mongo.UpdateResult, error)

	// Adds a new flag to the project by the project name.
	AddFlag(
		ctx context.Context,
		projectName string,
		userID string,
		flagID primitive.ObjectID,
	) (*mongo.UpdateResult, error)
}
