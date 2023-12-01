package flag

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Flag struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	ProjectID primitive.ObjectID `bson:"project_id"`
	CreatedBy primitive.ObjectID `bson:"created_by"`
	CreatedAt time.Time          `bson:"created_at"`
}

type FlagRepository interface {
	// Save creates a new `Flag` document.
	Save(ctx context.Context, flag *Flag) (*mongo.InsertOneResult, error)

	// GetByID gets a `Flag` by it's object ID.
	GetByID(ctx context.Context, flagID primitive.ObjectID) (*Flag, error)

	// GetByNameAndProjectID gets a `Flag` by it's name and the ID of the
	// project that it belongs to.
	GetByNameAndProjectID(
		ctx context.Context,
		flagName string,
		projectID primitive.ObjectID,
	) (*Flag, error)
}
