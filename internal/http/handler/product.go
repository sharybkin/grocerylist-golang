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
	userId, err := getUserId(c)
	if err != nil {
		//Response Body was formed inside getUserId
		return
	}

	listId := c.Param("id")
	productId := c.Param("product_id")

	var product model.Product
	if err := c.BindJSON(&product); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), "updateProduct")
		return
	}

	product.Id = productId

	if err := h.services.Product.UpdateProduct(userId, listId, product); err != nil {
		setErrorResponse(c, err, "updateProduct")
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) deleteProduct(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		//Response Body was formed inside getUserId
		return
	}

	listId := c.Param("id")
	productId := c.Param("product_id")

	if err := h.services.Product.DeleteProduct(userId, listId, productId); err != nil {
		setErrorResponse(c, err, "deleteProduct")
		return
	}

	c.Status(http.StatusOK)
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
