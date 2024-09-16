package project

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/waduhek/flagger/internal/user"
)

const ProjectCollection string = "projects"

// projectMongoModel is the MongoDB representation of the `Project` struct.
type projectMongoModel struct {
	ID           primitive.ObjectID   `bson:"_id,omitempty"`
	Key          string               `bson:"key"`
	Name         string               `bson:"name"`
	Environments []primitive.ObjectID `bson:"environments,omitempty"`
	Flags        []primitive.ObjectID `bson:"flags,omitempty"`
	FlagSettings []primitive.ObjectID `bson:"flag_settings,omitempty"`
	CreatedBy    primitive.ObjectID   `bson:"created_by"`
	CreatedAt    time.Time            `bson:"created_at"`
	UpdatedAt    time.Time            `bson:"updated_at"`
}

type MongoDataRepository struct {
	coll *mongo.Collection
}

func (p *MongoDataRepository) Save(
	ctx context.Context,
	project *Project,
) (string, error) {
	mappedEnvironmentIDs, mapEnvironmentIDErr :=
		mapStringSliceToObjectIDs(project.Environments)
	if mapEnvironmentIDErr != nil {
		log.Printf("could not map environment ids to object ids: %v", mapEnvironmentIDErr)
		return "", ErrCouldNotSave
	}

	mappedFlagIDs, mapFlagIDErr := mapStringSliceToObjectIDs(project.Flags)
	if mapFlagIDErr != nil {
		log.Printf("could not map flag ids to object ids: %v", mapFlagIDErr)
		return "", ErrCouldNotSave
	}

	mappedFlagSettingIDs, mapFlagSettingIDErr :=
		mapStringSliceToObjectIDs(project.FlagSettings)
	if mapFlagSettingIDErr != nil {
		log.Printf("could not map flag setting ids to object ids: %v", mapFlagSettingIDErr)
		return "", ErrCouldNotSave
	}

	createdByObjID, createdByObjIDErr := primitive.ObjectIDFromHex(project.CreatedBy)
	if createdByObjIDErr != nil {
		log.Printf("could not convert created by to object id: %v", createdByObjIDErr)
		return "", ErrCouldNotSave
	}

	projectToAdd := &projectMongoModel{
		Key:          project.Key,
		Name:         project.Name,
		Environments: mappedEnvironmentIDs,
		Flags:        mappedFlagIDs,
		FlagSettings: mappedFlagSettingIDs,
		CreatedBy:    createdByObjID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	saveResult, saveErr := p.coll.InsertOne(ctx, projectToAdd)
	if saveErr != nil {
		if mongo.IsDuplicateKeyError(saveErr) {
			log.Printf(
				"a project with name %q exists",
				project.Name,
			)
			return "", ErrNameTaken
		}

		log.Printf("error occurred while saving project: %v", saveErr)
		return "", ErrCouldNotSave
	}

	projectID, ok := saveResult.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Printf("could not assert saved project id as object id")
		return "", ErrCouldNotSave
	}

	return projectID.Hex(), nil
}

func (p *MongoDataRepository) GetByNameAndUserID(
	ctx context.Context,
	projectName string,
	userID string,
) (*Project, error) {
	userIDObjectID, userIDErr := primitive.ObjectIDFromHex(userID)
	if userIDErr != nil {
		log.Printf("could not convert user ID to ObjectID")
		return nil, user.ErrUserIDConvert
	}

	query := bson.D{
		{Key: "created_by", Value: userIDObjectID},
		{Key: "name", Value: projectName},
	}

	var decodedProject projectMongoModel
	findErr := p.coll.FindOne(ctx, query).Decode(&decodedProject)
	if findErr != nil {
		log.Printf(
			"error while getting project %q with user %q: %v",
			projectName,
			userID,
			findErr,
		)

		if errors.Is(findErr, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}

		return nil, ErrCouldNotFetch
	}

	project := &Project{
		ID:           decodedProject.ID.Hex(),
		Key:          decodedProject.Key,
		Name:         decodedProject.Name,
		Environments: mapObjectIDsToStringSlice(decodedProject.Environments),
		Flags:        mapObjectIDsToStringSlice(decodedProject.Flags),
		FlagSettings: mapObjectIDsToStringSlice(decodedProject.FlagSettings),
		CreatedBy:    decodedProject.CreatedBy.Hex(),
		CreatedAt:    decodedProject.CreatedAt,
		UpdatedAt:    decodedProject.UpdatedAt,
	}

	return project, nil
}

func (p *MongoDataRepository) AddEnvironment(
	ctx context.Context,
	projectID string,
	environmentID string,
) (uint, error) {
	projectIDObjID, projectIDConvertErr := primitive.ObjectIDFromHex(projectID)
	if projectIDConvertErr != nil {
		log.Printf("could not convert project ID to ObjectID")
		return 0, ErrProjectIDConvert
	}

	environmentIDObjID, environmentIDObjIDErr := primitive.ObjectIDFromHex(environmentID)
	if environmentIDObjIDErr != nil {
		log.Printf("could not convert environment id to object id")
		return 0, ErrProjectIDConvert
	}

	filterQuery := bson.D{
		{Key: "_id", Value: projectIDObjID},
	}
	updateQuery := bson.D{
		{
			Key:   "$push",
			Value: bson.D{{Key: "environments", Value: environmentIDObjID}},
		},
		{
			Key:   "$set",
			Value: bson.D{{Key: "updated_at", Value: time.Now()}},
		},
	}

	updateResult, updateErr := p.coll.UpdateOne(ctx, filterQuery, updateQuery)
	if updateErr != nil {
		log.Printf(
			"error while adding environment %q to project %q: %v",
			environmentID,
			projectIDObjID,
			updateErr,
		)
		return 0, ErrAddEnvironment
	}

	//nolint:gosec // ModifiedCount can't be a negative number.
	return uint(updateResult.ModifiedCount), nil
}

