package middleware

import (
	"crescendo-api/security"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	ContextUserID    = "userID"
	ContextUsername  = "username"
	ContextTokenType = "tokenType"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "missing authorization header",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "invalid authorization header format",
			})
			return
		}

		tokenString := parts[1]

		claims, err := security.ValidateLoginToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "invalid or expired token",
			})
			return
		}

		if claims.Type != "auth" {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "invalid token type",
			})
			return
		}

		c.Next()
	}
}
