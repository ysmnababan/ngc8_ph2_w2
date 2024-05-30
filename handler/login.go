package handler

import (
	"fmt"
	"log"
	"net/http"
	"ngc8/model"
	"ngc8/repo"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

func generateToken(u model.User) (string, error) {
	// create payload
	payload := jwt.MapClaims{
		"id":    u.ID,
		"email": u.Email,
		"name":  u.Name,
	}

	// define the method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	err := godotenv.Load()
	if err != nil {
		return "", fmt.Errorf("unable to get .env")
	}

	// get token string
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", fmt.Errorf("unable to get token String")
	}

	return tokenString, nil
}
func (h *UserHandler) Login(c *gin.Context) {
	var u model.User
	err := c.ShouldBindJSON(&u)
	if err != nil {
		handleError(err, c)
		return
	}

	if u.Email == "" || u.Pwd == "" {
		handleError(repo.ErrParam, c)
		return
	}

	newUser, err := h.Repo.Login(u)
	if err != nil {
		handleError(err, c)
		return
	}

	token, err := generateToken(newUser)
	if err != nil {
		log.Println("unable to generate token:", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "unable to generate token",
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "login succeed",
			"token":   token,
		},
	)
}
