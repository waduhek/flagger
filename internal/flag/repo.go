package flag

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const flagCollection string = "flags"

type flagRepository struct {
	coll *mongo.Collection
}

func (r *flagRepository) Save(
	ctx context.Context,
	flag *Flag,
) (*mongo.InsertOneResult, error) {
	return r.coll.InsertOne(ctx, flag)
}

func (r *flagRepository) GetByID(
	ctx context.Context,
	flagID primitive.ObjectID,
) (*Flag, error) {
	query := bson.D{{Key: "_id", Value: flagID}}

	var flag Flag

	err := r.coll.FindOne(ctx, query).Decode(&flag)
	if err != nil {
		return nil, err
	}

	return &flag, nil
}

func (r *flagRepository) GetByNameAndProjectID(
	ctx context.Context,
	flagName string,
	projectID primitive.ObjectID,
) (*Flag, error) {
	query := bson.D{
		{Key: "name", Value: flagName},
		{Key: "project_id", Value: projectID},
	}

	var flag Flag

	err := r.coll.FindOne(ctx, query).Decode(&flag)
	if err != nil {
		return nil, err
	}

	return &flag, nil
}

func setupCollIndexes(ctx context.Context, coll *mongo.Collection) error {
	flagProjectIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "name", Value: 1},
			{Key: "project_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	_, err := coll.Indexes().CreateOne(ctx, flagProjectIndexModel)

	return err
}

func NewFlagRepository(
	ctx context.Context,
	db *mongo.Database,
) (*flagRepository, error) {
	coll := db.Collection(flagCollection)

	err := setupCollIndexes(ctx, coll)
	if err != nil {
		return nil, err
	}

	return &flagRepository{coll}, nil
}
