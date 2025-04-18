package auth

import (
	"context"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"

	"github.com/waduhek/flagger/internal/logger"
)

const tokenDuration = 24 * time.Hour

type jwtClaimsKey struct{}

// FlaggerJWTClaims are the claims that a JWT must contain when authenticating
// with flagger.
type FlaggerJWTClaims struct {
	jwt.RegisteredClaims
}

// CreateJWT generates a new JWT with an expiry time that is 24 hours from
// now. It also adds the provided username to the `sub` field of the token.
func CreateJWT(logger logger.Logger, username string) (string, error) {
	jwtSecret := os.Getenv("FLAGGER_JWT_SECRET")

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
		logger.Error("could not sign jwt: %v", err)
		return "", ErrJWTSign
	}

	return token, nil
}

// VerifyJWT verifies that the provided string is valid and that it contains
// FlaggerJWTClaims as the claims. All errors returned from this function are
// GRPC compliant.
func VerifyJWT(logger logger.Logger, str string) (*FlaggerJWTClaims, error) {
	jwtSecret := os.Getenv("FLAGGER_JWT_SECRET")

	parsedToken, err := jwt.ParseWithClaims(
		str,
		&FlaggerJWTClaims{},
		func(t *jwt.Token) (interface{}, error) {
			signingMethod, ok := t.Method.(*jwt.SigningMethodHMAC)

			if !ok || signingMethod != jwt.SigningMethodHS512 {
				return nil, ErrInvalidJWT
			}

			return []byte(jwtSecret), nil
		},
	)
	if err != nil {
		logger.Error("error while parsing token: %v", err)
		return nil, err
	}

	if !parsedToken.Valid {
		logger.Error("token is not valid")
		return nil, ErrInvalidJWT
	}

	claims, ok := parsedToken.Claims.(*FlaggerJWTClaims)
	if !ok {
		logger.Error("could not parse token claims as FlaggerJWTClaims")
		return nil, ErrNoTokenClaims
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
