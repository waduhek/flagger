package middleware

import (
	"context"
	"log"
	"os"
	"regexp"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	jwt "github.com/golang-jwt/jwt/v5"
)

type jwtClaimsKey struct{}

// AuthoriseJWT takes a GRPC context and validates that the current request has
// the "authorization" header and is a valid JWT. If the token is present and
// valid, adds the claims to the context and returns a new context for the
// handler. All errors from this middleware will be GRPC compliant.
func AuthoriseJWT(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Printf("could not find incoming request metadata")
		return nil, status.Error(
			codes.InvalidArgument,
			"could not find incoming request metadata",
		)
	}

	authHeader, ok := md["authorization"]
	if !ok {
		log.Printf("could not find authorization header")
		return nil, status.Error(
			codes.Unauthenticated,
			"authorization header not found",
		)
	}

	if len(authHeader) != 1 {
		log.Printf(
			"authorization header was found to be of len %d which is not expected",
			len(authHeader),
		)
		return nil, status.Error(
			codes.InvalidArgument,
			"invalid authorization header value",
		)
	}

	bearerToken := authHeader[0]

	claims, err := validateJWT(bearerToken)
	if err != nil {
		return nil, err
	}

	claimCtx := context.WithValue(ctx, jwtClaimsKey{}, claims)

	return claimCtx, nil
}

// ClaimsFromContext takes a GRPC context and tries to find the claims added by
// the AuthoriseJWT middleware.
func ClaimsFromContext(ctx context.Context) (*jwt.RegisteredClaims, bool) {
	claims, ok := ctx.Value(jwtClaimsKey{}).(*jwt.RegisteredClaims)
	if !ok {
		return nil, false
	}

	return claims, true
}

// validateJWT accepts the token header value of the "authorization" header and
// validates it. If the token is valid, returns the claims from the body of the
// token. If an error occurs, will always return a GRPC compliant error.
func validateJWT(token string) (*jwt.RegisteredClaims, error) {
	bearerTokenRegEx := regexp.MustCompile(
		`^[b|B]earer [a-zA-Z0-9-_]+\.[a-zA-Z0-9-_]+\.[a-zA-Z0-9-_]+$`,
	)
	if !bearerTokenRegEx.MatchString(token) {
		log.Println("token header format does not match")
		return nil,
			status.Error(codes.InvalidArgument, "invalid bearer header format")
	}

	headerJWT := strings.Split(token, " ")[1]

	parsedToken, err := jwt.ParseWithClaims(
		headerJWT,
		&jwt.RegisteredClaims{},
		func(jwtToken *jwt.Token) (interface{}, error) {
			if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil,
					status.Error(codes.Unauthenticated, "invalid jwt")
			}

			jwtSecret, ok := os.LookupEnv("FLAGGER_JWT_SECRET")
			if !ok {
				log.Println("could not find jwt signing key in environment variables")
				return nil,
					status.Error(codes.Internal, "error while signing jwt")
			}

			return []byte(jwtSecret), nil
		},
	)
	if err != nil {
		log.Printf("error while parsing token: %v", err)
		return nil,
			status.Error(codes.Internal, "error while parsing token")
	}

	if !parsedToken.Valid {
		log.Println("token is not valid")
		return nil,
			status.Error(codes.Unauthenticated, "invalid jwt")
	}

	claims, ok := parsedToken.Claims.(*jwt.RegisteredClaims)
	if !ok {
		log.Println("could not parse token claims as RegisteredClaims")
		return nil,
			status.Error(codes.Internal, "could not parse token claims")
	}

	return claims, nil
}
