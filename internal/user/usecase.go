package user

import (
	"context"
	"errors"
	"time"
)

var (
	// ErrEmailAlreadyExists is returned when a user is not found
	ErrEmailAlreadyExists = errors.New("Email already exists")
	// ErrUserNotFound is returned when a user is not found
	ErrUserNotFound = errors.New("User not found")
	// ErrUserNotActive is returned when a user is not active
	ErrUserNotActive = errors.New("User not active")
	// ErrInvalidPassword is returned when the password is invalid
	ErrInvalidPassword = errors.New("Invalid password")
	// ErrEmptyPassword is returned when the email is empty
	ErrEmptyPassword = errors.New("Empty password")
	// ErrPasswordTooShort is returned when the password is too short
	ErrPasswordTooShort = errors.New("Password too short. At least 8 characters")
	// ErrPasswordTooLong is returned when the password is too long
	ErrPasswordTooLong = errors.New("Password too long. Maximum 100 characters")
)

type repository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	FindUserByEmail(ctx context.Context, email string) (*User, error)
	FindUserByID(ctx context.Context, id string) (*User, error)
	LogoutUser(ctx context.Context, userID, refreshToken string) error
	DeleteUser(ctx context.Context, userID string, password []byte) error
	UpdateUser(ctx context.Context, user *User) (*User, error)
	ResetPassword(ctx context.Context, userID string, password []byte) error
	ChangeEmail(ctx context.Context, userID, email string) error
}

type uuidGenerator interface {
	GenerateUUID() (string, error)
}

type passwordHasher interface {
	GeneratePassword(password string) ([]byte, error)
	ValidatePassword(hash []byte, clearword string) error
}

type tokenGenerator interface {
	GenerateToken(userID string, ttl time.Duration) (string, error)
}

// UseCase is the use case for creating a user
type UseCase struct {
	repo               repository
	uuidGen            uuidGenerator
	hash               passwordHasher
	tokenGen           tokenGenerator
	accessTokenExpire  time.Duration
	refreshTokenExpire time.Duration
}

// NewUseCase creates a new CreateUserUseCase
func NewUseCase(repo repository, uuidGen uuidGenerator, hash passwordHasher, tokenGen tokenGenerator, accessTokenExpire, refreshTokenExpire time.Duration) *UseCase {
	return &UseCase{
		repo:               repo,
		uuidGen:            uuidGen,
		hash:               hash,
		tokenGen:           tokenGen,
		accessTokenExpire:  accessTokenExpire,
		refreshTokenExpire: refreshTokenExpire,
	}
}

// CreateInput is the input for the crateUser use case
type CreateInput struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

type userCreator interface {
	CreateUser(ctx context.Context, input *CreateInput) error
}

// CreateUser is the interactor for creating a user
func (c *UseCase) CreateUser(ctx context.Context, input *CreateInput) error {
	// Need to check? Move to user entity?
	if input.Password == "" {
		return ErrEmptyPassword
	}
	if len(input.Password) < 8 {
		return ErrPasswordTooShort
	}
	if len(input.Password) > 100 {
		return ErrPasswordTooLong
	}

	pwdHash, err := c.hash.GeneratePassword(input.Password)
	if err != nil {
		return err
	}
	id, err := c.uuidGen.GenerateUUID()
	if err != nil {
		return err
	}

	user, err := NewUser(id, input.FirstName, input.LastName, "", input.Email, pwdHash)
	if err != nil {
		return err
	}
	// Test only! should send email to verify the user
	user.Aktiviert()

	if _, err = c.repo.CreateUser(ctx, user); err != nil {
		if errors.Is(err, ErrEmailAlreadyExists) {
			return ErrEmailAlreadyExists
		}
		return err
	}
	return nil
}

// LoginInput is the input for the login use case
type LoginInput struct {
	Email    string
	Password string
}

// LoginOutput is the output for the login use case
type LoginOutput struct {
	AccessToken  string
	RefreshToken string
}

type userAuthenticator interface {
	LoginUser(ctx context.Context, input *LoginInput) (*LoginOutput, error)
}

