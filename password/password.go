package password

import (
	passwordvalidator "github.com/wagslane/go-password-validator"
)

const DefaultEntropy = 60.0

func Check(password string) (bool, string) {
	return CheckEntropy(password, DefaultEntropy)
}

func CheckEntropy(password string, minEntropy float64) (bool, string) {
	err := passwordvalidator.Validate(password, minEntropy)
	if err != nil {
		return false, err.Error()
	}
	return true, ""
}
