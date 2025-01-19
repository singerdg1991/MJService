package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// GenerateCorrelationID generates a unique correlation ID.
func GenerateCorrelationID() string {
	// Generate a random byte array
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		// If random generation fails, fallback to timestamp
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	// Convert byte array to hexadecimal string
	return hex.EncodeToString(bytes)
}
