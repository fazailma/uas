package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"time"

	"UAS/app/models"

	"github.com/golang-jwt/jwt/v5"
)

// HashPassword hashes a password using SHA256
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// VerifyPassword verifies if password matches hash
func VerifyPassword(password string, hash string) bool {
	return HashPassword(password) == hash
}

// GenerateJWT generates a short-lived JWT access token (1 hour)
func GenerateJWT(user *models.User, role models.Role, permissions []string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key-change-in-production"
	}

	claims := jwt.MapClaims{
		"user_id":     user.ID,
		"username":    user.Username,
		"email":       user.Email,
		"role":        role.Name,
		"permissions": permissions,
		"exp":         time.Now().Add(time.Hour * 1).Unix(), // 1 hour expiry
		"iat":         time.Now().Unix(),
		"type":        "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// GenerateRefreshToken generates a refresh JWT token (7 days expiry)
func GenerateRefreshToken(user *models.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key-change-in-production"
	}

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days expiry
		"iat":      time.Now().Unix(),
		"type":     "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateRefreshToken validates and extracts claims from refresh token
func ValidateRefreshToken(tokenString string) (*models.User, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key-change-in-production"
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	// Check if it's a refresh token
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return nil, jwt.ErrTokenInvalidClaims
	}

	user := &models.User{
		ID:       claims["user_id"].(string),
		Username: claims["username"].(string),
	}

	return user, nil
}
