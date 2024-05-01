//go:build integrationtest

package startup

import "testing"

func TestSuccessfulRedisConnection(t *testing.T) {
	const connectionString string = "redis://localhost:6379/0"

	_, err := connectRedisWithConnectionString(connectionString)
	if err != nil {
		t.Errorf("did not expect error when connecting to redis: %v", err)
	}
}

func TestUnsuccessfulRedisConnection(t *testing.T) {
	const connectionString = "http://example.com"

	_, err := connectRedisWithConnectionString(connectionString)
	if err == nil {
		t.Error("expected error when connecting to redis")
	}
}
