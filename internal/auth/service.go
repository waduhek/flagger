package auth

import (
	"context"

	"github.com/waduhek/flagger/proto/authpb"

	"github.com/waduhek/flagger/internal/hash"
	"github.com/waduhek/flagger/internal/logger"
	"github.com/waduhek/flagger/internal/user"
)

type Server struct {
	authpb.UnimplementedAuthServer
	userDataRepo user.DataRepository
	logger       logger.Logger
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
		s.logger.Error("could not generate password hash: %v", err)
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

	s.logger.Info("created new user, username %s, id %s", username, newUserID)

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
		s.logger.Error("incorrect credentials for user \"%s\"", username)
		return nil, ErrIncorrectUsernameOrPassword
	}

	token, err := CreateJWT(s.logger, fetchedUser.Username)
	if err != nil {
		s.logger.Error("error while generating jwt: %v", err)
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
		s.logger.Error("could not find claims from token")
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
		s.logger.Error("incorrect current password for resetting password")
		return nil, ErrIncorrectUsernameOrPassword
	}

	newPasswordHash, err := hash.GeneratePasswordHash(newPassword)
	if err != nil {
		s.logger.Error("error while hashing password: %v", err)
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

	s.logger.Info("changed password for user %q", username)
	return &authpb.ChangePasswordResponse{}, nil
}

// NewServer creates a new server for the auth service.
func NewServer(userDataRepo user.DataRepository, logger logger.Logger) *Server {
	server := &Server{userDataRepo: userDataRepo, logger: logger}

	return server
}
