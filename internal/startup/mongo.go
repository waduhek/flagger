package startup

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectMongo establishes a connection with the MongoDB database cluster and
// returns the client.
func ConnectMongo() (*mongo.Client, error) {
	filePath := os.Getenv("FLAGGER_MONGODB_CONN_FILE_PATH")

	return connectMongoWithFile(filePath)
}

func connectMongoWithFile(filePath string) (*mongo.Client, error) {
	mongoConnectionString, getConnStringErr :=
		getMongoConnectionStringFromFile(filePath)
	if getConnStringErr != nil {
		log.Printf("could not get mongo connection string: %v", getConnStringErr)
		return nil, getConnStringErr
	}

	client, clientErr := getMongoClient(mongoConnectionString)
	if clientErr != nil {
		log.Printf("could not connect to mongo: %v", clientErr)
		return nil, clientErr
	}

	pingErr := pingMongo(client)
	if pingErr != nil {
		log.Printf("error while pinging mongo: %v", pingErr)
		return nil, pingErr
	}

	return client, nil
}

func getMongoConnectionStringFromFile(path string) (string, error) {
	mongoConnectionStringBytes, readFileErr := os.ReadFile(path)
	if readFileErr != nil {
		return "", readFileErr
	}

	return string(mongoConnectionStringBytes), nil
}

func getMongoClient(connectionString string) (*mongo.Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.
		Client().
		ApplyURI(connectionString).
		SetServerAPIOptions(serverAPI)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func pingMongo(client *mongo.Client) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	var result bson.M

	pingErr := client.
		Database("admin").
		RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).
		Decode(&result)
	if pingErr != nil {
		return pingErr
	}

	return nil
}
