package project

import "math/rand"

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const lettersLen = len(letters)

// generateProjectKey generates a new project key with a length specified by n.
func generateProjectKey(n uint) string {
	b := make([]byte, n)

	for i := range b {
		b[i] = letters[rand.Intn(lettersLen)]
	}

	return string(b)
}
