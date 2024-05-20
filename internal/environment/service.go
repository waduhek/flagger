package environment

import (
	"context"
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

type Server struct {
	environmentpb.UnimplementedEnvironmentServer
	mongoClient         *mongo.Client
	userDataRepo        user.DataRepository
	projectDataRepo     project.DataRepository
	flagSettingDataRepo flagsetting.DataRepository
	environmentDataRepo DataRepository
}

func (s *Server) CreateEnvironment(
	ctx context.Context,
	req *environmentpb.CreateEnvironmentRequest,
) (*environmentpb.CreateEnvironmentResponse, error) {
	jwtClaims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		log.Printf("could not find jwt claims in request context")
		return nil, auth.ErrNoTokenClaims
	}

	username := jwtClaims.Subject

	fetchedUser, err := s.userDataRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	projectName := req.GetProjectName()
	environmentName := req.GetEnvironmentName()

	// Check if the provided project exists with the user.
	fetchedProject, err := s.projectDataRepo.GetByNameAndUserID(
		ctx,
		projectName,
		fetchedUser.ID,
	)
	if err != nil {
		return nil, err
	}

	// Create a session that will initiate a transaction to save the details of
	// the environment.
	session, err := s.mongoClient.StartSession()
	if err != nil {
		log.Printf("could not create a new session: %v", err)
		return nil, ErrTxnSession
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
func (s *Server) handleCreateEnvrionment(
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

		envResult, envSaveErr := s.environmentDataRepo.Save(ctx, &newEnvironment)
		if envSaveErr != nil {
			if mongo.IsDuplicateKeyError(envSaveErr) {
				log.Printf(
					"an environment %q already exists for project %q",
					environmentName,
					fetchedProject.Name,
				)
				return nil, ErrNameTaken
			}

			log.Printf(
				"could not create environment %q for project %q",
				environmentName,
				fetchedProject.ID,
			)
			return nil, ErrCouldNotSave
		}

		// Cast the returned ID of the inserted environment as an ObjectID.
		environmentID, ok := envResult.InsertedID.(primitive.ObjectID)
		if !ok {
			log.Printf("environment ID is not of type ObjectID")
			return nil, ErrEnvironmentIDCast
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
			insertedFlagSettings, flagSettingSaveErr := s.flagSettingDataRepo.SaveMany(
				ctx,
				flagSettings,
			)
			if flagSettingSaveErr != nil {
				log.Printf("error while saving flag settings: %v", flagSettingSaveErr)
				return nil, flagsetting.ErrCouldNotSave
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
			_, projectFlagSettingErr := s.projectDataRepo.AddFlagSettings(
				ctx,
				fetchedProject.ID,
				insertedFlagSettingIDs...,
			)
			if projectFlagSettingErr != nil {
				return nil, projectFlagSettingErr
			}
		}

		// Add the environment to the project.
		_, projectUpdateErr := s.projectDataRepo.AddEnvironment(
			ctx,
			fetchedProject.ID,
			environmentID,
		)
		if projectUpdateErr != nil {
			return nil, projectUpdateErr
		}

		return nil, nil
	}
}

func NewEnvironmentServer(
	client *mongo.Client,
	userDataRepo user.DataRepository,
	projectDataRepo project.DataRepository,
	flagSettingDataRepo flagsetting.DataRepository,
	environmentDataRepo DataRepository,
) *Server {
	return &Server{
		mongoClient:         client,
		userDataRepo:        userDataRepo,
		projectDataRepo:     projectDataRepo,
		flagSettingDataRepo: flagSettingDataRepo,
		environmentDataRepo: environmentDataRepo,
	}
}
