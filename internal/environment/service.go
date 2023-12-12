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
	"github.com/waduhek/flagger/internal/flagsetting"
	"github.com/waduhek/flagger/internal/project"
	"github.com/waduhek/flagger/internal/user"
)

type EnvironmentServer struct {
	environmentpb.UnimplementedEnvironmentServer
	mongoClient     *mongo.Client
	userRepo        user.UserRepository
	projectRepo     project.ProjectRepository
	flagSettingRepo flagsetting.FlagSettingRepository
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

	// Create a session that will initiate a transaction to save the details of
	// the environment.
	session, err := s.mongoClient.StartSession()
	if err != nil {
		log.Printf("could not create a new session: %v", err)
		return nil, status.Error(
			codes.Internal,
			"could not create a transaction session",
		)
	}

	// Start the transaction to save the environment.
	_, txnErr := session.WithTransaction(
		ctx,
		s.handleCreateEnvrionment(req, user, project),
	)
	if txnErr != nil {
		log.Printf(
			"error while performing environment save transaction: %v",
			txnErr,
		)
		return nil, txnErr
	}

	log.Printf(
		"successfully created environment %q for project %q",
		req.EnvironmentName,
		req.ProjectName,
	)
	return &environmentpb.CreateEnvironmentResponse{}, nil
}

// handleCreateEnvrionment performs the transaction for creating the
// environment.
func (s *EnvironmentServer) handleCreateEnvrionment(
	req *environmentpb.CreateEnvironmentRequest,
	user *user.User,
	project *project.Project,
) func(mongo.SessionContext) (interface{}, error) {
	return func(ctx mongo.SessionContext) (interface{}, error) {
		// Create a new environment.
		newEnvironment := Environment{
			Name:      req.EnvironmentName,
			ProjectID: project.ID,
			CreatedBy: user.ID,
			CreatedAt: time.Now(),
		}

		envResult, err := s.environmentRepo.Save(ctx, &newEnvironment)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				log.Printf(
					"an environment %q already exists for project %q",
					req.EnvironmentName,
					project.Name,
				)
				return nil, status.Error(
					codes.AlreadyExists,
					"an environment with that name already exists",
				)
			}

			log.Printf(
				"could not create environment %q for project %q",
				req.EnvironmentName,
				project.ID,
			)
			return nil, status.Error(
				codes.Internal,
				"could not create a new environment",
			)
		}

		// Cast the returned ID of the inserted environment as an ObjectID.
		environmentID, ok := envResult.InsertedID.(primitive.ObjectID)
		if !ok {
			log.Printf("environment ID is not of type ObjectID")
			return nil, status.Error(
				codes.Internal,
				"could not create a new environment",
			)
		}

		// Create new flag settings for all the flags that are present in the
		// current project.
		var flagSettings []flagsetting.FlagSetting
		for _, flagID := range project.Flags {
			flagSetting := flagsetting.FlagSetting{
				FlagID:        flagID,
				ProjectID:     project.ID,
				EnvironmentID: environmentID,
				IsActive:      true,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}

			flagSettings = append(flagSettings, flagSetting)
		}

		// Save the flag settings to the collection.
		if len(flagSettings) > 0 {
			insertedFlagSettings, err := s.flagSettingRepo.SaveMany(
				ctx,
				flagSettings,
			)
			if err != nil {
				log.Printf("error while saving flag settings: %v", err)
				return nil, status.Error(
					codes.Internal,
					"could not save flag settings",
				)
			}

			// Cast the IDs of all the flag settings as ObjectIDs.
			var insertedFlagSettingIDs []primitive.ObjectID
			for _, id := range insertedFlagSettings.InsertedIDs {
				insertedFlagSettingIDs = append(
					insertedFlagSettingIDs,
					id.(primitive.ObjectID),
				)
			}

			// Update the project with the new flag settings.
			_, projectFlagSettingErr := s.projectRepo.AddFlagSettings(
				ctx,
				project.ID,
				insertedFlagSettingIDs...,
			)
			if projectFlagSettingErr != nil {
				log.Printf(
					"error while updating project with flag settings: %v",
					projectFlagSettingErr,
				)
				return nil, status.Error(
					codes.Internal,
					"error while updating project with flag settings",
				)
			}
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
			return nil, status.Error(
				codes.Internal,
				"could not create a new environment",
			)
		}

		return nil, nil
	}
}

func NewEnvironmentServer(
	client *mongo.Client,
	userRepo user.UserRepository,
	projectRepo project.ProjectRepository,
	flagSettingRepo flagsetting.FlagSettingRepository,
	environmentRepo EnvironmentRepository,
) *EnvironmentServer {
	return &EnvironmentServer{
		mongoClient:     client,
		userRepo:        userRepo,
		projectRepo:     projectRepo,
		flagSettingRepo: flagSettingRepo,
		environmentRepo: environmentRepo,
	}
}
