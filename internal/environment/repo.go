package environment

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const environmentCollection string = "environments"

type environmentRepository struct {
	coll *mongo.Collection
}

func (r *environmentRepository) Save(
	ctx context.Context,
	environment *Environment,
) (*mongo.InsertOneResult, error) {
	return r.coll.InsertOne(ctx, environment)
}

func (r *environmentRepository) GetByID(
	ctx context.Context,
	id primitive.ObjectID,
) (*Environment, error) {
	query := bson.D{{Key: "_id", Value: id}}

	var environment Environment

	err := r.coll.FindOne(ctx, query).Decode(&environment)
	if err != nil {
		return nil, err
	}

	return &environment, nil
}

func (r *environmentRepository) GetByNameAndProjectID(
	ctx context.Context,
	environmentName string,
	projectID primitive.ObjectID,
) (*Environment, error) {
	query := bson.D{
		{Key: "name", Value: environmentName},
		{Key: "project_id", Value: projectID},
	}

	var environment Environment

	err := r.coll.FindOne(ctx, query).Decode(&environment)
	if err != nil {
		return nil, err
	}

	return &environment, nil
}

func setupCollIndexes(
	ctx context.Context,
	coll *mongo.Collection,
) error {
	environmentProjectIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "name", Value: 1},
			{Key: "project_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	_, err := coll.Indexes().CreateOne(ctx, environmentProjectIndexModel)

	return err
}

func NewEnvironmentRepository(
	ctx context.Context,
	db *mongo.Database,
) (*environmentRepository, error) {
	coll := db.Collection(environmentCollection)

	err := setupCollIndexes(ctx, coll)
	if err != nil {
		return nil, err
	}

	return &environmentRepository{coll}, nil
}
