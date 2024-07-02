package flag

import (
	"context"
	"time"
)

type Flag struct {
	ID        string
	Name      string
	ProjectID string
	CreatedBy string
	CreatedAt time.Time
}

type DataRepository interface {
	// Save creates a new `Flag` document.
	Save(ctx context.Context, flag *Flag) (string, error)

	// GetByID gets a `Flag` by it's object ID.
	GetByID(ctx context.Context, flagID string) (*Flag, error)

	// GetByNameAndProjectID gets a `Flag` by it's name and the ID of the
	// project that it belongs to.
	GetByNameAndProjectID(
		ctx context.Context,
		flagName string,
		projectID string,
	) (*Flag, error)
}
