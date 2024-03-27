package user

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const userCollection string = "users"

type MongoDataRepository struct {
	coll *mongo.Collection
}

func (u *MongoDataRepository) Save(
	ctx context.Context,
	user *User,
) (*mongo.InsertOneResult, error) {
	result, err := u.coll.InsertOne(ctx, user)

	return result, err
}

func (u *MongoDataRepository) GetByUsername(
	ctx context.Context,
	username string,
) (*User, error) {
	query := bson.D{{Key: "username", Value: username}}

	var user User
	err := u.coll.FindOne(ctx, query).Decode(&user)

	return &user, err
}

func (u *MongoDataRepository) UpdatePassword(
	ctx context.Context,
	username string,
	password *Password,
) (*mongo.UpdateResult, error) {
	filter := bson.D{{Key: "username", Value: username}}
	update := bson.D{{
		Key: "$set",
		Value: bson.D{
			{Key: "hash", Value: password.Hash},
			{Key: "salt", Value: password.Salt},
		},
	}}

	return u.coll.UpdateOne(ctx, filter, update)
}

func setupUserCollIndexes(ctx context.Context, coll *mongo.Collection) error {
	usernameIndexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := coll.Indexes().CreateOne(ctx, usernameIndexModel)

	return err
}

// NewUserRepository creates a new repository for performing database operations
// in the users collection.
func NewUserRepository(
	ctx context.Context,
	db *mongo.Database,
) (*MongoDataRepository, error) {
	coll := db.Collection(userCollection)

	err := setupUserCollIndexes(ctx, coll)
	if err != nil {
		return nil, err
	}

	repo := &MongoDataRepository{coll}
	return repo, nil
}
