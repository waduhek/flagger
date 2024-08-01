package environment_test

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/waduhek/flagger/internal/environment"
)

const mongoDBConnectionString string = "mongodb://localhost:27017"

const dummyObjectID string = "66a4836693cca0acf7482f8b"

//nolint:gochecknoglobals // Only to be used to create an environment.
var dummyEnvironmentStruct = &environment.Environment{
	Name:      "test",
	ProjectID: dummyObjectID,
	CreatedBy: dummyObjectID,
	CreatedAt: time.Now(),
}

//nolint:gochecknoglobals // Initialised only in TestMain function.
var environmentCollection *mongo.Collection

//nolint:gochecknoglobals // Initialised only in TestMain function.
var environmentRepository environment.DataRepository

func TestSaveEnvironment(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, saveErr := environmentRepository.Save(ctx, dummyEnvironmentStruct)
	if saveErr != nil {
		t.Fatalf("error while saving environment: %v\n", saveErr)
	}

	cleanupEnvironment(t, dummyEnvironmentStruct)
}

func TestSaveEnvironment_InvalidIDs(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		Environment *environment.Environment
	}{
		{
			Name: "invalid_project_id",
			Environment: &environment.Environment{
				Name:      "test_name_invalid_project_id",
				ProjectID: "invalid_object_id",
				CreatedBy: dummyObjectID,
				CreatedAt: time.Now(),
			},
		},
		{
			Name: "invalid_created_by",
			Environment: &environment.Environment{
				Name:      "test_name_invalid_created_by",
				ProjectID: dummyObjectID,
				CreatedBy: "invalid_object_id",
				CreatedAt: time.Now(),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			cleanupEnvironment(t, testCase.Environment)

			_, saveErr := environmentRepository.Save(ctx, testCase.Environment)
			if !errors.Is(saveErr, environment.ErrCouldNotSave) {
				t.Fatal("expected ErrCouldNotSave error while saving")
			}
		})
	}
}

func TestSaveEnvironment_Duplicate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, firstSaveErr := environmentRepository.Save(ctx, dummyEnvironmentStruct)
	if firstSaveErr != nil {
		t.Fatalf("error while saving: %v", firstSaveErr)
	}

	cleanupEnvironment(t, dummyEnvironmentStruct)

	_, secondSaveErr := environmentRepository.Save(ctx, dummyEnvironmentStruct)
	if !errors.Is(secondSaveErr, environment.ErrNameTaken) {
		t.Fatal("expected error when saving duplicate environment")
	}
}

func TestGetByID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	savedID, saveErr := environmentRepository.Save(ctx, dummyEnvironmentStruct)
	if saveErr != nil {
		t.Fatalf("unexpected error when saving environment: %v", saveErr)
	}

	cleanupEnvironment(t, dummyEnvironmentStruct)

	fetchedEnvironment, getErr := environmentRepository.GetByID(ctx, savedID)
	if getErr != nil {
		t.Fatalf("error while fetching environment details: %v", getErr)
	}

	if fetchedEnvironment.ID != savedID {
		t.Fatalf("received unequal environment ids")
	}
}

func TestGetByID_InvalidID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, getErr := environmentRepository.GetByID(ctx, "invalid_object_id")
	if !errors.Is(getErr, environment.ErrCouldNotFetch) {
		t.Fatal("expected ErrCouldNotFetch while fetching environment details")
	}
}

func TestGetByID_NoEnvironment(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, getErr := environmentRepository.GetByID(ctx, dummyObjectID)
	if !errors.Is(getErr, environment.ErrNotFound) {
		t.Fatalf("expected not found error")
	}
}

func TestGetByNameAndProjectID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	savedID, saveErr := environmentRepository.Save(ctx, dummyEnvironmentStruct)
	if saveErr != nil {
		t.Fatalf("unexpected error when saving environment: %v", saveErr)
	}

	cleanupEnvironment(t, dummyEnvironmentStruct)

	fetchedEnvironment, getErr := environmentRepository.GetByNameAndProjectID(
		ctx,
		dummyEnvironmentStruct.Name,
		dummyEnvironmentStruct.ProjectID,
	)
	if getErr != nil {
		t.Fatalf("error while getting environment: %v", getErr)
	}

	if savedID != fetchedEnvironment.ID {
		t.Fatalf("received unequal environments")
	}
}

func TestGetByNameAndProjectID_InvalidProjectID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, getErr := environmentRepository.GetByNameAndProjectID(
		ctx,
		dummyEnvironmentStruct.Name,
		"invalid_object_id",
	)
	if !errors.Is(getErr, environment.ErrCouldNotFetch) {
		t.Fatalf("expected error while getting environment")
	}
}

func TestGetByNameAndProjectID_NoEnvironment(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, getErr := environmentRepository.GetByNameAndProjectID(
		ctx,
		dummyEnvironmentStruct.Name,
		dummyEnvironmentStruct.ProjectID,
	)
	if !errors.Is(getErr, environment.ErrNotFound) {
		t.Fatalf("expected ErrNotFound")
	}
}

func TestMain(m *testing.M) {
	client, mongoClientErr := getMongoClient()
	if mongoClientErr != nil {
		os.Exit(1)
	}

	mongoDatabase := client.Database("test")
	environmentCollection = mongoDatabase.Collection("environments")

	environmentRepo, repositoryErr := getEnvironmentRepository(mongoDatabase)
	if repositoryErr != nil {
		log.Fatalf("error while creating environment repository: %v\n", repositoryErr)
	}
	environmentRepository = environmentRepo

	code := m.Run()

	_ = disconnectFromMongoDB(client)

	os.Exit(code)
}

func getMongoClient() (*mongo.Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.
		Client().
		ApplyURI(mongoDBConnectionString).
		SetServerAPIOptions(serverAPI)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Printf("error while connecting to mongodb: %v\n", err)
		return nil, err
	}

	return client, nil
}

func disconnectFromMongoDB(client *mongo.Client) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	err := client.Disconnect(ctx)
	if err != nil {
		log.Printf("error while disconnecting from mongodb: %v", err)
		return err
	}

	return nil
}

func getEnvironmentRepository(mongoDatabase *mongo.Database) (environment.DataRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	environmentRepository, repositoryErr := environment.NewEnvironmentRepository(ctx, mongoDatabase)
	if repositoryErr != nil {
		return nil, repositoryErr
	}

	return environmentRepository, nil
}

func cleanupEnvironment(t *testing.T, e *environment.Environment) {
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		_, _ = environmentCollection.DeleteOne(
			ctx,
			bson.D{{Key: "name", Value: e.Name}},
		)
	})
}
