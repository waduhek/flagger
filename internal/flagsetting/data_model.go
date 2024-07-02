package flagsetting

import (
	"context"
	"time"
)

type FlagSetting struct {
	ID            string
	ProjectID     string
	EnvironmentID string
	FlagID        string
	IsActive      bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type DataRepository interface {
	// Save creates a new `FlagSetting`.
	Save(
		ctx context.Context,
		flagSetting *FlagSetting,
	) (string, error)

	// SaveMany creates multiple `FlagSetting`s.
	SaveMany(
		ctx context.Context,
		flagSettings []FlagSetting,
	) ([]string, error)

	// Get gets the flag setting by the project ID, environment ID and the
	// flag ID.
	Get(
		ctx context.Context,
		projectID string,
		environmentID string,
		flagID string,
	) (*FlagSetting, error)

	// UpdateIsActive updates a flag setting's `IsActive` field to the provided
	// value.
	UpdateIsActive(
		ctx context.Context,
		projectID string,
		environmentID string,
		flagID string,
		isActive bool,
	) (uint, error)
}
