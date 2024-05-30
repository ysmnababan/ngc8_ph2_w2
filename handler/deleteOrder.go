package handler

import (
	"net/http"
	"ngc8/repo"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	param_id := c.Param("id")
	id, err := strconv.Atoi(param_id)
	if err != nil || id <= 0 {
		handleError(repo.ErrInvalidId, c)
		return
	}

	err = h.Repo.DeleteProduct(id)
	if err != nil {
		handleError(err, c)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "data deleted successfuly",
		},
	)
}
