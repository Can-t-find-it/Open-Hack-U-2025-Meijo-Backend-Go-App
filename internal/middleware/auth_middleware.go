package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "認証トークンがありません"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "トークンの形式が不正です"})
			c.Abort()
			return
		}
		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			secret := os.Getenv("TOKEN_SECRET")
			if secret == "" {
				secret = "develop_secret_key"
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "無効なトークンです: " + err.Error()})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if ok && token.Valid {
			if UserIDStr, ok := claims["user_id"].(string); ok {
				c.Set("userID", UserIDStr)
			}else {
				// stringで取れなかった場合の保険 (古いトークン対策など)
				// floatで来る可能性も考慮するならここ分岐ですが、
				// DBリセット後はstring確定なのでこれでOK
				fmt.Println("Middleware Error: user_id is not a string")
			}
		}

		c.Next()
	}
}
