package common

import "fmt"

// WrapError into another error
func WrapError(original, aditional error) error {
	if aditional == nil {
		return original
	}

	if original == nil {
		return aditional
	}

	return fmt.Errorf("%w: %w", original, aditional)
}
