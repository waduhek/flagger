package auth

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/waduhek/flagger/proto/authpb"

	"github.com/waduhek/flagger/internal/hash"
	"github.com/waduhek/flagger/internal/user"
)

type AuthServer struct {
	authpb.UnimplementedAuthServer
	userRepo user.UserRepository
}

func (s *AuthServer) CreateNewUser(
	ctx context.Context,
	req *authpb.CreateNewUserRequest,
) (*authpb.CreateNewUserResponse, error) {
	username := req.GetUsername()
	password := req.GetPassword()
	name := req.GetName()
	email := req.GetEmail()

	passwordHash, err := hash.GeneratePasswordHash(password)
	if err != nil {
		log.Printf("could not generate password hash: %v", err)
		return nil, hash.ErrGenPasswordHash
	}

	newUser := user.User{
		Username: username,
		Name:     name,
		Email:    email,
		Password: user.Password{
			Hash: passwordHash.Hash,
			Salt: passwordHash.Salt,
		},
	}

	newUserResult, err := s.userRepo.Save(ctx, &newUser)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			log.Printf("a user with username %q already exists", username)
			return nil, user.ErrUsernameTaken
		}

		log.Printf("could not save user details: %v", err)
		return nil, user.ErrNotSaved
	}

	log.Printf(
		"created new user, username %s, id %s",
		username, newUserResult.InsertedID,
	)

	response := authpb.CreateNewUserResponse{}

	return &response, nil
}

func (s *AuthServer) Login(
	ctx context.Context,
	req *authpb.LoginRequest,
) (*authpb.LoginResponse, error) {
	username := req.GetUsername()
	password := req.GetPassword()

	fetchedUser, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		log.Printf("could not get details of user by username: %v", err)
		return nil, user.ErrCouldNotFetch
	}

	if !hash.VerifyPasswordHash(
		password,
		fetchedUser.Password.Hash,
		fetchedUser.Password.Salt,
	) {
		log.Printf("incorrect credentials for user \"%s\"", username)
		return nil, ErrIncorrectUsernameOrPassword
	}

	token, err := CreateJWT(fetchedUser.Username)
	if err != nil {
		log.Printf("error while generating jwt: %v", err)
		return nil, err
	}

	response := &authpb.LoginResponse{Token: token}

	return response, nil
}

func (s *AuthServer) ChangePassword(
	ctx context.Context,
	req *authpb.ChangePasswordRequest,
) (*authpb.ChangePasswordResponse, error) {
	claims, ok := ClaimsFromContext(ctx)
	if !ok {
		log.Printf("could not find claims from token")
		return nil, ErrNoTokenClaims
	}

	username := claims.Subject

	fetchedUser, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		log.Printf("error while fetching user %q: %v", username, err)
		return nil, user.ErrCouldNotFetch
	}

	currentPassword := req.GetCurrentPassword()
	newPassword := req.GetNewPassword()

	if !hash.VerifyPasswordHash(
		currentPassword,
		fetchedUser.Password.Hash,
		fetchedUser.Password.Salt,
	) {
		log.Printf("incorrect current password for resetting password")
		return nil, ErrIncorrectUsernameOrPassword
	}

	newPasswordHash, err := hash.GeneratePasswordHash(newPassword)
	if err != nil {
		log.Printf("error while hashing password: %v", err)
		return nil, hash.ErrGenPasswordHash
	}

	password := user.Password{
		Hash: newPasswordHash.Hash,
		Salt: newPasswordHash.Salt,
	}

	_, updateErr := s.userRepo.UpdatePassword(ctx, username, &password)
	if updateErr != nil {
		log.Printf("error while saving new password: %v", updateErr)
		return nil, user.ErrPasswordUpdate
	}

	log.Printf("changed password for user %q", username)
	return &authpb.ChangePasswordResponse{}, nil
}

// NewAuthServer creates a new server for the auth service.
func NewAuthServer(userRepo user.UserRepository) *AuthServer {
	server := &AuthServer{userRepo: userRepo}

	return server
}
