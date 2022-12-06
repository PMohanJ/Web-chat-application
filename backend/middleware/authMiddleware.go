package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pmohanj/web-chat-app/helpers"
)

// Authenticate acts as authorization middleware that receives the client request
// and performs validation of the provided token
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		token := strings.Split(header, " ")[1]

		claims, err := helpers.ValidateToken(token)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenMalformed) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token malformed"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			}
			c.Abort()
			return
		}

		c.Set("_id", claims.ID)
		c.Set("name", claims.Name)
		c.Set("email", claims.Email)
		c.Next()
	}
}
