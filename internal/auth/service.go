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
	passwordHash, err := hash.GeneratePasswordHash(req.Password)
	if err != nil {
		log.Printf("could not generate password hash: %v", err)
		return nil, hash.EHashGenPasswordHash
	}

	newUser := user.User{
		Username: req.Username,
		Name:     req.Name,
		Email:    req.Email,
		Password: user.Password{
			Hash: passwordHash.Hash,
			Salt: passwordHash.Salt,
		},
	}

	newUserResult, err := s.userRepo.Save(ctx, &newUser)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			log.Printf("a user with username %q already exists", req.Username)
			return nil, user.EUsernameTaken
		}

		log.Printf("could not save user details: %v", err)
		return nil, user.EUserNotSaved
	}

	log.Printf(
		"created new user, username %s, id %s",
		req.Username, newUserResult.InsertedID,
	)

	response := authpb.CreateNewUserResponse{}

	return &response, nil
}

func (s *AuthServer) Login(
	ctx context.Context,
	req *authpb.LoginRequest,
) (*authpb.LoginResponse, error) {
	fetchedUser, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		log.Printf("could not get details of user by username: %v", err)
		return nil, user.ECouldNotFetchUser
	}

	if !hash.VerifyPasswordHash(
		req.Password,
		fetchedUser.Password.Hash,
		fetchedUser.Password.Salt,
	) {
		log.Printf("incorrect credentials for user \"%s\"", req.Username)
		return nil, EIncorrectUsernameOrPassword
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
		return nil, ENoTokenClaims
	}

	username := claims.Subject

	fetchedUser, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		log.Printf("error while fetching user %q: %v", username, err)
		return nil, user.ECouldNotFetchUser
	}

	if !hash.VerifyPasswordHash(
		req.CurrentPassword,
		fetchedUser.Password.Hash,
		fetchedUser.Password.Salt,
	) {
		log.Printf("incorrect current password for resetting password")
		return nil, EIncorrectUsernameOrPassword
	}

	newPasswordHash, err := hash.GeneratePasswordHash(req.NewPassword)
	if err != nil {
		log.Printf("error while hashing password: %v", err)
		return nil, hash.EHashGenPasswordHash
	}

	password := user.Password{
		Hash: newPasswordHash.Hash,
		Salt: newPasswordHash.Salt,
	}

	_, updateErr := s.userRepo.UpdatePassword(ctx, username, &password)
	if updateErr != nil {
		log.Printf("error while saving new password: %v", updateErr)
		return nil, user.EPasswordUpdate
	}

	log.Printf("changed password for user %q", username)
	return &authpb.ChangePasswordResponse{}, nil
}

// NewAuthServer creates a new server for the auth service.
func NewAuthServer(userRepo user.UserRepository) *AuthServer {
	server := &AuthServer{userRepo: userRepo}

	return server
}
