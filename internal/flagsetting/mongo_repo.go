package flagsetting

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const flagSettingCollection string = "flag_settings"

type MongoDataRepository struct {
	coll *mongo.Collection
}

func (r *MongoDataRepository) Save(
	ctx context.Context,
	flagSetting *FlagSetting,
) (*mongo.InsertOneResult, error) {
	return r.coll.InsertOne(ctx, flagSetting)
}

func (r *MongoDataRepository) SaveMany(
	ctx context.Context,
	flagSettings []FlagSetting,
) (*mongo.InsertManyResult, error) {
	var toInsert []interface{}
	for _, setting := range flagSettings {
		toInsert = append(toInsert, setting)
	}

	return r.coll.InsertMany(ctx, toInsert)
}

func (r *MongoDataRepository) Get(
	ctx context.Context,
	projectID primitive.ObjectID,
	environmentID primitive.ObjectID,
	flagID primitive.ObjectID,
) (*FlagSetting, error) {
	query := bson.D{
		{Key: "project_id", Value: projectID},
		{Key: "environment_id", Value: environmentID},
		{Key: "flag_id", Value: flagID},
	}

	var flagSetting FlagSetting

	err := r.coll.FindOne(ctx, query).Decode(&flagSetting)
	if err != nil {
		return nil, err
	}

	return &flagSetting, err
}

func (r *MongoDataRepository) UpdateIsActive(
	ctx context.Context,
	projectID primitive.ObjectID,
	environmentID primitive.ObjectID,
	flagID primitive.ObjectID,
	isActive bool,
) (*mongo.UpdateResult, error) {
	filter := bson.D{
		{Key: "project_id", Value: projectID},
		{Key: "environment_id", Value: environmentID},
		{Key: "flag_id", Value: flagID},
	}

	update := bson.D{{
		Key: "$set", Value: bson.D{{Key: "is_active", Value: isActive}},
	}}

	return r.coll.UpdateOne(ctx, filter, update)
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
) (*MongoDataRepository, error) {
	coll := db.Collection(flagSettingCollection)

	err := setupIndexes(ctx, coll)
	if err != nil {
		return nil, err
	}

	return &MongoDataRepository{coll}, nil
}
