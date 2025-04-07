package auth_test

import (
	"testing"

	"github.com/waduhek/flagger/internal/auth"
	"github.com/waduhek/flagger/internal/logger"
)

const jwtSecretEnv string = "FLAGGER_JWT_SECRET"
const jwtSecretEnvValue string = "testing_secret"

var stubLogger = &logger.StubLogger{}

func TestCreateJWT(t *testing.T) {
	t.Setenv(jwtSecretEnv, jwtSecretEnvValue)

	_, tokenErr := auth.CreateJWT(stubLogger, "test_username")
	if tokenErr != nil {
		t.Errorf("did not expect error while generating token: %v", tokenErr)
	}
}

func TestVerifyJWT(t *testing.T) {
	t.Setenv(jwtSecretEnv, jwtSecretEnvValue)

	t.Run("valid_jwt", func(subT *testing.T) {
		const validJWT string = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9." +
			"eyJzdWIiOiJ0ZXN0X3VzZXIiLCJpYXQiOjE1MTYyMzkwMjIsImV4cCI6MzgxODQ0ODAwMH0." +
			"TgIBdgdQ7yYARUTqnsmcKJdtKtoH5lEyj-di012kkuxGWW3PCxpwWUOj8kbUR26rRwr3ThDXY1kQsoiaaXdwVQ"

		_, verifyErr := auth.VerifyJWT(stubLogger, validJWT)
		if verifyErr != nil {
			subT.Errorf("did not expect error for valid JWT: %v", verifyErr)
		}
	})

	t.Run("incorrect_signing_method", func(subT *testing.T) {
		const incorrectSignedJWT string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
			"eyJzdWIiOiJ0ZXN0X3VzZXIiLCJpYXQiOjE1MTYyMzkwMjIsImV4cCI6MzgxODQ0ODAwMH0." +
			"In_XjaJevdTToyqb9pf_iKPjVU0cPRZwybBGy-ncUFM"

		_, verifyErr := auth.VerifyJWT(stubLogger, incorrectSignedJWT)
		if verifyErr == nil {
			subT.Error("expected error when verifying incorrect signing method in JWT")
		}
	})

	t.Run("expired_jwt", func(subT *testing.T) {
		const expiredJWT string = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9." +
			"eyJzdWIiOiJ0ZXN0X3VzZXIiLCJpYXQiOjE1MTYyMzkwMjIsImV4cCI6MH0." +
			"V6HuGYsLsd9L1HhxNqZREF_hru5J0CyqQ6qN7Oy6vZ9Y6UVALO_H17FcXkKs2ZolkCJ5M12r-a2e0sE7zlQ3jg"

		_, verifyErr := auth.VerifyJWT(stubLogger, expiredJWT)
		if verifyErr == nil {
			subT.Error("expected error when verifying expired JWT")
		}
	})
}
