package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) getAllProducts(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		//Response Body was formed inside getUserId
		return
	}

	listId := c.Param("id")

	products, err := h.services.Product.GetAllProducts(userId, listId)

	if err != nil {
		setErrorResponse(c, err, "getAllProducts")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"products": products,
	})

}

func (h *Handler) updateProduct(c *gin.Context) {

}

func (h *Handler) deleteProduct(c *gin.Context) {

}

func (h *Handler) createProduct(c *gin.Context) {

}
