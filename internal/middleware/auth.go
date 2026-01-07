package middleware

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"pwa-backend/internal/repositories"
)

func AuthMiddleware(jwtSecret string, sessionRepo *repositories.UserSessionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		var tokenString string
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			tokenString = authHeader
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userID := claims["user_id"]
			sessionID := claims["sid"]

			// Validate session if present
			if sessionID != nil {
				isValid, err := sessionRepo.IsSessionValid(sessionID.(string))
				if err != nil && err != sql.ErrNoRows {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate session"})
					c.Abort()
					return
				}

				if !isValid {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Session invalid or revoked"})
					c.Abort()
					return
				}

				c.Set("session_id", sessionID.(string))
			}

			c.Set("user_id", userID)
			c.Set("username", claims["username"])
			c.Set("role", claims["role"])
		}

		c.Next()
	}
}