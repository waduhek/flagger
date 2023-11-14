package environment

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/waduhek/flagger/proto/environmentpb"

	"github.com/waduhek/flagger/internal/auth"
	"github.com/waduhek/flagger/internal/project"
	"github.com/waduhek/flagger/internal/user"
)

type EnvironmentServer struct {
	environmentpb.UnimplementedEnvironmentServer
	userRepo        user.UserRepository
	projectRepo     project.ProjectRepository
	environmentRepo EnvironmentRepository
}

func (s *EnvironmentServer) CreateEnvironment(
	ctx context.Context,
	req *environmentpb.CreateEnvironmentRequest,
) (*environmentpb.CreateEnvironmentResponse, error) {
	jwtClaims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		log.Printf("could not find jwt claims in request context")
		return nil, status.Error(codes.Internal, "could not find token claims")
	}

	username := jwtClaims.Subject

	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		log.Printf("error while fetching user %q: %v", username, err)
		return nil, status.Error(codes.Internal, "error while fetching user")
	}

	// Check if the provided project exists with the user.
	project, err := s.projectRepo.GetByNameAndUserID(
		ctx,
		req.ProjectName,
		user.ID,
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf(
				"no projects were found with name %q with user %q",
				req.ProjectName,
				username,
			)
			return nil,
				status.Error(
					codes.NotFound,
					"no project was found with that name",
				)
		}

		log.Printf("error occurred while fetching projects: %v", err)
		return nil, status.Error(codes.Internal, "could not fetch projects")
	}

	// Create a new environment
	newEnvironment := Environment{
		Name:      req.EnvironmentName,
		ProjectID: project.ID,
		CreatedBy: user.ID,
		CreatedAt: time.Now(),
	}

	envResult, envSaveErr := s.environmentRepo.Save(ctx, &newEnvironment)
	if envSaveErr != nil {
		if mongo.IsDuplicateKeyError(envSaveErr) {
			log.Printf(
				"an environment %q already exists for project %q",
				req.EnvironmentName,
				project.Name,
			)
			return nil,
				status.Error(
					codes.AlreadyExists,
					"an environment with that name already exists",
				)
		}

		log.Printf(
			"could not create environment %q for project %q",
			req.EnvironmentName,
			project.ID,
		)
		return nil,
			status.Error(codes.Internal, "could not create a new environment")
	}

	environmentID, ok := envResult.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Printf("environment ID is not of type ObjectID")
		return nil,
			status.Error(codes.Internal, "could not create a new environment")
	}

	// Add the environment to the project.
	_, projectUpdateErr := s.projectRepo.AddEnvironment(
		ctx,
		project.ID,
		environmentID,
	)
	if projectUpdateErr != nil {
		log.Printf(
			"error while adding environment %q to project %q: %v",
			environmentID,
			project.ID,
			projectUpdateErr,
		)
		return nil,
			status.Error(codes.Internal, "could not create a new environment")
	}

	log.Printf(
		"successfully created environment %q for project %q",
		req.EnvironmentName,
		req.ProjectName,
	)
	return &environmentpb.CreateEnvironmentResponse{}, nil
}

func NewEnvironmentServer(
	userRepo user.UserRepository,
	projectRepo project.ProjectRepository,
	environmentRepo EnvironmentRepository,
) *EnvironmentServer {
	return &EnvironmentServer{
		userRepo:        userRepo,
		projectRepo:     projectRepo,
		environmentRepo: environmentRepo,
	}
}