// LoginUser is the interactor for logging in a user
func (c *UseCase) LoginUser(ctx context.Context, input *LoginInput) (*LoginOutput, error) {
	user, err := c.repo.FindUserByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	if !user.IstAktiv() {
		return nil, ErrUserNotActive
	}

	if err := c.hash.ValidatePassword(user.Passwort(), input.Password); err != nil {
		return nil, ErrInvalidPassword
	}

	accessToken, err := c.tokenGen.GenerateToken(user.ID(), time.Second*c.accessTokenExpire)
	if err != nil {
		return nil, err
	}

	refreshToken, err := c.tokenGen.GenerateToken(user.ID(), time.Second*c.refreshTokenExpire)
	if err != nil {
		return nil, err
	}

	return &LoginOutput{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

// LogoutInput is the input for the logout use case
type LogoutInput struct {
	UserID       string
	RefreshToken string
}

type userSessionCloser interface {
	LogoutUser(ctx context.Context, input *LogoutInput) error
}

// LogoutUser is the interactor for logging out a user
func (c *UseCase) LogoutUser(ctx context.Context, input *LogoutInput) error {
	return c.repo.LogoutUser(ctx, input.UserID, input.RefreshToken)
}

// DeleteInput is the input for the delete user use case
type DeleteInput struct {
	UserID   string
	Password string
}

type userRemover interface {
	DeleteUser(ctx context.Context, input *DeleteInput) error
}

// DeleteUser is the interactor for deleting a user
func (c *UseCase) DeleteUser(ctx context.Context, input *DeleteInput) error {
	return c.repo.DeleteUser(ctx, input.UserID, []byte(input.Password))
}

// UpdateInput is the input for the update user use case
type UpdateInput struct {
	UserID          string
	Email           *string
	FirstName       *string
	LastName        *string
	CurrentPassword *string
	NewPassword     *string
	ConfirmPassword *string
}

// UpdateOutput is the output for the update user use case
type UpdateOutput struct {
	Email     string
	FirstName string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type userUpdater interface {
	UpdateUser(ctx context.Context, input *UpdateInput) (*UpdateOutput, error)
}

// UpdateUser is the interactor for updating a user
func (c *UseCase) UpdateUser(ctx context.Context, input *UpdateInput) (*UpdateOutput, error) {
	user, err := c.repo.FindUserByID(ctx, input.UserID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// TODO: own usecase for email
	if input.Email != nil {
		err := user.AktualisiereEmail(*input.Email)
		if err != ErrInvalideEmail {
			return nil, ErrInvalideEmail
		}
	}

	if input.FirstName != nil {
		err := user.AktualisiereVorname(*input.FirstName)
		if err != nil {
			switch err {
			case ErrLeererVorname:
				return nil, ErrLeererVorname
			case ErrVornameZuKurz:
				return nil, ErrVornameZuKurz
			case ErrVornameZuLang:
				return nil, ErrVornameZuLang
			default:
				return nil, err
			}
		}
	}

	if input.LastName != nil {
		err := user.AktualisiereNachname(*input.LastName)
		if err != nil {
			switch err {
			case ErrLeererNachname:
				return nil, ErrLeererNachname
			case ErrNachnameZuKurz:
				return nil, ErrNachnameZuKurz
			case ErrNachnameZuLang:
				return nil, ErrNachnameZuLang
			default:
				return nil, err
			}
		}
	}

	// TODO: own usecase for password
	if input.NewPassword != nil {
		if err := c.hash.ValidatePassword(user.Passwort(), *input.CurrentPassword); err != nil {
			return nil, ErrInvalidPassword
		}
		pwdHash, err := c.hash.GeneratePassword(*input.NewPassword)
		if err != nil {
			return nil, err
		}
		user.AktualisierePasswort([]byte(pwdHash))
	}

	if _, err = c.repo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	userOuput := &UpdateOutput{
		user.Email(),
		user.Vorname(),
		user.Nachname(),
		user.ErstelltAm(),
		user.AktualisiertAm(),
	}
	return userOuput, nil
}
