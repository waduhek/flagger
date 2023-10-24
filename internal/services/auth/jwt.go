package auth

import (
	"errors"
	"log"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

// createNewJWT generates a new JWT with an expiry time that is 24 hours from
// now. It also adds the provided username to the `sub` field of the token.
func createNewJWT(username string) (string, error) {
	jwtSecret, ok := os.LookupEnv("FLAGGER_JWT_SECRET")
	if !ok {
		log.Print("could not find jwt signing key in environment variables")
		return "", errors.New("error while signing jwt")
	}

	now := time.Now()
	expTime := now.Add(24 * time.Hour)

	claims := jwt.RegisteredClaims{
		IssuedAt: &jwt.NumericDate{
			Time: now,
		},
		ExpiresAt: &jwt.NumericDate{
			Time: expTime,
		},
		Subject: username,
	}

	tokenGenerator := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	token, err := tokenGenerator.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Printf("could not sign jwt: %v", err)
		return "", errors.New("error while signing jwt")
	}

	return token, nil
}
