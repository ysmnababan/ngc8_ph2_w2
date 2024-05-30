package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"ngc8/repo"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	Repo repo.ProductRepo
}

func handleError(err error, ctx *gin.Context) {
	log.Println(err)
	status := http.StatusOK
	message := ""
	switch {
	case errors.Is(err, repo.ErrQuery):
		fallthrough
	case errors.Is(err, repo.ErrScan):
		fallthrough
	case errors.Is(err, repo.ErrRowsAffected):
		fallthrough
	case errors.Is(err, repo.ErrLastInsertId):
		fallthrough
	case errors.Is(err, repo.ErrNoAffectedRow):
		status = http.StatusInternalServerError
		message = "Internal Server Error"
	case errors.Is(err, repo.ErrNoRows):
		status = http.StatusNotFound
		message = "No row found"
	case errors.Is(err, repo.ErrParam):
		status = http.StatusBadRequest
		message = "error or missing param"
	case errors.Is(err, repo.ErrBindJSON):
		status = http.StatusBadRequest
		message = "Bad request"
	case errors.Is(err, repo.ErrInvalidId):
		status = http.StatusBadRequest
		message = "Invalid ID"
	case errors.Is(err, repo.ErrCredential):
		status = http.StatusBadRequest
		message = "Incorrect credential"
	case errors.Is(err, repo.ErrUserExists):
		status = http.StatusBadRequest
		message = "User Already Exists"
	case errors.Is(err, repo.ErrNoUpdate):
		status = http.StatusBadRequest
		message = "Data is the same"
	default:
		status = http.StatusInternalServerError
		message = "Unknown error"
	}

	ctx.JSON(
		status,
		gin.H{
			"message": message,
		},
	)
}

func (h *ProductHandler) GetProducts(ctx *gin.Context) {
	product, err := h.Repo.GetAllProducts()
	if err != nil {
		handleError(err, ctx)
		return
	}

	fmt.Println(product)
	ctx.JSON(http.StatusOK, product)
}