func (p *MongoDataRepository) AddFlag(
	ctx context.Context,
	projectID string,
	flagID string,
) (uint, error) {
	projectIDObjID, projectIDObjIDErr := primitive.ObjectIDFromHex(projectID)
	if projectIDObjIDErr != nil {
		log.Printf("error while converting project id to object id: %v", projectIDObjIDErr)
		return 0, ErrAddFlag
	}

	flagIDObjID, flagIDObjIDErr := primitive.ObjectIDFromHex(flagID)
	if flagIDObjIDErr != nil {
		log.Printf("error while converting flag id to object id: %v", flagIDObjIDErr)
		return 0, ErrAddFlag
	}

	filterQuery := bson.D{{Key: "_id", Value: projectIDObjID}}
	updateQuery := bson.D{
		{
			Key:   "$push",
			Value: bson.D{{Key: "flags", Value: flagIDObjID}},
		},
		{
			Key:   "$set",
			Value: bson.D{{Key: "updated_at", Value: time.Now()}},
		},
	}

	updateResult, updateErr := p.coll.UpdateOne(ctx, filterQuery, updateQuery)
	if updateErr != nil {
		log.Printf(
			"error while updating project with the flag: %v",
			updateErr,
		)

		return 0, ErrAddFlag
	}

	//nolint:gosec // ModifiedCount can't be a negative number.
	return uint(updateResult.ModifiedCount), nil
}

func (p *MongoDataRepository) AddFlagSettings(
	ctx context.Context,
	projectID string,
	flagSettingIDs ...string,
) (uint, error) {
	projectIDObjID, projectIDObjIDErr := primitive.ObjectIDFromHex(projectID)
	if projectIDObjIDErr != nil {
		log.Printf("error while converting project id to object id: %v", projectIDObjIDErr)
		return 0, ErrAddFlagSetting
	}

	flagSettingIDObjIDs := make([]primitive.ObjectID, 0, len(flagSettingIDs))
	for _, flagSettingID := range flagSettingIDs {
		flagSettingIDObjID, flagSettingIDObjIDErr := primitive.ObjectIDFromHex(flagSettingID)
		if flagSettingIDObjIDErr != nil {
			log.Printf("error while converting flag setting id to object id: %v", flagSettingIDObjIDErr)
			return 0, ErrAddFlagSetting
		}

		flagSettingIDObjIDs = append(flagSettingIDObjIDs, flagSettingIDObjID)
	}

	filter := bson.D{{Key: "_id", Value: projectIDObjID}}
	update := bson.D{
		{
			Key: "$push",
			Value: bson.D{
				{
					Key:   "flag_settings",
					Value: bson.D{{Key: "$each", Value: flagSettingIDObjIDs}},
				},
			},
		},
		{
			Key:   "$set",
			Value: bson.D{{Key: "updated_at", Value: time.Now()}},
		},
	}

	updateResult, updateErr := p.coll.UpdateOne(ctx, filter, update)
	if updateErr != nil {
		log.Printf(
			"error while updating project with flag settings: %v",
			updateErr,
		)
		return 0, ErrAddFlagSetting
	}

	//nolint:gosec // ModifiedCount can't be a negative number.
	return uint(updateResult.ModifiedCount), nil
}

func setupProjectCollIndexes(
	ctx context.Context,
	coll *mongo.Collection,
) error {
	projectUserIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "name", Value: 1},
			{Key: "created_by", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	projectKeyIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "key", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	_, err := coll.Indexes().CreateMany(
		ctx,
		[]mongo.IndexModel{
			projectUserIndexModel,
			projectKeyIndexModel,
		},
	)

	return err
}

func mapStringSliceToObjectIDs(s []string) ([]primitive.ObjectID, error) {
	mappedObjectIDs := make([]primitive.ObjectID, 0, len(s))

	for _, id := range s {
		currentID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return []primitive.ObjectID{}, err
		}

		mappedObjectIDs = append(mappedObjectIDs, currentID)
	}

	return mappedObjectIDs, nil
}

func mapObjectIDsToStringSlice(s []primitive.ObjectID) []string {
	mappedStrings := make([]string, 0, len(s))

	for _, id := range s {
		mappedStrings = append(mappedStrings, id.Hex())
	}

	return mappedStrings
}

func NewProjectRepository(
	ctx context.Context,
	db *mongo.Database,
) (*MongoDataRepository, error) {
	coll := db.Collection(ProjectCollection)

	err := setupProjectCollIndexes(ctx, coll)
	if err != nil {
		return nil, err
	}

	return &MongoDataRepository{coll}, nil
}
