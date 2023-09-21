package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func JWTAuthenStore() gin.HandlerFunc {
	return func(c *gin.Context) {
		hmacSampleSecretStore := []byte(os.Getenv("JWT_SECRET_KEY_STORE"))
		header := c.Request.Header.Get("Authorization")
		tokenString := strings.Replace(header, "Bearer ", "", 1)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return hmacSampleSecretStore, nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("storeId", claims["storeId"])

		} else {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"status": "forbidden", "message": err.Error()})
			return
		}
		c.Next()
	}
}
