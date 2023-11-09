package hash

import (
	"bytes"
	"crypto/rand"

	"golang.org/x/crypto/argon2"
)

type PasswordHashDetails struct {
	Hash []byte
	Salt []byte
}

// The number of bytes for the salt
const saltLength uint16 = 16

// The time parameter for the hashing algorithm
const hashTime uint32 = 1

// The memory to allocate for the hashing algorithm
const memory uint32 = 64 * 1024

// The number of threads for the hashing algorithm
const threads uint8 = 4

// The key length of the generated hash
const keyLen uint32 = 64

// generateRandomSalt generates a random salt of a set byte size for password
// hashing.
func generateRandomSalt() ([]byte, error) {
	salt := make([]byte, saltLength)

	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}

	return salt, nil
}

// generatePasswordHash hashes the provided password returning the hash
// generated and the random salt that was used for hashing the password.
func GeneratePasswordHash(password string) (PasswordHashDetails, error) {
	salt, err := generateRandomSalt()
	if err != nil {
		return PasswordHashDetails{}, err
	}

	passwordHash := argon2.IDKey(
		[]byte(password),
		[]byte(salt),
		hashTime,
		memory,
		threads,
		keyLen,
	)

	return PasswordHashDetails{Hash: passwordHash, Salt: salt}, nil
}

// verifyPasswordHash verifies the provided plain text password to verify if the
// provided plain text password generates a hash that will match the expected
// hash.
func VerifyPasswordHash(plainPassword string, expectedHash []byte, salt []byte) bool {
	generatedHash := argon2.IDKey(
		[]byte(plainPassword),
		salt,
		hashTime,
		memory,
		threads,
		keyLen,
	)

	return bytes.Equal(generatedHash, expectedHash)
}
