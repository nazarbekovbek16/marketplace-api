// internal/api/middleware/auth_middleware.go

package middleware

import (
	"fmt"
	"marketplace-api/internal/api/handlers"
	"marketplace-api/internal/models"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware middleware for JWT authentication
func AuthMiddleware(jwtSecret string, role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}

		tokenString := parts[1]

		token, err := jwt.ParseWithClaims(tokenString, &handlers.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})
		claims, ok := token.Claims.(*handlers.CustomClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		user := models.User{
			ID:    claims.ID,
			Email: claims.Email,
			Role:  claims.Role,
		}
		c.Set("user", user)
		c.Set("user_id", user.ID)
		if err != nil {
			fmt.Println(err.Error(), jwtSecret)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "failed to parse token"})
			return
		}
		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		if role != user.Role && role != "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Not Permitted"})
			return
		}

		c.Next()
	}
}

func AuthorizeRoleMiddleware(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			fmt.Println(user)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		userRole := user.(*models.User).Role
		if userRole != role {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
			return
		}

		c.Next()
	}
}

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
