package middleware

import (
	"backend/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// AuthRequired is a middleware to validate the JWT token
func AuthRequired(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization token is required"})
		c.Abort()
		return
	}

	token, err := utils.ParseJWT(tokenString)
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		c.Abort()
		return
	}

	// You can set user info from token claims if needed
	claims := token.Claims.(jwt.MapClaims)
	c.Set("user_id", claims["user_id"])
	c.Set("role", claims["role"])

	c.Next()
}
