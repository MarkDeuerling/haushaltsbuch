package id

import "github.com/google/uuid"

// UUID enables the Signatur.
type UUID struct{}

// NewUUID creates a new UUID.
func NewUUID() *UUID {
	return &UUID{}
}

// GenerateUUID Signatur.
func (id *UUID) GenerateUUID() (string, error) {
	return uuid.New().String(), nil
}
