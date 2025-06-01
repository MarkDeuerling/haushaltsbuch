package user

import (
	"context"
	"errors"
	"net/mail"
	"time"
)

var (
	// ErrInvalidEmail is returned when the email is invalid
	ErrInvalidEmail = errors.New("Invalid email format")
	// ErrEmailAlreadyExists is returned when a user is not found
	ErrEmailAlreadyExists = errors.New("Email already exists")
	// ErrEmailTooLong is returned when the email is too long
	ErrEmailTooLong = errors.New("Email too long. Maximum 256 characters")
	// ErrUserNotFound is returned when a user is not found
	ErrUserNotFound = errors.New("User not found")
	// ErrUserNotActive is returned when a user is not active
	ErrUserNotActive = errors.New("User not active")
	// ErrUserAlreadyActivated is returned when a user is already verified
	ErrUserAlreadyActivated = errors.New("User already verified")
	// ErrInvalidPassword is returned when the password is invalid
	ErrInvalidPassword = errors.New("Invalid password")
	// ErrPasswordTooShort is returned when the password is too short
	ErrPasswordTooShort = errors.New("Password too short. At least 8 characters")
	// ErrFirstNameTooLong is returned when the first name is too long
	ErrFirstNameTooLong = errors.New("First name too long. Maximum 50 characters")
	// ErrLastNameTooLong is returned when the last name is too long
	ErrLastNameTooLong = errors.New("Last name too long. Maximum 50 characters")
)

const (
	maxFirstNameLength = 128
	maxLastNameLength  = 128
	maxEmailLength     = 256
	minPasswordLength  = 8
)

type repository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	FindUserByEmail(ctx context.Context, email string) (*User, error)
	FindUserByID(ctx context.Context, id string) (*User, error)
	LogoutUser(ctx context.Context, userID, refreshToken string) error
	DeleteUser(ctx context.Context, userID string, password []byte) error
	UpdateUser(ctx context.Context, user *User) (*User, error)
	ChangePassword(ctx context.Context, userID string, password []byte) error
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
	GenerateAccessToken(userID string, ttl time.Duration) (string, error)
	GenerateRefreshToken(userID string, ttl time.Duration) (string, error)
}

type emailSender interface {
	SendVerificationEmail(to, subject, token string) error
}

// UseCase is the use case for creating a user
type UseCase struct {
	repo                    repository
	uuidGen                 uuidGenerator
	hash                    passwordHasher
	mailer                  emailSender
	tokenGen                tokenGenerator
	accessTokenExpire       time.Duration
	refreshTokenExpire      time.Duration
	verificationTokenExpire time.Duration
}

// NewUseCase creates a new CreateUserUseCase
func NewUseCase(repo repository, uuidGen uuidGenerator, hash passwordHasher, mailer emailSender, tokenGen tokenGenerator, accessTokenExpire, refreshTokenExpire, verificationTokenExpire time.Duration) *UseCase {
	return &UseCase{
		repo:                    repo,
		uuidGen:                 uuidGen,
		hash:                    hash,
		mailer:                  mailer,
		tokenGen:                tokenGen,
		accessTokenExpire:       accessTokenExpire,
		refreshTokenExpire:      refreshTokenExpire,
		verificationTokenExpire: verificationTokenExpire,
	}
}

// CreateInput is the input for the crateUser use case
type CreateInput struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

func (i *CreateInput) validate() error {
	if len(i.FirstName) > maxFirstNameLength {
		return ErrFirstNameTooLong
	}
	if len(i.LastName) > maxLastNameLength {
		return ErrLastNameTooLong
	}
	if len(i.Email) > maxEmailLength {
		return ErrEmailTooLong
	}
	if _, err := mail.ParseAddress(i.Email); err != nil {
		return ErrInvalidEmail
	}
	if len(i.Password) < minPasswordLength {
		return ErrPasswordTooShort
	}
	return nil
}

