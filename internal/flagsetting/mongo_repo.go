package flagsetting

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/waduhek/flagger/internal/logger"
)

const flagSettingCollection string = "flag_settings"

// flagSettingMongoModel is the MongoDB representation of the `FlagSetting`
// struct.
type flagSettingMongoModel struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	ProjectID     primitive.ObjectID `bson:"project_id"`
	EnvironmentID primitive.ObjectID `bson:"environment_id"`
	FlagID        primitive.ObjectID `bson:"flag_id"`
	IsActive      bool               `bson:"is_active"`
	CreatedAt     time.Time          `bson:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at"`
}

type MongoDataRepository struct {
	coll   *mongo.Collection
	logger logger.Logger
}

func (r *MongoDataRepository) Save(
	ctx context.Context,
	flagSetting *FlagSetting,
) (string, error) {
	flagSettingToSave, createErr := r.createNewFlagSetting(flagSetting)
	if createErr != nil {
		return "", createErr
	}

	saveResult, saveErr := r.coll.InsertOne(ctx, flagSettingToSave)
	if saveErr != nil {
		if mongo.IsDuplicateKeyError(saveErr) {
			r.logger.Error("got duplicate key error when saving flag setting: %v", saveErr)
		} else {
			r.logger.Error("error while saving flag setting: %v", saveErr)
		}

		return "", ErrCouldNotSave
	}

	flagSettingID, ok := saveResult.InsertedID.(primitive.ObjectID)
	if !ok {
		r.logger.Error("could not assert saved flag setting id as object id")
		return "", ErrCouldNotSave
	}

	return flagSettingID.Hex(), nil
}

func (r *MongoDataRepository) SaveMany(
	ctx context.Context,
	flagSettings []FlagSetting,
) ([]string, error) {
	toInsert := make([]interface{}, 0, len(flagSettings))
	for _, setting := range flagSettings {
		mongoModel, err := r.createNewFlagSetting(&setting)
		if err != nil {
			return []string{}, err
		}

		toInsert = append(toInsert, mongoModel)
	}

	saveResult, saveErr := r.coll.InsertMany(ctx, toInsert)
	if saveErr != nil {
		if mongo.IsDuplicateKeyError(saveErr) {
			r.logger.Error("got duplicate key error while saving flag settings: %v", saveErr)
		} else {
			r.logger.Error("error while saving flag settings: %v", saveErr)
		}

		return []string{}, ErrCouldNotSave
	}

	savedIDs := make([]string, 0, len(saveResult.InsertedIDs))
	for _, id := range saveResult.InsertedIDs {
		flagSettingID, ok := id.(primitive.ObjectID)
		if !ok {
			r.logger.Error("could not assert save flag setting id as object id")
			return []string{}, ErrCouldNotSave
		}

		savedIDs = append(savedIDs, flagSettingID.Hex())
	}

	return savedIDs, nil
}

