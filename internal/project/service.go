package project

import (
	"context"
	"math/rand"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const lettersLen = len(letters)

type projectTokenKey struct{}

// generateProjectKey generates a new project key with a length specified by n.
func generateProjectKey(n uint) string {
	b := make([]byte, n)

	for i := range b {
		b[i] = letters[rand.Intn(lettersLen)]
	}

	return string(b)
}

// injectProjectTokenIntoContext creates a new context with the project token
// in the value.
func injectProjectTokenIntoContext(
	ctx context.Context,
	token string,
) context.Context {
	return context.WithValue(ctx, projectTokenKey{}, token)
}

// ProjectTokenFromContext returns the project token from the provided context
// if it is available.
func ProjectTokenFromContext(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(projectTokenKey{}).(string)
	if !ok {
		return "", false
	}

	return token, true
}
