package flag_test

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

	"github.com/waduhek/flagger/internal/flag"
	"github.com/waduhek/flagger/internal/logger"
)

const mongoDBConnectionString string = "mongodb://localhost:27017"

const dummyObjectID string = "66a4836693cca0acf7482f8b"

var dummyFlag = &flag.Flag{
	Name:      "test",
	ProjectID: dummyObjectID,
	CreatedBy: dummyObjectID,
	CreatedAt: time.Now(),
}

var flagCollection *mongo.Collection

var flagRepository flag.DataRepository

func TestSave(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, saveErr := flagRepository.Save(ctx, dummyFlag)
	if saveErr != nil {
		t.Fatalf("error while saving flag")
	}

	cleanupFlag(t, dummyFlag)
}

func TestSave_InvalidID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name string
		Flag *flag.Flag
	}{
		{
			Name: "invalid_project_id",
			Flag: &flag.Flag{
				Name:      "test_invalid_project_id",
				ProjectID: "invalid_object_id",
				CreatedBy: dummyObjectID,
				CreatedAt: time.Now(),
			},
		},
		{
			Name: "invalid_created_by",
			Flag: &flag.Flag{
				Name:      "test_invalid_object_id",
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

			_, saveErr := flagRepository.Save(ctx, testCase.Flag)
			if !errors.Is(saveErr, flag.ErrCouldNotSave) {
				t.Fatalf("expected ErrCouldNotSave error")
			}
		})
	}
}

func TestSave_Duplicate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, firstSaveErr := flagRepository.Save(ctx, dummyFlag)
	if firstSaveErr != nil {
		t.Fatalf("error while saving flag: %v", firstSaveErr)
	}

	cleanupFlag(t, dummyFlag)

	_, secondSaveError := flagRepository.Save(ctx, dummyFlag)
	if !errors.Is(secondSaveError, flag.ErrNameTaken) {
		t.Fatal("expected ErrNameTaken when saving duplicate flag")
	}
}

func TestGetByID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	savedID, saveErr := flagRepository.Save(ctx, dummyFlag)
	if saveErr != nil {
		t.Fatalf("error while saving flag: %v", saveErr)
	}

	cleanupFlag(t, dummyFlag)

	gotFlag, getErr := flagRepository.GetByID(ctx, savedID)
	if getErr != nil {
		t.Fatalf("error while getting flag details: %v", getErr)
	}

	if gotFlag.ID != savedID {
		t.Fatal("got different flag ids")
	}
}

func TestGetByID_InvalidFlagID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, fetchErr := flagRepository.GetByID(ctx, "invalid_object_id")
	if !errors.Is(fetchErr, flag.ErrCouldNotFetch) {
		t.Fatal("expected ErrCouldNotFetch error when getting flag details")
	}
}

func TestGetByID_NoFlag(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, getErr := flagRepository.GetByID(ctx, dummyObjectID)
	if !errors.Is(getErr, flag.ErrNotFound) {
		t.Fatal("expected ErrNotFound error")
	}
}

func TestGetByNameAndProjectID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	savedID, saveErr := flagRepository.Save(ctx, dummyFlag)
	if saveErr != nil {
		t.Fatalf("error while saving flag: %v", saveErr)
	}

	cleanupFlag(t, dummyFlag)

	gotFlag, getErr := flagRepository.GetByNameAndProjectID(ctx, dummyFlag.Name, dummyFlag.ProjectID)
	if getErr != nil {
		t.Fatalf("error while getting flag details: %v", getErr)
	}

	if gotFlag.ID != savedID {
		t.Fatal("got different flag IDs")
	}
}

func TestGetByNameAndProjectID_NoFlag(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, getErr := flagRepository.GetByNameAndProjectID(ctx, dummyFlag.Name, dummyFlag.ProjectID)
	if !errors.Is(getErr, flag.ErrNotFound) {
		t.Fatal("expected ErrNotFound error")
	}
}

func TestGetByNameAndProjectID_InvalidProjectID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, getErr := flagRepository.GetByNameAndProjectID(ctx, dummyFlag.Name, "invalid_object_id")
	if !errors.Is(getErr, flag.ErrCouldNotFetch) {
		t.Fatalf("expected ErrCouldNotFetch error for invalid object id")
	}
}

func TestMain(m *testing.M) {
	client, mongoClientErr := getMongoClient()
	if mongoClientErr != nil {
		os.Exit(1)
	}

	mongoDatabase := client.Database("test")
	flagCollection = mongoDatabase.Collection("flags")

	flagRepo, repositoryErr := getFlagRepository(mongoDatabase)
	if repositoryErr != nil {
		log.Fatalf("error while creating flag repository: %v\n", repositoryErr)
	}
	flagRepository = flagRepo

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

func getFlagRepository(mongoDatabase *mongo.Database) (flag.DataRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	flagRepository, repositoryErr := flag.NewFlagRepository(
		ctx,
		mongoDatabase,
		&logger.StubLogger{},
	)
	if repositoryErr != nil {
		return nil, repositoryErr
	}

	return flagRepository, nil
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

func cleanupFlag(t *testing.T, flag *flag.Flag) {
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		_, _ = flagCollection.DeleteOne(
			ctx,
			bson.D{{Key: "name", Value: flag.Name}},
		)
	})
}
