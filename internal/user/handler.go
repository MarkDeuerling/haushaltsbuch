package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.com/shingeki-no-kyojin/ymir/config"
	"gitlab.com/shingeki-no-kyojin/ymir/internal/middleware"
	"gitlab.com/shingeki-no-kyojin/ymir/internal/presenter"
	"gitlab.com/shingeki-no-kyojin/ymir/pkg/logger"
)

type usecase interface {
	CreateUser(ctx context.Context, input *CreateInput) error
	LoginUser(ctx context.Context, input *LoginInput) (*LoginOutput, error)
	LogoutUser(ctx context.Context, input *LogoutInput) error
	DeleteUser(ctx context.Context, input *DeleteInput) error
	UpdateUser(ctx context.Context, input *UpdateInput) (*UpdateOutput, error)
}

// Controller is the controller for the user usecase.
type Controller struct {
	log     logger.Logger
	config  *config.Config
	usecase usecase
}

// NewController creates a new controller for the user usecase.
func NewController(log logger.Logger, config *config.Config, usecase usecase) *Controller {
	return &Controller{
		log:     log,
		config:  config,
		usecase: usecase,
	}
}

// CreateUserRequest is a serializable struct for the user creation request body.
type CreateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// CreateUser handles the user creation request.
func (c *Controller) CreateUser(w http.ResponseWriter, r *http.Request) {
	var body CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		c.log.Error(fmt.Sprintf("failed to decode request body. %v", err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	input := &CreateInput{
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Email:     body.Email,
		Password:  body.Password,
	}
	if err := c.usecase.CreateUser(r.Context(), input); err != nil {
		switch err {
		case ErrEmailAlreadyExists:
			c.log.Error("email already exists")
			http.Error(w, "email already exists", http.StatusConflict)
		case ErrInvalideEmail:
			c.log.Error(fmt.Sprintf("email is invalid. %v", err))
			http.Error(w, "email is invalid", http.StatusInternalServerError)
		case ErrLeererVorname:
			c.log.Error(fmt.Sprintf("first name is empty. %v", err))
			http.Error(w, "first name is empty", http.StatusBadRequest)
		case ErrLeererNachname:
			c.log.Error(fmt.Sprintf("last name is empty. %v", err))
			http.Error(w, "last name is empty", http.StatusBadRequest)
		case ErrVornameZuKurz:
			c.log.Error(fmt.Sprintf("first name is too short. %v", err))
			http.Error(w, "first name is too short", http.StatusBadRequest)
		case ErrNachnameZuKurz:
			c.log.Error(fmt.Sprintf("last name is too short. %v", err))
			http.Error(w, "last name is too short", http.StatusBadRequest)
		case ErrVornameZuLang:
			c.log.Error(fmt.Sprintf("first name is too long. %v", err))
			http.Error(w, "first name is too long", http.StatusBadRequest)
		case ErrNachnameZuLang:
			c.log.Error(fmt.Sprintf("last name is too long. %v", err))
			http.Error(w, "last name is too long", http.StatusBadRequest)
		default:
			c.log.Error(fmt.Sprintf("failed to create user. %v", err))
			http.Error(w, "failed to create user", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// LoginUserRequest is a serializable struct for the login request body.
type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginUser handles the user login request.
func (c *Controller) LoginUser(w http.ResponseWriter, r *http.Request) {
	var body LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		c.log.Error(fmt.Sprintf("failed to decode request body. %v", err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	input := &LoginInput{
		Email:    body.Email,
		Password: body.Password,
	}
	tokens, err := c.usecase.LoginUser(r.Context(), input)
	if err != nil {
		c.log.Error(fmt.Sprintf("failed to login user. %v", err))
		http.Error(w, "failed to login user", http.StatusInternalServerError)
		return
	}
	presenter.NewJSONPresenter(w).Successful(tokens)
}

// LogoutUser handles the user logout request.
func (c *Controller) LogoutUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserID).(string)
	if !ok {
		c.log.Error("User ID not found in context")
		http.Error(w, "User ID not found", http.StatusUnauthorized)
		return
	}
	refreshToken, ok := r.Context().Value(middleware.Token).(string)
	if refreshToken == "" {
		c.log.Error("refresh token not found")
		http.Error(w, "refresh token not found", http.StatusUnauthorized)
		return
	}
	input := &LogoutInput{
		UserID:       userID,
		RefreshToken: refreshToken,
	}
	if err := c.usecase.LogoutUser(r.Context(), input); err != nil {
		c.log.Error(fmt.Sprintf("failed to logout user. %v", err))
		http.Error(w, "failed to logout user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DeleteUserRequest is a serializable struct for the delete user request body.
type DeleteUserRequest struct {
	Password string `json:"password"`
}

// DeleteUser handles the user deletion request.
func (c *Controller) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserID).(string)
	if !ok {
		c.log.Error("User ID not found in context")
		http.Error(w, "User ID not found", http.StatusUnauthorized)
		return
	}
	var body DeleteUserRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		c.log.Error(fmt.Sprintf("failed to decode request body. %v", err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	input := &DeleteInput{
		UserID:   userID,
		Password: body.Password,
	}
	if err := c.usecase.DeleteUser(r.Context(), input); err != nil {
		c.log.Error(fmt.Sprintf("failed to delete user. %v", err))
		http.Error(w, "failed to delete user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// UpdateUserRequest is a serializable struct for the user update request body.
type UpdateUserRequest struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// UpdateUserResponse is a serializable struct for the user update response body.
type UpdateUserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// UpdateUser handles the user update request.
func (c *Controller) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserID).(string)
	if !ok {
		c.log.Error("User ID not found in context")
		http.Error(w, "User ID not found", http.StatusUnauthorized)
		return
	}
	var body UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		c.log.Error(fmt.Sprintf("failed to decode request body. %v", err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	input := &UpdateInput{
		UserID:    userID,
		Email:     &body.Email,
		FirstName: &body.FirstName,
		LastName:  &body.LastName,
	}
	output, err := c.usecase.UpdateUser(r.Context(), input)
	if err != nil {
		switch err {
		case ErrUserNotFound:
			c.log.Error("user not found")
			http.Error(w, "user not found", http.StatusNotFound)
		case
			ErrInvalideEmail,
			ErrLeererVorname,
			ErrLeererNachname,
			ErrVornameZuKurz,
			ErrNachnameZuKurz,
			ErrVornameZuLang,
			ErrNachnameZuLang:
			c.log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			c.log.Error(fmt.Sprintf("failed to update user. %v", err))
			http.Error(w, "failed to update user", http.StatusInternalServerError)
		}
		return
	}
	response := &UpdateUserResponse{
		ID:        userID,
		Email:     output.Email,
		FirstName: output.FirstName,
		LastName:  output.LastName,
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	presenter.NewJSONPresenter(w).Successful(response)
}
