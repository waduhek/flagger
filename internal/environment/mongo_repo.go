package environment

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// environmentMongoModel is the MongoDB representation of the `Environment`
// struct.
type environmentMongoModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	ProjectID primitive.ObjectID `bson:"project_id"`
	CreatedBy primitive.ObjectID `bson:"created_by"`
	CreatedAt time.Time          `bson:"created_at"`
}

const environmentCollection string = "environments"

type MongoDataRepository struct {
	coll *mongo.Collection
}

func (r *MongoDataRepository) Save(
	ctx context.Context,
	environment *Environment,
) (string, error) {
	projectIDObjID, projectIDObjIDErr := primitive.ObjectIDFromHex(environment.ProjectID)
	if projectIDObjIDErr != nil {
		log.Printf("could not convert project id to object id: %v", projectIDObjIDErr)
		return "", ErrCouldNotSave
	}

	createdByObjID, createdByObjIDErr := primitive.ObjectIDFromHex(environment.CreatedBy)
	if createdByObjIDErr != nil {
		log.Printf("could not convert created by id to object id: %v", createdByObjIDErr)
		return "", ErrCouldNotSave
	}

	environmentToAdd := &environmentMongoModel{
		Name:      environment.Name,
		ProjectID: projectIDObjID,
		CreatedBy: createdByObjID,
		CreatedAt: time.Now(),
	}

	saveResult, saveErr := r.coll.InsertOne(ctx, environmentToAdd)
	if saveErr != nil {
		if mongo.IsDuplicateKeyError(saveErr) {
			log.Printf("an environment with name %q exists", environment.Name)
			return "", ErrNameTaken
		}

		log.Printf("error occurred while saving environment: %v", saveErr)
		return "", ErrCouldNotSave
	}

	environmentID, ok := saveResult.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Printf("could not assert saved environment id as object id")
		return "", ErrCouldNotSave
	}

	return environmentID.Hex(), nil
}

func (r *MongoDataRepository) GetByID(
	ctx context.Context,
	id string,
) (*Environment, error) {
	environmentID, environmentIDErr := primitive.ObjectIDFromHex(id)
	if environmentIDErr != nil {
		log.Printf("error while converting environment id to object id: %v", environmentIDErr)
		return nil, ErrCouldNotFetch
	}

	query := bson.D{{Key: "_id", Value: environmentID}}

	var decodedEnvironment environmentMongoModel

	err := r.coll.FindOne(ctx, query).Decode(&decodedEnvironment)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("no environment with the id %q was found: %v", id, err)
			return nil, ErrNotFound
		}

		log.Printf("error occurred while fetching environment: %v", err)
		return nil, ErrCouldNotFetch
	}

	return mapMongoModelToStruct(&decodedEnvironment), nil
}

func (r *MongoDataRepository) GetByNameAndProjectID(
	ctx context.Context,
	environmentName string,
	projectID string,
) (*Environment, error) {
	projectIDObjID, projectIDObjIDErr := primitive.ObjectIDFromHex(projectID)
	if projectIDObjIDErr != nil {
		log.Printf("could not convert project id to object id: %v", projectIDObjIDErr)
		return nil, ErrCouldNotFetch
	}

	query := bson.D{
		{Key: "name", Value: environmentName},
		{Key: "project_id", Value: projectIDObjID},
	}

	var decodedEnvironment environmentMongoModel

	err := r.coll.FindOne(ctx, query).Decode(&decodedEnvironment)
	if err != nil {
		return nil, err
	}

	return mapMongoModelToStruct(&decodedEnvironment), nil
}

func mapMongoModelToStruct(decodedEnvironment *environmentMongoModel) *Environment {
	return &Environment{
		ID:        decodedEnvironment.ID.Hex(),
		Name:      decodedEnvironment.Name,
		ProjectID: decodedEnvironment.ID.Hex(),
		CreatedBy: decodedEnvironment.CreatedBy.Hex(),
		CreatedAt: decodedEnvironment.CreatedAt,
	}
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
) (*MongoDataRepository, error) {
	coll := db.Collection(environmentCollection)

	err := setupCollIndexes(ctx, coll)
	if err != nil {
		return nil, err
	}

	return &MongoDataRepository{coll}, nil
}
