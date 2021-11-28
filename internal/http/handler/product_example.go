package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) getAllProductExamples(c *gin.Context) {

	productExamples := h.services.ProductExample.GetProductExamples()

	c.JSON(http.StatusOK, productExamples)
}
