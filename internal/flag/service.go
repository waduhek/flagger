package flag

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/waduhek/flagger/proto/flagpb"

	"github.com/waduhek/flagger/internal/auth"
	"github.com/waduhek/flagger/internal/environment"
	"github.com/waduhek/flagger/internal/flagsetting"
	"github.com/waduhek/flagger/internal/project"
	"github.com/waduhek/flagger/internal/user"
)

type FlagServer struct {
	flagpb.UnimplementedFlagServer
	mongoClient     *mongo.Client
	userRepo        user.UserRepository
	projectRepo     project.ProjectRepository
	environmentRepo environment.EnvironmentRepository
	flagRepo        FlagRepository
	flagSettingRepo flagsetting.FlagSettingRepository
}

type mongoTxnCallback func(ctx mongo.SessionContext) (interface{}, error)

func (s *FlagServer) CreateFlag(
	ctx context.Context,
	req *flagpb.CreateFlagRequest,
) (*flagpb.CreateFlagResponse, error) {
	// Get the details of the currently authenticated user from the JWT.
	jwtClaims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		log.Println("could not get token claims")
		return nil, auth.ENoTokenClaims
	}

	username := jwtClaims.Subject

	fetchedUser, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		log.Printf("error while fetching user %q: %v", username, err)
		return nil, user.ECouldNotFetchUser
	}

	projectName := req.GetProjectName()
	flagName := req.GetFlagName()

	// Get the project that the flag is to be added to. If the project does not
	// belong to the currently authenticated user, or if the project doesn't
	// exist, return an error.
	fetchedProject, err := s.projectRepo.GetByNameAndUserID(
		ctx,
		projectName,
		fetchedUser.ID,
	)
	if err != nil {
		log.Printf(
			"error while getting project %q with user %q: %v",
			projectName,
			username,
			err,
		)

		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, project.EProjectNotFound
		}

		return nil, project.EProjectFetch
	}

	// If there are no environments configured for the project, don't allow any
	// flags to be created.
	if len(fetchedProject.Environments) == 0 {
		log.Printf(
			"no environments are configured for the project %q cannot create flag",
			projectName,
		)
		return nil, ENoEnvironments
	}

	// Start a transaction to save the flag.
	txnSession, err := s.mongoClient.StartSession()
	if err != nil {
		log.Printf("could not start transaction to save flag: %v", err)
		return nil, EFlagSaveTxn
	}

	_, txnErr := txnSession.WithTransaction(
		ctx,
		s.handleCreateFlag(req, fetchedUser, fetchedProject),
	)
	if txnErr != nil {
		log.Printf("could not complete flag save transaction: %v", err)
		return nil, err
	}

	log.Printf(
		"successfully created the flag %q in project %q",
		flagName,
		projectName,
	)
	return &flagpb.CreateFlagResponse{}, nil
}

// handleCreateFlag performs the transaction for saving the flag in the DB.
func (s *FlagServer) handleCreateFlag(
	req *flagpb.CreateFlagRequest,
	user *user.User,
	fetchedProject *project.Project,
) mongoTxnCallback {
	return func(ctx mongo.SessionContext) (interface{}, error) {
		flagName := req.GetFlagName()

		// Save the details of the flag.
		newFlag := &Flag{
			Name:      flagName,
			ProjectID: fetchedProject.ID,
			CreatedBy: user.ID,
			CreatedAt: time.Now(),
		}

		flagSaveResult, err := s.flagRepo.Save(ctx, newFlag)
		if err != nil {
			log.Printf("error while saving the flag: %v", err)

			if mongo.IsDuplicateKeyError(err) {
				return nil, EFlagNameTaken
			}

			return nil, EFlagSave
		}

		// Cast the returned ID of the saved flag as an ObjectID.
		savedFlagID, _ := flagSaveResult.InsertedID.(primitive.ObjectID)

		// Get all the environments that the project has.
		projectEnvIDs := fetchedProject.Environments

		// Create flag settings in all the environments created in the project.
		// By default the flags will be enabled in all the environments.
		var flagSettings []flagsetting.FlagSetting

		for _, envID := range projectEnvIDs {
			setting := flagsetting.FlagSetting{
				ProjectID:     fetchedProject.ID,
				EnvironmentID: envID,
				FlagID:        savedFlagID,
				IsActive:      true,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}

			flagSettings = append(flagSettings, setting)
		}

		flagSettingSaveResult, err := s.flagSettingRepo.SaveMany(
			ctx,
			flagSettings,
		)
		if err != nil {
			log.Printf("error while saving flag settings: %v", err)

			return nil, flagsetting.EFlagSettingSave
		}

		// Cast the IDs of the flag settings as ObjectIDs.
		var flagSettingIDs []primitive.ObjectID
		for _, id := range flagSettingSaveResult.InsertedIDs {
			flagSettingIDs = append(flagSettingIDs, id.(primitive.ObjectID))
		}

		// Add the flag's object ID to the list of flags in the project.
		_, projectFlagUpdateErr := s.projectRepo.AddFlag(
			ctx,
			fetchedProject.ID,
			savedFlagID,
		)
		if projectFlagUpdateErr != nil {
			log.Printf(
				"error while updating project with the flag: %v",
				projectFlagUpdateErr,
			)

			return nil, project.EProjectAddFlag
		}

		// Add all the object IDs of the flag settings to the project.
		_, projectSettingErr := s.projectRepo.AddFlagSettings(
			ctx,
			fetchedProject.ID,
			flagSettingIDs...,
		)
		if projectSettingErr != nil {
			log.Printf(
				"error while updating project with flag settings: %v",
				projectSettingErr,
			)

			return nil, project.EProjectAddFlagSetting
		}

		return nil, nil
	}
}

