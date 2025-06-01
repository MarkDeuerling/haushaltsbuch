package user

import (
	"golang.org/x/crypto/bcrypt"
)

// PasswordHasherFunc is a function type that generates and validates passwords.
type PasswordHasherFunc struct {
	Generate func(string) ([]byte, error)
	Validate func([]byte, string) error
}

// GeneratePassword is a function type that generates a password.
func (f PasswordHasherFunc) GeneratePassword(clearword string) ([]byte, error) {
	return f.Generate(clearword)
}

// ValidatePassword checks if the password is valid.
func (f PasswordHasherFunc) ValidatePassword(hash []byte, clearword string) error {
	return f.Validate(hash, clearword)
}

// GeneratePassword Signatur.
func GeneratePassword(clearword string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(clearword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

// ValidatePassword checks if the password is valid.
func ValidatePassword(hash []byte, clearword string) error {
	err := bcrypt.CompareHashAndPassword(hash, []byte(clearword))
	return err
}
