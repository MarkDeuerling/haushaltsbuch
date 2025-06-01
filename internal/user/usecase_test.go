package user_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/shingeki-no-kyojin/ymir/internal/user"
)

type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) CreateUser(ctx context.Context, u *user.User) (*user.User, error) {
	args := m.Called(ctx, u)
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *mockUserRepository) FindUserByEmail(ctx context.Context, email string) (*user.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *mockUserRepository) FindUserByID(ctx context.Context, id string) (*user.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *mockUserRepository) LogoutUser(ctx context.Context, userID, refreshToken string) error {
	args := m.Called(ctx, userID, refreshToken)
	return args.Error(0)
}

func (m *mockUserRepository) DeleteUser(ctx context.Context, userID string, password []byte) error {
	args := m.Called(ctx, userID, password)
	return args.Error(1)
}

func (m *mockUserRepository) UpdateUser(ctx context.Context, u *user.User) (*user.User, error) {
	args := m.Called(ctx, u)
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *mockUserRepository) ChangePassword(ctx context.Context, userID string, password []byte) error {
	args := m.Called(ctx, userID, password)
	return args.Error(0)
}

func (m *mockUserRepository) ChangeEmail(ctx context.Context, userID, email string) error {
	args := m.Called(ctx, userID, email)
	return args.Error(0)
}

type mockUUIDGenerator struct {
	mock.Mock
}

func (m *mockUUIDGenerator) GenerateUUID() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

type mockPasswordHasher struct {
	mock.Mock
}

func (m *mockPasswordHasher) GeneratePassword(password string) ([]byte, error) {
	args := m.Called(password)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *mockPasswordHasher) ValidatePassword(hash []byte, clearword string) error {
	args := m.Called(hash, clearword)
	return args.Error(0)
}

type mockMailer struct {
	mock.Mock
}

func (m *mockMailer) SendVerificationEmail(to, subject, body string) error {
	args := m.Called(to, subject, body)
	return args.Error(0)
}

type mockTokenGenerator struct {
	mock.Mock
}

func (m *mockTokenGenerator) GenerateAccessToken(userID string, ttl time.Duration) (string, error) {
	args := m.Called(userID, ttl)
	return args.String(0), args.Error(1)
}

func (m *mockTokenGenerator) GenerateRefreshToken(userID string, ttl time.Duration) (string, error) {
	args := m.Called(userID, ttl)
	return args.String(0), args.Error(1)
}

func TestCreateUser(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name       string
		input      *user.CreateInput
		setupMocks func(*mockUserRepository, *mockUUIDGenerator, *mockPasswordHasher, *mockMailer, *mockTokenGenerator)
		expectErr  error
	}{
		{
			name: "Valid User Input Successful",
			input: &user.CreateInput{
				FirstName: "Max",
				LastName:  "Mustermann",
				Email:     "max.mustermann@gmail.de",
				Password:  "password",
			},
			setupMocks: func(repo *mockUserRepository, uuidGen *mockUUIDGenerator, hasher *mockPasswordHasher, mailer *mockMailer, tokenGen *mockTokenGenerator) {
				uuidGen.On("GenerateUUID").Return("12345", nil)
				hasher.On("GeneratePassword", "password").Return([]byte("password"), nil)
				u := user.NewUser("123", "Max", "Mustermann", "max.mustermann@gmx.de", []byte("password"), time.Now(), time.Now())
				u.Aktiviert()
				repo.On("CreateUser", ctx, mock.AnythingOfType("*user.User")).Return(u, nil)
				tokenGen.On("GenerateAccessToken", "12345", mock.Anything).Return("accesstokensecret12345", nil)
				tokenGen.On("GenerateRefreshToken", "12345", mock.Anything).Return("refreshtokensecret12345", nil)
			},
			expectErr: nil,
		},
		{
			name: "Invalid User Input - Empty firstname",
			input: &user.CreateInput{
				FirstName: "",
				LastName:  "Mustermann",
				Email:     "max.mustermann@gmail.de",
				Password:  "password",
			},
			setupMocks: func(repo *mockUserRepository, uuidGen *mockUUIDGenerator, hasher *mockPasswordHasher, mailer *mockMailer, tokenGen *mockTokenGenerator) {
				uuidGen.On("GenerateUUID").Return("12345", nil)
				hasher.On("GeneratePassword", "password").Return([]byte("password"), nil)
				u := user.NewUser("123", "Max", "Mustermann", "max.mustermann@gmx.de", []byte("password"), time.Now(), time.Now())
				u.Aktiviert()
				repo.On("CreateUser", ctx, mock.AnythingOfType("*user.User")).Return(u, nil)
				tokenGen.On("GenerateAccessToken", "12345", mock.Anything).Return("accesstokensecret12345", nil)
				tokenGen.On("GenerateRefreshToken", "12345", mock.Anything).Return("refreshtokensecret12345", nil)
			},
			expectErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mockUserRepository)
			uuidGen := new(mockUUIDGenerator)
			hasher := new(mockPasswordHasher)
			tokenGen := new(mockTokenGenerator)
			mailer := new(mockMailer)
			uc := user.NewUseCase(repo, uuidGen, hasher, mailer, tokenGen, time.Millisecond*99999, time.Millisecond*99999, time.Millisecond*99999) // Mock dependencies as needed

			tt.setupMocks(repo, uuidGen, hasher, mailer, tokenGen)

			err := uc.CreateUser(ctx, tt.input)

			if tt.expectErr != nil {
				assert.ErrorIs(t, err, tt.expectErr)
			}

			if tt.expectErr == nil {
				assert.NoError(t, err)
			}
		})
	}
}
