package security

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gitlab.com/shingeki-no-kyojin/ymir/internal/middleware"
	"gitlab.com/shingeki-no-kyojin/ymir/internal/user"
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

// GenerateToken Signatur
func (t *JWT) GenerateToken(userID user.ID, ttl time.Duration) (string, error) {
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

// GenerateAccessToken Signatur; add  scope []string later; Not used yet
func (t *JWT) GenerateAccessToken(userID user.ID, ttl time.Duration) (string, error) {
	secret := t.AccessSecret
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": userID,
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(ttl).Unix(),
			// "role":        scope[0],
			// "permissions": scope[1],
		},
	)
	jwt, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return jwt, nil
}

// GenerateRefreshToken Signatur; not used yet
func (t *JWT) GenerateRefreshToken(userID user.ID, ttl time.Duration) (string, error) {
	secret := t.RefreshSecret
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": userID,
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(ttl).Unix(),
			// "jit": id.NewUUID().GenerateUUID(),
		},
	)
	jwt, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return jwt, nil
}

// Parse Signatur parse JWT to extract the claims and validate the token.
func (t *JWT) Parse(tokenString string) (*middleware.Claims, error) {
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
		permissions, _ := claims["permissions"].([]string)
		jit, _ := claims["jit"].(string)

		return &middleware.Claims{
			Sub:         sub,
			Iat:         iat,
			Exp:         exp,
			Role:        role,
			Permissions: permissions,
			Jit:         jit,
		}, nil
	}
	return nil, errors.New("Can not parse claims")
}
