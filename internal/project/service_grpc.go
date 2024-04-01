package project

import (
	"context"
	"errors"
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

type Server struct {
	projectpb.UnimplementedProjectServer
	projectDataRepo DataRepository
	userDataRepo    user.DataRepository
}

func (p *Server) CreateNewProject(
	ctx context.Context,
	req *projectpb.CreateNewProjectRequest,
) (*projectpb.CreateNewProjectResponse, error) {
	jwtClaims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		log.Printf("could not find claims from token")
		return nil, auth.ErrNoTokenClaims
	}

	username := jwtClaims.Subject

	fetchedUser, err := p.userDataRepo.GetByUsername(ctx, username)
	if err != nil {
		log.Printf("error while fetching user %q: %v", username, err)
		return nil, user.ErrCouldNotFetch
	}

	projectName := req.GetProjectName()

	newProject := Project{
		Name:      projectName,
		Key:       generateProjectKey(projectKeyLen),
		CreatedBy: fetchedUser.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	var projectErr error
	for i := projectSaveRetries; i > 0; i-- {
		_, projectErr = p.projectDataRepo.Save(ctx, &newProject)
		// If the save was successful, then break out of the loop.
		if projectErr == nil {
			break
		}

		// If there was an error in saving the project, try again with a new
		// project key.
		newProject.Key = generateProjectKey(projectKeyLen)
	}
	// If all attempts to save the project were unsuccessful, then return the
	// error.
	if projectErr != nil {
		if mongo.IsDuplicateKeyError(projectErr) {
			log.Printf(
				"a project with name %q exists for user %q",
				projectName,
				username,
			)
			return nil, ErrNameTaken
		}

		log.Printf("error while creating new project: %v", projectErr)
		return nil, ErrCouldNotSave
	}

	return &projectpb.CreateNewProjectResponse{}, nil
}

func (p *Server) GetProjectKey(
	ctx context.Context,
	req *projectpb.GetProjectKeyRequest,
) (*projectpb.GetProjectKeyResponse, error) {
	jwtClaims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		log.Printf("could not find claims from token")
		return nil, auth.ErrNoTokenClaims
	}

	username := jwtClaims.Subject

	fetchedUser, err := p.userDataRepo.GetByUsername(ctx, username)
	if err != nil {
		log.Printf("error while fetching user %q: %v", username, err)
		return nil, user.ErrCouldNotFetch
	}

	projectName := req.GetProjectName()

	project, err := p.projectDataRepo.GetByNameAndUserID(
		ctx,
		projectName,
		fetchedUser.ID,
	)
	if err != nil {
		log.Printf(
			"error while fetching project %q with user %q: %v",
			projectName,
			username,
			err,
		)

		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}

		return nil, ErrCouldNotFetch
	}

	response := projectpb.GetProjectKeyResponse{ProjectKey: project.Key}

	return &response, nil
}

// NewProjectServer creates a new server for the project service.
func NewProjectServer(
	projectDataRepo DataRepository,
	userDataRepo user.DataRepository,
) *Server {
	return &Server{projectDataRepo: projectDataRepo, userDataRepo: userDataRepo}
}
