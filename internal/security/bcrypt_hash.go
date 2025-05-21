package security

import (
	"golang.org/x/crypto/bcrypt"
)

// BcryptHash enables the Signatur.
type BcryptHash struct{}

// NewBcryptHash creates a new Authentificator.
func NewBcryptHash() *BcryptHash {
	return &BcryptHash{}
}

// GeneratePassword Signatur.
func (a *BcryptHash) GeneratePassword(clearword string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(clearword), 12)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

// ValidatePassword checks if the password is valid.
func (a *BcryptHash) ValidatePassword(hash []byte, clearword string) error {
	err := bcrypt.CompareHashAndPassword(hash, []byte(clearword))
	return err
}
