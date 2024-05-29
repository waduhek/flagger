package environment

import (
	"context"
	"time"
)

type Environment struct {
	ID        string
	Name      string
	ProjectID string
	CreatedBy string
	CreatedAt time.Time
}

type DataRepository interface {
	// Creates a new `Environment` document.
	Save(
		ctx context.Context,
		environment *Environment,
	) (string, error)

	// Gets an `Environment` by it's object ID.
	GetByID(ctx context.Context, id string) (*Environment, error)

	// Gets an environment by it's name and the project ID.
	GetByNameAndProjectID(
		ctx context.Context,
		environmentName string,
		projectID string,
	) (*Environment, error)
}
