package middleware

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get token from header
		tokenString := ctx.Request.Header.Get("auth")
		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, "Unauthorized")
			ctx.Abort()
			return
		}

		//get secret code from .env
		err := godotenv.Load()
		if err != nil {
			log.Println("error while getting secret key")
			ctx.Abort()
			return
		}

		// get token string
		key := os.Getenv("SECRET_KEY")

		// get token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return []byte(key), nil
		})

		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, "unauthorized")
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
