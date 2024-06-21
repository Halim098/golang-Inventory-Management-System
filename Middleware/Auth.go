package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var privateKey = []byte(os.Getenv("JWT_PRIVATE_KEY"))

func Auth(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		// remove the Bearer prefix
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return privateKey, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if claims["role"] == role {
				c.Set("user_id", uint(claims["id"].(float64)))
				c.Next()
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized sepsial"})
				c.Abort()
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
		}
	}
}

func AdminAuth() gin.HandlerFunc {
	return Auth("admin")
}

func UserAuth() gin.HandlerFunc {
	return Auth("user")
}
