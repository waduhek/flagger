package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/waduhek/flagger/proto/authpb"
	"github.com/waduhek/flagger/proto/environmentpb"
	"github.com/waduhek/flagger/proto/flagpb"
	"github.com/waduhek/flagger/proto/projectpb"
	"github.com/waduhek/flagger/proto/providerpb"

	"github.com/waduhek/flagger/internal/auth"
	"github.com/waduhek/flagger/internal/environment"
	"github.com/waduhek/flagger/internal/flag"
	"github.com/waduhek/flagger/internal/flagsetting"
	"github.com/waduhek/flagger/internal/project"
	"github.com/waduhek/flagger/internal/provider"
	"github.com/waduhek/flagger/internal/user"
)

var mongoConnectionString string = os.Getenv("FLAGGER_MONGO_URI")

var flaggerDB string = os.Getenv("FLAGGER_DB")

var serverPort, _ = strconv.ParseUint(os.Getenv("FLAGGER_PORT"), 10, 16)

func initAuthServer(db *mongo.Database) *auth.AuthServer {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	userRepo, err := user.NewUserRepository(ctx, db)
	if err != nil {
		log.Panicf("could not initialise user repository: %v", err)
	}

	authServer := auth.NewAuthServer(userRepo)

	return authServer
}

func initProjectServer(db *mongo.Database) *project.ProjectServer {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	projectRepo, err := project.NewProjectRepository(ctx, db)
	if err != nil {
		log.Panicf("could not initialise project repository: %v", err)
	}

	userRepo, err := user.NewUserRepository(ctx, db)
	if err != nil {
		log.Panicf("could not initialise user repository: %v", err)
	}

	return project.NewProjectServer(projectRepo, userRepo)
}

func initEnvironmentServer(
	client *mongo.Client,
	db *mongo.Database,
) *environment.EnvironmentServer {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	userRepo, err := user.NewUserRepository(ctx, db)
	if err != nil {
		log.Panicf("could not initialise user repository: %v", err)
	}

	projectRepo, err := project.NewProjectRepository(ctx, db)
	if err != nil {
		log.Panicf("could not initialise project repository: %v", err)
	}

	flagSettingRepo, err := flagsetting.NewFlagSettingRepository(ctx, db)
	if err != nil {
		log.Panicf("could not initialise flag setting repository: %v", err)
	}

	environmentRepo, err := environment.NewEnvironmentRepository(ctx, db)
	if err != nil {
		log.Panicf("could not initialise environment repository: %v", err)
	}

	return environment.NewEnvironmentServer(
		client,
		userRepo,
		projectRepo,
		flagSettingRepo,
		environmentRepo,
	)
}

func initFlagServer(client *mongo.Client, db *mongo.Database) *flag.FlagServer {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	userRepo, err := user.NewUserRepository(ctx, db)
	if err != nil {
		log.Panicf("could not initialise user repository: %v", err)
	}

	projectRepo, err := project.NewProjectRepository(ctx, db)
	if err != nil {
		log.Panicf("could not initialise project repository: %v", err)
	}

	flagRepo, err := flag.NewFlagRepository(ctx, db)
	if err != nil {
		log.Panicf("could not initialise flag repository: %v", err)
	}

	flagSettingRepo, err := flagsetting.NewFlagSettingRepository(ctx, db)
	if err != nil {
		log.Panicf("could not initialise flag setting repository: %v", err)
	}

	environmentRepo, err := environment.NewEnvironmentRepository(ctx, db)
	if err != nil {
		log.Panicf("could not initialise environment repository: %v", err)
	}

	return flag.NewFlagServer(
		client,
		userRepo,
		projectRepo,
		environmentRepo,
		flagRepo,
		flagSettingRepo,
	)
}

func initFlagProviderServer() *provider.FlagProviderServer {
	return provider.NewFlagProviderServer()
}

func connectMongo() *mongo.Client {
	serverApi := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.
		Client().
		ApplyURI(mongoConnectionString).
		SetServerAPIOptions(serverApi)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Panicf("could not connect to mongo: %v", err)
	}

	var result bson.M

	pingErr := client.
		Database("admin").
		RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).
		Decode(&result)
	if pingErr != nil {
		log.Panicf("error while pinging mongo: %v", pingErr)
	}

	return client
}

func gracefulShutdown(cleanup func()) {
	sig := make(chan os.Signal, 1)

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		recvSig := <-sig
		log.Printf("shutting down due to signal %v", recvSig)
		cleanup()
		log.Printf("goodbye")
		os.Exit(0)
	}()
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", serverPort))
	if err != nil {
		log.Panicf("could not listen on port %d: %v", serverPort, err)
	}

	mongoClient := connectMongo()
	db := mongoClient.Database(flaggerDB)

	// Initialising all the servers
	authServer := initAuthServer(db)
	projectServer := initProjectServer(db)
	environmentServer := initEnvironmentServer(mongoClient, db)
	flagServer := initFlagServer(mongoClient, db)
	flagProviderServer := initFlagProviderServer()

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			auth.AuthServerUnaryInterceptor,
			auth.AuthoriseRequestInterceptor("/projectpb.Project/"),
			auth.AuthoriseRequestInterceptor(
				"/environmentpb.Environment/",
			),
			auth.AuthoriseRequestInterceptor("/flagpb.Flag/"),
			project.ProjectKeyUnaryInterceptor("/providerpb.FlagProvider/"),
		),
	)
	// Registering servers
	authpb.RegisterAuthServer(grpcServer, authServer)
	projectpb.RegisterProjectServer(grpcServer, projectServer)
	environmentpb.RegisterEnvironmentServer(grpcServer, environmentServer)
	flagpb.RegisterFlagServer(grpcServer, flagServer)
	providerpb.RegisterFlagProviderServer(grpcServer, flagProviderServer)

	// GRPC reflection
	reflection.Register(grpcServer)

	gracefulShutdown(func() {
		ctx := context.Background()

		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Panicf("could not disconnect from mongodb: %v", err)
		}

		grpcServer.GracefulStop()
	})

	log.Printf("flagger server listening at %q", lis.Addr().String())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("could not serve: %v", err)
	}
}
