package startup

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis() (*redis.Client, error) {
	connectionString := os.Getenv("FLAGGER_REDIS_URI")

	return connectRedisWithConnectionString(connectionString)
}

func connectRedisWithConnectionString(connString string) (*redis.Client, error) {
	opt, parseErr := redis.ParseURL(connString)
	if parseErr != nil {
		log.Printf("error while parsing redis connection string: %v", parseErr)
		return nil, parseErr
	}

	client := redis.NewClient(opt)

	pingErr := pingRedis(client)
	if pingErr != nil {
		log.Printf("error while pinging redis: %v", pingErr)
		return nil, pingErr
	}

	return client, nil
}

func pingRedis(client *redis.Client) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	pingErr := client.Ping(ctx).Err()
	if pingErr != nil {
		return pingErr
	}

	return nil
}
