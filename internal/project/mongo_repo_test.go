package project_test

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

	"github.com/waduhek/flagger/internal/project"
	"github.com/waduhek/flagger/internal/user"
)

const mongoDBConnectionString string = "mongodb://localhost:27017"

const dummyObjectID string = "66a4836693cca0acf7482f8b"

var dummyProject = &project.Project{
	Name:         "test",
	Key:          "test",
	Flags:        []string{},
	Environments: []string{},
	FlagSettings: []string{},
	CreatedBy:    dummyObjectID,
	CreatedAt:    time.Now(),
	UpdatedAt:    time.Now(),
}

var projectCollection *mongo.Collection

var projectRepository project.DataRepository

func TestSave(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, saveErr := projectRepository.Save(ctx, dummyProject)
	if saveErr != nil {
		t.Fatalf("error while saving project: %v", saveErr)
	}

	cleanupProject(t, dummyProject)
}

func TestSave_Duplicate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, firstSaveErr := projectRepository.Save(ctx, dummyProject)
	if firstSaveErr != nil {
		t.Fatalf("error while saving project: %v", firstSaveErr)
	}

	cleanupProject(t, dummyProject)

	_, secondSaveErr := projectRepository.Save(ctx, dummyProject)
	if !errors.Is(secondSaveErr, project.ErrNameTaken) {
		t.Fatalf("expected ErrNameTaken error but got %v", secondSaveErr)
	}
}

func TestSave_InvalidIDs(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name    string
		Project *project.Project
	}{
		{
			Name: "invalid_environment_id",
			Project: &project.Project{
				Name:         "test",
				Key:          "test",
				Environments: []string{"invalid_object_id"},
				Flags:        []string{dummyObjectID},
				FlagSettings: []string{dummyObjectID},
				CreatedBy:    dummyObjectID,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
		},
		{
			Name: "invalid_flag_id",
			Project: &project.Project{
				Name:         "test",
				Key:          "test",
				Environments: []string{dummyObjectID},
				Flags:        []string{"invalid_object_id"},
				FlagSettings: []string{dummyObjectID},
				CreatedBy:    dummyObjectID,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
		},
		{
			Name: "invalid_flagsetting_id",
			Project: &project.Project{
				Name:         "test",
				Key:          "test",
				Environments: []string{dummyObjectID},
				Flags:        []string{dummyObjectID},
				FlagSettings: []string{"invalid_object_id"},
				CreatedBy:    dummyObjectID,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
		},
		{
			Name: "invalid_createdby_id",
			Project: &project.Project{
				Name:         "test",
				Key:          "test",
				Environments: []string{dummyObjectID},
				Flags:        []string{dummyObjectID},
				FlagSettings: []string{dummyObjectID},
				CreatedBy:    "invalid_object_id",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			_, saveErr := projectRepository.Save(ctx, testCase.Project)
			if !errors.Is(saveErr, project.ErrCouldNotSave) {
				t.Fatalf("expected ErrCouldNotSave error while saving project")
			}
		})
	}
}

func TestGetByNameAndUserID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	savedID, saveErr := projectRepository.Save(ctx, dummyProject)
	if saveErr != nil {
		t.Fatalf("error while saving project: %v", saveErr)
	}

	cleanupProject(t, dummyProject)

	gotProject, getErr := projectRepository.GetByNameAndUserID(
		ctx,
		dummyProject.Name,
		dummyProject.CreatedBy,
	)
	if getErr != nil {
		t.Fatalf("error while getting project details: %v", getErr)
	}

	if savedID != gotProject.ID {
		t.Fatalf("unexpected project returned")
	}
}

func TestGetByNameAndUserID_InvalidUserID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, getErr := projectRepository.GetByNameAndUserID(
		ctx,
		dummyProject.Name,
		"invalid_object_id",
	)
	if !errors.Is(getErr, user.ErrUserIDConvert) {
		t.Fatalf("expected ErrUserIDConvert error but got %v", getErr)
	}
}

func TestGetByNameAndUserID_NotFound(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, getErr := projectRepository.GetByNameAndUserID(
		ctx,
		dummyProject.Name,
		dummyProject.CreatedBy,
	)
	if !errors.Is(getErr, project.ErrNotFound) {
		t.Fatalf("expected ErrNotFound but got %v", getErr)
	}
}

