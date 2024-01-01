package project

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/waduhek/flagger/proto/projectpb"

	"github.com/waduhek/flagger/internal/auth"
	"github.com/waduhek/flagger/internal/user"
)

type ProjectServer struct {
	projectpb.UnimplementedProjectServer
	projectRepo ProjectRepository
	userRepo    user.UserRepository
}

func (p *ProjectServer) CreateNewProject(
	ctx context.Context,
	req *projectpb.CreateNewProjectRequest,
) (*projectpb.CreateNewProjectResponse, error) {
	jwtClaims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		log.Printf("could not find claims from token")
		return nil, auth.ENoTokenClaims
	}

	username := jwtClaims.Subject

	fetchedUser, err := p.userRepo.GetByUsername(ctx, username)
	if err != nil {
		log.Printf("error while fetching user %q: %v", username, err)
		return nil, user.ECouldNotFetchUser
	}

	newProject := Project{
		Name:      req.ProjectName,
		CreatedBy: fetchedUser.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, projectErr := p.projectRepo.Save(ctx, &newProject)
	if projectErr != nil {
		if mongo.IsDuplicateKeyError(projectErr) {
			log.Printf(
				"a project with name %q exists for user %q",
				req.ProjectName,
				username,
			)
			return nil,
				status.Error(codes.AlreadyExists, "project already exists")
		}

		log.Printf("error while creating new project: %v", projectErr)
		return nil,
			status.Error(codes.Internal, "could not create a new project")
	}

	return &projectpb.CreateNewProjectResponse{}, nil
}

// NewProjectServer creates a new server for the project service.
func NewProjectServer(
	projectRepo ProjectRepository,
	userRepo user.UserRepository,
) *ProjectServer {
	return &ProjectServer{projectRepo: projectRepo, userRepo: userRepo}
}
