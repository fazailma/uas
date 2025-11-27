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

// GenerateJWT generates a JWT token
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
		"exp":         time.Now().Add(time.Hour * 24).Unix(),
		"iat":         time.Now().Unix(),
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
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
