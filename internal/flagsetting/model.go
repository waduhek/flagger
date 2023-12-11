package flagsetting

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FlagSetting struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	ProjectID     primitive.ObjectID `bson:"project_id"`
	EnvironmentID primitive.ObjectID `bson:"environment_id"`
	FlagID        primitive.ObjectID `bson:"flag_id"`
	IsActive      bool               `bson:"is_active"`
	CreatedAt     time.Time          `bson:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at"`
}

type FlagSettingRepository interface {
	// Save creates a new `FlagSetting`.
	Save(
		ctx context.Context,
		flagSetting *FlagSetting,
	) (*mongo.InsertOneResult, error)

	// SaveMany creates multiple `FlagSetting`s.
	SaveMany(
		ctx context.Context,
		flagSettings []FlagSetting,
	) (*mongo.InsertManyResult, error)

	// Get gets the flag setting by the project ID, environment ID and the
	// flag ID.
	Get(
		ctx context.Context,
		projectID primitive.ObjectID,
		environmentID primitive.ObjectID,
		flagID primitive.ObjectID,
	) (*FlagSetting, error)

	// UpdateIsActive updates a flag setting's `IsActive` field to the provided
	// value.
	UpdateIsActive(
		ctx context.Context,
		flagSettingID primitive.ObjectID,
		isActive bool,
	) (*mongo.UpdateResult, error)
}
