package user

import "context"

// Password is the password information for a user. This struct will be inlined
// with the User struct.
type Password struct {
	Hash []byte `bson:"hash" json:"hash"`
	Salt []byte `bson:"salt" json:"salt"`
}

// User represents the details of a user.
type User struct {
	ID       string
	Name     string
	Email    string
	Username string
	Password *Password
}

// DataRepository is an interface to the operations that can be performed on the
// users collection.
type DataRepository interface {
	// Save saves the details of the user to the collection and returns the ID
	// of the saved user.
	Save(ctx context.Context, user *User) (string, error)

	// GetByUsername gets the details of the user by their username.
	GetByUsername(ctx context.Context, username string) (*User, error)

	// UpdatePassword updates the password of the user by the username and
	// returns the number of records that were updated.
	UpdatePassword(
		ctx context.Context,
		username string,
		password *Password,
	) (uint, error)
}
