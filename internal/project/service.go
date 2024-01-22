package project

import (
	"context"
	"math/rand"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const lettersLen = len(letters)

// projectKey is the key for storing the project key in the request context.
type projectKey struct{}

// generateProjectKey generates a new project key with a length specified by n.
func generateProjectKey(n uint) string {
	b := make([]byte, n)

	for i := range b {
		b[i] = letters[rand.Intn(lettersLen)]
	}

	return string(b)
}

// injectProjectKeyIntoContext creates a new context with the project key in the
// value.
func injectProjectKeyIntoContext(
	ctx context.Context,
	token string,
) context.Context {
	return context.WithValue(ctx, projectKey{}, token)
}

// ProjectKeyFromContext returns the project key from the provided context if it
// is available.
func ProjectKeyFromContext(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(projectKey{}).(string)
	if !ok {
		return "", false
	}

	return token, true
}
