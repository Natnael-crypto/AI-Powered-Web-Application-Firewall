package utils

import (
	"backend/internal/config"
	"backend/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// GenerateJWT generates a new JWT token
func GenerateJWT(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.UserID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.JWTSecretKey)
}

// ParseJWT parses the JWT token and returns the claims
func ParseJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("invalid signing method", jwt.ValidationErrorSignatureInvalid)
		}
		return config.JWTSecretKey, nil
	})
}
