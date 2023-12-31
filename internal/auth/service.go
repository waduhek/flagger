package auth

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
		return nil, status.Error(codes.Internal, "could not generate a hash")
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
			return nil, status.Error(codes.AlreadyExists, "username is taken")
		}
		log.Printf("could not save user details: %v", err)
		return nil, status.Error(codes.Internal, "could not save the user details")
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
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		log.Printf("could not get details of user by username: %v", err)
		return nil, status.Error(codes.Internal, "could not get details of the username")
	}

	if !hash.VerifyPasswordHash(
		req.Password,
		user.Password.Hash,
		user.Password.Salt,
	) {
		log.Printf("incorrect credentials for user \"%s\"", req.Username)
		return nil, status.Error(codes.Unauthenticated, "incorrect username or password")
	}

	token, err := CreateJWT(user.Username)
	if err != nil {
		log.Printf("error while generating jwt: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
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
		return nil, status.Error(codes.Internal, "could not get details of the username")
	}

	if !hash.VerifyPasswordHash(
		req.CurrentPassword,
		fetchedUser.Password.Hash,
		fetchedUser.Password.Salt,
	) {
		log.Printf("incorrect current password for resetting password")
		return nil, status.Error(codes.Unauthenticated, "incorrect current password")
	}

	newPasswordHash, err := hash.GeneratePasswordHash(req.NewPassword)
	if err != nil {
		log.Printf("error while hashing password: %v", err)
		return nil, status.Error(codes.Internal, "could not hash new password")
	}

	password := user.Password{
		Hash: newPasswordHash.Hash,
		Salt: newPasswordHash.Salt,
	}

	_, updateErr := s.userRepo.UpdatePassword(ctx, username, &password)
	if updateErr != nil {
		log.Printf("error while saving new password: %v", updateErr)
		return nil, status.Error(codes.Internal, "could not save password")
	}

	log.Printf("changed password for user %q", username)
	return &authpb.ChangePasswordResponse{}, nil
}

// NewAuthServer creates a new server for the auth service.
func NewAuthServer(userRepo user.UserRepository) *AuthServer {
	server := &AuthServer{userRepo: userRepo}

	return server
}