func (s *FlagServer) UpdateFlagStatus(
	ctx context.Context,
	req *flagpb.UpdateFlagStatusRequest,
) (*flagpb.UpdateFlagStatusResponse, error) {
	// Get the details of the currently authenticated user from the JWT.
	jwtClaims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		log.Println("could not get token claims")
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
	flagName := req.GetFlagName()
	isActive := req.GetIsActive()

	// Get the project that the has to be updated.
	fetchedProject, err := s.projectRepo.GetByNameAndUserID(
		ctx,
		projectName,
		fetchedUser.ID,
	)
	if err != nil {
		log.Printf("error while fetching project %q: %v", projectName, err)

		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, project.EProjectNotFound
		}

		return nil, project.EProjectFetch
	}

	// Get the environment that is to be updated.
	fetchedEnvironment, err := s.environmentRepo.GetByNameAndProjectID(
		ctx,
		environmentName,
		fetchedProject.ID,
	)
	if err != nil {
		log.Printf(
			"error while fetching environment %q: %v",
			environmentName,
			err,
		)

		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, environment.EEnvironmentNotFound
		}

		return nil, environment.EEnvironmentFetch
	}

	// Get the flag to update.
	flag, err := s.flagRepo.GetByNameAndProjectID(
		ctx,
		flagName,
		fetchedProject.ID,
	)
	if err != nil {
		log.Printf("error while fetching flag %q: %v", flagName, err)

		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, EFlagNotFound
		}

		return nil, EFlagFetch
	}

	// Update the flag setting to the desired value.
	updateResult, err := s.flagSettingRepo.UpdateIsActive(
		ctx,
		fetchedProject.ID,
		fetchedEnvironment.ID,
		flag.ID,
		isActive,
	)
	if err != nil {
		log.Printf("error while updating flag setting: %v", err)

		return nil, flagsetting.EFlagSettingStatusUpdate
	}

	if updateResult.ModifiedCount == 0 {
		log.Printf(
			"no flag settings were updated for project %q, environment %q, flag %q",
			projectName,
			environmentName,
			flagName,
		)

		return nil, EUpdateFlagStatus
	}

	return &flagpb.UpdateFlagStatusResponse{}, nil
}

// NewFlagServer creates a new `FlagServer` for serving GRPC requests.
func NewFlagServer(
	mongoClient *mongo.Client,
	userRepo user.UserRepository,
	projectRepo project.ProjectRepository,
	environmentRepo environment.EnvironmentRepository,
	flagRepo FlagRepository,
	flagSettingRepo flagsetting.FlagSettingRepository,
) *FlagServer {
	return &FlagServer{
		mongoClient:     mongoClient,
		userRepo:        userRepo,
		projectRepo:     projectRepo,
		environmentRepo: environmentRepo,
		flagRepo:        flagRepo,
		flagSettingRepo: flagSettingRepo,
	}
}
