package utils

import (
	"math/rand"
	"time"
)

func GenerateRandomNumber(min, max int) int {
	// Use rand.NewSource to create a new source of random numbers
	src := rand.NewSource(time.Now().UnixNano())

	// Create a new rand.Rand using the source
	r := rand.New(src)

	return r.Intn(max-min+1) + min
}
