package validator

import "fmt"

func ValidateString(value string, minLength int, maxLength int) error {
	if len(value) < minLength || len(value) > maxLength {
		return fmt.Errorf("string length must be between %d and %d characters", minLength, maxLength)
	}
	return nil
}