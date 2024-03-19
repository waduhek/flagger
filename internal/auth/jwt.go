package auth

import (
	"context"
	"log"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

var jwtSecret = os.Getenv("FLAGGER_JWT_SECRET")

const tokenDuration = 24 * time.Hour

type jwtClaimsKey struct{}

// FlaggerJWTClaims are the claims that a JWT must contain when authenticating
// with flagger.
type FlaggerJWTClaims struct {
	jwt.RegisteredClaims
}

// CreateJWT generates a new JWT with an expiry time that is 24 hours from
// now. It also adds the provided username to the `sub` field of the token.
func CreateJWT(username string) (string, error) {
	now := time.Now()
	expTime := now.Add(tokenDuration)

	claims := FlaggerJWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: &jwt.NumericDate{
				Time: now,
			},
			ExpiresAt: &jwt.NumericDate{
				Time: expTime,
			},
			Subject: username,
		},
	}

	tokenGenerator := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	token, err := tokenGenerator.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Printf("could not sign jwt: %v", err)
		return "", EJWTSign
	}

	return token, nil
}

// VerifyJWT verifies that the provided string is valid and that it contains
// FlaggerJWTClaims as the claims. All errors returned from this function are
// GRPC compliant.
func VerifyJWT(str string) (*FlaggerJWTClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(
		str,
		&FlaggerJWTClaims{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, EInvalidJWT
			}

			return []byte(jwtSecret), nil
		},
	)
	if err != nil {
		log.Printf("error while parsing token: %v", err)
		return nil, err
	}

	if !parsedToken.Valid {
		log.Println("token is not valid")
		return nil, EInvalidJWT
	}

	claims, ok := parsedToken.Claims.(*FlaggerJWTClaims)
	if !ok {
		log.Println("could not parse token claims as RegisteredClaims")
		return nil, ENoTokenClaims
	}

	return claims, nil
}

// InjectClaimsIntoContext creates a new context with the claims of a JWT.
func InjectClaimsIntoContext(
	ctx context.Context,
	claims jwt.Claims,
) context.Context {
	return context.WithValue(ctx, jwtClaimsKey{}, claims)
}

// ClaimsFromContext takes a context and tries to find the claims added by the
// authentication middleware.
func ClaimsFromContext(ctx context.Context) (*FlaggerJWTClaims, bool) {
	claims, ok := ctx.Value(jwtClaimsKey{}).(*FlaggerJWTClaims)
	if !ok {
		return nil, false
	}

	return claims, true
}
