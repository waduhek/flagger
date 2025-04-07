package flagsetting_test

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"
	"time"

	"github.com/waduhek/flagger/internal/flagsetting"
	"github.com/waduhek/flagger/internal/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const mongoDBConnectionString string = "mongodb://localhost:27017"

const dummyObjectID string = "66a4836693cca0acf7482f8b"

var flagsettingCollection *mongo.Collection

var flagsettingRepository flagsetting.DataRepository

var dummyFlagSetting = &flagsetting.FlagSetting{
	ProjectID:     dummyObjectID,
	EnvironmentID: dummyObjectID,
	FlagID:        dummyObjectID,
	IsActive:      true,
	CreatedAt:     time.Now(),
	UpdatedAt:     time.Now(),
}

func TestSave(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, saveErr := flagsettingRepository.Save(ctx, dummyFlagSetting)
	if saveErr != nil {
		t.Fatalf("error while saving flag setting: %v", saveErr)
	}

	cleanupFlagSetting(t, dummyFlagSetting)
}

func TestSave_Duplicate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, firstSaveErr := flagsettingRepository.Save(ctx, dummyFlagSetting)
	if firstSaveErr != nil {
		t.Fatalf("error while saving flag setting: %v", firstSaveErr)
	}

	cleanupFlagSetting(t, dummyFlagSetting)

	_, secondSaveErr := flagsettingRepository.Save(ctx, dummyFlagSetting)
	if !errors.Is(secondSaveErr, flagsetting.ErrCouldNotSave) {
		t.Fatal("expected ErrCouldNotSave error when saving duplicate flag setting")
	}
}

func TestSave_InvalidID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		FlagSetting *flagsetting.FlagSetting
	}{
		{
			Name: "invalid_project_id",
			FlagSetting: &flagsetting.FlagSetting{
				ProjectID:     "invalid_object_id",
				EnvironmentID: dummyObjectID,
				FlagID:        dummyObjectID,
				IsActive:      true,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
		},
		{
			Name: "invalid_environment_id",
			FlagSetting: &flagsetting.FlagSetting{
				ProjectID:     dummyObjectID,
				EnvironmentID: "invalid_object_id",
				FlagID:        dummyObjectID,
				IsActive:      true,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
		},
		{
			Name: "invalid_flag_id",
			FlagSetting: &flagsetting.FlagSetting{
				ProjectID:     dummyObjectID,
				EnvironmentID: dummyObjectID,
				FlagID:        "invalid_object_id",
				IsActive:      true,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			_, saveErr := flagsettingRepository.Save(ctx, testCase.FlagSetting)
			if !errors.Is(saveErr, flagsetting.ErrCouldNotSave) {
				t.Fatalf("expected ErrCouldNotSave error")
			}
		})
	}
}

func TestSaveMany(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, saveErr := flagsettingRepository.SaveMany(
		ctx,
		[]flagsetting.FlagSetting{*dummyFlagSetting},
	)
	if saveErr != nil {
		t.Fatalf("got error when saving flag settings: %v", saveErr)
	}

	cleanupFlagSetting(t, dummyFlagSetting)
}

func TestSaveMany_Duplicate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, saveErr := flagsettingRepository.SaveMany(
		ctx,
		[]flagsetting.FlagSetting{*dummyFlagSetting, *dummyFlagSetting},
	)
	if !errors.Is(saveErr, flagsetting.ErrCouldNotSave) {
		t.Fatal("expected error when saving duplicate flag settings")
	}

	cleanupFlagSetting(t, dummyFlagSetting)
}

func TestGet(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	savedID, saveErr := flagsettingRepository.Save(ctx, dummyFlagSetting)
	if saveErr != nil {
		t.Fatalf("error while saving flag setting: %v", saveErr)
	}

	cleanupFlagSetting(t, dummyFlagSetting)

	gotFlagSetting, getErr := flagsettingRepository.Get(
		ctx,
		dummyFlagSetting.ProjectID,
		dummyFlagSetting.EnvironmentID,
		dummyFlagSetting.FlagID,
	)
	if getErr != nil {
		t.Fatalf("error while getting flag setting: %v", getErr)
	}

	if savedID != gotFlagSetting.ID {
		t.Fatalf("got unequal flag setting")
	}
}

