package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sharybkin/grocerylist-golang/internal/model"
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

func (h *Handler) addProduct(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		//Response Body was formed inside getUserId
		return
	}

	listId := c.Param("id")

	var product model.Product

	if err := c.BindJSON(&product); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), "addProduct")
		return
	}

	productId, err := h.services.Product.AddProduct(userId, listId, product)

	if err != nil {
		setErrorResponse(c, err, "addProduct")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"productId": productId,
	})
}
