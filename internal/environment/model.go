package environment

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Environment struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	ProjectID primitive.ObjectID `bson:"project_id"`
	CreatedBy primitive.ObjectID `bson:"created_by"`
	CreatedAt time.Time          `bson:"created_at"`
}

type EnvironmentRepository interface {
	// Creates a new `Environment` document.
	Save(
		ctx context.Context,
		environment *Environment,
	) (*mongo.InsertOneResult, error)

	// Gets an `Environment` by it's object ID.
	GetByID(ctx context.Context, id primitive.ObjectID) (*Environment, error)

	// Gets an environment by it's name and the project ID.
	GetByNameAndProjectID(
		ctx context.Context,
		environmentName string,
		projectID primitive.ObjectID,
	) (*Environment, error)
}
