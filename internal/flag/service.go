package flag

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

	// Get the project that the flag is to be added to. If the project does not
	// belong to the currently authenticated user, or if the project doesn't
	// exist, return an error.
	fetchedProject, err := s.projectRepo.GetByNameAndUserID(
		ctx,
		req.ProjectName,
		fetchedUser.ID,
	)
	if err != nil {
		log.Printf(
			"error while getting project %q with user %q: %v",
			req.ProjectName,
			username,
			err,
		)

		if err == mongo.ErrNoDocuments {
			return nil, project.EProjectNotFound
		}

		return nil, project.EProjectFetch
	}

	// If there are no environments configured for the project, don't allow any
	// flags to be created.
	if len(fetchedProject.Environments) == 0 {
		log.Printf(
			"no environments are configured for the project %q cannot create flag",
			req.ProjectName,
		)
		return nil, status.Error(
			codes.FailedPrecondition,
			"no environments have been configured for the project",
		)
	}

	// Start a transaction to save the flag.
	txnSession, err := s.mongoClient.StartSession()
	if err != nil {
		log.Printf("could not start transaction to save flag: %v", err)
		return nil, status.Error(
			codes.Internal,
			"could start the transaction to save flag",
		)
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
		req.FlagName,
		req.ProjectName,
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
		// Save the details of the flag.
		newFlag := &Flag{
			Name:      req.FlagName,
			ProjectID: fetchedProject.ID,
			CreatedBy: user.ID,
			CreatedAt: time.Now(),
		}

		flagSaveResult, err := s.flagRepo.Save(ctx, newFlag)
		if err != nil {
			log.Printf("error while saving the flag: %v", err)

			if mongo.IsDuplicateKeyError(err) {
				return nil, status.Error(
					codes.AlreadyExists,
					"a flag with that name already exists",
				)
			}
		}

		// Cast the returned ID of the saved flag as an ObjectID.
		savedFlagID := flagSaveResult.InsertedID.(primitive.ObjectID)

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

			return nil,
				status.Error(codes.Internal, "could not save flag settings")
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

	// Get the project that the has to be updated.
	fetchedProject, err := s.projectRepo.GetByNameAndUserID(
		ctx,
		req.ProjectName,
		fetchedUser.ID,
	)
	if err != nil {
		log.Printf("error while fetching project %q: %v", req.ProjectName, err)

		if err == mongo.ErrNoDocuments {
			return nil, project.EProjectNotFound
		}

		return nil, project.EProjectFetch
	}

	// Get the environment that is to be updated.
	fetchedEnvironment, err := s.environmentRepo.GetByNameAndProjectID(
		ctx,
		req.EnvironmentName,
		fetchedProject.ID,
	)
	if err != nil {
		log.Printf(
			"error while fetching environment %q: %v",
			req.EnvironmentName,
			err,
		)

		if err == mongo.ErrNoDocuments {
			return nil, environment.EEnvironmentNotFound
		}

		return nil, environment.EEnvironmentFetch
	}

	// Get the flag to update.
	flag, err := s.flagRepo.GetByNameAndProjectID(
		ctx,
		req.FlagName,
		fetchedProject.ID,
	)
	if err != nil {
		log.Printf("error while fetching flag %q: %v", req.FlagName, err)

		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, "flag not found")
		}

		return nil, status.Error(
			codes.Internal,
			"error occurred while fetching flag",
		)
	}

	// Update the flag setting to the desired value.
	updateResult, err := s.flagSettingRepo.UpdateIsActive(
		ctx,
		fetchedProject.ID,
		fetchedEnvironment.ID,
		flag.ID,
		req.IsActive,
	)
	if err != nil {
		log.Printf("error while updating flag setting: %v", err)

		return nil, status.Error(
			codes.Internal,
			"error occurred while updating the flag setting",
		)
	}

	if updateResult.ModifiedCount == 0 {
		log.Printf(
			"no flag settings were updated for project %q, environment %q, flag %q",
			req.ProjectName,
			req.EnvironmentName,
			req.FlagName,
		)

		return nil, status.Error(
			codes.Internal,
			"no flag settings could be updated",
		)
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
