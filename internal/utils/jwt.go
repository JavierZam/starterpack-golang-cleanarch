package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID   string `json:"user_id"`
	TenantID string `json:"tenant_id"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID, tenantID, role string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	expiresInMinutes := os.Getenv("JWT_EXPIRES_IN_MINUTES")

	if jwtSecret == "" || expiresInMinutes == "" {
		return "", fmt.Errorf("JWT_SECRET or JWT_EXPIRES_IN_MINUTES not set")
	}

	expMinutes, err := time.ParseDuration(expiresInMinutes + "m")
	if err != nil {
		return "", fmt.Errorf("invalid JWT_EXPIRES_IN_MINUTES format: %w", err)
	}

	claims := Claims{
		UserID:   userID,
		TenantID: tenantID,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expMinutes)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "go-starterpack",
			Subject:   userID,
			Audience:  jwt.ClaimStrings{"go-users"},
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}
	return tokenString, nil
}

func GenerateRefreshToken(userID string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	refreshExpiresInHours := os.Getenv("REFRESH_TOKEN_EXPIRES_IN_HOURS")

	if jwtSecret == "" || refreshExpiresInHours == "" {
		return "", fmt.Errorf("JWT_SECRET or REFRESH_TOKEN_EXPIRES_IN_HOURS not set")
	}

	expHours, err := time.ParseDuration(refreshExpiresInHours + "h")
	if err != nil {
		return "", fmt.Errorf("invalid REFRESH_TOKEN_EXPIRES_IN_HOURS format: %w", err)
	}

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expHours)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "go-starterpack",
			Subject:   userID,
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}
	return tokenString, nil
}

func ValidateToken(tokenString string) (*Claims, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET not set")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
