package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWT is a struct that holds the JWT secret keys.
type JWT struct {
	AccessSecret  []byte
	RefreshSecret []byte
}

// NewJWT creates a new Token.
func NewJWT(accessSecret, refreshSecret string) *JWT {
	return &JWT{
		AccessSecret:  []byte(accessSecret),
		RefreshSecret: []byte(refreshSecret),
	}
}

// GenerateAccessToken Signatur
func (t *JWT) GenerateAccessToken(userID string, ttl time.Duration) (string, error) {
	secret := t.AccessSecret
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": userID,
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(ttl).Unix(),
		},
	)
	jwt, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return jwt, nil
}

// GenerateRefreshToken Signatur
func (t *JWT) GenerateRefreshToken(userID string, ttl time.Duration) (string, error) {
	secret := t.RefreshSecret
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub":  userID,
			"iat":  time.Now().Unix(),
			"exp":  time.Now().Add(ttl).Unix(),
			"type": "refresh",
		},
	)
	jwt, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return jwt, nil
}

// Parse Signatur parse JWT to extract the claims and validate the token.
func (t *JWT) Parse(tokenString string) (*Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return t.AccessSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("Token Invalid")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		sub, _ := claims["sub"].(string)
		iat, _ := claims["iat"].(time.Time)
		exp, _ := claims["exp"].(time.Time)
		role, _ := claims["role"].(string)
		// permissions, _ := claims["permissions"].([]string)
		// jit, _ := claims["jit"].(string)

		return &Claims{
			Sub:  sub,
			Iat:  iat,
			Exp:  exp,
			Role: role,
			// Permissions: permissions,
			// Jit:         jit,
		}, nil
	}
	return nil, errors.New("Can not parse claims")
}
