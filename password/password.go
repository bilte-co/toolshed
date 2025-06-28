// Package password provides password strength validation using entropy calculation.
// It uses the go-password-validator library to ensure passwords meet minimum security requirements
// by calculating the entropy (randomness) of the password.
//
// Example usage:
//
//	// Check password with default entropy requirement (60.0)
//	valid, msg := password.Check("mySecurePassword123!")
//	if !valid {
//		fmt.Printf("Password validation failed: %s\n", msg)
//	}
//
//	// Check password with custom entropy requirement
//	valid, msg = password.CheckEntropy("myPassword", 50.0)
//	if !valid {
//		fmt.Printf("Password validation failed: %s\n", msg)
//	}
package password

import (
	passwordvalidator "github.com/wagslane/go-password-validator"
)

// DefaultEntropy is the default minimum entropy requirement for password validation.
// A value of 60.0 provides strong security while remaining reasonable for users.
const DefaultEntropy = 60.0

// Check validates a password using the default entropy requirement.
// Returns true and empty string if valid, false and error message if invalid.
func Check(password string) (bool, string) {
	return CheckEntropy(password, DefaultEntropy)
}

// CheckEntropy validates a password against a custom minimum entropy requirement.
// Higher entropy values require more complex passwords. Returns true and empty string
// if the password meets the requirement, false and error message if it doesn't.
func CheckEntropy(password string, minEntropy float64) (bool, string) {
	err := passwordvalidator.Validate(password, minEntropy)
	if err != nil {
		return false, err.Error()
	}
	return true, ""
}
