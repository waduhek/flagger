package flag

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/waduhek/flagger/proto/flagpb"

	"github.com/waduhek/flagger/internal/auth"
	"github.com/waduhek/flagger/internal/environment"
	"github.com/waduhek/flagger/internal/flagsetting"
	"github.com/waduhek/flagger/internal/project"
	"github.com/waduhek/flagger/internal/user"
)

type Server struct {
	flagpb.UnimplementedFlagServer
	mongoClient         *mongo.Client
	userDataRepo        user.DataRepository
	projectDataRepo     project.DataRepository
	environmentDataRepo environment.DataRepository
	flagDataRepo        DataRepository
	flagSettingDataRepo flagsetting.DataRepository
}

type mongoTxnCallback func(ctx mongo.SessionContext) (interface{}, error)

func (s *Server) CreateFlag(
	ctx context.Context,
	req *flagpb.CreateFlagRequest,
) (*flagpb.CreateFlagResponse, error) {
	// Get the details of the currently authenticated user from the JWT.
	jwtClaims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		log.Println("could not get token claims")
		return nil, auth.ErrNoTokenClaims
	}

	username := jwtClaims.Subject

	fetchedUser, err := s.userDataRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	projectName := req.GetProjectName()
	flagName := req.GetFlagName()

	// Get the project that the flag is to be added to. If the project does not
	// belong to the currently authenticated user, or if the project doesn't
	// exist, return an error.
	fetchedProject, err := s.projectDataRepo.GetByNameAndUserID(
		ctx,
		projectName,
		fetchedUser.ID,
	)
	if err != nil {
		return nil, err
	}

	// If there are no environments configured for the project, don't allow any
	// flags to be created.
	if len(fetchedProject.Environments) == 0 {
		log.Printf(
			"no environments are configured for the project %q cannot create flag",
			projectName,
		)
		return nil, ErrNoEnvironments
	}

	// Start a transaction to save the flag.
	txnSession, err := s.mongoClient.StartSession()
	if err != nil {
		log.Printf("could not start transaction to save flag: %v", err)
		return nil, ErrTxnSession
	}

	_, txnErr := txnSession.WithTransaction(
		ctx,
		s.handleCreateFlag(req, fetchedUser, fetchedProject),
	)
	if txnErr != nil {
		log.Printf("could not complete flag save transaction: %v", err)
		return nil, txnErr
	}

	log.Printf(
		"successfully created the flag %q in project %q",
		flagName,
		projectName,
	)
	return &flagpb.CreateFlagResponse{}, nil
}

// handleCreateFlag performs the transaction for saving the flag in the DB.
func (s *Server) handleCreateFlag(
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

		savedFlagID, err := s.flagDataRepo.Save(ctx, newFlag)
		if err != nil {
			return nil, err
		}

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

		flagSettingIDs, err := s.flagSettingDataRepo.SaveMany(
			ctx,
			flagSettings,
		)
		if err != nil {
			log.Printf("error while saving flag settings: %v", err)
			return nil, err
		}

		// Add the flag's object ID to the list of flags in the project.
		_, projectFlagUpdateErr := s.projectDataRepo.AddFlag(
			ctx,
			fetchedProject.ID,
			savedFlagID,
		)
		if projectFlagUpdateErr != nil {
			return nil, projectFlagUpdateErr
		}

		// Add all the object IDs of the flag settings to the project.
		_, projectSettingErr := s.projectDataRepo.AddFlagSettings(
			ctx,
			fetchedProject.ID,
			flagSettingIDs...,
		)
		if projectSettingErr != nil {
			return nil, projectSettingErr
		}

		return nil, nil
	}
}

func (s *Server) UpdateFlagStatus(
	ctx context.Context,
	req *flagpb.UpdateFlagStatusRequest,
) (*flagpb.UpdateFlagStatusResponse, error) {
	// Get the details of the currently authenticated user from the JWT.
	jwtClaims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		log.Println("could not get token claims")
		return nil, auth.ErrNoTokenClaims
	}

	username := jwtClaims.Subject

	fetchedUser, err := s.userDataRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	projectName := req.GetProjectName()
	environmentName := req.GetEnvironmentName()
	flagName := req.GetFlagName()
	isActive := req.GetIsActive()

	// Get the project that the has to be updated.
	fetchedProject, err := s.projectDataRepo.GetByNameAndUserID(
		ctx,
		projectName,
		fetchedUser.ID,
	)
	if err != nil {
		return nil, err
	}

	// Get the environment that is to be updated.
	fetchedEnvironment, err := s.environmentDataRepo.GetByNameAndProjectID(
		ctx,
		environmentName,
		fetchedProject.ID,
	)
	if err != nil {
		return nil, err
	}

	// Get the flag to update.
	flag, err := s.flagDataRepo.GetByNameAndProjectID(
		ctx,
		flagName,
		fetchedProject.ID,
	)
	if err != nil {
		return nil, err
	}

	// Update the flag setting to the desired value.
	updatedCount, err := s.flagSettingDataRepo.UpdateIsActive(
		ctx,
		fetchedProject.ID,
		fetchedEnvironment.ID,
		flag.ID,
		isActive,
	)
	if err != nil {
		return nil, err
	}

	if updatedCount == 0 {
		log.Printf(
			"no flag settings were updated for project %q, environment %q, flag %q",
			projectName,
			environmentName,
			flagName,
		)

		return nil, ErrUpdateStatus
	}

	return &flagpb.UpdateFlagStatusResponse{}, nil
}

// NewFlagServer creates a new `FlagServer` for serving GRPC requests.
func NewFlagServer(
	mongoClient *mongo.Client,
	userDataRepo user.DataRepository,
	projectDataRepo project.DataRepository,
	environmentDataRepo environment.DataRepository,
	flagDataRepo DataRepository,
	flagSettingDataRepo flagsetting.DataRepository,
) *Server {
	return &Server{
		mongoClient:         mongoClient,
		userDataRepo:        userDataRepo,
		projectDataRepo:     projectDataRepo,
		environmentDataRepo: environmentDataRepo,
		flagDataRepo:        flagDataRepo,
		flagSettingDataRepo: flagSettingDataRepo,
	}
}
