package startup

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/waduhek/flagger/internal/logger"
)

func ConnectRedis(logger logger.Logger) (*redis.Client, error) {
	connectionString := os.Getenv("FLAGGER_REDIS_URI")

	return connectRedisWithConnectionString(logger, connectionString)
}

func connectRedisWithConnectionString(
	logger logger.Logger,
	connString string,
) (*redis.Client, error) {
	opt, parseErr := redis.ParseURL(connString)
	if parseErr != nil {
		logger.Error("error while parsing redis connection string: %v", parseErr)
		return nil, parseErr
	}

	client := redis.NewClient(opt)

	pingErr := pingRedis(client)
	if pingErr != nil {
		logger.Error("error while pinging redis: %v", pingErr)
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
