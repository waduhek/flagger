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

	"github.com/redis/go-redis/v9"

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

func initAuthServer(db *mongo.Database) *auth.Server {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	userRepo, err := user.NewUserRepository(ctx, db)
	if err != nil {
		log.Panicf("could not initialise user repository: %v", err)
	}

	authServer := auth.NewServer(userRepo)

	return authServer
}

func initProjectServer(db *mongo.Database) *project.Server {
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
) *environment.Server {
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

func initFlagProviderServer(
	db *mongo.Database,
	redisClient *redis.Client,
) *provider.FlagProviderServer {
	providerRepo := provider.NewProviderRepository(db)
	cacheRepo := provider.NewProviderCacheRepository(redisClient)

	return provider.NewFlagProviderServer(providerRepo, cacheRepo)
}

func connectMongo() *mongo.Client {
	mongoConnectionString := os.Getenv("FLAGGER_MONGO_URI")

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.
		Client().
		ApplyURI(mongoConnectionString).
		SetServerAPIOptions(serverAPI)

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

func connectRedis() *redis.Client {
	redisConnectionString := os.Getenv("FLAGGER_REDIS_URI")

	opt, err := redis.ParseURL(redisConnectionString)
	if err != nil {
		log.Panicf("could not parse redis connection string: %v", err)
	}

	client := redis.NewClient(opt)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	pingErr := client.Ping(ctx).Err()
	if pingErr != nil {
		log.Panicf("error while pinging redis: %v", pingErr)
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
	flaggerDB := os.Getenv("FLAGGER_DB")
	serverPort, _ := strconv.ParseUint(os.Getenv("FLAGGER_PORT"), 10, 16)

	lis, lisErr := net.Listen("tcp", fmt.Sprintf(":%d", serverPort))
	if lisErr != nil {
		log.Panicf("could not listen on port %d: %v", serverPort, lisErr)
	}

	mongoClient := connectMongo()
	mongoDB := mongoClient.Database(flaggerDB)

	redisClient := connectRedis()

	// Initialising all the servers
	authServer := initAuthServer(mongoDB)
	projectServer := initProjectServer(mongoDB)
	environmentServer := initEnvironmentServer(mongoClient, mongoDB)
	flagServer := initFlagServer(mongoClient, mongoDB)
	flagProviderServer := initFlagProviderServer(mongoDB, redisClient)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			auth.UnaryServerInterceptor,
			auth.AuthoriseRequestInterceptor("/projectpb.Project/"),
			auth.AuthoriseRequestInterceptor(
				"/environmentpb.Environment/",
			),
			auth.AuthoriseRequestInterceptor("/flagpb.Flag/"),
			project.KeyUnaryInterceptor("/providerpb.FlagProvider/"),
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

		if disconnectErr := mongoClient.Disconnect(ctx); disconnectErr != nil {
			log.Panicf("could not disconnect from mongodb: %v", disconnectErr)
		}

		grpcServer.GracefulStop()
	})

	log.Printf("flagger server listening at %q", lis.Addr().String())
	if serveErr := grpcServer.Serve(lis); serveErr != nil {
		log.Fatalf("could not serve: %v", serveErr)
	}
}
