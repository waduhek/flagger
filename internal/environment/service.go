package environment

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

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
		return nil, auth.ENoTokenClaims
	}

	username := jwtClaims.Subject

	fetchedUser, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		log.Printf("error while fetching user %q: %v", username, err)
		return nil, user.ECouldNotFetchUser
	}

	projectName := req.GetProjectName()
	environmentName := req.GetEnvironmentName()

	// Check if the provided project exists with the user.
	fetchedProject, err := s.projectRepo.GetByNameAndUserID(
		ctx,
		projectName,
		fetchedUser.ID,
	)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf(
				"no projects were found with name %q with user %q",
				projectName,
				username,
			)
			return nil, project.EProjectNotFound
		}

		log.Printf("error occurred while fetching projects: %v", err)
		return nil, project.EProjectFetch
	}

	// Create a session that will initiate a transaction to save the details of
	// the environment.
	session, err := s.mongoClient.StartSession()
	if err != nil {
		log.Printf("could not create a new session: %v", err)
		return nil, EEnvironmentTxn
	}

	// Start the transaction to save the environment.
	_, txnErr := session.WithTransaction(
		ctx,
		s.handleCreateEnvrionment(req, fetchedUser, fetchedProject),
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
		environmentName,
		projectName,
	)
	return &environmentpb.CreateEnvironmentResponse{}, nil
}

// handleCreateEnvrionment performs the transaction for creating the
// environment.
func (s *EnvironmentServer) handleCreateEnvrionment(
	req *environmentpb.CreateEnvironmentRequest,
	user *user.User,
	fetchedProject *project.Project,
) func(mongo.SessionContext) (interface{}, error) {
	return func(ctx mongo.SessionContext) (interface{}, error) {
		environmentName := req.GetEnvironmentName()

		// Create a new environment.
		newEnvironment := Environment{
			Name:      environmentName,
			ProjectID: fetchedProject.ID,
			CreatedBy: user.ID,
			CreatedAt: time.Now(),
		}

		envResult, envSaveErr := s.environmentRepo.Save(ctx, &newEnvironment)
		if envSaveErr != nil {
			if mongo.IsDuplicateKeyError(envSaveErr) {
				log.Printf(
					"an environment %q already exists for project %q",
					environmentName,
					fetchedProject.Name,
				)
				return nil, EEnvironmentNameTaken
			}

			log.Printf(
				"could not create environment %q for project %q",
				environmentName,
				fetchedProject.ID,
			)
			return nil, EEnvironmentSave
		}

		// Cast the returned ID of the inserted environment as an ObjectID.
		environmentID, ok := envResult.InsertedID.(primitive.ObjectID)
		if !ok {
			log.Printf("environment ID is not of type ObjectID")
			return nil, EEnvironmentIDCast
		}

		// Create new flag settings for all the flags that are present in the
		// current project.
		var flagSettings []flagsetting.FlagSetting
		for _, flagID := range fetchedProject.Flags {
			flagSetting := flagsetting.FlagSetting{
				FlagID:        flagID,
				ProjectID:     fetchedProject.ID,
				EnvironmentID: environmentID,
				IsActive:      true,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}

			flagSettings = append(flagSettings, flagSetting)
		}

		// Save the flag settings to the collection.
		if len(flagSettings) > 0 {
			insertedFlagSettings, flagSettingSaveErr := s.flagSettingRepo.SaveMany(
				ctx,
				flagSettings,
			)
			if flagSettingSaveErr != nil {
				log.Printf("error while saving flag settings: %v", flagSettingSaveErr)
				return nil, flagsetting.EFlagSettingSave
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
				fetchedProject.ID,
				insertedFlagSettingIDs...,
			)
			if projectFlagSettingErr != nil {
				log.Printf(
					"error while updating project with flag settings: %v",
					projectFlagSettingErr,
				)
				return nil, project.EProjectAddFlagSetting
			}
		}

		// Add the environment to the project.
		_, projectUpdateErr := s.projectRepo.AddEnvironment(
			ctx,
			fetchedProject.ID,
			environmentID,
		)
		if projectUpdateErr != nil {
			log.Printf(
				"error while adding environment %q to project %q: %v",
				environmentID,
				fetchedProject.ID,
				projectUpdateErr,
			)
			return nil, project.EProjectAddEnvironment
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
