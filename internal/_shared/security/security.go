package security

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashPassword hashes a password using SHA-256 and returns the hashed password as a string.
func HashPassword(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))
	return hashedPassword
}

// ValidatePassword checks if the provided password matches the stored hashed password.
func ValidatePassword(inputPassword, hashedPassword string) bool {
	return HashPassword(inputPassword) == hashedPassword
}
