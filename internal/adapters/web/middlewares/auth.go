package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"wappi/internal/platform/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// Context key for user data (stores the full payload as a map, similar to assistant-ia-api)
const ContextKeyUser = "user"

// CustomClaims matches the auth-api-be JWT structure
type CustomClaims struct {
	UserID       string `json:"user_id"`
	TokenVersion uint   `json:"token_version"`
	jwt.StandardClaims
}

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(tokenStr, jwtSecret string) (*CustomClaims, error) {
	if jwtSecret == "" {
		return nil, errors.New("JWT secret not configured")
	}

	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

// AuthMiddleware validates JWT tokens from Authorization header
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		cfg := config.GetInstance()

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "Missing authorization header",
			})
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "Invalid authorization format. Use: Bearer <token>",
			})
			return
		}

		tokenStr := tokenParts[1]

		claims, err := ValidateToken(tokenStr, cfg.JWTSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "Invalid or expired token",
			})
			return
		}

		userPayload := map[string]interface{}{
			"user_id":       claims.UserID,
			"token_version": claims.TokenVersion,
		}

		c.Set(ContextKeyUser, userPayload)

		c.Next()
	}
}

// OptionalAuthMiddleware tries to validate token but doesn't block if missing
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := config.GetInstance()

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			c.Next()
			return
		}

		tokenStr := tokenParts[1]

		claims, err := ValidateToken(tokenStr, cfg.JWTSecret)
		if err != nil {
			c.Next()
			return
		}

		userPayload := map[string]interface{}{
			"user_id":       claims.UserID,
			"token_version": claims.TokenVersion,
		}

		c.Set(ContextKeyUser, userPayload)

		c.Next()
	}
}

// GetUserIDFromContext extracts user ID from gin context (similar to assistant-ia-api pattern)
func GetUserIDFromContext(c *gin.Context) (string, bool) {
	user, exists := c.Get(ContextKeyUser)
	if !exists {
		return "", false
	}

	userMap, ok := user.(map[string]interface{})
	if !ok {
		return "", false
	}

	userID, ok := userMap["user_id"].(string)
	if !ok {
		return "", false
	}

	return userID, true
}