func TestGet_InvalidID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name          string
		ProjectID     string
		EnvironmentID string
		FlagID        string
	}{
		{
			Name:          "invalid_project_id",
			ProjectID:     "invalid_object_id",
			EnvironmentID: dummyObjectID,
			FlagID:        dummyObjectID,
		},
		{
			Name:          "invalid_environment_id",
			ProjectID:     dummyObjectID,
			EnvironmentID: "invalid_object_id",
			FlagID:        dummyObjectID,
		},
		{
			Name:          "invalid_flag_id",
			ProjectID:     dummyObjectID,
			EnvironmentID: dummyObjectID,
			FlagID:        "invalid_object_id",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			_, saveErr := flagsettingRepository.Get(
				ctx,
				testCase.ProjectID,
				testCase.EnvironmentID,
				testCase.FlagID,
			)
			if !errors.Is(saveErr, flagsetting.ErrCouldNotGet) {
				t.Fatalf("expected ErrCouldNotGet error when getting flag setting")
			}
		})
	}
}

func TestUpdateIsActive(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, saveErr := flagsettingRepository.Save(ctx, dummyFlagSetting)
	if saveErr != nil {
		t.Fatalf("error while saving flag setting: %v", saveErr)
	}

	cleanupFlagSetting(t, dummyFlagSetting)

	updatedCount, updateErr := flagsettingRepository.UpdateIsActive(
		ctx,
		dummyFlagSetting.ProjectID,
		dummyFlagSetting.EnvironmentID,
		dummyFlagSetting.FlagID,
		false,
	)
	if updateErr != nil {
		t.Fatalf("error while updating active status of flag setting: %v", updateErr)
	}

	if updatedCount != 1 {
		t.Fatalf("expected 1 flag setting to be updated but got %v", updatedCount)
	}
}

func TestUpdateIsActive_InvalidID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name          string
		ProjectID     string
		EnvironmentID string
		FlagID        string
	}{
		{
			Name:          "invalid_project_id",
			ProjectID:     "invalid_object_id",
			EnvironmentID: dummyObjectID,
			FlagID:        dummyObjectID,
		},
		{
			Name:          "invalid_environment_id",
			ProjectID:     dummyObjectID,
			EnvironmentID: "invalid_object_id",
			FlagID:        dummyObjectID,
		},
		{
			Name:          "invalid_flag_id",
			ProjectID:     dummyObjectID,
			EnvironmentID: dummyObjectID,
			FlagID:        "invalid_object_id",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			_, updateErr := flagsettingRepository.UpdateIsActive(
				ctx,
				testCase.ProjectID,
				testCase.EnvironmentID,
				testCase.FlagID,
				true,
			)
			if !errors.Is(updateErr, flagsetting.ErrStatusUpdate) {
				t.Fatalf("expected ErrStatusUpdate error when updating setting active status")
			}
		})
	}
}

func TestMain(m *testing.M) {
	client, mongoClientErr := getMongoClient()
	if mongoClientErr != nil {
		os.Exit(1)
	}

	mongoDatabase := client.Database("test")
	flagsettingCollection = mongoDatabase.Collection("flag_settings")

	flagRepo, repositoryErr := getFlagSettingRepository(mongoDatabase)
	if repositoryErr != nil {
		log.Fatalf("error while creating flag repository: %v\n", repositoryErr)
	}
	flagsettingRepository = flagRepo

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

func getFlagSettingRepository(mongoDatabase *mongo.Database) (flagsetting.DataRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	flagRepository, repositoryErr := flagsetting.NewFlagSettingRepository(
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

func cleanupFlagSetting(t *testing.T, flagSetting *flagsetting.FlagSetting) {
	projectObjectID, _ := primitive.ObjectIDFromHex(flagSetting.ProjectID)
	environmentObjectID, _ := primitive.ObjectIDFromHex(flagSetting.EnvironmentID)
	flagObjectID, _ := primitive.ObjectIDFromHex(flagSetting.FlagID)

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		_, _ = flagsettingCollection.DeleteOne(
			ctx,
			bson.D{
				{Key: "project_id", Value: projectObjectID},
				{Key: "environment_id", Value: environmentObjectID},
				{Key: "flag_id", Value: flagObjectID},
			},
		)
	})
}
