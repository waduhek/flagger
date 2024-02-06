package provider

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/waduhek/flagger/internal/environment"
	"github.com/waduhek/flagger/internal/flag"
	"github.com/waduhek/flagger/internal/flagsetting"
)

// FlagDetails contains the details of a particular flag in a particular
// environment.
type FlagDetails struct {
	ID          primitive.ObjectID      `bson:"_id"`
	Key         string                  `bson:"key"`
	Name        string                  `bson:"name"`
	Environment environment.Environment `bson:"environment"`
	Flag        flag.Flag               `bson:"flag"`
	FlagSetting flagsetting.FlagSetting `bson:"flag_setting"`
	CreatedBy   primitive.ObjectID      `bson:"created_by"`
	CreatedAt   time.Time               `bson:"created_at"`
	UpdatedAt   time.Time               `bson:"updated_at"`
}

type ProviderRepository interface {
	// GetFlagDetailsByProjectKey gets the details of a flag by the project key,
	// environment name and the flag name.
	GetFlagDetailsByProjectKey(
		ctx context.Context,
		projectKey string,
		environmentName string,
		flagName string,
	) ([]FlagDetails, error)
}
