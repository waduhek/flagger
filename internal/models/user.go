package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Password is the password information for a user. This struct will be inlined
// with the User struct.
type Password struct {
	Hash []byte `bson:"hash" json:"hash"`
	Salt []byte `bson:"salt" json:"salt"`
}

// The details of the user.
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name" json:"name"`
	Email    string             `bson:"email" json:"email"`
	Username string             `bson:"username" json:"username"`
	Password Password           `bson:"inline" json:"password"`
}

// UserRepository is an interface to the operations that can be performed on the
// users collection.
type UserRepository interface {
	// Save saves the details of the user to the collection.
	Save(ctx context.Context, user *User) (*mongo.InsertOneResult, error)

	// GetByUsername fetches the details of the user by their username.
	GetByUsername(ctx context.Context, username string) (*User, error)

	// UpdatePassword updates the password of the user by the username.
	UpdatePassword(
		ctx context.Context,
		username string,
		password *Password,
	) (*mongo.UpdateResult, error)
}
