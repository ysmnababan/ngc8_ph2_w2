package handler

import (
	"net/http"
	"ngc8/repo"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *ProductHandler) GetProductById(ctx *gin.Context) {
	param_id := ctx.Param("id")

	id, err := strconv.Atoi(param_id)
	if err != nil || id <= 0 {
		handleError(repo.ErrInvalidId, ctx)
		return
	}

	product, err := h.Repo.GetProductById(uint(id))
	if err != nil {
		handleError(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, product)
}
