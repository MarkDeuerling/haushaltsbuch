package id

import "github.com/google/uuid"

// UUIDGeneratorFunc is a function type that generates a UUID.
type UUIDGeneratorFunc func() (string, error)

// GenerateUUID Signatur.
func (f UUIDGeneratorFunc) GenerateUUID() (string, error) {
	return f()
}

// GenerateUUID generates a new UUID as a string.
func GenerateUUID() (string, error) {
	return uuid.New().String(), nil
}
