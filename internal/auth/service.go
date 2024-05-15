package auth

import (
	"context"
	"log"

	"github.com/waduhek/flagger/proto/authpb"

	"github.com/waduhek/flagger/internal/hash"
	"github.com/waduhek/flagger/internal/user"
)

type Server struct {
	authpb.UnimplementedAuthServer
	userDataRepo user.DataRepository
}

func (s *Server) CreateNewUser(
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
		Password: &user.Password{
			Hash: passwordHash.Hash,
			Salt: passwordHash.Salt,
		},
	}

	newUserID, err := s.userDataRepo.Save(ctx, &newUser)
	if err != nil {
		return nil, err
	}

	log.Printf("created new user, username %s, id %s", username, newUserID)

	response := authpb.CreateNewUserResponse{}

	return &response, nil
}

func (s *Server) Login(
	ctx context.Context,
	req *authpb.LoginRequest,
) (*authpb.LoginResponse, error) {
	username := req.GetUsername()
	password := req.GetPassword()

	fetchedUser, err := s.userDataRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
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

func (s *Server) ChangePassword(
	ctx context.Context,
	req *authpb.ChangePasswordRequest,
) (*authpb.ChangePasswordResponse, error) {
	claims, ok := ClaimsFromContext(ctx)
	if !ok {
		log.Printf("could not find claims from token")
		return nil, ErrNoTokenClaims
	}

	username := claims.Subject

	fetchedUser, err := s.userDataRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
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

	_, updateErr := s.userDataRepo.UpdatePassword(ctx, username, &password)
	if updateErr != nil {
		return nil, updateErr
	}

	log.Printf("changed password for user %q", username)
	return &authpb.ChangePasswordResponse{}, nil
}

// NewServer creates a new server for the auth service.
func NewServer(userDataRepo user.DataRepository) *Server {
	server := &Server{userDataRepo: userDataRepo}

	return server
}