func (r *MongoDataRepository) createNewFlagSetting(flagSetting *FlagSetting) (*flagSettingMongoModel, error) {
	projectIDObjID, projectIDErr := primitive.ObjectIDFromHex(flagSetting.ProjectID)
	if projectIDErr != nil {
		r.logger.Error("could not convert project id to object id: %v", projectIDErr)
		return nil, ErrCouldNotSave
	}

	environmentIDObjID, environmentIDErr := primitive.ObjectIDFromHex(flagSetting.EnvironmentID)
	if environmentIDErr != nil {
		r.logger.Error("could not convert environment id to object id: %v", environmentIDErr)
		return nil, ErrCouldNotSave
	}

	flagIDObjID, flagIDErr := primitive.ObjectIDFromHex(flagSetting.FlagID)
	if flagIDErr != nil {
		r.logger.Error("could not convert flag id to object id: %v", flagIDErr)
		return nil, ErrCouldNotSave
	}

	flagSettingToSave := &flagSettingMongoModel{
		ProjectID:     projectIDObjID,
		EnvironmentID: environmentIDObjID,
		FlagID:        flagIDObjID,
		IsActive:      flagSetting.IsActive,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	return flagSettingToSave, nil
}

func (r *MongoDataRepository) Get(
	ctx context.Context,
	projectID string,
	environmentID string,
	flagID string,
) (*FlagSetting, error) {
	projectIDObjID, projectIDErr := primitive.ObjectIDFromHex(projectID)
	if projectIDErr != nil {
		r.logger.Error("could not convert project id to object id: %v", projectIDErr)
		return nil, ErrCouldNotGet
	}

	environmentIDObjID, environmentIDErr := primitive.ObjectIDFromHex(environmentID)
	if environmentIDErr != nil {
		r.logger.Error("could not convert environment id to object id: %v", environmentIDErr)
		return nil, ErrCouldNotGet
	}

	flagIDObjID, flagIDErr := primitive.ObjectIDFromHex(flagID)
	if flagIDErr != nil {
		r.logger.Error("could not convert flag id to object id: %v", flagIDErr)
		return nil, ErrCouldNotGet
	}

	query := bson.D{
		{Key: "project_id", Value: projectIDObjID},
		{Key: "environment_id", Value: environmentIDObjID},
		{Key: "flag_id", Value: flagIDObjID},
	}

	var decodedFlagSetting flagSettingMongoModel

	err := r.coll.FindOne(ctx, query).Decode(&decodedFlagSetting)
	if err != nil {
		return nil, err
	}

	flagSetting := &FlagSetting{
		ID:            decodedFlagSetting.ID.Hex(),
		ProjectID:     decodedFlagSetting.ProjectID.Hex(),
		EnvironmentID: decodedFlagSetting.EnvironmentID.Hex(),
		FlagID:        decodedFlagSetting.FlagID.Hex(),
		IsActive:      decodedFlagSetting.IsActive,
		CreatedAt:     decodedFlagSetting.CreatedAt,
		UpdatedAt:     decodedFlagSetting.UpdatedAt,
	}

	return flagSetting, err
}

func (r *MongoDataRepository) UpdateIsActive(
	ctx context.Context,
	projectID string,
	environmentID string,
	flagID string,
	isActive bool,
) (uint, error) {
	projectIDObjID, projectIDErr := primitive.ObjectIDFromHex(projectID)
	if projectIDErr != nil {
		r.logger.Error("could not convert project id to object id: %v", projectIDErr)
		return 0, ErrStatusUpdate
	}

	environmentIDObjID, environmentIDErr := primitive.ObjectIDFromHex(environmentID)
	if environmentIDErr != nil {
		r.logger.Error("could not convert environment id to object id: %v", environmentIDErr)
		return 0, ErrStatusUpdate
	}

	flagIDObjID, flagIDErr := primitive.ObjectIDFromHex(flagID)
	if flagIDErr != nil {
		r.logger.Error("could not convert flag id to object id: %v", flagIDErr)
		return 0, ErrStatusUpdate
	}

	filter := bson.D{
		{Key: "project_id", Value: projectIDObjID},
		{Key: "environment_id", Value: environmentIDObjID},
		{Key: "flag_id", Value: flagIDObjID},
	}

	update := bson.D{{
		Key: "$set", Value: bson.D{{Key: "is_active", Value: isActive}},
	}}

	updateResult, updateErr := r.coll.UpdateOne(ctx, filter, update)
	if updateErr != nil {
		r.logger.Error("could not update status of flag setting: %v", updateErr)
		return 0, ErrStatusUpdate
	}

	//nolint:gosec // ModifiedCount can't be a negative number.
	return uint(updateResult.ModifiedCount), nil
}

func setupIndexes(ctx context.Context, coll *mongo.Collection) error {
	projectEnvFlagIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "project_id", Value: 1},
			{Key: "environment_id", Value: 1},
			{Key: "flag_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	_, err := coll.Indexes().CreateOne(ctx, projectEnvFlagIndexModel)

	return err
}

// NewFlagSettingRepository creates a new repository that implements the
// `flagsetting.FlagSettingRepository` interface.
func NewFlagSettingRepository(
	ctx context.Context,
	db *mongo.Database,
	logger logger.Logger,
) (*MongoDataRepository, error) {
	coll := db.Collection(flagSettingCollection)

	err := setupIndexes(ctx, coll)
	if err != nil {
		return nil, err
	}

	return &MongoDataRepository{coll, logger}, nil
}
