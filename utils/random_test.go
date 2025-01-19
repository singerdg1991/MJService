package utils_test

import (
	"github.com/hoitek/Maja-Service/utils"
	"testing"
)

// TestGenerateRandomNumber tests the functionality of utils.GenerateRandomNumber function.
// This function generates two random numbers within the range [1, 100] and checks if they are different.
// If the generated numbers are the same, it fails the test using t.Error.
//
// Parameters:
// - t *testing.T: A pointer to a testing.T object, used for reporting test failures.
//
// Return:
// - None
func TestGenerateRandomNumber(t *testing.T) {
	firstRandom := utils.GenerateRandomNumber(1, 100)
	secondRandom := utils.GenerateRandomNumber(1, 100)

	if firstRandom == secondRandom {
		t.Error("Generated random numbers should be different")
	}
}

// TestGenerateRandomNumberInRange tests the functionality of utils.GenerateRandomNumber function with a range of values.
//
// The function generates a random number within the range [1, 100] and checks if it falls within the expected range.
// If the generated number is not within the range, it fails the test using t.Errorf.
//
// Parameters:
// - t *testing.T: A pointer to a testing.T object, used for reporting test failures.
//
// Return:
// - None
func TestGenerateRandomNumberInRange(t *testing.T) {
	min := 1
	max := 100
	random := utils.GenerateRandomNumber(min, max)

	if random < min || random > max {
		t.Errorf("Generated random number %d is not within the range [%d, %d]", random, min, max)
	}
}

// TestGenerateRandomNumberMinRange tests the functionality of utils.GenerateRandomNumber function with a minimum range value.
//
// The function generates a random number within the range [100, 100] and checks if it is equal to the minimum range value.
// If the generated number is not equal to the minimum range value, it fails the test using t.Errorf.
//
// Parameters:
// - t *testing.T: A pointer to a testing.T object, used for reporting test failures.
//
// Return:
// - None
func TestGenerateRandomNumberMinRange(t *testing.T) {
	min := 100
	max := 100
	random := utils.GenerateRandomNumber(min, max)

	if random != min {
		t.Errorf("Generated random number %d is not equal to the minimum range value %d", random, min)
	}
}

// TestGenerateRandomNumberMaxRange tests the functionality of utils.GenerateRandomNumber function
// with a maximum range value.
//
// The function generates a random number within the range [100, 100] and checks if it is equal to the maximum range value.
// If the generated number is not equal to the maximum range value, it fails the test using t.Errorf.
//
// Parameters:
// - t *testing.T: A pointer to a testing.T object, used for reporting test failures.
//
// Return:
// - None
func TestGenerateRandomNumberMaxRange(t *testing.T) {
	min := 100
	max := 100
	random := utils.GenerateRandomNumber(min, max)

	if random != max {
		t.Errorf("Generated random number %d is not equal to the maximum range value %d", random, max)
	}
}

// TestGenerateRandomNumberNegativeRange tests the functionality of utils.GenerateRandomNumber function
// with a negative range value.
//
// The function generates a random number within the range [-100, -1] and checks if it falls within
// the expected range. If the generated number is not within the range, it fails the test using
// t.Errorf.
//
// Parameters:
// - t *testing.T: A pointer to a testing.T object, used for reporting test failures.
//
// Return:
// - None
func TestGenerateRandomNumberNegativeRange(t *testing.T) {
	min := -100
	max := -1
	random := utils.GenerateRandomNumber(min, max)

	if random > max || random < min {
		t.Errorf("Generated random number %d is not within the range [%d, %d]", random, min, max)
	}
}

// TestGenerateRandomNumberLargeRange tests the functionality of utils.GenerateRandomNumber function
// with a large range of values.
//
// Parameters:
// - t *testing.T: A pointer to a testing.T object, used for reporting test failures.
//
// The function generates a random number within the range [1, 1000000] and checks if it falls within
// the expected range. If the generated number is not within the range, it fails the test using
// t.Errorf.
func TestGenerateRandomNumberLargeRange(t *testing.T) {
	min := 1
	max := 1000000
	random := utils.GenerateRandomNumber(min, max)

	if random < min || random > max {
		t.Errorf("Generated random number %d is not within the range [%d, %d]", random, min, max)
	}
}

// TestGenerateRandomNumberZeroRange tests the functionality of utils.GenerateRandomNumber function
// with a range where both minimum and maximum values are zero.
//
// Parameters:
// - t *testing.T: A pointer to a testing.T object, used for reporting test failures.
//
// The function generates a random number within the range [0, 0] and checks if it falls within
// the expected range. Since the range is only one value, the generated number should always be
// equal to the minimum range value. If the generated number is not equal to the minimum range
// value, it fails the test using t.Errorf.
func TestGenerateRandomNumberZeroRange(t *testing.T) {
	min := 0
	max := 0
	random := utils.GenerateRandomNumber(min, max)

	if random != min {
		t.Errorf("Generated random number %d is not equal to the minimum range value %d", random, min)
	}
}

// TestGenerateRandomNumberLargeMin tests the functionality of utils.GenerateRandomNumber function
// with a large minimum range value.
//
// Parameters:
// - t *testing.T: A pointer to a testing.T object, used for reporting test failures.
//
// The function generates a random number within the range [1000000, 10000000] and checks if it falls within
// the expected range. If the generated number is not within the range, it fails the test using
// t.Errorf.
func TestGenerateRandomNumberLargeMin(t *testing.T) {
	min := 1000000
	max := 10000000
	random := utils.GenerateRandomNumber(min, max)

	if random < min || random > max {
		t.Errorf("Generated random number %d is not within the range [%d, %d]", random, min, max)
	}
}

// TestGenerateRandomNumberLargeMax tests the functionality of utils.GenerateRandomNumber function
// with a large maximum range value.
//
// Parameters:
// - t *testing.T: A pointer to a testing.T object, used for reporting test failures.
//
// The function generates a random number within the range [1, 1000000000] and checks if it falls within
// the expected range. If the generated number is not within the range, it fails the test using
// t.Errorf.
func TestGenerateRandomNumberLargeMax(t *testing.T) {
	min := 1
	max := 1000000000
	random := utils.GenerateRandomNumber(min, max)

	if random < min || random > max {
		t.Errorf("Generated random number %d is not within the range [%d, %d]", random, min, max)
	}
}

// TestGenerateRandomNumberSameRange tests the functionality of utils.GenerateRandomNumber function
// with the same minimum and maximum range values.
//
// Parameters:
// - t *testing.T: A pointer to a testing.T object, used for reporting test failures.
//
// The function generates two random numbers within the range [1, 100] and checks if they are different.
// If the generated numbers are the same, it fails the test using t.Error.
//
// Return:
// - None
func TestGenerateRandomNumberSameRange(t *testing.T) {
	min := 1
	max := 100
	firstRandom := utils.GenerateRandomNumber(min, max)
	secondRandom := utils.GenerateRandomNumber(min, max)

	if firstRandom == secondRandom {
		t.Error("Generated random numbers should be different")
	}
}
