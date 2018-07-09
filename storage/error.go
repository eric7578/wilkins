package storage

import (
	"errors"
	"fmt"
)

func newErrNotAllowed(operation string) error {
	return fmt.Errorf("%s not allowed", operation)
}

func newErrOperationNotAllowed() error {
	return errors.New("operation not allowed")
}

// IsErrNotAllowed check if error is ErrNotAllowed
func IsErrNotAllowed(err error) bool {
	return isErrMessageContains(err, "not allowed")
}

func newErrNotFound(resource string) error {
	return fmt.Errorf("%s not found", resource)
}

// IsErrNotFound check if error is ErrNotFound
func IsErrNotFound(err error) bool {
	return isErrMessageContains(err, "not found")
}
