package startup_test

import (
	"testing"

	"github.com/waduhek/flagger/internal/logger"
	"github.com/waduhek/flagger/internal/startup"
)

func TestSuccessfulRedisConnection(t *testing.T) {
	const connectionString string = "redis://localhost:6379/0"
	t.Setenv("FLAGGER_REDIS_URI", connectionString)

	_, err := startup.ConnectRedis(&logger.StubLogger{})
	if err != nil {
		t.Errorf("did not expect error when connecting to redis: %v", err)
	}
}

func TestUnsuccessfulRedisConnection(t *testing.T) {
	const connectionString = "http://example.com"
	t.Setenv("FLAGGER_REDIS_URI", connectionString)

	_, err := startup.ConnectRedis(&logger.StubLogger{})
	if err == nil {
		t.Error("expected error when connecting to redis")
	}
}
