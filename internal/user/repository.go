package user

import (
	"context"
	"errors"
	"sync"
)

// InMemoryUserRepository implements UserRepository with an in-memory store.
type InMemoryUserRepository struct {
	users         map[string]*User
	emailToID     map[string]string
	refreshTokens map[string][]string
	mutex         sync.RWMutex
}

// NewInMemoryUserRepository creates a new InMemoryUserRepository.
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users:         make(map[string]*User),
		emailToID:     make(map[string]string),
		refreshTokens: make(map[string][]string),
	}
}

// CreateUser adds a new user to the repository.
func (r *InMemoryUserRepository) CreateUser(ctx context.Context, user *User) (*User, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		r.mutex.Lock()
		defer r.mutex.Unlock()

		if _, exists := r.emailToID[user.Email()]; exists {
			return nil, errors.New("email already exists")
		}

		if _, exists := r.users[user.ID()]; exists {
			return nil, errors.New("user ID already exists")
		}

		user.Aktualisert()
		user.ErstelltAm()

		r.users[user.ID()] = user
		r.emailToID[user.Email()] = user.ID()
		return user, nil
	}
}

// FindUserByEmail retrieves a user by their email.
func (r *InMemoryUserRepository) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		r.mutex.RLock()
		defer r.mutex.RUnlock()

		userID, exists := r.emailToID[email]
		if !exists {
			return nil, errors.New("user not found")
		}

		user, exists := r.users[userID]
		if !exists {
			return nil, errors.New("user not found")
		}
		return user, nil
	}
}

// FindUserByID retrieves a user by their ID.
func (r *InMemoryUserRepository) FindUserByID(ctx context.Context, id string) (*User, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		r.mutex.RLock()
		defer r.mutex.RUnlock()

		user, exists := r.users[id]
		if !exists {
			return nil, errors.New("user not found")
		}
		return user, nil
	}
}

// LogoutUser removes a specific refresh token for a user.
func (r *InMemoryUserRepository) LogoutUser(ctx context.Context, userID, refreshToken string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		r.mutex.Lock()
		defer r.mutex.Unlock()

		if _, exists := r.users[userID]; !exists {
			return errors.New("user not found")
		}

		tokens, exists := r.refreshTokens[userID]
		if !exists {
			return nil
		}

		for i, token := range tokens {
			if token == refreshToken {
				// r.refreshTokens[userID] = slices.Delete(tokens, i, i+1)
				r.refreshTokens[userID] = append(tokens[:i], tokens[i+1:]...)
				if len(r.refreshTokens[userID]) == 0 {
					delete(r.refreshTokens, userID)
				}
				return nil
			}
		}
		return nil
	}
}

// DeleteUser removes a user if the provided password matches.
func (r *InMemoryUserRepository) DeleteUser(ctx context.Context, userID string, password []byte) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		r.mutex.Lock()
		defer r.mutex.Unlock()

		user, exists := r.users[userID]
		if !exists {
			return errors.New("user not found")
		}

		if !bytesEqual(user.Passwort(), password) {
			return errors.New("invalid password")
		}

		delete(r.users, userID)
		delete(r.emailToID, user.Email())
		delete(r.refreshTokens, userID)
		return nil
	}
}

// UpdateUser updates an existing user's information.
func (r *InMemoryUserRepository) UpdateUser(ctx context.Context, user *User) (*User, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		r.mutex.Lock()
		defer r.mutex.Unlock()

		existingUser, exists := r.users[user.ID()]
		if !exists {
			return nil, errors.New("user not found")
		}

		if user.Email() != existingUser.Email() {
			if otherID, exists := r.emailToID[user.Email()]; exists && otherID != user.ID() {
				return nil, errors.New("email already exists")
			}
		}

		user.Aktualisert()

		r.users[user.ID()] = user
		if user.Email() != existingUser.Email() {
			delete(r.emailToID, existingUser.Email())
			r.emailToID[user.Email()] = user.ID()
		}
		return user, nil
	}
}

// ChangePassword updates the user's password.
func (r *InMemoryUserRepository) ChangePassword(ctx context.Context, userID string, password []byte) error {
	return errors.New("not implemented")
}

// ChangeEmail changes the email of a user.
func (r *InMemoryUserRepository) ChangeEmail(ctx context.Context, userID, email string) error {
	return errors.New("not implemented")
}

// bytesEqual compares two byte slices for equality.
func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
