package utils_test

import (
	"testing"

	"github.com/waduhek/flagger/internal/utils"
)

func TestGeneratePasswordHash(t *testing.T) {
	t.Run("some_passsword", func(subT *testing.T) {
		// password := "super_s3cR3t"
		password := "this_is anUnsualLY l0ng_and_stong p@a$$w0Rd!"

		got, err := utils.GeneratePasswordHash(password)
		if err != nil {
			subT.Errorf("error while hashing %q: %v", password, err)
		}

		if len(got.Hash) <= 0 {
			subT.Errorf(
				"hash length of %q is less than or equal to 0",
				password,
			)
		}
		if len(got.Salt) <= 0 {
			subT.Errorf(
				"salt length of %q is less than or equal to 0",
				password,
			)
		}
	})

	t.Run("empty_password", func(subT *testing.T) {
		password := ""

		got, err := utils.GeneratePasswordHash(password)
		if err != nil {
			subT.Errorf("error while hashing %q: %v", password, err)
		}

		if len(got.Hash) <= 0 {
			subT.Errorf(
				"hash length of %q is less than or equal to 0",
				password,
			)
		}
		if len(got.Salt) <= 0 {
			subT.Errorf(
				"salt length of %q is less than or equal to 0",
				password,
			)
		}

	})
}

func TestVerifyPasswordHash(t *testing.T) {
	hash := []byte{
		208, 164, 67, 67, 145, 147, 230, 181, 67, 48, 137, 90, 66, 117, 1, 53,
		246, 79, 122, 165, 148, 120, 25, 20, 178, 98, 132, 132, 75, 198, 151,
		67, 242, 167, 255, 221, 87, 51, 208, 193, 26, 20, 162, 36, 77, 100, 86,
		80, 231, 69, 236, 14, 136, 170, 92, 86, 52, 16, 209, 13, 170, 210, 166,
		209,
	}
	salt := []byte{
		21, 6, 98, 248, 134, 145, 11, 71, 44, 193, 242, 190, 102, 139, 65, 186,
	}
	plain := "super_s3cR3t"

	ok := utils.VerifyPasswordHash(plain, hash, salt)
	if !ok {
		t.Error("password verification failed")
	}
}

func BenchmarkGeneratePasswordHash(b *testing.B) {
	password := "this_is anUnsualLY l0ng_and_stong p@a$$w0Rd!"

	for i := 0; i < b.N; i++ {
		_, err := utils.GeneratePasswordHash(password)
		if err != nil {
			b.Errorf("hashing failed at i=%d", i)
		}
	}
}

func BenchmarkVerifyPasswordHash(b *testing.B) {
	hash := []byte{
		85, 149, 191, 225, 195, 168, 153, 215, 15, 201, 120, 235, 243, 188, 122,
		51, 144, 241, 36, 237, 80, 73, 111, 41, 74, 167, 63, 92, 189, 255, 46,
		122, 42, 116, 39, 91, 24, 114, 249, 100, 37, 150, 252, 113, 189, 174,
		121, 214, 211, 186, 155, 190, 59, 150, 221, 102, 18, 104, 108, 253, 237,
		12, 42, 55,
	}
	salt := []byte{
		127, 207, 82, 119, 226, 34, 30, 30, 85, 15, 155, 126, 240, 169, 2, 113,
	}
	plain := "this_is anUnsualLY l0ng_and_stong p@a$$w0Rd!"

	for i := 0; i < b.N; i++ {
		utils.VerifyPasswordHash(plain, hash, salt)
	}
}