func TestAddEnvironment(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	savedID, saveErr := projectRepository.Save(ctx, dummyProject)
	if saveErr != nil {
		t.Fatalf("got error while saving project: %v", saveErr)
	}

	cleanupProject(t, dummyProject)

	_, updateErr := projectRepository.AddEnvironment(ctx, savedID, dummyObjectID)
	if updateErr != nil {
		t.Fatalf("got error while adding environment to project: %v", updateErr)
	}
}

func TestAddEnvironment_InvalidIDs(t *testing.T) {
	t.Parallel()

	testCases := []struct{
		Name string
		ProjectID string
		EnvironmentID string
	}{
		{
			Name: "invalid_project_id",
			ProjectID: "invalid_object_id",
			EnvironmentID: dummyObjectID,
		},
		{
			Name: "invalid_environment_id",
			ProjectID: dummyObjectID,
			EnvironmentID: "invalid_environment_id",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			_, updateErr := projectRepository.AddEnvironment(
				ctx,
				testCase.ProjectID,
				testCase.EnvironmentID,
			)

			if !errors.Is(updateErr, project.ErrProjectIDConvert) {
				t.Fatalf("expected ErrProjectIDConvert error but got %v", updateErr)
			}
		})
	}
}

func TestAddFlag(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	savedID, saveErr := projectRepository.Save(ctx, dummyProject)
	if saveErr != nil {
		t.Fatalf("got error while saving project: %v", saveErr)
	}

	cleanupProject(t, dummyProject)

	_, updateErr := projectRepository.AddFlag(ctx, savedID, dummyObjectID)
	if updateErr != nil {
		t.Fatalf("got error while adding flag to project: %v", updateErr)
	}
}

func TestAddFlag_InvalidIDs(t *testing.T) {
	t.Parallel()

	testCases := []struct{
		Name string
		ProjectID string
		FlagID string
	}{
		{
			Name: "invalid_project_id",
			ProjectID: "invalid_object_id",
			FlagID: dummyObjectID,
		},
		{
			Name: "invalid_environment_id",
			ProjectID: dummyObjectID,
			FlagID: "invalid_environment_id",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			_, updateErr := projectRepository.AddFlag(ctx, testCase.ProjectID, testCase.FlagID)
			if !errors.Is(updateErr, project.ErrAddFlag) {
				t.Fatalf("expected ErrAddFlag error but got %v", updateErr)
			}
		})
	}
}

func TestAddFlagSettings(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	savedID, saveErr := projectRepository.Save(ctx, dummyProject)
	if saveErr != nil {
		t.Fatalf("got error while saving project: %v", saveErr)
	}

	cleanupProject(t, dummyProject)

	_, updateErr := projectRepository.AddFlagSettings(ctx, savedID, dummyObjectID)
	if updateErr != nil {
		t.Fatalf("got error while adding flag settings to project: %v", updateErr)
	}
}

func TestAddFlagSettings_InvalidIDs(t *testing.T) {
	t.Parallel()

	testCases := []struct{
		Name string
		ProjectID string
		FlagSettingIDs []string
	}{
		{
			Name: "invalid_project_id",
			ProjectID: "invalid_object_id",
			FlagSettingIDs: []string{dummyObjectID},
		},
		{
			Name: "invalid_environment_id",
			ProjectID: dummyObjectID,
			FlagSettingIDs: []string{"invalid_environment_id"},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			_, updateErr := projectRepository.AddFlagSettings(
				ctx,
				testCase.ProjectID,
				testCase.FlagSettingIDs...,
			)
			if !errors.Is(updateErr, project.ErrAddFlagSetting) {
				t.Fatalf("expected ErrAddFlagSetting error but got %v", updateErr)
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
	projectCollection = mongoDatabase.Collection(project.ProjectCollection)

	projectRepo, repositoryErr := getProjectRepository(mongoDatabase)
	if repositoryErr != nil {
		log.Fatalf("error while creating flag repository: %v\n", repositoryErr)
	}
	projectRepository = projectRepo

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

func getProjectRepository(mongoDatabase *mongo.Database) (project.DataRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	flagRepository, repositoryErr := project.NewProjectRepository(ctx, mongoDatabase)
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

func cleanupProject(t *testing.T, project *project.Project) {
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		_, _ = projectCollection.DeleteOne(ctx, bson.D{{Key: "key", Value: project.Key}})
	})
}