type userCreator interface {
	CreateUser(ctx context.Context, input *CreateInput) error
}

// CreateUser is the interactor for creating a user
func (c *UseCase) CreateUser(ctx context.Context, input *CreateInput) error {
	if err := input.validate(); err != nil {
		return err
	}

	pwdHash, err := c.hash.GeneratePassword(input.Password)
	if err != nil {
		return err
	}

	id, err := c.uuidGen.GenerateUUID()
	if err != nil {
		return err
	}

	now := time.Now()
	user := NewUser(id, input.FirstName, input.LastName, input.Email, pwdHash, now, now)
	if _, err = c.repo.CreateUser(ctx, user); err != nil {
		return err
	}

	verificationToken, err := c.tokenGen.GenerateAccessToken(user.ID(), time.Second*c.verificationTokenExpire)
	if err != nil {
		return err
	}

	if err := c.mailer.SendVerificationEmail(user.Email(), "Account Verification", verificationToken); err != nil {
		return err
	}

	return nil
}

type userActivator interface {
	ActivateUser(ctx context.Context, userID string) error
}

// ActivateUser is the interactor for verifying a user
func (c *UseCase) ActivateUser(ctx context.Context, userID string) error {
	user, err := c.repo.FindUserByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}
	if user.IstAktiv() {
		return ErrUserAlreadyActivated
	}
	user.Aktiviert()
	if _, err := c.repo.UpdateUser(ctx, user); err != nil {
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

	accessToken, err := c.tokenGen.GenerateAccessToken(user.ID(), time.Second*c.accessTokenExpire)
	if err != nil {
		return nil, err
	}

	refreshToken, err := c.tokenGen.GenerateRefreshToken(user.ID(), time.Second*c.refreshTokenExpire)
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
	UpdateUser(context.Context, *UpdateInput) (*UpdateOutput, error)
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
		user.NeueEmail(*input.Email)
	}

	if input.FirstName != nil {
		user.NeuerVorname(*input.FirstName)
	}

	if input.LastName != nil {
		user.NeuerNachname(*input.LastName)
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
		user.NeuesPasswort(pwdHash)
	}

	if _, err = c.repo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	userOuput := &UpdateOutput{
		Email:     user.Email(),
		FirstName: user.Vorname(),
		LastName:  user.Nachname(),
		CreatedAt: user.ErstelltAm(),
		UpdatedAt: user.AktualisiertAm(),
	}
	return userOuput, nil
}

// ChangePasswordInput is the input for the change password use case
type ChangePasswordInput struct {
	UserID   string
	Password []byte
}

type passwordChanger interface {
	ChangePassword(ctx context.Context, input *ChangePasswordInput) error
}

// ChangePassword is the interactor for changing a user's password
func (c *UseCase) ChangePassword(ctx context.Context, input *ChangePasswordInput) error {
	if len(input.Password) < 8 {
		return ErrPasswordTooShort
	}

	pwdHash, err := c.hash.GeneratePassword(string(input.Password))
	if err != nil {
		return err
	}

	return c.repo.ChangePassword(ctx, input.UserID, pwdHash)
}

// ChangeEmailInput is the input for the change email use case
type ChangeEmailInput struct {
	UserID string
	Email  string
}

type emailChanger interface {
	ChangeEmail(context.Context, *ChangeEmailInput) error
}

// ChangeEmail is the interactor for changing a user's email
func (c *UseCase) ChangeEmail(ctx context.Context, input *ChangeEmailInput) error {
	return c.repo.ChangeEmail(ctx, input.UserID, input.Email)
}

type passwordResetter interface {
	ResetPassword(ctx context.Context, email string) error
}

// ResetPassword is the interactor for resetting a user's password
func (c *UseCase) ResetPassword(ctx context.Context, email string) error {
	if _, err := c.repo.FindUserByEmail(ctx, email); err != nil {
		return ErrUserNotFound
	}

	// TODO: Implement password reset logic, e.g., sending a reset link via email
	return nil
}
