package auth

import (
	"errors"
	"fmt"
	"time"

	"uuid"

	"github.com/golang-jwt/jwt/v5"
)

// AccessTokenTTL is the lifetime of an access token. Kept short-lived since
// there is no refresh-token flow yet.
const AccessTokenTTL = 15 * time.Minute

// Claims are the JWT claims carried by an access token.
type Claims struct {
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

// GenerateAccessToken creates a signed HS256 access token for the given user and permissions.
func GenerateAccessToken(secret []byte, userID uuid.UUID, permissions []string) (string, error) {
	now := time.Now()
	claims := Claims{
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(AccessTokenTTL)),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("sign access token: %w", err)
	}

	return token, nil
}

// ParseAccessToken validates the token signature and expiration and returns its claims.
func ParseAccessToken(secret []byte, tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse access token: %w", err)
	}
	if !token.Valid {
		return nil, errors.New("invalid access token")
	}

	return claims, nil
}
