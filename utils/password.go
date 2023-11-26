package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/scrypt"
	"strings"
)

// HashPassword generates a secure hash of a password using scrypt.
// Returns the hashed password as a hex-encoded string or an error if hashing fails.
func HashPassword(password string) (string, error) {
	// example for making salt - https://play.golang.org/p/_Aw6WeWC42I
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	// using recommended cost parameters from - https://godoc.org/golang.org/x/crypto/scrypt
	shash, err := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
	if err != nil {
		return "", err
	}

	// return hex-encoded string with salt appended to password
	hashedPW := fmt.Sprintf("%s.%s", hex.EncodeToString(shash), hex.EncodeToString(salt))

	return hashedPW, nil
}

// ComparePasswords checks if the supplied password matches the stored hashed password.
// Returns true if they match or an error if there's a problem in the comparison process.
func ComparePasswords(storedPassword, suppliedPassword string) (bool, error) {
	pwsalt := strings.Split(storedPassword, ".")

	// check supplied password salted with hash
	salt, err := hex.DecodeString(pwsalt[1])
	if err != nil {
		return false, fmt.Errorf("error decoding password salt: %w", err)
	}

	shash, err := scrypt.Key([]byte(suppliedPassword), salt, 32768, 8, 1, 32)
	if err != nil {
		return false, fmt.Errorf("error generating password hash: %w", err)
	}

	return hex.EncodeToString(shash) == pwsalt[0], nil
}
