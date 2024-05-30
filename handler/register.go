package handler

import (
	"net/http"
	"ngc8/model"
	"ngc8/repo"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Repo repo.UserRepo
}

func (h *UserHandler) Register(ctx *gin.Context) {
	var u model.User

	if err := ctx.ShouldBindJSON(&u); err != nil {
		handleError(repo.ErrBindJSON, ctx)
		return
	}

	// validate user
	if u.Name == "" || u.Email == "" || u.Pwd == "" {
		handleError(repo.ErrParam, ctx)
		return
	}

	newUser, err := h.Repo.Register(u)
	if err != nil {
		handleError(err, ctx)
		return
	}

	ctx.JSON(
		http.StatusCreated,
		gin.H{
			"message": "new user added",
			"user":    newUser,
		},
	)
}
