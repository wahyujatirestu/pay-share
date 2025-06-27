package security

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"strings"

	"golang.org/x/crypto/argon2"
)


func GeneratePasswordHash(password string) (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := strings.Join([]string{"argon2id", "v=19", b64Salt, b64Hash}, "$")

	return encoded, nil
}


func VerifyPasswordHash(encodedHash, password string) (bool, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 4 {
		return false, errors.New("invalid password hash format")
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[2])
	if err != nil {
		return false, err
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[3])
	if err != nil {
		return false, err
	}

	computedHash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	return string(computedHash) == string(expectedHash), nil
}
