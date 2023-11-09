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
	"github.com/waduhek/flagger/proto/projectpb"

	"github.com/waduhek/flagger/internal/auth"
	"github.com/waduhek/flagger/internal/interceptors"
	"github.com/waduhek/flagger/internal/project"
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

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.AuthServerUnaryInterceptor,
			interceptors.ProjectServerUnaryInterceptor,
		),
	)
	// Registering servers
	authpb.RegisterAuthServer(grpcServer, authServer)
	projectpb.RegisterProjectServer(grpcServer, projectServer)

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
