package flag

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

const flagCollection string = "flags"

// flagMongoModel is the MongoDB representation of the `Flag` structure.
type flagMongoModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	ProjectID primitive.ObjectID `bson:"project_id"`
	CreatedBy primitive.ObjectID `bson:"created_by"`
	CreatedAt time.Time          `bson:"created_at"`
}

type MongoDataRepository struct {
	coll *mongo.Collection
}

func (r *MongoDataRepository) Save(
	ctx context.Context,
	flag *Flag,
) (string, error) {
	projectIDObjID, projectIDErr := primitive.ObjectIDFromHex(flag.ProjectID)
	if projectIDErr != nil {
		log.Printf("could not convert project id to object id: %v", projectIDErr)
		return "", ErrCouldNotSave
	}

	userIDObjID, userIDErr := primitive.ObjectIDFromHex(flag.CreatedBy)
	if userIDErr != nil {
		log.Printf("could not convert user id to object id: %v", userIDErr)
		return "", ErrCouldNotSave
	}

	flagToSave := &flagMongoModel{
		Name:      flag.Name,
		ProjectID: projectIDObjID,
		CreatedBy: userIDObjID,
		CreatedAt: time.Now(),
	}

	saveResult, saveErr := r.coll.InsertOne(ctx, flagToSave)
	if saveErr != nil {
		if mongo.IsDuplicateKeyError(saveErr) {
			log.Printf("flag name is taken: %v", saveErr)
			return "", ErrNameTaken
		}

		log.Printf("error while saving flag: %v", saveErr)
		return "", ErrCouldNotSave
	}

	savedFlagID, ok := saveResult.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Printf("could not assert saved flag id as object id")
		return "", ErrCouldNotSave
	}

	return savedFlagID.Hex(), nil
}

func (r *MongoDataRepository) GetByID(
	ctx context.Context,
	flagID string,
) (*Flag, error) {
	flagIDObjID, flagIDErr := primitive.ObjectIDFromHex(flagID)
	if flagIDErr != nil {
		log.Printf("could not convert flag id to object id: %v", flagIDErr)
		return nil, ErrCouldNotFetch
	}

	query := bson.D{{Key: "_id", Value: flagIDObjID}}

	var decodedFlag flagMongoModel

	err := r.coll.FindOne(ctx, query).Decode(&decodedFlag)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("flag not found with id %v", flagID)
			return nil, ErrNotFound
		}

		return nil, ErrCouldNotFetch
	}

	return mapDecodedFlag(&decodedFlag), nil
}

func (r *MongoDataRepository) GetByNameAndProjectID(
	ctx context.Context,
	flagName string,
	projectID string,
) (*Flag, error) {
	projectIDObjID, projectIDErr := primitive.ObjectIDFromHex(projectID)
	if projectIDErr != nil {
		log.Printf("could not convert project id to object id: %v", projectIDErr)
		return nil, ErrCouldNotFetch
	}

	query := bson.D{
		{Key: "name", Value: flagName},
		{Key: "project_id", Value: projectIDObjID},
	}

	var decodedFlag flagMongoModel

	err := r.coll.FindOne(ctx, query).Decode(&decodedFlag)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("no flag found with name %q and project id %q", flagName, projectID)
			return nil, ErrNotFound
		}

		log.Printf("error while getting flag: %v", err)
		return nil, err
	}

	return mapDecodedFlag(&decodedFlag), nil
}

func mapDecodedFlag(decodedFlag *flagMongoModel) *Flag {
	return &Flag{
		ID:        decodedFlag.ID.Hex(),
		Name:      decodedFlag.Name,
		ProjectID: decodedFlag.ProjectID.Hex(),
		CreatedBy: decodedFlag.CreatedBy.Hex(),
		CreatedAt: decodedFlag.CreatedAt,
	}
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
) (*MongoDataRepository, error) {
	coll := db.Collection(flagCollection)

	err := setupCollIndexes(ctx, coll)
	if err != nil {
		return nil, err
	}

	return &MongoDataRepository{coll}, nil
}
