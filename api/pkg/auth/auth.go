package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID     int    `json:"user_id"`
	BusinessID int    `json:"business_id"`
	Role       string `json:"role"`
	jwt.RegisteredClaims
}

type TokenManager struct {
	signingKey string
}

func NewTokenManager(signingKey string) (*TokenManager, error) {
	if signingKey == "" {
		return nil, fmt.Errorf("empty signing key")
	}

	return &TokenManager{signingKey: signingKey}, nil
}

func (m *TokenManager) GenerateToken(userID, businessID int, role string) (string, error) {
	claims := Claims{
		UserID:     userID,
		BusinessID: businessID,
		Role:       role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.signingKey))
}

func (m *TokenManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.signingKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}
