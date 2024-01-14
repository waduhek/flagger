package project

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/waduhek/flagger/proto/projectpb"

	"github.com/waduhek/flagger/internal/auth"
	"github.com/waduhek/flagger/internal/user"
)

// The number of characters in the project key.
const projectKeyLen uint = 32

// The number of retries for saving the project before an error is returned.
const projectSaveRetries uint = 5

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
		Key:       generateProjectKey(projectKeyLen),
		CreatedBy: fetchedUser.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	var projectErr error = nil
	for i := projectSaveRetries; i > 0; i-- {
		_, projectErr = p.projectRepo.Save(ctx, &newProject)
		// If the save was successful, then break out of the loop.
		if projectErr == nil {
			break
		} else {
			// If there was an error in saving the project, try again with a new
			// project key.
			newProject.Key = generateProjectKey(projectKeyLen)
		}
	}
	// If all attempts to save the project were unsuccessful, then return the
	// error.
	if projectErr != nil {
		if mongo.IsDuplicateKeyError(projectErr) {
			log.Printf(
				"a project with name %q exists for user %q",
				req.ProjectName,
				username,
			)
			return nil, EProjectNameTaken
		}

		log.Printf("error while creating new project: %v", projectErr)
		return nil, EProjectSave
	}

	return &projectpb.CreateNewProjectResponse{}, nil
}

func (p *ProjectServer) GetProjectKey(
	ctx context.Context,
	req *projectpb.GetProjectKeyRequest,
) (*projectpb.GetProjectKeyResponse, error) {
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

	project, err := p.projectRepo.GetByNameAndUserID(
		ctx,
		req.ProjectName,
		fetchedUser.ID,
	)
	if err != nil {
		log.Printf(
			"error while fetching project %q with user %q: %v",
			req.ProjectName,
			username,
			err,
		)

		if err == mongo.ErrNoDocuments {
			return nil, EProjectNotFound
		}

		return nil, EProjectFetch
	}

	response := projectpb.GetProjectKeyResponse{ProjectKey: project.Key}

	return &response, nil
}

// NewProjectServer creates a new server for the project service.
func NewProjectServer(
	projectRepo ProjectRepository,
	userRepo user.UserRepository,
) *ProjectServer {
	return &ProjectServer{projectRepo: projectRepo, userRepo: userRepo}
}
