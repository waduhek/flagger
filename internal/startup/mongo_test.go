package startup_test

import (
	"os"
	"testing"

	"github.com/waduhek/flagger/internal/startup"
)

func TestSuccessfulMongoConnection(t *testing.T) {
	connStringFilePath, _ := createConnectionStringFile("mongodb://localhost:27017")
	t.Setenv("FLAGGER_MONGODB_CONN_FILE_PATH", connStringFilePath)

	_, clientErr := startup.ConnectMongo()
	if clientErr != nil {
		t.Error("did not expect error when connecting to mongo")
	}

	deleteConnectionStringFile(connStringFilePath)
}

func TestUnsuccessfulMongoConnection(t *testing.T) {
	connStringFilePath, _ := createConnectionStringFile("http://example.com")
	t.Setenv("FLAGGER_MONGODB_CONN_FILE_PATH", connStringFilePath)

	_, clientErr := startup.ConnectMongo()
	if clientErr == nil {
		t.Error("expected error when connecting to mongo")
	}

	deleteConnectionStringFile(connStringFilePath)
}

func createConnectionStringFile(connString string) (string, error) {
	tempFile, tempFileErr := os.CreateTemp("", "mongo-conn-string-")
	if tempFileErr != nil {
		return "", tempFileErr
	}
	defer tempFile.Close()

	_, writeErr := tempFile.WriteString(connString)
	if writeErr != nil {
		return "", writeErr
	}

	return tempFile.Name(), nil
}

func deleteConnectionStringFile(filePath string) {
	os.Remove(filePath)
}
