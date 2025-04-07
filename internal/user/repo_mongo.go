package user

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/waduhek/flagger/internal/logger"
)

const userCollection string = "users"

// userMongoModel is the MongoDB representation of the `User` struct.
type userMongoModel struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username"`
	Name     string             `bson:"name"`
	Email    string             `bson:"email"`
	Password *Password          `bson:"password,inline"`
}

type MongoDataRepository struct {
	coll   *mongo.Collection
	logger logger.Logger
}

func (u *MongoDataRepository) Save(
	ctx context.Context,
	user *User,
) (string, error) {
	userToAdd := &userMongoModel{
		Username: user.Username,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}

	result, err := u.coll.InsertOne(ctx, userToAdd)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			u.logger.Error("a user with username %q already exists", user.Username)
			return "", ErrUsernameTaken
		}

		u.logger.Error("could not save user details: %v", err)
		return "", ErrNotSaved
	}

	userID, userIDOk := result.InsertedID.(primitive.ObjectID)
	if !userIDOk {
		u.logger.Error("could not cast mongodb user id as objectid")
		return "", ErrNotSaved
	}

	return userID.Hex(), nil
}

func (u *MongoDataRepository) GetByUsername(
	ctx context.Context,
	username string,
) (*User, error) {
	query := bson.D{{Key: "username", Value: username}}

	var decodedUser userMongoModel
	err := u.coll.FindOne(ctx, query).Decode(&decodedUser)
	if err != nil {
		u.logger.Error("error while fetching user %q: %v", username, err)
		return nil, ErrCouldNotFetch
	}

	user := &User{
		ID:       decodedUser.ID.Hex(),
		Username: decodedUser.Username,
		Name:     decodedUser.Name,
		Email:    decodedUser.Email,
		Password: decodedUser.Password,
	}

	return user, nil
}

func (u *MongoDataRepository) UpdatePassword(
	ctx context.Context,
	username string,
	password *Password,
) (uint, error) {
	filter := bson.D{{Key: "username", Value: username}}
	update := bson.D{{
		Key: "$set",
		Value: bson.D{
			{Key: "hash", Value: password.Hash},
			{Key: "salt", Value: password.Salt},
		},
	}}

	updateResult, updateErr := u.coll.UpdateOne(ctx, filter, update)
	if updateErr != nil {
		u.logger.Error("error while saving new password: %v", updateErr)
		return 0, ErrPasswordUpdate
	}

	//nolint:gosec // ModifiedCount can't be a negative number.
	return uint(updateResult.ModifiedCount), nil
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
	logger logger.Logger,
) (*MongoDataRepository, error) {
	coll := db.Collection(userCollection)

	err := setupUserCollIndexes(ctx, coll)
	if err != nil {
		return nil, err
	}

	repo := &MongoDataRepository{coll, logger}
	return repo, nil
}
