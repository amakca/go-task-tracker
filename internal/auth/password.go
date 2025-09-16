package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"golang.org/x/crypto/argon2"
)

const (
	saltLength  = 16
	timeCost    = 1
	memoryCost  = 64 * 1024 // 64MB
	parallelism = 2
	keyLength   = 32
)

// HashPassword returns encoded hash string containing params, salt and hash.
func HashPassword(plain string) (string, error) {
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(plain), salt, timeCost, memoryCost, parallelism, keyLength)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)
	// format: argon2id$t=1$m=65536$p=2$salt$hash
	encoded := fmt.Sprintf("argon2id$t=%d$m=%d$p=%d$%s$%s", timeCost, memoryCost, parallelism, b64Salt, b64Hash)
	return encoded, nil
}

// VerifyPassword compares plain password against encoded hash.
func VerifyPassword(plain string, encoded string) (bool, error) {
	var t, m, p uint32
	var b64Salt, b64Hash string
	_, err := fmt.Sscanf(encoded, "argon2id$t=%d$m=%d$p=%d$%s$%s", &t, &m, &p, &b64Salt, &b64Hash)
	if err != nil {
		return false, errors.New("invalid hash format")
	}
	salt, err := base64.RawStdEncoding.DecodeString(b64Salt)
	if err != nil {
		return false, err
	}
	decodedHash, err := base64.RawStdEncoding.DecodeString(b64Hash)
	if err != nil {
		return false, err
	}
	hash := argon2.IDKey([]byte(plain), salt, t, m, uint8(p), uint32(len(decodedHash)))
	return base64.RawStdEncoding.EncodeToString(hash) == b64Hash, nil
}
