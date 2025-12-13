package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	// JWT secret keys from environment variables
	accessTokenSecret  = []byte(getEnv("JWT_ACCESS_SECRET"))
	refreshTokenSecret = []byte(getEnv("JWT_REFRESH_SECRET"))

	// Token durations
	accessTokenDuration  = 15 * time.Minute   // 15 minutes
	refreshTokenDuration = 7 * 24 * time.Hour // 7 days
)

// JWTClaims represents the claims stored in JWT tokens
type JWTClaims struct {
	UserID   uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
	jwt.RegisteredClaims
}

// GenerateAccessToken creates a short-lived access token
func GenerateAccessToken(userID uuid.UUID, email, username, role string) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		Email:    email,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(accessTokenSecret)
}

// GenerateRefreshToken creates a long-lived refresh token
func GenerateRefreshToken(userID uuid.UUID) (string, time.Time, error) {
	expiresAt := time.Now().Add(refreshTokenDuration)

	claims := JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(refreshTokenSecret)
	return tokenString, expiresAt, err
}

// ValidateAccessToken validates and parses an access token
func ValidateAccessToken(tokenString string) (*JWTClaims, error) {
	return validateToken(tokenString, accessTokenSecret)
}

// ValidateRefreshToken validates and parses a refresh token
func ValidateRefreshToken(tokenString string) (*JWTClaims, error) {
	return validateToken(tokenString, refreshTokenSecret)
}

// validateToken is a helper function to validate JWT tokens
func validateToken(tokenString string, secret []byte) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		// Additional validation: check if token is expired
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			return nil, errors.New("token has expired")
		}
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// getEnv retrieves environment variable or returns empty string if not set
func getEnv(key string) string {
	return os.Getenv(key)
}

// GetAccessTokenDuration returns the access token duration (for cookie max age)
func GetAccessTokenDuration() time.Duration {
	return accessTokenDuration
}

// GetRefreshTokenDuration returns the refresh token duration (for cookie max age)
func GetRefreshTokenDuration() time.Duration {
	return refreshTokenDuration
}
