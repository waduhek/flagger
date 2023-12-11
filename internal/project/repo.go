package project

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const projectCollection string = "projects"

type projectRepository struct {
	coll *mongo.Collection
}

func (p *projectRepository) Save(
	ctx context.Context,
	project *Project,
) (*mongo.InsertOneResult, error) {
	return p.coll.InsertOne(ctx, project)
}

func (p *projectRepository) GetByNameAndUserID(
	ctx context.Context,
	projectName string,
	userID primitive.ObjectID,
) (*Project, error) {
	query := bson.D{
		{Key: "created_by", Value: userID},
		{Key: "name", Value: projectName},
	}

	var project Project
	err := p.coll.FindOne(ctx, query).Decode(&project)

	return &project, err
}

func (p *projectRepository) AddEnvironment(
	ctx context.Context,
	projectID primitive.ObjectID,
	environmentID primitive.ObjectID,
) (*mongo.UpdateResult, error) {
	filterQuery := bson.D{
		{Key: "_id", Value: projectID},
	}
	updateQuery := bson.D{
		{
			Key:   "$push",
			Value: bson.D{{Key: "environments", Value: environmentID}},
		},
		{
			Key:   "$set",
			Value: bson.D{{Key: "updated_at", Value: time.Now()}},
		},
	}

	return p.coll.UpdateOne(ctx, filterQuery, updateQuery)
}

func (p *projectRepository) AddFlag(
	ctx context.Context,
	projectID primitive.ObjectID,
	flagID primitive.ObjectID,
) (*mongo.UpdateResult, error) {
	filterQuery := bson.D{{Key: "_id", Value: projectID}}
	updateQuery := bson.D{
		{
			Key:   "$push",
			Value: bson.D{{Key: "flags", Value: flagID}},
		},
		{
			Key:   "$set",
			Value: bson.D{{Key: "updated_at", Value: time.Now()}},
		},
	}

	return p.coll.UpdateOne(ctx, filterQuery, updateQuery)
}

func (p *projectRepository) AddFlagSettings(
	ctx context.Context,
	projectID primitive.ObjectID,
	flagSettingIDs ...primitive.ObjectID,
) (*mongo.UpdateResult, error) {
	filter := bson.D{{Key: "_id", Value: projectID}}
	update := bson.D{
		{
			Key: "$push",
			Value: bson.D{
				{
					Key: "flag_settings",
					Value: bson.D{{Key: "$each", Value: flagSettingIDs}},
				},
			},
		},
		{
			Key:   "$set",
			Value: bson.D{{Key: "updated_at", Value: time.Now()}},
		},
	}

	return p.coll.UpdateOne(ctx, filter, update)
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
	_, err := coll.Indexes().CreateOne(ctx, projectUserIndexModel)

	return err
}

func NewProjectRepository(
	ctx context.Context,
	db *mongo.Database,
) (*projectRepository, error) {
	coll := db.Collection(projectCollection)

	err := setupProjectCollIndexes(ctx, coll)
	if err != nil {
		return nil, err
	}

	return &projectRepository{coll}, nil
}
