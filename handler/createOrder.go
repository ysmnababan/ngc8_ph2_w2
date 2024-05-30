package handler

import (
	"log"
	"net/http"
	"ngc8/model"
	"ngc8/repo"

	"github.com/gin-gonic/gin"
)

func (h *ProductHandler) CreateProduct(ctx *gin.Context) {
	var p model.ProductDB

	if err := ctx.ShouldBindJSON(&p); err != nil {
		log.Println(err)
		handleError(repo.ErrBindJSON, ctx)
		return
	}

	// validate product
	if p.ID == 1 {
		handleError(repo.ErrParam, ctx)
		return
	}

	newProducts, err := h.Repo.CreateProduct(p)
	if err != nil {
		handleError(err, ctx)
		return
	}
	ctx.JSON(http.StatusCreated, newProducts)
}
