package auth

import (
	"context"
	"net/http"
	"strings"
	"time"
)

type contextkey string

var (
	// UserID is the key for the user ID in the context
	UserID contextkey = "userID"
	// Token is the key for the token in the context
	Token contextkey = "token"
)

// Claims ...
type Claims struct {
	Sub         string
	Iat         time.Time
	Exp         time.Time
	Role        string
	Permissions []string
	Jit         string
}

type tokenParser interface {
	Parse(tokenString string) (*Claims, error)
}

// Authorization is a middleware that checks the authorization token
type Authorization struct {
	tokenAuth tokenParser
}

// NewAuthorization creates a new Authorization middleware
func NewAuthorization(tokenAuth tokenParser) *Authorization {
	return &Authorization{tokenAuth: tokenAuth}
}

// Authorize verifies the token and extracts the user ID from it
func (a *Authorization) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqToken := r.Header.Get("Authorization")
		token := strings.TrimPrefix(reqToken, "Bearer ")
		claim, err := a.tokenAuth.Parse(token)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID := claim.Sub
		ctx := context.WithValue(r.Context(), UserID, userID)
		ctx = context.WithValue(ctx, Token, token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
